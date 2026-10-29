package pcb

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

var constraintAgentItems = []edakit.Item{
	{KP: 1, Q: "In a board placer built on the four-layer engine model, what belongs in the immutable core?",
		A:   "Courtyard non-overlap, keep-out and board-outline legality",
		D:   []string{"The weight between ratsnest length and thermal spreading", "A rule that keeps decoupling capacitors within 2 mm of power pins", "A hook that streams live congestion maps to an agent"},
		Why: "Legality is the invariant every placement must satisfy. Objectives, rules and hooks vary per design and sit above it."},
	{KP: 1, Q: "Which item is a parameter (configuration) of a board placer?",
		A:   "The weight between ratsnest length and a thermal keep-out margin around regulators",
		D:   []string{"The courtyard overlap test", "The algorithm that legalizes overlapping parts", "The parser for the netlist format"},
		Why: "Weights and thresholds change with no recompile. They are where most customer requests land."},
	{KP: 1, Q: "What should an observation hook in a board router expose to an agent?",
		A:   "Live state such as unrouted nets, congestion by region and current length-match errors",
		D:   []string{"The router core's source code", "The legalizer's internal data structures for direct editing", "The fabricator's price list"},
		Why: "Hooks are read points into running state, so an agent can decide what to change without touching the core."},
	{KP: 1, Q: "Which choice belongs to the policy layer of a board autorouter?",
		A:   "The rip-up ordering and the cost schedule for congested regions",
		D:   []string{"The definition of a legal clearance", "The Gerber writer", "The netlist parser"},
		Why: "Policies decide how the router searches; the clearance model and file I/O are fixed infrastructure."},
	{KP: 2, Q: "A customer wants every DQ in a DDR byte lane matched within 5 ps of its strobe. Which change tier is this?",
		A:   "Tier 1: a value in the existing length-match constraint table",
		D:   []string{"Tier 2: a new plug-in rule checker", "Tier 3: a core edit to the router", "No change, since routers match lengths automatically"},
		Why: "Length-match groups and tolerances are standard constraint types, so the request is configuration."},
	{KP: 2, Q: "A customer wants a new rule: no signal trace on any layer within 3 mm of a switching inductor's outline. The constraint engine has no rule of this kind. Which tier?",
		A:   "Tier 2: a plug-in rule checker at the constraint extension point",
		D:   []string{"Tier 1: a value in an existing table", "Tier 3: a rewrite of the DRC core", "No change, since clearance rules already cover it"},
		Why: "The rule needs new logic (component outline to trace on all layers) but plugs into the existing rule interface without touching the core."},
	{KP: 2, Q: "A request asks to change how the placement legalizer resolves overlapping courtyards. Which tier?",
		A:   "Tier 3: a core edit that needs human approval and full regression",
		D:   []string{"Tier 1: a parameter change", "Tier 2: a plug-in", "No change, since overlaps are allowed"},
		Why: "Legalization is the core invariant. Changing it can break every design, so it is rare and gated."},
	{KP: 2, Q: "An agent reads a memory controller datasheet and fills in length and impedance constraints. What must happen before routing uses them?",
		A:   "A check of each value against the datasheet tables and a human review, since an extracted number can be wrong",
		D:   []string{"Nothing, since the agent read the datasheet", "Only a DRC run", "Only a Gerber export"},
		Why: "Extracted constraints are claims. A wrong tolerance silently produces a board that passes every check and fails in the lab."},
	{KP: 3, Q: "An agent proposes a placement that cuts ratsnest length by 20%. What evidence should gate acceptance?",
		A:   "Courtyard and keep-out checks clean, critical-net constraints met, and thermal and assembly checks no worse than the baseline",
		D:   []string{"The ratsnest reduction alone", "The agent's own confidence score", "A screenshot of the new placement"},
		Why: "Ratsnest length is a proxy. Acceptance needs the real constraints checked by deterministic tools against a baseline."},
	{KP: 3, Q: "Why should an agent's constraint edits be versioned and diffable?",
		A:   "A reviewer can see exactly which rules changed and revert one without losing the rest",
		D:   []string{"Versioning makes the router faster", "Diffs are required by the Gerber format", "It hides changes from the designer"},
		Why: "Small, reviewable, reversible changes are what make agent edits safe to accept."},
	{KP: 3, Q: "An agent loop alternates fast placement with a quick congestion estimate and runs a full autoroute only on promising candidates. Which idea is this?",
		A:   "Fast and accurate engines in one loop: the fast estimator prunes and the accurate router confirms",
		D:   []string{"Tool-grounded reinforcement learning on signoff data", "A Tier 3 core change", "Golden testing of the router core"},
		Why: "The fast model is cheap enough to explore with, and the agent recurses to the accurate engine where the decision matters."},
	{KP: 3, Q: "To reach 100% completion, the agent proposes relaxing the minimum clearance rule. What should the harness do?",
		A:   "Block it and escalate: clearance comes from the fabricator's capability, not from the agent's objective",
		D:   []string{"Accept it, since completion improved", "Accept it if DRC then passes", "Halve the clearance on half the board as a compromise"},
		Why: "Changing the rule that defines legality makes every downstream check meaningless. Such rules are outside the agent's change surface."},
}

