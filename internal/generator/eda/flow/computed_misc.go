package flow

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// lutCount: configuration bits of a LUT fabric, or LUTs needed for a wide AND, as computed
// multiple choice (the FPGA concept's second knowledge point).
func lutCount(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Pick([]int{4, 5, 6})
	if edakit.Pick([]bool{true, false}) {
		n := 1000 * edakit.Int(1, 400)
		ans := n << k
		q := edakit.Sprintf("An FPGA region holds %s %d-input LUTs. How many LUT configuration bits does it contain, not counting routing?", comma(n), k)
		return edakit.Choice(q, comma(ans),
			edakit.Sprintf("Each %d-LUT stores one bit per truth-table row, $2^{%d} = %d$ bits, so $%s \\cdot %d = %s$.", k, k, 1<<k, comma(n), 1<<k, comma(ans)),
			comma(n*k), comma(n<<(k-1)), comma(n*k*k), comma(n<<(k+1)))
	}
	m := edakit.Int(k+2, 40)
	ans := (m - 2 + k - 1) / (k - 1) // ceil((m-1)/(k-1))
	q := edakit.Sprintf("A %d-input AND must be built from %d-input LUTs arranged as a tree. What is the fewest LUTs that can implement it?", m, k)
	return edakit.Choice(q, edakit.IntStr(ans),
		edakit.Sprintf("Each LUT reduces the number of signals by at most %d, so $\\lceil (%d - 1)/(%d - 1) \\rceil = %d$ LUTs.", k-1, m, k, ans),
		edakit.Distinct(ans, 3, (m+k-1)/k, ans+1, ans-1, ans+2)...)
}

type tool struct{ name, vendor, class string }

var commercialTools = []tool{
	{"VCS", "Synopsys", "Logic simulation of RTL and gate netlists"},
	{"Xcelium", "Cadence", "Logic simulation of RTL and gate netlists"},
	{"Questa", "Siemens EDA", "Logic simulation of RTL and gate netlists"},
	{"JasperGold", "Cadence", "Formal property verification"},
	{"Formality", "Synopsys", "Logic equivalence checking"},
	{"Conformal", "Cadence", "Logic equivalence checking"},
	{"Design Compiler", "Synopsys", "Logic synthesis from RTL to gates"},
	{"Genus", "Cadence", "Logic synthesis from RTL to gates"},
	{"Tessent", "Siemens EDA", "Scan insertion, test compression and ATPG"},
	{"Innovus", "Cadence", "Placement and routing of a gate netlist"},
	{"IC Compiler II", "Synopsys", "Placement and routing of a gate netlist"},
	{"Fusion Compiler", "Synopsys", "Synthesis and place-and-route in one database"},
	{"PrimeTime", "Synopsys", "Static timing analysis signoff"},
	{"Tempus", "Cadence", "Static timing analysis signoff"},
	{"StarRC", "Synopsys", "Parasitic extraction"},
	{"Quantus", "Cadence", "Parasitic extraction"},
	{"Calibre", "Siemens EDA", "Physical verification: DRC and LVS"},
	{"IC Validator", "Synopsys", "Physical verification: DRC and LVS"},
	{"Pegasus", "Cadence", "Physical verification: DRC and LVS"},
	{"Voltus", "Cadence", "Power-integrity and IR-drop analysis"},
	{"RedHawk", "", "Power-integrity and IR-drop analysis"},
	{"Virtuoso", "Cadence", "Custom and analog schematic and layout design"},
	{"Spectre", "Cadence", "Transistor-level circuit simulation"},
	{"HSPICE", "Synopsys", "Transistor-level circuit simulation"},
	{"Vivado", "AMD", "FPGA synthesis, place and route and bitstream generation"},
	{"Quartus", "Altera", "FPGA synthesis, place and route and bitstream generation"},
	{"Allegro", "Cadence", "PCB and package layout"},
	{"Xpedition", "Siemens EDA", "PCB and package layout"},
	{"Advanced Design System (ADS)", "Keysight", "RF and microwave circuit and electromagnetic design"},
}

func toolClasses(ts []tool) []string {
	seen := map[string]bool{}
	var out []string
	for _, t := range ts {
		if !seen[t.class] {
			seen[t.class] = true
			out = append(out, t.class)
		}
	}
	return out
}

