package infra

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerFlow) }

func registerFlow(reg *generator.Registry) {
	reg.Register("infra.flow.yosys_abc", &edakit.Bank{Items: yosysAbcItems})
	reg.Register("infra.flow.openroad_tools", &edakit.Bank{Items: openroadToolsItems})
	reg.Register("infra.flow.orfs", edakit.Func(orfs))
	reg.Register("infra.flow.gcd_first_run", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: gcdFirstRunItems}, edakit.Func(gcdFirstRun)}})
	reg.Register("infra.flow.openlane", &edakit.Bank{Items: openlaneItems})
	reg.Register("infra.flow.layout_tools", &edakit.Bank{Items: layoutToolsItems})
	reg.Register("infra.flow.simulation_tools", &edakit.Bank{Items: simulationToolsItems})
}

// flowStep is one in-stage step of the OpenROAD-flow-scripts flow, in flow order.
type flowStep struct {
	name, stage string
}

// flowSteps lists steps whose relative order is fixed by data dependence.
var flowSteps = []flowStep{
	{"synthesis with Yosys and ABC", "synth"},
	{"floorplan initialization from utilization", "floorplan"},
	{"macro placement", "floorplan"},
	{"power distribution network generation", "floorplan"},
	{"global placement", "place"},
	{"repair design (resizing and buffering)", "place"},
	{"detailed placement", "place"},
	{"clock tree synthesis", "cts"},
	{"post-CTS timing repair", "cts"},
	{"global routing", "route"},
	{"detailed routing", "route"},
	{"final reports and the GDS merge", "finish"},
}

var flowStages = []string{"synth", "floorplan", "place", "cts", "route", "finish"}

var flowWhy = map[string]string{
	"synth":     "synthesis must produce the gate-level netlist before anything physical can happen",
	"floorplan": "the floorplan fixes the die, core rows, macros and power grid that placement needs",
	"place":     "cells must have locations before a clock tree can be built to their clock pins",
	"cts":       "the clock tree's buffers and wires must exist before signal routing claims the tracks",
	"route":     "only routed wires give the geometry that final fill, reports and GDS need",
	"finish":    "the finish stage works on the fully routed design",
}

// flowDraw picks n distinct step indices and returns them sorted.
func flowDraw(n int) []int {
	idx := edakit.Shuffle([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})[:n]
	for i := 1; i < len(idx); i++ {
		for j := i; j > 0 && idx[j] < idx[j-1]; j-- {
			idx[j], idx[j-1] = idx[j-1], idx[j]
		}
	}
	return idx
}

func flowOrderText(idx []int) string {
	parts := make([]string, len(idx))
	for i, k := range idx {
		parts[i] = flowSteps[k].name
	}
	return strings.Join(parts, " → ")
}