type cand struct {
	name   string
	drc    int
	skew   float64
	z      float64
	vias   int
	wl     float64
	reward float64
	siPass bool
}

// autorouterGolden: count candidates that pass DRC and skew (KP 1), plus impedance (KP 2), or
// count unintended changes in a differential delay report (KP 3).
func autorouterGolden(ctx generator.GeneratorContext) generator.Problem {
	kp := kpOf(ctx)
	if kp <= 2 {
		n := edakit.Scaled(ctx, 4, 7)
		skewMax := float64(edakit.Pick([]int{2, 3, 5}))
		zt := float64(edakit.Pick([]int{50, 85, 90, 100}))
		pass := 0
		var rows, names []string
		for i := 0; i < n; i++ {
			c := cand{name: edakit.Sprintf("R%d", i+1)}
			if edakit.Int(1, 3) == 1 {
				c.drc = edakit.Int(1, 6)
			}
			c.skew = float64(edakit.Int(2, 2*int(skewMax)+6)) / 2
			c.z = zt * (1 + float64(edakit.Int(-16, 16))/100)
			ok := c.drc == 0 && c.skew <= skewMax+1e-9
			row := edakit.Sprintf("%s: %s, worst intra-pair skew %s ps", c.name, plural(c.drc, "DRC violation", "DRC violations"), edakit.Num(c.skew, 1))
			if kp == 2 {
				ok = ok && absf(c.z-zt)/zt <= 0.10+1e-9
				row += edakit.Sprintf(", worst simulated impedance %s ohms", edakit.Num(c.z, 1))
			}
			if ok {
				pass++
				names = append(names, c.name)
			}
			rows = append(rows, row)
		}
		rule := edakit.Sprintf("zero DRC violations and worst intra-pair skew at most %s ps", edakit.Num(skewMax, 0))
		if kp == 2 {
			rule += edakit.Sprintf(" and every controlled net within 10%% of its %s ohm target in the SI simulation", edakit.Num(zt, 0))
		}
		q := "An AI autorouter produced candidate routings of one board. A candidate is accepted only with " + rule + ". The results are: " + strings.Join(rows, "; ") + ". How many candidates are accepted?"
		exp := "Apply every gate to every candidate; failing any one rejects it."
		if pass > 0 {
			exp += " Accepted: " + strings.Join(names, ", ") + "."
		}
		exp += " The checks are deterministic tools, so the router cannot talk its way past them."
		return edakit.Numeric(q, edakit.IntStr(pass), exp, map[string]string{"accepted": edakit.IntStr(pass)})
	}
	n := edakit.Scaled(ctx, 6, 10)
	tol := edakit.Pick([]float64{0.5, 1.0})
	bad := 0
	var rows, names []string
	for i := 0; i < n; i++ {
		ddr := edakit.Int(1, 2) == 1
		name := edakit.Sprintf("SIG%d", i)
		grp := "other"
		if ddr {
			name = edakit.Sprintf("DQ%d", i)
			grp = "DDR"
		}
		base := float64(edakit.Int(1500, 4000)) / 10
		delta := float64(edakit.Int(-30, 30)) / 10
		if !ddr && edakit.Int(1, 2) == 1 {
			delta = float64(edakit.Int(-4, 4)) / 10
		}
		nw := base + delta
		if !ddr && absf(nw-base) > tol+1e-9 {
			bad++
			names = append(names, name)
		}
		rows = append(rows, edakit.Sprintf("%s (%s) %s to %s ps", name, grp, edakit.Num(base, 1), edakit.Num(nw, 1)))
	}
	q := edakit.Sprintf("An agent's routing change was meant to retune only the DDR group. A differential delay report against the golden baseline lists each net's delay before and after: %s. Any net outside the DDR group whose delay moved by more than %s ps is an unintended change. How many unintended changes are there?", strings.Join(rows, "; "), edakit.Num(tol, 1))
	exp := "Ignore the DDR nets, which were supposed to move, and compare each other net's before and after delays."
	if bad > 0 {
		exp += " Unintended: " + strings.Join(names, ", ") + "."
	} else {
		exp += " None moved beyond the tolerance."
	}
	return edakit.Numeric(q, edakit.IntStr(bad), exp, map[string]string{"unintended": edakit.IntStr(bad)})
}

