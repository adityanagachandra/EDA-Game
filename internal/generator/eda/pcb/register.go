// Package pcb generates problems on printed-circuit-board design automation: the board flow and
// its libraries, stackups and controlled impedance, vias and HDI, board placement and routing,
// differential pairs, board-level signal and power integrity, assembly, verification and
// manufacturing outputs, and the agent-first future of board design.
package pcb

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func bank(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }

// mix serves a computed part on one knowledge point and a bank on the other two.
func mix(kp int, f edakit.Func, b *edakit.Bank) *edakit.Mix {
	parts := []generator.Generator{b, b, b}
	parts[kp-1] = f
	return &edakit.Mix{Parts: parts}
}

// Register adds the PCB design generators to reg.
func Register(reg *generator.Registry) {
	// flow
	reg.Register("pcb.flow.schematic_to_layout", bank(flowItems))
	reg.Register("pcb.flow.symbols_footprints", bank(footprintItems))
	reg.Register("pcb.flow.bom", edakit.Func(bom))
	reg.Register("pcb.flow.netlist_handoff", bank(netlistItems))
	reg.Register("pcb.flow.erc", edakit.Func(erc))
	reg.Register("pcb.flow.tools", bank(toolItems))
	// stackup
	reg.Register("pcb.stackup.layers", edakit.Func(stackLayers))
	reg.Register("pcb.stackup.copper_resistance", edakit.Func(copperResistance))
	reg.Register("pcb.stackup.current_capacity", edakit.Func(currentCapacity))
	reg.Register("pcb.stackup.controlled_impedance", edakit.Func(controlledImpedance))
	reg.Register("pcb.stackup.prop_delay", edakit.Func(propDelay))
	reg.Register("pcb.stackup.flex_rigid_flex", mix(2, bendRadius, bank(flexItems)))
	// via
	reg.Register("pcb.via.types", mix(3, stubResonance, bank(viaTypeItems)))
	reg.Register("pcb.via.aspect_ratio", edakit.Func(aspectRatio))
	reg.Register("pcb.via.hdi", edakit.Func(hdi))
	// layout
	reg.Register("pcb.layout.placement", mix(2, ratsnest, bank(placementItems)))
	reg.Register("pcb.layout.autorouter", bank(autorouterItems))
	reg.Register("pcb.layout.bga_escape", edakit.Func(bgaEscape))
	// diff
	reg.Register("pcb.diff.odd_even", bank(oddEvenItems))
	reg.Register("pcb.diff.pairs", edakit.Func(diffPairs))
	reg.Register("pcb.diff.length_matching", edakit.Func(lengthMatching))
	// si
	reg.Register("pcb.si.reflection", edakit.Func(reflection))
	reg.Register("pcb.si.ibis", mix(2, ibisImpedance, bank(ibisItems)))
	reg.Register("pcb.si.touchstone", edakit.Func(touchstone))
	reg.Register("pcb.si.ibis_ami", mix(2, ffeOutput, bank(amiItems)))
	reg.Register("pcb.si.pdn_target", edakit.Func(pdnTarget))
	// assembly
	reg.Register("pcb.assembly.tht_smt", bank(thtSmtItems))
	reg.Register("pcb.assembly.dfa", bank(dfaItems))
	reg.Register("pcb.assembly.pick_place", edakit.Func(pickPlace))
	reg.Register("pcb.assembly.reflow", edakit.Func(reflow))
	// verify
	reg.Register("pcb.verify.drc", edakit.Func(boardDRC))
	reg.Register("pcb.verify.dft", edakit.Func(boardDFT))
	reg.Register("pcb.verify.bringup", bank(bringupItems))
	// output
	reg.Register("pcb.output.gerber_odb", mix(2, gerberCoord, bank(gerberItems)))
	reg.Register("pcb.output.ipc_standards", bank(ipcItems))
	// frontier
	reg.Register("pcb.frontier.constraint_agent", bank(constraintAgentItems))
	reg.Register("pcb.frontier.autorouter_golden", edakit.Func(autorouterGolden))
	reg.Register("pcb.frontier.schematic_review", mix(2, reviewPrecision, bank(reviewItems)))
	reg.Register("pcb.frontier.routing_reward", edakit.Func(routingReward))
	reg.Register("pcb.frontier.package_codesign", bank(codesignItems))
}

// kpOf returns the requested knowledge point, or a uniform one when the request is unguided.
func kpOf(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

func sqrt(x float64) float64   { return math.Sqrt(x) }
func pow(x, y float64) float64 { return math.Pow(x, y) }
func ln(x float64) float64     { return math.Log(x) }
func log10(x float64) float64  { return math.Log10(x) }
func absf(x float64) float64   { return math.Abs(x) }
func ceilInt(x float64) int    { return int(math.Ceil(x - 1e-9)) }
func floorInt(x float64) int   { return int(math.Floor(x + 1e-9)) }

// options formats a correct value and candidate wrong values at the given precision, dropping
// any candidate that renders the same as the answer or as another option, and topping the list
// up with nearby values so that at least three distractors remain.
func options(correct float64, places int, unit string, wrong ...float64) (string, []string) {
	f := func(x float64) string {
		s := edakit.Num(x, places)
		if unit != "" {
			s += " " + unit
		}
		return s
	}
	ans := f(correct)
	seen := map[string]bool{ans: true}
	var out []string
	for _, w := range wrong {
		if correct > 0 && w <= 0 {
			continue
		}
		s := f(w)
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	step := math.Abs(correct) * 0.15
	if step == 0 {
		step = 1
	}
	for k := 1; len(out) < 3; k++ {
		for _, w := range []float64{correct + float64(k)*step, correct * (1 + 0.37*float64(k))} {
			s := f(w)
			if !seen[s] && len(out) < 3 {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return ans, out
}

// plural returns "1 thing" or "n things".
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return edakit.Sprintf("%d %s", n, many)
}

// an returns the indefinite article for a spoken number such as 8 or 11.
func an(n int) string {
	s := edakit.IntStr(n)
	if s[0] == '8' || n == 11 || n == 18 || (n >= 11000 && n < 12000) {
		return "an"
	}
	return "a"
}

// paren wraps a negative number in parentheses for use inside an expression.
func paren(s string) string {
	if strings.HasPrefix(s, "-") {
		return "(" + s + ")"
	}
	return s
}

// joinAnd joins items as "a, b and c".
func joinAnd(xs []string) string {
	if len(xs) <= 1 {
		return strings.Join(xs, "")
	}
	return strings.Join(xs[:len(xs)-1], ", ") + " and " + xs[len(xs)-1]
}