// orfs asks stage-membership (KP 1), full-ordering (KP 2) and immediate-successor or
// first-step (KP 3) questions over drawn ORFS steps.
func orfs(ctx generator.GeneratorContext) generator.Problem {
	switch variant(ctx) {
	case 1:
		k := edakit.Int(0, len(flowSteps)-1)
		st := flowSteps[k]
		var d []string
		for _, s := range flowStages {
			if s != st.stage {
				d = append(d, s)
			}
		}
		d = edakit.Shuffle(d)[:4]
		q := edakit.Sprintf("In OpenROAD-flow-scripts, the Makefile runs the stages synth, floorplan, place, cts, route and finish. In which stage does the step \"%s\" run?", st.name)
		return edakit.Choice(q, st.stage, edakit.Sprintf("The step \"%s\" belongs to the %s stage: %s.", st.name, st.stage, flowWhy[st.stage]), d...)
	case 2:
		n := edakit.Scaled(ctx, 4, 6)
		idx := flowDraw(n)
		correct := flowOrderText(idx)
		var d []string
		seen := map[string]bool{correct: true}
		for tries := 0; len(d) < 3 && tries < 100; tries++ {
			p := append([]int(nil), idx...)
			i := edakit.Int(0, n-2)
			j := edakit.Int(i+1, n-1)
			if tries%2 == 0 {
				j = i + 1
			}
			p[i], p[j] = p[j], p[i]
			s := flowOrderText(p)
			if !seen[s] {
				seen[s] = true
				d = append(d, s)
			}
		}
		shown := edakit.Shuffle(idx)
		var listed []string
		for _, k := range shown {
			listed = append(listed, "\""+flowSteps[k].name+"\"")
		}
		q := edakit.Sprintf("An OpenROAD-flow-scripts run performs these steps, listed here out of order: %s. Which ordering is the one the flow uses?", strings.Join(listed, ", "))
		var why []string
		for _, k := range idx {
			why = append(why, flowSteps[k].name+" ("+flowSteps[k].stage+")")
		}
		return edakit.Choice(q, correct, "Ordered by stage and by data dependence: "+strings.Join(why, ", then ")+".", d...)
	}
	n := edakit.Scaled(ctx, 4, 6)
	idx := flowDraw(n)
	shown := edakit.Shuffle(idx)
	var opts []string
	for _, k := range shown {
		opts = append(opts, flowSteps[k].name)
	}
	if edakit.Int(0, 2) == 0 {
		ans := flowSteps[idx[0]].name
		var d []string
		for _, k := range idx[1:] {
			d = append(d, flowSteps[k].name)
		}
		q := edakit.Sprintf("Among the steps %s, which runs first in an OpenROAD-flow-scripts run?", strings.Join(flowQuote(opts), ", "))
		return edakit.Choice(q, ans, edakit.Sprintf("\"%s\" is in the %s stage, the earliest stage among these steps, and %s.", ans, flowSteps[idx[0]].stage, flowWhy[flowSteps[idx[0]].stage]), d...)
	}
	pos := edakit.Int(0, n-2)
	x, next := flowSteps[idx[pos]], flowSteps[idx[pos+1]]
	var d []string
	for i, k := range idx {
		if i != pos+1 && i != pos {
			d = append(d, flowSteps[k].name)
		}
	}
	d = append(d, x.name)
	q := edakit.Sprintf("Considering only the steps %s, which of them is the next to run after \"%s\" in an OpenROAD-flow-scripts run?", strings.Join(flowQuote(opts), ", "), x.name)
	return edakit.Choice(q, next.name, edakit.Sprintf("In flow order the drawn steps run %s, so \"%s\" follows \"%s\".", flowOrderText(idx), next.name, x.name), d...)
}

func flowQuote(xs []string) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = "\"" + x + "\""
	}
	return out
}

// gcdFirstRun sizes a square core from a synthesized cell area and CORE_UTILIZATION.
func gcdFirstRun(ctx generator.GeneratorContext) generator.Problem {
	type plat struct {
		name string
		row  float64
	}
	p := edakit.Pick([]plat{{"nangate45", 1.4}, {"sky130hd", 2.72}, {"asap7", 0.27}})
	var area float64
	switch p.name {
	case "nangate45":
		area = float64(50 * edakit.Int(8, 30))
	case "sky130hd":
		area = float64(200 * edakit.Int(10, 40))
	default:
		area = float64(5 * edakit.Int(6, 30))
	}
	util := float64(5 * edakit.Int(5, 14))
	core := area / (util / 100)
	side := math.Sqrt(core)
	rows := int(math.Floor(side / p.row))
	base := edakit.Sprintf("A first ORFS run of the GCD design on the %s platform synthesizes to %s square microns of standard-cell area, and config.mk sets CORE_UTILIZATION to %s. The floorplan step makes a square core. ",
		p.name, edakit.Num(area, 0), edakit.Num(util, 0))
	_ = rows
	switch edakit.Int(0, 2) {
	case 0:
		ans, d := numOptions(core, 1, "square microns", area*util/100, area*(1+util/100), core*util/100)
		return edakit.Choice(base+"What is the core area?", ans,
			edakit.Sprintf("Utilization is cell area over core area, so core area = %s / %s = %s square microns.", edakit.Num(area, 0), edakit.Num(util/100, 2), edakit.Num(core, 1)), d...)
	case 1:
		ans, d := numOptions(side, 2, "microns", math.Sqrt(area), core/4, math.Sqrt(area*util/100))
		return edakit.Choice(base+"What is the side length of the square core? (Round to 2 decimal places.)", ans,
			edakit.Sprintf("Core area = %s / %s = %s square microns; its square root is %s microns.", edakit.Num(area, 0), edakit.Num(util/100, 2), edakit.Num(core, 1), edakit.Num(side, 2)), d...)
	}
	cands := []int{int(math.Ceil(side / p.row)), int(math.Floor(math.Sqrt(area) / p.row)), int(math.Floor(core / p.row / 10)), rows + 2}
	ans := edakit.IntStr(rows)
	d := edakit.Distinct(rows, 3, cands...)
	return edakit.Choice(base+edakit.Sprintf("The %s standard-cell row height is %s microns. How many whole rows fit in the core height?", p.name, edakit.Num(p.row, 2)), ans,
		edakit.Sprintf("Core side = sqrt(%s / %s) = %s microns, and floor(%s / %s) = %d rows.", edakit.Num(area, 0), edakit.Num(util/100, 2), edakit.Num(side, 2), edakit.Num(side, 2), edakit.Num(p.row, 2), rows), d...)
}