// reviewPrecision is the computed KP 2 of schematic_review: precision or recall of an agent's findings.
func reviewPrecision(ctx generator.GeneratorContext) generator.Problem {
	seeded := edakit.Int(10, 40)
	tp := edakit.Int(seeded/3, seeded)
	fp := edakit.Int(1, 30)
	flagged := tp + fp
	prec := 100 * float64(tp) / float64(flagged)
	rec := 100 * float64(tp) / float64(seeded)
	wantPrec := edakit.Int(1, 2) == 1
	var correct float64
	var what string
	var wrong []float64
	if wantPrec {
		correct, what = prec, "precision"
		wrong = []float64{rec, 100 * float64(fp) / float64(flagged), 100 * float64(tp) / float64(flagged+seeded)}
	} else {
		correct, what = rec, "recall"
		wrong = []float64{prec, 100 * float64(seeded-tp) / float64(seeded), 100 * float64(tp) / float64(flagged+seeded)}
	}
	ans, ds := options(correct, 1, "percent", wrong...)
	q := edakit.Sprintf("A schematic-review agent is evaluated on reference designs with %d seeded defects. It raises %d findings, of which %d match seeded defects and %d are false alarms. What is its %s?", seeded, flagged, tp, fp, what)
	return edakit.Choice(q, ans, edakit.Sprintf("Precision is true findings over all findings, $%d/%d$; recall is true findings over seeded defects, $%d/%d$. The %s is %s.", tp, flagged, tp, seeded, what, ans), ds...)
}