// toolLandscape: tool to problem class (KP 1), class to tool (KP 2), tool to vendor (KP 3).
func toolLandscape(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		t := edakit.Pick(commercialTools)
		q := "What class of problem does " + t.name + " solve?"
		return edakit.Choice(q, t.class, t.name+" is used for: "+lowerFirst(t.class)+".", pickOthers(toolClasses(commercialTools), 3, t.class)...)
	case 2:
		t := edakit.Pick(commercialTools)
		var others []string
		for _, o := range edakit.Shuffle(commercialTools) {
			if o.class != t.class {
				others = append(others, o.name)
			}
		}
		q := "A team needs a tool for this job: " + lowerFirst(t.class) + ". Which product is built for it?"
		return edakit.Choice(q, t.name, t.name+" solves that class of problem; the others serve different stages.", others[:3]...)
	default:
		var withVendor []tool
		for _, t := range commercialTools {
			if t.vendor != "" {
				withVendor = append(withVendor, t)
			}
		}
		t := edakit.Pick(withVendor)
		q := "Which company develops " + t.name + "?"
		return edakit.Choice(q, t.vendor, t.name+" is a "+t.vendor+" product.",
			pickOthers([]string{"Synopsys", "Cadence", "Siemens EDA", "AMD", "Altera", "Keysight", "Ansys"}, 3, t.vendor)...)
	}
}

type otool struct{ name, role, group, input string }

var openTools = []otool{
	{"Yosys", "RTL synthesis to a gate-level netlist", "yosys", "Verilog RTL"},
	{"ABC", "AIG-based logic optimization and technology mapping", "abc", "An and-inverter graph or logic network"},
	{"OpenROAD", "Floorplanning, placement, clock tree synthesis and routing", "openroad", "A gate netlist with LEF, Liberty and SDC"},
	{"OpenSTA", "Static timing analysis", "opensta", "A netlist with Liberty, SDC and optionally SPEF"},
	{"KLayout", "Viewing and editing mask layout and running DRC scripts", "klayout", "GDSII or OASIS layout"},
	{"Magic", "Interactive custom layout editing with continuous DRC and extraction", "magic", ""},
	{"Netgen", "LVS comparison of two netlists", "netgen", "Two netlists, one extracted from layout"},
	{"ngspice", "SPICE circuit simulation", "spice", "A SPICE netlist with device models"},
	{"Xyce", "Parallel SPICE circuit simulation for large circuits", "spice", ""},
	{"Verilator", "Compiling Verilog into a fast cycle-based C++ model", "verilator", ""},
	{"Icarus Verilog", "Event-driven Verilog simulation", "icarus", ""},
	{"cocotb", "Writing hardware testbenches in Python", "cocotb", ""},
	{"SymbiYosys", "Formal property checking of assertions", "sby", ""},
	{"nextpnr", "FPGA placement and routing", "fpga", "A synthesized netlist for a specific FPGA device"},
	{"VPR", "FPGA placement and routing research platform of the VTR project", "fpga", ""},
}

// openSource: tool to role (KP 1), role to tool (KP 2), tool to main input (KP 3).
func openSource(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx, 3) {
	case 1:
		t := edakit.Pick(openTools)
		var roles []string
		for _, o := range openTools {
			if o.group != t.group {
				roles = append(roles, o.role)
			}
		}
		q := "What does the open-source tool " + t.name + " do?"
		return edakit.Choice(q, t.role, t.name+": "+lowerFirst(t.role)+".", pickOthers(roles, 3)...)
	case 2:
		t := edakit.Pick(openTools)
		var names []string
		for _, o := range openTools {
			if o.group != t.group {
				names = append(names, o.name)
			}
		}
		q := "Which open-source tool would you use for this job: " + lowerFirst(t.role) + "?"
		return edakit.Choice(q, t.name, t.name+" is built for that job.", pickOthers(names, 3)...)
	default:
		var withIn []otool
		for _, o := range openTools {
			if o.input != "" {
				withIn = append(withIn, o)
			}
		}
		t := edakit.Pick(withIn)
		var ins []string
		for _, o := range withIn {
			if o.name != t.name {
				ins = append(ins, o.input)
			}
		}
		q := "What is the main input of " + t.name + " in an open-source flow?"
		return edakit.Choice(q, t.input, t.name+" reads "+lowerFirst(t.input)+".", pickOthers(ins, 3)...)
	}
}