var yosysAbcItems = []edakit.Item{
	{KP: 1, Q: "In a Yosys script, what does the proc pass do?",
		A: "Turns behavioral processes into multiplexers and flip-flop cells", D: []string{"Maps gates onto Liberty cells", "Writes the final netlist", "Checks the design hierarchy for a top module"},
		Why: "proc converts always blocks (processes in RTLIL) into explicit multiplexers, flip-flops and latches, which later passes can optimize and map."},
	{KP: 1, Q: "A synthesis script runs hierarchy -top gcd first. What is this step for?",
		A: "To elaborate the design from the named top module and find missing modules", D: []string{"To flatten every module into one", "To map flip-flops to the library", "To insert scan chains"},
		Why: "hierarchy sets the top module, elaborates parameters and reports undefined modules before optimization begins."},
	{KP: 1, Q: "In an open flow, which pass maps the design's flip-flops onto the library's sequential cells?",
		A: "dfflibmap with the Liberty file", D: []string{"abc with the Liberty file", "opt_clean", "write_verilog"},
		Why: "dfflibmap maps internal flip-flop cells to Liberty flip-flops. abc handles the combinational logic."},
	{KP: 1, Q: "Why does a Yosys script usually end with opt_clean before write_verilog?",
		A: "To remove unused wires and cells left by earlier passes", D: []string{"To add timing constraints to the netlist", "To convert the netlist into GDS", "To merge the Liberty file into the netlist"},
		Why: "Mapping leaves dangling wires and unused cells. opt_clean removes them so the written netlist is clean."},
	{KP: 2, Q: "What does ABC receive from Yosys during abc -liberty?",
		A: "The combinational logic, which ABC optimizes and maps to library gates", D: []string{"The full design including flip-flops and memories", "The routed layout", "The SDC constraints only"},
		Why: "Yosys hands ABC the combinational portions. ABC optimizes them (for example as AIGs) and maps them to Liberty cells, then Yosys stitches the result back."},
	{KP: 2, Q: "ABC reads the Liberty file during mapping. What does it take from it?",
		A: "Cell functions, areas and delay information", D: []string{"Cell layout geometry", "Routing layer directions", "Clock definitions"},
		Why: "To choose gates ABC needs each cell's Boolean function, area and delay. Geometry is in LEF and clocks are in SDC."},
	{KP: 2, Q: "A flow passes a target clock period to ABC. What does ABC do with it?",
		A: "Steers mapping toward delay-oriented choices to meet the target", D: []string{"Writes it into the SDC file", "Builds a clock tree", "Ignores it, since ABC is area-only"},
		Why: "With a delay target, ABC's mapping and resynthesis favor faster structures, typically at some area cost."},
	{KP: 2, Q: "After abc finishes, the netlist contains only gates from the Liberty file. What representation did ABC most likely optimize internally?",
		A: "An and-inverter graph", D: []string{"A placement density map", "A SPEF parasitic network", "A GDS polygon set"},
		Why: "ABC's core optimizations work on AIGs, which are then mapped to library cells."},
	{KP: 3, Q: "Yosys stops with an error that module sky130_sram is not defined. What is the likeliest fix?",
		A: "Provide a blackbox or library model for the macro", D: []string{"Raise the clock period", "Lower the placement density", "Run abc twice"},
		Why: "A hard macro has no RTL to synthesize. Yosys needs a blackbox declaration or a Liberty model so it can keep the instance."},
	{KP: 3, Q: "The synthesized netlist has far fewer cells than expected and outputs are tied to constants. What is a likely cause?",
		A: "Logic with no observable outputs was optimized away", D: []string{"The Liberty file had too many corners", "ABC refused to map flip-flops", "write_verilog skipped the top module"},
		Why: "Logic that drives nothing observable is removed. Unconnected outputs or a missing keep attribute often explain a shrunken netlist."},
	{KP: 3, Q: "A netlist from Yosys contains assign statements that the downstream tool rejects. Which kind of pass addresses this?",
		A: "A pass that splits nets or inserts buffers for direct assignments", D: []string{"A pass that adds scan flip-flops", "A pass that sets the clock period", "A pass that flattens the Liberty file"},
		Why: "Some physical tools dislike plain assigns between ports. Flows add buffers or reshape nets so every connection goes through a cell or a simple net."},
	{KP: 3, Q: "The same RTL synthesizes to very different areas with two Liberty files from one PDK. Why?",
		A: "The libraries differ in cell sizes and available functions", D: []string{"Yosys parses Verilog differently per library", "Area is random in ABC", "The second library ignores flip-flops"},
		Why: "Mapping depends on the cells offered. Different libraries offer different areas and functions, so the mapped result changes."},
}