var reviewItems = []edakit.Item{
	{KP: 1, Q: "Which schematic check is best left to a deterministic rule engine rather than the language model in a review agent?",
		A:   "Whether every input has a driver and every power input has a source (ERC)",
		D:   []string{"Whether a filter's design intent matches the vendor's application note", "Whether a block's comments describe it well", "Whether an alternative part would be cheaper"},
		Why: "Connectivity rules are exact and cheap to compute. The model's value is in the judgment calls the rules cannot express."},
	{KP: 1, Q: "What should a review agent attach to each finding?",
		A:   "The evidence: the net or part, the rule or datasheet section it relies on, and the observed values",
		D:   []string{"A confidence score and nothing else", "A rewritten version of the whole schematic", "No supporting detail, to keep reports short"},
		Why: "A finding a human can verify in seconds is useful; one that must be re-derived is noise."},
	{KP: 1, Q: "An agent checks that every I2C bus has pull-up resistors. Which representation lets it do so reliably?",
		A:   "The netlist with pin types and part values, queried through tool calls",
		D:   []string{"A screenshot of the schematic PDF", "The Gerber files", "The BOM without connectivity"},
		Why: "Structured connectivity can be queried exactly. Reading an image invites misread nets and missed sheets."},
	{KP: 1, Q: "Which finding needs datasheet knowledge beyond ERC?",
		A:   "A regulator's output capacitor is smaller than the minimum its datasheet requires for stability",
		D:   []string{"A net that has only one pin", "Two outputs driving one net", "An input pin with no driver"},
		Why: "ERC knows pin types, not each part's application requirements; that is where retrieval and reasoning help."},
	{KP: 3, Q: "How is a schematic-review agent's quality best measured?",
		A:   "Seed known defects into reference schematics and measure the precision and recall of its findings",
		D:   []string{"Count how many findings it produces", "Measure the length of its reports", "Ask the agent to grade itself"},
		Why: "Seeded defects give ground truth. Precision and recall separate noise from coverage."},
	{KP: 3, Q: "An agent reports 40 findings per board, most of them style remarks. What is the risk?",
		A:   "Low precision buries the real defects, and engineers stop reading the reports",
		D:   []string{"High recall always compensates", "More findings always mean a safer board", "None, since every remark is harmless"},
		Why: "Review output competes for human attention. A review that is mostly noise has negative value."},
	{KP: 3, Q: "An agent cites a datasheet requirement that the datasheet does not contain. Which harness feature catches this?",
		A:   "Grounding: each cited requirement is matched against the retrieved datasheet text before the finding is shown",
		D:   []string{"A longer prompt", "A higher sampling temperature", "Running DRC on the layout"},
		Why: "A citation check turns an unverifiable claim into a test the harness can run."},
	{KP: 3, Q: "How should a review agent's proposed fixes reach the design?",
		A:   "As schematic edits in a reviewable diff that the designer approves",
		D:   []string{"Applied silently to the master schematic", "Applied directly to the Gerber files", "Sent to the fabricator first"},
		Why: "The designer owns the circuit. Agent edits are proposals, and the diff makes them cheap to accept or reject."},
	{KP: 3, Q: "Why keep the seeded-defect evaluation designs out of the data used to tune the agent's prompts and skills?",
		A:   "Tuning on the evaluation set inflates measured recall without improving real reviews",
		D:   []string{"The designs are too large to reuse", "Seeded defects expire over time", "It makes the agent run faster"},
		Why: "An evaluation set that leaks into tuning measures memorization rather than review skill."},
}

func rewardOf(c cand) float64 {
	return 100 - 20*float64(c.drc) - 0.1*float64(c.vias) - 2*c.wl - 5*c.skew
}

const rewardText = "A routing agent is trained with the tool-grounded reward $R = 100 - 20 V - 0.1 N - 2 W - 5 S$, where $V$ is the DRC violation count, $N$ the via count, $W$ the wirelength excess over the baseline in percent, and $S$ the worst intra-pair skew in ps."

func randCand(name string) cand {
	c := cand{name: name, vias: edakit.Int(20, 200), wl: float64(edakit.Int(-20, 80)) / 10, skew: float64(edakit.Int(0, 12)) / 2}
	if edakit.Int(1, 3) == 1 {
		c.drc = edakit.Int(1, 3)
	}
	c.reward = rewardOf(c)
	return c
}

func (c cand) text() string {
	return edakit.Sprintf("%s has %s, %d vias, %s%% wirelength excess and %s ps worst skew", c.name, plural(c.drc, "DRC violation", "DRC violations"), c.vias, edakit.Num(c.wl, 1), edakit.Num(c.skew, 1))
}