// fastAccurate: loop time, speedup and the endpoints a fast timer cannot decide.
func fastAccurate(ctx generator.GeneratorContext) generator.Problem {
	n := 10 * edakit.Int(2, 10)
	tf := edakit.Int(2, 15)
	ta := tf * edakit.Int(5, 30)
	k := edakit.Pick([]int{5, 10, 20})
	total := n*tf + (n/k)*ta
	setup := edakit.Sprintf("An agent runs %d iterations of a fast place-and-time loop, each taking %d minutes, and every %dth iteration it confirms the result with an accurate run taking %d minutes.", n, tf, k, ta)
	switch kp(ctx, 3) {
	case 1:
		q := setup + " How many minutes does the whole loop take?"
		return edakit.Numeric(q, edakit.IntStr(total),
			edakit.Sprintf("$%d \\cdot %d + \\lfloor %d/%d \\rfloor \\cdot %d = %d$ minutes.", n, tf, n, k, ta, total), nil)
	case 2:
		ans := edakit.Num(float64(n*ta)/float64(total), 2)
		q := setup + " How many times faster is this than running all iterations with the accurate engine? Round to 2 decimal places."
		return edakit.Numeric(q, ans,
			edakit.Sprintf("All-accurate: $%d \\cdot %d = %d$ minutes. Mixed: %d minutes. Ratio %s.", n, ta, n*ta, total, ans), nil)
	default:
		e := edakit.Pick([]int{5, 10, 15, 20, 25})
		m := edakit.Scaled(ctx, 6, 10)
		var sl []int
		cnt := 0
		for len(sl) < m || cnt == 0 {
			if len(sl) == m {
				sl, cnt = nil, 0
			}
			v := edakit.Int(-60, 60)
			if v == e || v == -e {
				continue
			}
			sl = append(sl, v)
			if v > -e && v < e {
				cnt++
			}
		}
		q := edakit.Sprintf("A fast timer is calibrated to be within %d ps of signoff timing. It reports endpoint setup slacks of %s ps. An endpoint can be decided from the fast run only if its slack is at least %d ps (surely passing) or below -%d ps (surely failing). How many endpoints must be re-checked with the accurate timer?", e, joinInts(sl), e, e)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Count the slacks strictly between -%d and %d ps: %d.", e, e, cnt), nil)
	}
}

// differentialQoR: endpoints moved, TNS change, and moves outside the intended set.
func differentialQoR(ctx generator.GeneratorContext) generator.Problem {
	m := edakit.Scaled(ctx, 5, 8)
	tol := edakit.Pick([]int{1, 2, 3, 5})
	base := make([]int, m)
	nw := make([]int, m)
	moved := make([]bool, m)
	for i := range base {
		base[i] = edakit.Int(-30, 40)
		if edakit.Int(0, 2) == 0 {
			d := edakit.Int(tol+1, 20)
			if edakit.Pick([]bool{true, false}) {
				d = -d
			}
			nw[i] = base[i] + d
		} else {
			nw[i] = base[i] + edakit.Int(-tol, tol)
		}
		moved[i] = abs(nw[i]-base[i]) > tol
	}
	names := make([]string, m)
	for i := range names {
		names[i] = string(rune('A' + i))
	}
	table := ""
	for i := range base {
		table += edakit.Sprintf("%s: %d to %d; ", names[i], base[i], nw[i])
	}
	table = strings.TrimSuffix(table, "; ")
	switch kp(ctx, 3) {
	case 1:
		cnt := 0
		for _, b := range moved {
			if b {
				cnt++
			}
		}
		q := edakit.Sprintf("A differential timing check compares endpoint setup slacks (in ps, baseline to new) after a flow change: %s. With a tolerance of %d ps, how many endpoints moved?", table, tol)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Count endpoints whose slack changed by more than %d ps: %d.", tol, cnt), nil)
	case 2:
		tb, tn := 0, 0
		for i := range base {
			if base[i] < 0 {
				tb += base[i]
			}
			if nw[i] < 0 {
				tn += nw[i]
			}
		}
		q := edakit.Sprintf("Endpoint setup slacks (in ps, baseline to new) after a flow change are: %s. By how much did total negative slack change? Give new TNS minus baseline TNS, in ps.", table)
		return edakit.Numeric(q, edakit.IntStr(tn-tb),
			edakit.Sprintf("Baseline TNS is %d ps and new TNS is %d ps, so the change is %d ps.", tb, tn, tn-tb), nil)
	default:
		var intended []string
		in := make([]bool, m)
		for i := range names {
			if edakit.Int(0, 1) == 0 {
				in[i] = true
				intended = append(intended, names[i])
			}
		}
		if len(intended) == 0 {
			in[0] = true
			intended = []string{names[0]}
		}
		cnt := 0
		for i := range moved {
			if moved[i] && !in[i] {
				cnt++
			}
		}
		q := edakit.Sprintf("A derating plug-in was meant to affect only endpoints %s. Endpoint setup slacks (in ps, baseline to new) are: %s. With a tolerance of %d ps, how many endpoints outside the intended set moved?", strings.Join(intended, ", "), table, tol)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("Among endpoints outside %s, %d changed by more than %d ps. Any such move is a side effect that must be explained before approval.", strings.Join(intended, ", "), cnt, tol), nil)
	}
}

func abs(x int) int { return int(math.Abs(float64(x))) }