var openroadToolsItems = []edakit.Item{
	{KP: 1, Q: "Which OpenROAD tool creates the die and core areas and the standard-cell rows?",
		A: "ifp, the floorplan initializer", D: []string{"gpl, the global placer", "dpl, the detailed placer", "grt, the global router"},
		Why: "ifp initializes the floorplan from a utilization or explicit die area and builds rows from the library site."},
	{KP: 1, Q: "Which OpenROAD tool assigns top-level IO pins to locations on the die boundary?",
		A: "ppl, the pin placer", D: []string{"pdn, the power-grid generator", "cts, the clock tree tool", "drt, the detailed router"},
		Why: "ppl places IO pins on the boundary, typically to reduce wirelength to the cells they connect to."},
	{KP: 1, Q: "Which OpenROAD tool computes overlapping, analytical cell locations using an electrostatic density model?",
		A: "gpl, the global placer", D: []string{"dpl, the detailed placer", "ifp, the floorplan initializer", "rsz, the resizer"},
		Why: "gpl comes from the RePlAce line of electrostatics-based placers. Its output still has overlaps that legalization removes."},
	{KP: 1, Q: "After global placement, which OpenROAD tool snaps cells to legal sites in rows without overlap?",
		A: "dpl, the detailed placer", D: []string{"gpl, the global placer", "ppl, the pin placer", "grt, the global router"},
		Why: "dpl legalizes and refines placement so each cell sits on a site in a row with no overlaps."},
	{KP: 2, Q: "Which OpenROAD tool sizes gates and inserts buffers to repair slew, capacitance and timing violations?",
		A: "rsz, the resizer", D: []string{"cts, the clock tree tool", "dpl, the detailed placer", "OpenRCX, the extractor"},
		Why: "The resizer runs repair design and repair timing, sizing cells and inserting buffers."},
	{KP: 2, Q: "Which tool in OpenROAD builds the clock distribution network?",
		A: "cts, from TritonCTS", D: []string{"grt, from FastRoute", "drt, from TritonRoute", "pdn, the power-grid generator"},
		Why: "cts builds a buffered clock tree toward balanced latency at the sinks."},
	{KP: 2, Q: "Which OpenROAD tool assigns nets to gcells and layers and estimates congestion before exact wires exist?",
		A: "grt, the global router", D: []string{"drt, the detailed router", "gpl, the global placer", "OpenSTA, the timer"},
		Why: "grt is global routing from the FastRoute line. It produces route guides and congestion data."},
	{KP: 2, Q: "The detailed router reports remaining DRC violations. Which tool produced that report?",
		A: "drt, from TritonRoute", D: []string{"grt, from FastRoute", "dpl, the detailed placer", "ifp, the floorplan initializer"},
		Why: "drt turns route guides into exact DRC-aware wires and vias and reports any violations it cannot fix."},
	{KP: 3, Q: "Which OpenROAD component writes parasitics for the routed nets in SPEF form?",
		A: "OpenRCX, the parasitic extractor", D: []string{"OpenSTA, the static timer", "rsz, the resizer", "ppl, the pin placer"},
		Why: "OpenRCX extracts resistance and capacitance from routed geometry, and the timer reads them back for signoff-style analysis."},
	{KP: 3, Q: "After resizing, the repair step queries slack without writing any file. How can it do that?",
		A: "OpenSTA runs inside the same process on the same database", D: []string{"It reruns synthesis to recompute slack", "It reads slack from the last DEF file", "It estimates slack from wirelength only"},
		Why: "OpenROAD links its tools into one application over one OpenDB, so the timer sees edits immediately."},
	{KP: 3, Q: "Between which two OpenROAD tools do route guides pass?",
		A: "From grt to drt", D: []string{"From drt to grt", "From gpl to ifp", "From cts to ppl"},
		Why: "Global routing produces guides that confine the detailed router's search."},
	{KP: 3, Q: "Why can OpenROAD run gpl with timing-driven weighting during placement?",
		A: "The placer can call the shared timer between iterations", D: []string{"Timing is stored inside LEF files", "The placer reads slack from SPEF only", "Placement ignores timing entirely"},
		Why: "Because tools share the database and process, the placer can ask the timer for criticality and reweight nets as it goes."},
}