// routingReward: one reward (KP 1), a difference (KP 2), reward hacking among the top candidates (KP 3).
func routingReward(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		c := randCand("Candidate A")
		ans := edakit.Num(c.reward, 1)
		return edakit.Numeric(rewardText+" "+c.text()+". What is its reward? Round to 1 decimal place.", ans,
			edakit.Sprintf("$100 - 20 \\cdot %d - 0.1 \\cdot %d - 2 \\cdot %s - 5 \\cdot %s = %s$.", c.drc, c.vias, paren(edakit.Num(c.wl, 1)), edakit.Num(c.skew, 1), ans), nil)
	case 2:
		a, b := randCand("A"), randCand("B")
		d := a.reward - b.reward
		ans := edakit.Num(d, 1)
		return edakit.Numeric(rewardText+" Candidate "+a.text()+"; candidate "+b.text()+". What is $R_A - R_B$? Round to 1 decimal place.", ans,
			edakit.Sprintf("$R_A = %s$ and $R_B = %s$, so the difference is %s. One DRC violation outweighs 200 vias, which encodes that legality comes first.", edakit.Num(a.reward, 1), edakit.Num(b.reward, 1), ans), nil)
	}
	for a := 0; ; a++ {
		n := edakit.Scaled(ctx, 5, 7)
		k := edakit.Int(2, 3)
		var cs []cand
		for i := 0; i < n; i++ {
			c := randCand(edakit.Sprintf("C%d", i+1))
			c.siPass = edakit.Int(1, 2) == 1
			cs = append(cs, c)
		}
		sorted := append([]cand(nil), cs...)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i].reward > sorted[j].reward })
		if a < generator.MaxDrawAttempts && edakit.Round(sorted[k-1].reward, 1) == edakit.Round(sorted[k].reward, 1) {
			continue
		}
		fails := 0
		var top []string
		for _, c := range sorted[:k] {
			top = append(top, c.name)
			if !c.siPass {
				fails++
			}
		}
		var rows []string
		for _, c := range cs {
			s := "passes"
			if !c.siPass {
				s = "fails"
			}
			rows = append(rows, c.text()+edakit.Sprintf(", and %s the eye-mask SI signoff", s))
		}
		q := edakit.Sprintf("%s The reward has no signal-integrity term. Candidates: %s. How many of the %d highest-reward candidates fail SI signoff?", rewardText, strings.Join(rows, "; "), k)
		return edakit.Numeric(q, edakit.IntStr(fails), edakit.Sprintf("The top %d by reward are %s; %d of them fail signoff. A policy trained on this proxy learns to win the reward while the real objective is unmeasured, so the reward needs the SI check or a gate on it.", k, strings.Join(top, ", "), fails), map[string]string{"top": strings.Join(top, ",")})
	}
}

