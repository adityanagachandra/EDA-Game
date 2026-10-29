// Package techmap generates technology-mapping problems: subject graphs and library patterns,
// tree and DAG covering, cut enumeration and area recovery, FPGA LUT mapping, Boolean matching,
// load-dependent mapping, post-mapping sizing and buffering, the open-source mapping commands,
// and the agent-first frontier of mapping.
//
// The curriculum key for this domain is "map", which is a Go keyword, so the package lives in
// internal/generator/eda/techmap and is named techmap.
package techmap

import (
	"math/bits"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the technology-mapping generators to reg.
func Register(reg *generator.Registry) {
	// subject graphs
	reg.Register("map.subject.covering_problem", &edakit.Bank{Items: coveringProblemItems})
	reg.Register("map.subject.decomposition", edakit.Func(decomposition))
	reg.Register("map.subject.pattern_graphs", edakit.Func(patternGraphs))
	// tree and DAG covering
	reg.Register("map.tree.partition", edakit.Func(treePartition))
	reg.Register("map.tree.matching", edakit.Func(treeMatching))
	reg.Register("map.tree.dp_covering", edakit.Func(dpCovering))
	reg.Register("map.tree.delay_covering", edakit.Func(delayCovering))
	reg.Register("map.tree.dag_covering", edakit.Func(dagCovering))
	reg.Register("map.tree.limits", &edakit.Bank{Items: treeLimitsItems})
	// cuts
	reg.Register("map.cut.k_feasible", edakit.Func(kFeasible))
	reg.Register("map.cut.enumeration", edakit.Func(cutEnumeration))
	reg.Register("map.cut.priority_cuts", edakit.Func(priorityCuts))
	reg.Register("map.cut.area_flow", edakit.Func(areaFlow))
	reg.Register("map.cut.exact_area", edakit.Func(exactArea))
	// LUT mapping
	reg.Register("map.lut.flowmap", edakit.Func(flowmap))
	reg.Register("map.lut.cover_selection", edakit.Func(coverSelection))
	reg.Register("map.lut.area_recovery", edakit.Func(lutAreaRecovery))
	reg.Register("map.lut.wide_functions", edakit.Func(wideFunctions))
	// Boolean matching
	reg.Register("map.match.npn", edakit.Func(npnClasses))
	reg.Register("map.match.canonical", edakit.Func(canonicalForms))
	reg.Register("map.match.equivalence_check", edakit.Func(npnEquivalence))
	reg.Register("map.match.supergates", edakit.Func(supergates))
	// load-dependent mapping
	reg.Register("map.load.liberty_delay", edakit.Func(libertyDelay))
	reg.Register("map.load.crossover", edakit.Func(loadCrossover))
	reg.Register("map.load.gain_based", edakit.Func(gainBased))
	// sizing
	reg.Register("map.size.gate_sizing", edakit.Func(gateSizing))
	reg.Register("map.size.lagrangian", edakit.Func(lagrangianSizing))
	reg.Register("map.size.vt_swap", edakit.Func(vtSwap))
	reg.Register("map.size.pin_swap", edakit.Func(pinSwap))
	// buffering
	reg.Register("map.buffer.candidates", edakit.Func(bufferCandidates))
	reg.Register("map.buffer.pruning", edakit.Func(bufferPruning))
	reg.Register("map.buffer.van_ginneken", edakit.Func(vanGinneken))
	// tools
	reg.Register("map.tools.abc", &edakit.Bank{Items: abcItems})
	reg.Register("map.tools.yosys", &edakit.Bank{Items: yosysItems})
	// frontier
	reg.Register("map.frontier.cost_seam", &edakit.Bank{Items: costSeamItems})
	reg.Register("map.frontier.recipe_selection", &edakit.Bank{Items: recipeItems})
	reg.Register("map.frontier.golden_qor", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(goldenVerdict), edakit.Func(goldenDelta), &edakit.Bank{Items: goldenItems}}})
	reg.Register("map.frontier.change_tiers", &edakit.Bank{Items: changeTierItems})
}

// kp returns the knowledge point to serve: the requested one, or a draw weighted by difficulty
// when the request is unguided.
func kp(ctx generator.GeneratorContext) int {
	if ctx.KPOrdinal >= 1 && ctx.KPOrdinal <= 3 {
		return ctx.KPOrdinal
	}
	return edakit.Int(1, 3)
}

func popcount(m uint32) int { return bits.OnesCount32(m) }

// piName names primary input i: a, b, c, ...
func piName(i int) string { return string(rune('a' + i)) }

// joinInts formats integers as "1, 2, 3".
func joinInts(xs []int) string {
	s := make([]string, len(xs))
	for i, x := range xs {
		s[i] = edakit.IntStr(x)
	}
	return strings.Join(s, ", ")
}

// distinctStrings keeps the first occurrence of each string.
func distinctStrings(xs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}

// numChoices builds distinct decimal distractors for a computed multiple-choice question.
func numChoices(correct float64, places int, cands ...float64) []string {
	c := edakit.Num(correct, places)
	var out []string
	seen := map[string]bool{c: true}
	for _, x := range cands {
		s := edakit.Num(x, places)
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	step := 1.0
	for i := 0; i < places; i++ {
		step /= 10
	}
	for d := 1; len(out) < 3; d++ {
		for _, x := range []float64{correct + float64(d)*step*5, correct - float64(d)*step*5} {
			s := edakit.Num(x, places)
			if !seen[s] && len(out) < 3 {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}

func sortedKeys(m map[int]bool) []int {
	var out []int
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}