var gcdFirstRunItems = []edakit.Item{
	{KP: 1, Q: "In ORFS, where does the GCD design's configuration for a given platform live?",
		A: "In designs/<platform>/gcd/config.mk", D: []string{"In the platform's Liberty file", "In the Yosys source tree", "In results/<platform>/gcd/final.gds"},
		Why: "Each design has a config.mk per platform under designs/. It names the design, the platform, the Verilog and the SDC."},
	{KP: 1, Q: "Which config.mk variable names the timing constraints file for GCD?",
		A: "SDC_FILE", D: []string{"VERILOG_FILES", "PLACE_DENSITY", "DESIGN_NAME"},
		Why: "SDC_FILE points to the constraints; VERILOG_FILES to the RTL; DESIGN_NAME names the top module."},
	{KP: 1, Q: "A new user changes PLACE_DENSITY in config.mk. What does that knob affect?",
		A: "The target cell density that global placement spreads toward", D: []string{"The clock period of the design", "The number of metal layers", "The Liberty corner"},
		Why: "PLACE_DENSITY sets the global placer's target density, which trades wirelength against spreading."},
	{KP: 1, Q: "What does CORE_UTILIZATION set in a first GCD run?",
		A: "The share of the core area the standard cells should fill", D: []string{"The fraction of CPU cores ORFS uses", "The share of the clock period spent in logic", "The fraction of routing tracks reserved"},
		Why: "The floorplan derives core area from the synthesized cell area divided by this utilization."},
	{KP: 1, Q: "Why is GCD a good first design?",
		A: "It is small, so a full run finishes quickly and failures are easy to read", D: []string{"It contains every kind of macro", "It needs no SDC", "It has no clock"},
		Why: "A small design exercises every stage in minutes, so the environment can be validated before larger designs."},
	{KP: 3, Q: "After a GCD run, where should a user look first for the final timing and area numbers?",
		A: "The reports directory for the design and platform", D: []string{"The objects directory", "The platform's LEF files", "The Yosys source"},
		Why: "ORFS writes logs, reports and results per platform and design; final metrics are in reports and logs."},
	{KP: 3, Q: "The first GCD run fails in global routing with overflow. What is the most sensible first change?",
		A: "Lower CORE_UTILIZATION to give the router more space", D: []string{"Raise CORE_UTILIZATION to shorten wires", "Delete the SDC file", "Switch the Verilog to VHDL"},
		Why: "Overflow means demand exceeds track supply. A bigger core adds supply; a denser one makes it worse."},
	{KP: 3, Q: "The final GCD report shows negative setup slack. Which input most directly sets the target being missed?",
		A: "The clock period in the SDC file", D: []string{"The PLACE_DENSITY value", "The platform's GDS layer map", "The number of routing threads"},
		Why: "Slack is measured against the clock defined in the SDC. Relaxing or tightening it moves every setup slack."},
	{KP: 3, Q: "Which directory holds the intermediate ODB databases that let a user reopen GCD after placement?",
		A: "The results directory", D: []string{"The reports directory", "The platform directory", "The designs directory"},
		Why: "Results hold stage outputs such as databases, netlists and the final GDS, so a later stage or GUI session can load them."},
	{KP: 3, Q: "Two GCD runs with identical inputs produce identical results. Why does that matter for a first run?",
		A: "It establishes a baseline that later changes can be compared against", D: []string{"It proves the design is fabricated", "It shows the tools ignore inputs", "It means no reports are needed"},
		Why: "A reproducible baseline lets every later edit be judged by a clean difference."},
}