var codesignItems = []edakit.Item{
	{KP: 1, Q: "What does package and board co-design optimize that neither tool can alone?",
		A:   "The ball map: which signal sits on which package ball, so that the substrate and the board both escape cleanly",
		D:   []string{"The die's standard-cell placement", "The board's silkscreen layout", "The BOM line count"},
		Why: "A ball assignment that is easy for the package can force crossings and extra layers on the board, and the reverse."},
	{KP: 1, Q: "Which data must stay consistent across the die, package and board databases?",
		A:   "Net names and the mapping from die bumps to package balls to board pads",
		D:   []string{"The Gerber aperture tables", "The BOM quantities", "The silkscreen fonts"},
		Why: "The three databases describe one connected system. A mismatch in the mapping is an open or short across the boundary."},
	{KP: 1, Q: "Why is FPGA pin assignment a co-design lever?",
		A:   "Many I/O can be swapped within bank and I/O-standard rules, so the board can pick pins that shorten and untangle routes",
		D:   []string{"FPGA pins have no electrical rules", "Pin assignment sets the FPGA's clock frequency", "The board netlist cannot change after capture"},
		Why: "Swappable pins turn a fixed obstacle into a free variable for the board router, within the vendor's rules."},
	{KP: 1, Q: "A high-speed channel for signoff spans which pieces?",
		A:   "Die pad, package, ball, board traces and vias, connector, and on to the far device",
		D:   []string{"Only the board traces", "Only the package substrate", "Only the connector"},
		Why: "Loss and reflections add along the whole cascade, so each piece's S-parameters are joined for channel simulation."},
	{KP: 2, Q: "An agent swaps two FPGA pins in one I/O bank to remove a crossing. What must approve the swap?",
		A:   "The FPGA vendor's pin rules, enforced by the FPGA tool, and an updated board netlist back-annotated to the schematic",
		D:   []string{"Nothing, since pins in a bank are always equivalent", "Only the board DRC", "Only the BOM"},
		Why: "Bank voltage, I/O standard and clock-capable pin rules decide legality. The swap must also flow back to every view."},
	{KP: 2, Q: "A request asks to rerun ball-map optimization with more weight on board layer count. Which tier is this?",
		A:   "Tier 1: a parameter of the existing optimizer",
		D:   []string{"Tier 2: a new plug-in", "Tier 3: a core edit", "No change is possible"},
		Why: "The objective weights are configuration; nothing new is compiled."},
	{KP: 2, Q: "An agent evaluates hundreds of candidate ball maps. Which models belong in its inner loop?",
		A:   "Fast estimates such as escape-routability counts and quasi-static channel models, with full-wave extraction only for finalists",
		D:   []string{"Full-wave 3D extraction of every candidate", "Fabrication of a test board per candidate", "No models, only the agent's judgment"},
		Why: "The fast engines prune; the accurate engine confirms the few decisions that matter."},
	{KP: 2, Q: "A package supplier ships an encrypted package model. What does that mean for the agent?",
		A:   "The agent can simulate with it through the tool but cannot read or export its contents",
		D:   []string{"The agent should decrypt it to improve accuracy", "The model cannot be used in any simulation", "The agent must rebuild the package from photos"},
		Why: "The IP boundary is architectural: tools run the model, and the agent sees only results."},
	{KP: 3, Q: "After an agent changes the ball map, which checks show the change is safe?",
		A:   "Cross-domain connectivity checks over die, package and board, plus SI and PI re-simulation of the affected channels and rails",
		D:   []string{"The agent's summary of what it changed", "A board-only DRC run", "A BOM comparison"},
		Why: "A ball-map edit can open a connection across the boundary or degrade a channel, and each risk has a deterministic check."},
	{KP: 3, Q: "An agent improves board routability by moving power balls to the package edge. What must be checked?",
		A:   "Power-delivery inductance to the die, since longer paths raise supply noise against the target impedance",
		D:   []string{"Nothing, since power balls carry DC only", "Only the silkscreen", "Only the drill table"},
		Why: "Power and ground balls under the die keep the PDN loop small. Moving them trades board routability for PI margin."},
	{KP: 3, Q: "Why should co-design tools exchange machine-readable geometry and connectivity rather than reports and screenshots?",
		A:   "Agents and checkers must act on and verify the exact data, which images and prose do not carry",
		D:   []string{"Screenshots are larger files", "Reports cannot be stored", "Images are forbidden by IPC"},
		Why: "An agent loop is only as reliable as the data it reads; structured exchange makes each step checkable."},
	{KP: 3, Q: "A co-design agent reports that the board can drop two layers. What should the human reviewer see before accepting?",
		A:   "The ball-map diff, the routing evidence, and SI and PI results compared with the baselines on both sides",
		D:   []string{"The agent's estimate of saved cost alone", "A list of the layers removed", "The final Gerber files only"},
		Why: "Approval needs the change and the proof together, exactly as in an agentic ECO loop."},
}