var openlaneItems = []edakit.Item{
	{KP: 1, Q: "How does OpenLane 2 structure an RTL-to-GDS flow?",
		A: "As a sequence of steps composed into a flow in Python", D: []string{"As one monolithic C++ binary", "As a single Tcl script per design", "As a fixed Makefile with no extension points"},
		Why: "OpenLane 2 is Python infrastructure in which each step wraps a tool action, and flows are compositions of steps."},
	{KP: 1, Q: "A team wants to insert a custom checker between placement and CTS in OpenLane 2. What is the natural approach?",
		A: "Write a new step and add it to a custom flow", D: []string{"Patch the OpenROAD C++ core", "Edit the PDK's Liberty files", "Rename the GDS output"},
		Why: "Steps are the extension unit, so a new check is a new step placed where it belongs."},
	{KP: 1, Q: "What passes between OpenLane 2 steps?",
		A: "A state of design views such as netlist, ODB, DEF and reports", D: []string{"Only a single log file", "The RTL source alone", "Nothing; each step starts from RTL"},
		Why: "Each step consumes and produces design views, so later steps start from the earlier outputs."},
	{KP: 1, Q: "Why does a step-based design make a flow easier to reuse?",
		A: "Steps can be reordered, replaced or reused in new flows", D: []string{"Steps hide every intermediate file", "Steps remove the need for a PDK", "Steps make tools run in parallel automatically"},
		Why: "Well-defined inputs and outputs per step let developers build new flows from existing pieces."},
	{KP: 2, Q: "Which tools does OpenLane 2 call for synthesis and for place and route?",
		A: "Yosys for synthesis and OpenROAD for implementation", D: []string{"ABC alone for both", "Magic for synthesis and Netgen for routing", "Verilator for synthesis and ngspice for routing"},
		Why: "OpenLane orchestrates the same engines: Yosys (with ABC) for synthesis and OpenROAD tools for floorplan through routing."},
	{KP: 2, Q: "In OpenLane, which tools typically handle the final DRC and LVS checks?",
		A: "Magic or KLayout for DRC and Netgen for LVS", D: []string{"OpenSTA for DRC and Yosys for LVS", "grt for DRC and gpl for LVS", "cocotb for both"},
		Why: "Physical verification uses layout tools: Magic and KLayout run rule decks, and Netgen compares netlists for LVS."},
	{KP: 2, Q: "How does OpenLane 2 relate to a PDK such as SKY130?",
		A: "It reads the PDK's views and per-PDK settings to configure each tool", D: []string{"It contains its own process technology", "It replaces the PDK with generic rules", "It works only without a PDK"},
		Why: "The flow is process-neutral. The PDK supplies LEF, Liberty, GDS, decks and settings."},
	{KP: 2, Q: "A bug appears in detailed routing during an OpenLane run. Which project's code most likely holds the bug?",
		A: "OpenROAD, whose router the step calls", D: []string{"OpenLane's configuration parser", "The SKY130 Liberty files", "The Yosys frontend"},
		Why: "OpenLane wraps engines. Routing behavior comes from the OpenROAD router the step invokes."},
	{KP: 3, Q: "A group already runs ORFS Makefiles and wants a Python API to compose custom flows. Which choice fits?",
		A: "OpenLane 2, built around Python steps", D: []string{"Only ORFS, since it is Python-based", "Neither, since flows cannot be composed", "Xschem, which composes digital flows"},
		Why: "ORFS is Makefile and script driven; OpenLane 2 offers a Python step and flow API."},
	{KP: 3, Q: "What do ORFS and OpenLane 2 have in common?",
		A: "Both orchestrate Yosys and OpenROAD on open PDKs", D: []string{"Both are forks of the same commercial tool", "Both ship their own placer", "Both replace OpenDB with files"},
		Why: "They are two orchestration layers over the same open engines."},
	{KP: 3, Q: "A team needs to compare QoR from ORFS and OpenLane 2 on one design. Why can the numbers differ?",
		A: "The flows use different default steps, settings and scripts", D: []string{"The engines produce random numbers each run", "One of them ignores the SDC", "OpenLane uses a different process"},
		Why: "Same engines, different recipes: defaults for density, repair, routing and checks differ, so QoR differs."},
	{KP: 3, Q: "When is ORFS the better place to start for engine development?",
		A: "When you want a direct, script-level view of OpenROAD's stages", D: []string{"When you need a graphical schematic editor", "When you want to avoid OpenROAD", "When you need an analog simulator"},
		Why: "ORFS is maintained alongside OpenROAD and exposes its stage scripts directly."},
}

var layoutToolsItems = []edakit.Item{
	{KP: 1, Q: "In ORFS, which tool merges the standard-cell GDS views into the routed design to make the final GDS?",
		A: "KLayout", D: []string{"Netgen", "Yosys", "OpenSTA"},
		Why: "The finish stage uses KLayout to replace cell abstracts with real layouts and write the final GDS."},
	{KP: 1, Q: "How are KLayout DRC rules usually written?",
		A: "As scripted rule decks run by KLayout's DRC engine", D: []string{"As Liberty attributes", "As SDC commands", "As Verilog assertions"},
		Why: "KLayout DRC decks are scripts that express layer operations and checks."},
	{KP: 1, Q: "A user opens a GDS in KLayout but sees meaningless layer numbers. What is missing?",
		A: "A layer-properties file mapping numbers to named layers", D: []string{"The SPEF file", "The Liberty timing file", "The testbench"},
		Why: "GDS stores layers as numbers. A layer map gives them names and colors."},
	{KP: 1, Q: "Why is KLayout useful for inspecting very large GDS files?",
		A: "It reads and draws hierarchical layout efficiently", D: []string{"It converts GDS to RTL", "It simulates the circuit", "It reroutes the design"},
		Why: "KLayout is built as a fast hierarchical layout viewer and editor."},
	{KP: 2, Q: "What is Magic?",
		A: "A VLSI layout editor with built-in DRC and circuit extraction", D: []string{"A Verilog simulator", "A logic synthesizer", "A Liberty characterization tool"},
		Why: "Magic edits layout, checks rules interactively and extracts SPICE netlists from layout."},
	{KP: 2, Q: "An analog designer wants a SPICE netlist with parasitics from a hand-drawn SKY130 layout. Which tool fits?",
		A: "Magic, using its extraction", D: []string{"Netgen, using its LVS", "Verilator, using its compiler", "Yosys, using ABC"},
		Why: "Magic extracts devices and parasitics from layout into a SPICE netlist."},
	{KP: 2, Q: "Why do many SKY130 flows run DRC in both Magic and KLayout?",
		A: "Two independent engines catch rule gaps either one misses", D: []string{"Each checks a different half of the die", "Magic cannot read GDS", "KLayout only checks metal1"},
		Why: "Cross-checking with independent decks reduces the chance a rule gap slips through."},
	{KP: 2, Q: "Magic reports a spacing violation on met2 that KLayout does not. What is the right first step?",
		A: "Inspect the geometry and compare the two decks' rules", D: []string{"Ignore it, since KLayout passed", "Delete met2 from the layout", "Rerun synthesis"},
		Why: "Disagreement means one deck encodes the rule differently. The geometry and the rule text settle it."},
	{KP: 3, Q: "What does Netgen compare during LVS?",
		A: "The netlist extracted from layout against the source netlist", D: []string{"Two GDS files polygon by polygon", "Timing reports from two corners", "RTL against a testbench"},
		Why: "LVS matches devices and connectivity between layout and schematic or gate netlist."},
	{KP: 3, Q: "Netgen reports a mismatch where two nets are merged in the layout. What physical error is likely?",
		A: "A short between the two nets", D: []string{"An open on one net", "A missing SDC constraint", "A wrong clock period"},
		Why: "Two source nets appearing as one layout net is the signature of a short."},
	{KP: 3, Q: "Why does LVS need the extraction step before Netgen runs?",
		A: "Netgen compares netlists, so the layout must first become a netlist", D: []string{"Netgen reads only GDS polygons", "Extraction fixes DRC errors", "Extraction writes the SDC"},
		Why: "Netgen works on netlists. Magic or KLayout first extracts one from the layout."},
	{KP: 3, Q: "An LVS run reports a pin mismatch only on a power net. What should be checked first?",
		A: "Whether the power pin names and connections match in both netlists", D: []string{"The clock tree depth", "The ABC script", "The placement density"},
		Why: "Power pins often differ in naming or are implicit in one netlist, a common LVS failure."},
}

var simulationToolsItems = []edakit.Item{
	{KP: 1, Q: "What does Verilator do with synthesizable Verilog?",
		A: "Compiles it into a C++ model that runs cycle by cycle", D: []string{"Interprets it event by event with four-state values", "Maps it to Liberty cells", "Converts it into SPICE"},
		Why: "Verilator translates the design into optimized C++ for fast cycle-based simulation."},
	{KP: 1, Q: "A test depends on seeing X values propagate from an uninitialized register. Why might Verilator not show them?",
		A: "Verilator simulates two-state logic by default", D: []string{"Verilator cannot simulate registers", "X values only appear in synthesis", "Verilator runs only analog models"},
		Why: "Verilator's default model is two-state, so unknowns become 0 or 1, possibly randomized."},
	{KP: 1, Q: "What is cocotb?",
		A: "A framework for writing testbenches as Python coroutines", D: []string{"A Verilog synthesis tool", "A formal model checker", "A layout viewer"},
		Why: "cocotb drives a simulator from Python, awaiting clock edges and signals in coroutines."},
	{KP: 1, Q: "A team wants Python-based random stimulus with the speed of a compiled simulator. Which pairing fits?",
		A: "cocotb driving Verilator", D: []string{"SymbiYosys driving Netgen", "Xschem driving KLayout", "ngspice driving Yosys"},
		Why: "cocotb supplies the Python testbench and Verilator the fast compiled model."},
	{KP: 2, Q: "What does SymbiYosys add on top of Yosys?",
		A: "A front end that runs formal checks with SAT and SMT engines", D: []string{"A cycle-based C++ simulator", "A detailed router", "A SPICE solver"},
		Why: "SymbiYosys reads designs and properties through Yosys and runs bounded or unbounded proofs with solvers."},
	{KP: 2, Q: "A SymbiYosys BMC run with depth 20 finds no counterexample. What does that prove?",
		A: "No violation exists within 20 cycles of reset", D: []string{"The property holds for all time", "The design is equivalent to its RTL", "Simulation coverage is complete"},
		Why: "Bounded model checking is complete only up to its depth. An unbounded proof needs induction or similar."},
	{KP: 2, Q: "Which SymbiYosys mode asks for an example trace that reaches a cover statement?",
		A: "Cover mode", D: []string{"Prove mode", "BMC mode", "Equivalence mode"},
		Why: "Cover mode searches for a trace that hits each cover property, showing it is reachable."},
	{KP: 2, Q: "A property passes in SymbiYosys prove mode using k-induction. What has been shown?",
		A: "The property holds in every reachable state", D: []string{"The property holds only for the first k cycles", "The design synthesizes", "The property is untestable"},
		Why: "A successful induction proof establishes the property for all reachable states, not just a bounded depth."},
	{KP: 3, Q: "Which open tool solves a transistor-level circuit from a SPICE netlist?",
		A: "ngspice", D: []string{"Verilator", "SymbiYosys", "Netgen"},
		Why: "ngspice is a SPICE-class circuit simulator; the others are digital simulation, formal and LVS tools."},
	{KP: 3, Q: "Why might a team pick Xyce over ngspice for a very large post-layout netlist?",
		A: "Xyce is designed to run in parallel on large circuits", D: []string{"Xyce is a digital simulator", "Xyce draws schematics", "ngspice cannot read SPICE"},
		Why: "Xyce, from Sandia, was built for parallel circuit simulation of large problems."},
	{KP: 3, Q: "What is Xschem used for?",
		A: "Drawing schematics and producing SPICE netlists for simulation", D: []string{"Running DRC on GDS", "Compiling Verilog to C++", "Proving assertions"},
		Why: "Xschem is schematic capture that exports netlists for ngspice or Xyce."},
	{KP: 3, Q: "A designer must check that a comparator settles within 2 ns across corners. Which tool answers this?",
		A: "A SPICE simulator such as ngspice with corner models", D: []string{"Verilator with a two-state model", "SymbiYosys in cover mode", "KLayout's DRC engine"},
		Why: "Analog settling depends on device physics, so it needs circuit simulation with the PDK's corner models."},
}
