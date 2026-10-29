package em

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func registerCircuit(reg *generator.Registry) {
	reg.Register("em.circuit.charge_current_voltage", edakit.Func(chargeCurrentVoltage))
	reg.Register("em.circuit.kcl_kvl", edakit.Func(kclKvl))
	reg.Register("em.circuit.resistor_networks", edakit.Func(resistorNetworks))
	reg.Register("em.circuit.nodal_analysis", edakit.Func(nodalAnalysis))
	reg.Register("em.circuit.mna_stamps", edakit.Func(mnaStamps))
	reg.Register("em.circuit.mesh_analysis", edakit.Func(meshAnalysis))
	reg.Register("em.circuit.thevenin_norton", edakit.Func(theveninNorton))
	reg.Register("em.circuit.capacitor_inductor", edakit.Func(capacitorInductor))
	reg.Register("em.circuit.rc_time_constant", edakit.Func(rcTimeConstant))
	reg.Register("em.circuit.rise_time", edakit.Func(riseTime))
	reg.Register("em.circuit.rlc", edakit.Func(rlc))
	reg.Register("em.circuit.distributed_rc", edakit.Func(distributedRC))
}

var supplies = []float64{0.8, 0.9, 1.0, 1.1, 1.2, 1.8, 2.5, 3.3}

func chargeCurrentVoltage(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		i, t := edakit.Int(2, 40), edakit.Int(2, 30)
		q := float64(i*t) / 2
		return edakit.Numeric(
			sp("A driver's output current ramps linearly from 0 to %d mA over %d ns and then stops. How much charge does it deliver during the ramp, in picocoulombs? Round to 1 decimal place.", i, t),
			num(q, 1),
			sp("Charge is the integral of current: the area under a triangle, %d mA times %d ns over 2, which is %s pC (one mA for one ns is one pC).", i, t, num(q, 1)),
			facts("charge_pc", num(q, 1)))
	case 2:
		v := edakit.Pick(supplies)
		i, t := edakit.Int(1, 60), edakit.Int(2, 100)
		e := v * float64(i) * float64(t)
		return edakit.Numeric(
			sp("A block draws a steady %d mA from a %s V supply for %d ns. How much energy does it take from the supply, in picojoules? Round to 1 decimal place.", i, num(v, 1), t),
			num(e, 1),
			sp("Power is $P = VI = %s \\cdot %d = %s$ mW, and energy is power times time: %s mW for %d ns is %s pJ.", num(v, 1), i, num(v*float64(i), 2), num(v*float64(i), 2), t, num(e, 1)),
			facts("power_mw", num(v*float64(i), 2), "energy_pj", num(e, 1)))
	default:
		r := float64(10 * edakit.Int(1, 50))
		if edakit.Int(0, 1) == 0 {
			v := rnd(0.5, 0.1, 30)
			p := v * v / r * 1000
			return edakit.Numeric(
				sp("A %s V drop appears across a %s ohm resistor. What power does the resistor dissipate, in milliwatts? Round to 2 decimal places.", num(v, 1), num(r, 0)),
				num(p, 2),
				sp("$P = V^2/R = %s^2 / %s = %s$ mW.", num(v, 1), num(r, 0), num(p, 2)),
				facts("power_mw", num(p, 2)))
		}
		i := edakit.Int(1, 40)
		p := float64(i*i) * r / 1000
		return edakit.Numeric(
			sp("A current of %d mA flows through a %s ohm resistor. What power does the resistor dissipate, in milliwatts? Round to 2 decimal places.", i, num(r, 0)),
			num(p, 2),
			sp("$P = I^2 R = (%d \\times 10^{-3})^2 \\cdot %s$ W, which is %s mW.", i, num(r, 0), num(p, 2)),
			facts("power_mw", num(p, 2)))
	}
}

func kclKvl(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 5)
		var parts []string
		in := 0
		for b := 1; b <= n; b++ {
			c := edakit.Int(1, 25)
			if edakit.Int(0, 1) == 0 {
				parts = append(parts, sp("%d mA enters through branch %d", c, b))
				in += c
			} else {
				parts = append(parts, sp("%d mA leaves through branch %d", c, b))
				in -= c
			}
		}
		return edakit.Numeric(
			sp("At a node, %s. One more branch, branch %d, connects to the node. What current leaves the node through branch %d, in mA? A negative answer means the current actually enters.", strings.Join(parts, ", "), n+1, n+1),
			edakit.IntStr(in),
			sp("KCL says the currents leaving a node sum to zero, so branch %d must carry away the net inflow of the other branches: %d mA.", n+1, in),
			facts("net_inflow_ma", edakit.IntStr(in)))
	case 2:
		vs := rnd(1, 0.5, 8)
		n := edakit.Scaled(ctx, 2, 4)
		var parts []string
		sum := 0.0
		for k := 1; k <= n; k++ {
			d := rnd(0.1, 0.1, 9)
			sum += d
			parts = append(parts, sp("element %d drops %s V", k, num(d, 1)))
		}
		x := vs - sum
		return edakit.Numeric(
			sp("A single loop contains a %s V source and %d elements in series with one unknown element. Going around the loop in the direction of the source's rise, %s. What voltage drops across the unknown element, in volts? Round to 1 decimal place; a negative answer means a rise.", num(vs, 1), n, strings.Join(parts, ", ")),
			num(x, 1),
			sp("KVL: the rise of the source equals the sum of the drops around the loop, so the unknown drop is $%s - %s = %s$ V.", num(vs, 1), num(sum, 1), num(x, 1)),
			facts("sum_known_drops", num(sum, 1)))
	default:
		nodes := edakit.Scaled(ctx, 4, 12)
		b := nodes + edakit.Int(0, 8)
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A connected circuit has %d nodes (counting the reference node) and %d branches. How many independent KVL (loop) equations does it have?", nodes, b),
				edakit.IntStr(b-nodes+1),
				sp("A spanning tree has %d branches; each of the remaining $b - n + 1 = %d - %d + 1 = %d$ branches closes one independent loop.", nodes-1, b, nodes, b-nodes+1),
				facts("loops", edakit.IntStr(b-nodes+1)))
		}
		return edakit.Numeric(
			sp("A connected circuit has %d nodes (counting the reference node) and %d branches. How many independent KCL equations does it have?", nodes, b),
			edakit.IntStr(nodes-1),
			sp("The KCL equations of all %d nodes sum to zero, so one is redundant: $n - 1 = %d$ are independent. The branch count does not matter.", nodes, nodes-1),
			facts("kcl", edakit.IntStr(nodes-1)))
	}
}

func resistorNetworks(ctx generator.GeneratorContext) generator.Problem {
	r1 := float64(10 * edakit.Int(1, 50))
	r2 := float64(10 * edakit.Int(1, 50))
	r3 := float64(10 * edakit.Int(1, 50))
	switch kp(ctx) {
	case 1:
		eq := r1 + par(r2, r3)
		return edakit.Numeric(
			sp("A resistor of %s ohm is in series with the parallel pair of %s ohm and %s ohm. What is the equivalent resistance, in ohms? Round to 1 decimal place.", num(r1, 0), num(r2, 0), num(r3, 0)),
			num(eq, 1),
			sp("The parallel pair is $%s \\cdot %s/(%s + %s) = %s$ ohm; adding the series resistor gives %s ohm.", num(r2, 0), num(r3, 0), num(r2, 0), num(r3, 0), num(par(r2, r3), 2), num(eq, 1)),
			facts("parallel_ohm", num(par(r2, r3), 2)))
	case 2:
		vs := edakit.Pick(supplies)
		if edakit.Int(0, 1) == 0 {
			v := vs * r2 / (r1 + r2)
			return edakit.Numeric(
				sp("A %s V source drives a divider: %s ohm from the source to the output node, %s ohm from the output node to ground. What is the output voltage, in volts? Round to 3 decimal places.", num(vs, 1), num(r1, 0), num(r2, 0)),
				num(v, 3),
				sp("$V_{out} = V_s R_2/(R_1 + R_2) = %s \\cdot %s/%s = %s$ V.", num(vs, 1), num(r2, 0), num(r1+r2, 0), num(v, 3)),
				facts("vout", num(v, 3)))
		}
		i := float64(edakit.Int(1, 50))
		i1 := i * r2 / (r1 + r2)
		return edakit.Numeric(
			sp("A current of %s mA splits between two parallel resistors, %s ohm and %s ohm. How much current flows in the %s ohm resistor, in mA? Round to 3 decimal places.", num(i, 0), num(r1, 0), num(r2, 0), num(r1, 0)),
			num(i1, 3),
			sp("The current divider sends the fraction $R_2/(R_1+R_2)$ through $R_1$: $%s \\cdot %s/%s = %s$ mA.", num(i, 0), num(r2, 0), num(r1+r2, 0), num(i1, 3)),
			facts("i1_ma", num(i1, 3)))
	default:
		vs := edakit.Pick(supplies)
		rl := float64(10 * edakit.Int(1, 50))
		p := par(r2, rl)
		v := vs * p / (r1 + p)
		return edakit.Numeric(
			sp("A divider of %s ohm (top) and %s ohm (bottom) is fed by %s V. A load of %s ohm is connected from the output to ground. What is the loaded output voltage, in volts? Round to 3 decimal places.", num(r1, 0), num(r2, 0), num(vs, 1), num(rl, 0)),
			num(v, 3),
			sp("The load parallels the bottom resistor: $%s \\parallel %s = %s$ ohm. Then $V_{out} = %s \\cdot %s/(%s + %s) = %s$ V, below the unloaded %s V.", num(r2, 0), num(rl, 0), num(p, 2), num(vs, 1), num(p, 2), num(r1, 0), num(p, 2), num(v, 3), num(vs*r2/(r1+r2), 3)),
			facts("bottom_parallel", num(p, 2), "unloaded", num(vs*r2/(r1+r2), 3)))
	}
}

var kohms = []float64{0.5, 1, 1.5, 2, 2.5, 3, 4, 5, 6, 8, 10}

func nodalAnalysis(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		vs := edakit.Pick(supplies)
		r1, r2 := edakit.Pick(kohms), edakit.Pick(kohms)
		i := rnd(0.1, 0.1, 9)
		v := (vs/r1 + i) / (1/r1 + 1/r2)
		return edakit.Numeric(
			sp("A %s V source connects through $R_1 = %s$ kohm to node A. Node A has $R_2 = %s$ kohm to ground and a current source injecting %s mA into it. What is the voltage at node A, in volts? Round to 3 decimal places.", num(vs, 1), num(r1, 1), num(r2, 1), num(i, 1)),
			num(v, 3),
			sp("KCL at A: $(V - %s)/%s + V/%s = %s$. So $V = (%s/%s + %s)/(1/%s + 1/%s) = %s$ V.", num(vs, 1), num(r1, 1), num(r2, 1), num(i, 1), num(vs, 1), num(r1, 1), num(i, 1), num(r1, 1), num(r2, 1), num(v, 3)),
			facts("v_a", num(v, 3)))
	case 2:
		vs := edakit.Pick(supplies)
		r1, r2, r3, r4 := edakit.Pick(kohms), edakit.Pick(kohms), edakit.Pick(kohms), edakit.Pick(kohms)
		g1, g2, g3, g4 := 1/r1, 1/r2, 1/r3, 1/r4
		x := solve([][]float64{{g1 + g2 + g3, -g2}, {-g2, g2 + g4}}, []float64{vs * g1, 0})
		which := edakit.Int(1, 2)
		return edakit.Numeric(
			sp("A %s V source drives $R_1 = %s$ kohm into node 1. Node 1 has $R_3 = %s$ kohm to ground and $R_2 = %s$ kohm to node 2. Node 2 has $R_4 = %s$ kohm to ground. What is the voltage at node %d, in volts? Round to 3 decimal places.", num(vs, 1), num(r1, 1), num(r3, 1), num(r2, 1), num(r4, 1), which),
			num(x[which-1], 3),
			sp("The node equations are $(G_1+G_2+G_3)V_1 - G_2 V_2 = G_1 V_s$ and $-G_2 V_1 + (G_2+G_4)V_2 = 0$, with $G = 1/R$ in mS. Solving gives $V_1 = %s$ V and $V_2 = %s$ V.", num(x[0], 3), num(x[1], 3)),
			facts("v1", num(x[0], 3), "v2", num(x[1], 3)))
	default:
		vs := edakit.Pick(supplies)
		r := []float64{edakit.Pick(kohms), edakit.Pick(kohms), edakit.Pick(kohms)}
		sh := []float64{edakit.Pick(kohms), edakit.Pick(kohms), edakit.Pick(kohms)}
		g := func(x float64) float64 { return 1 / x }
		a := [][]float64{
			{g(r[0]) + g(r[1]) + g(sh[0]), -g(r[1]), 0},
			{-g(r[1]), g(r[1]) + g(r[2]) + g(sh[1]), -g(r[2])},
			{0, -g(r[2]), g(r[2]) + g(sh[2])},
		}
		x := solve(a, []float64{vs * g(r[0]), 0, 0})
		which := edakit.Int(1, 3)
		return edakit.Numeric(
			sp("A ladder: a %s V source feeds %s kohm to node 1, then %s kohm from node 1 to node 2, then %s kohm from node 2 to node 3. Nodes 1, 2 and 3 have %s, %s and %s kohm to ground. What is the voltage at node %d, in volts? Round to 3 decimal places.", num(vs, 1), num(r[0], 1), num(r[1], 1), num(r[2], 1), num(sh[0], 1), num(sh[1], 1), num(sh[2], 1), which),
			num(x[which-1], 3),
			sp("Writing KCL at each node gives a tridiagonal 3 by 3 conductance system. Its solution is $V_1 = %s$, $V_2 = %s$, $V_3 = %s$ V.", num(x[0], 3), num(x[1], 3), num(x[2], 3)),
			facts("v1", num(x[0], 3), "v2", num(x[1], 3), "v3", num(x[2], 3)))
	}
}

func mnaStamps(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		vals := []float64{0.5, 1, 2, 4, 5, 10}
		type res struct {
			a, b int
			r    float64
		}
		var rs []res
		var desc []string
		pairs := [][2]int{{1, 0}, {2, 0}, {3, 0}, {1, 2}, {2, 3}, {1, 3}}
		for _, p := range edakit.Shuffle(pairs)[:edakit.Scaled(ctx, 3, 6)] {
			r := edakit.Pick(vals)
			rs = append(rs, res{p[0], p[1], r})
			if p[1] == 0 {
				desc = append(desc, sp("%s kohm from node %d to ground", num(r, 1), p[0]))
			} else {
				desc = append(desc, sp("%s kohm between nodes %d and %d", num(r, 1), p[0], p[1]))
			}
		}
		i, j := edakit.Int(1, 3), edakit.Int(1, 3)
		val := 0.0
		for _, x := range rs {
			g := 1 / x.r
			if i == j && (x.a == i || x.b == i) {
				val += g
			}
			if i != j && ((x.a == i && x.b == j) || (x.a == j && x.b == i)) {
				val -= g
			}
		}
		return edakit.Numeric(
			sp("A resistive network with ground as the reference node has %s. In the nodal conductance matrix $G$ (rows and columns ordered by node 1, 2, 3), what is the entry $G_{%d%d}$, in millisiemens? Round to 2 decimal places.", strings.Join(desc, ", "), i, j),
			num(val, 2),
			sp("Each conductance stamps $+g$ on the diagonals of its two nodes and $-g$ on the off-diagonal pair between them (a branch to ground stamps only its diagonal). Summing the stamps that touch entry (%d, %d) gives %s mS.", i, j, num(val, 2)),
			facts("entry_ms", num(val, 2)))
	case 2:
		n := edakit.Scaled(ctx, 3, 40)
		vs, ls := edakit.Int(1, 4), edakit.Int(0, 6)
		rs, cs, is := edakit.Int(3, 60), edakit.Int(0, 30), edakit.Int(0, 3)
		dim := n - 1 + vs + ls
		return edakit.Numeric(
			sp("A netlist has %d nodes including ground, %d resistors, %d capacitors, %d independent current sources, %d independent voltage sources and %d inductors. In modified nodal analysis, where each voltage source and each inductor carries its branch current as an extra unknown, how many unknowns does the system have?", n, rs, cs, is, vs, ls),
			edakit.IntStr(dim),
			sp("Unknowns are the %d non-ground node voltages plus one branch current per voltage source (%d) and per inductor (%d): %d. Resistors, capacitors and current sources stamp into existing rows.", n-1, vs, ls, dim),
			facts("dim", edakit.IntStr(dim)))
	default:
		c := float64(5 * edakit.Int(1, 40))
		h := pick(0.5, 1, 2, 5, 10, 20)
		vp := rnd(0.1, 0.1, 11)
		g := c / h
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A transient simulator discretizes a %s fF capacitor with backward Euler at a time step of %s ps. What conductance does the capacitor's companion model stamp, in millisiemens? Round to 3 decimal places.", num(c, 0), num(h, 1)),
				num(g, 3),
				sp("Backward Euler replaces $i = C\\,dv/dt$ with $i_n = (C/h)(v_n - v_{n-1})$: a conductance $C/h = %s/%s$ fF/ps, which is %s mS.", num(c, 0), num(h, 1), num(g, 3)),
				facts("g_ms", num(g, 3)))
		}
		return edakit.Numeric(
			sp("A transient simulator discretizes a %s fF capacitor with backward Euler at a time step of %s ps. The capacitor voltage at the previous step was %s V. What is the magnitude of the history current source in its companion model, in mA? Round to 4 decimal places.", num(c, 0), num(h, 1), num(vp, 1)),
			num(g*vp, 4),
			sp("The companion model is a conductance $C/h = %s$ mS in parallel with a current source $(C/h) v_{n-1} = %s \\cdot %s = %s$ mA.", num(g, 3), num(g, 3), num(vp, 1), num(g*vp, 4)),
			facts("g_ms", num(g, 3), "i_hist", num(g*vp, 4)))
	}
}

func meshAnalysis(ctx generator.GeneratorContext) generator.Problem {
	v1 := rnd(1, 0.5, 8)
	v2 := rnd(0.5, 0.5, 8)
	if v2 == v1 {
		v2 = v1 - 0.5
	}
	r1, r2, r3 := edakit.Pick(kohms), edakit.Pick(kohms), edakit.Pick(kohms)
	if kp(ctx) == 1 {
		i := (v1 - v2) / (r1 + r2)
		return edakit.Numeric(
			sp("A single loop contains a %s V source aiding the clockwise direction, a %s V source opposing it, and resistors of %s kohm and %s kohm. What is the clockwise loop current, in mA? Round to 3 decimal places.", num(v1, 1), num(v2, 1), num(r1, 1), num(r2, 1)),
			num(i, 3),
			sp("KVL around the loop: $%s - %s = (%s + %s) I$, so $I = %s$ mA (volts over kilohms gives mA).", num(v1, 1), num(v2, 1), num(r1, 1), num(r2, 1), num(i, 3)),
			facts("i_ma", num(i, 3)))
	}
	aid := edakit.Int(0, 1) == 0
	s2 := -v2
	word := "opposes"
	if aid {
		s2, word = v2, "aids"
	}
	x := solve([][]float64{{r1 + r3, -r3}, {-r3, r2 + r3}}, []float64{v1, s2})
	stem := sp("Two meshes share one resistor $R_3 = %s$ kohm. Mesh 1 holds a %s V source that aids its clockwise current $I_1$ and $R_1 = %s$ kohm. Mesh 2 holds $R_2 = %s$ kohm and a %s V source that %s its clockwise current $I_2$.", num(r3, 1), num(v1, 1), num(r1, 1), num(r2, 1), num(v2, 1), word)
	expl := sp("Mesh equations: $(R_1+R_3)I_1 - R_3 I_2 = %s$ and $-R_3 I_1 + (R_2+R_3) I_2 = %s$. Solving gives $I_1 = %s$ mA and $I_2 = %s$ mA.", num(v1, 1), num(s2, 1), num(x[0], 3), num(x[1], 3))
	f := facts("i1", num(x[0], 3), "i2", num(x[1], 3))
	if kp(ctx) == 2 {
		which := edakit.Int(1, 2)
		return edakit.Numeric(stem+sp(" What is $I_%d$, in mA? Round to 3 decimal places.", which), num(x[which-1], 3), expl, f)
	}
	return edakit.Numeric(stem+" What current flows through $R_3$ in the direction of $I_1$, in mA? Round to 3 decimal places.",
		num(x[0]-x[1], 3), expl+sp(" The shared resistor carries $I_1 - I_2 = %s$ mA.", num(x[0]-x[1], 3)), f)
}

func theveninNorton(ctx generator.GeneratorContext) generator.Problem {
	vs := edakit.Pick(supplies)
	r1, r2, r3 := edakit.Pick(kohms), edakit.Pick(kohms), edakit.Pick(kohms)
	vth := vs * r2 / (r1 + r2)
	rth := par(r1, r2) + r3
	stem := sp("A %s V source drives $R_1 = %s$ kohm into node a; $R_2 = %s$ kohm runs from a to ground, and $R_3 = %s$ kohm runs from a to the output terminal.", num(vs, 1), num(r1, 1), num(r2, 1), num(r3, 1))
	f := facts("vth", num(vth, 3), "rth", num(rth, 3))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+" What is the Thevenin (open-circuit) voltage at the output terminal, in volts? Round to 3 decimal places.",
			num(vth, 3), sp("With the output open, no current flows in $R_3$, so the output sits at the divider voltage $%s \\cdot %s/%s = %s$ V.", num(vs, 1), num(r2, 1), num(r1+r2, 1), num(vth, 3)), f)
	case 2:
		return edakit.Numeric(stem+" What is the Thevenin resistance seen at the output terminal, in kohm? Round to 3 decimal places.",
			num(rth, 3), sp("Zero the source (short it): $R_1 \\parallel R_2 = %s$ kohm, in series with $R_3$, gives %s kohm.", num(par(r1, r2), 3), num(rth, 3)), f)
	default:
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(stem+" What is the Norton (short-circuit) current at the output terminal, in mA? Round to 3 decimal places.",
				num(vth/rth, 3), sp("$I_N = V_{th}/R_{th} = %s/%s = %s$ mA.", num(vth, 3), num(rth, 3), num(vth/rth, 3)), f)
		}
		p := vth * vth / (4 * rth)
		return edakit.Numeric(stem+" What is the largest power any resistive load on the output can draw, in milliwatts? Round to 4 decimal places.",
			num(p, 4), sp("Maximum power transfer happens at $R_L = R_{th} = %s$ kohm, giving $V_{th}^2/(4R_{th}) = %s^2/(4 \\cdot %s) = %s$ mW.", num(rth, 3), num(vth, 3), num(rth, 3), num(p, 4)), f)
	}
}

func capacitorInductor(ctx generator.GeneratorContext) generator.Problem {
	switch kp(ctx) {
	case 1:
		c := float64(edakit.Int(1, 50))
		dv := edakit.Pick(supplies)
		dt := rnd(0.1, 0.1, 29)
		i := c * dv / dt
		return edakit.Numeric(
			sp("A %s pF load is slewed by %s V in %s ns at a constant rate. What current does the driver supply during the ramp, in mA? Round to 2 decimal places.", num(c, 0), num(dv, 1), num(dt, 1)),
			num(i, 2),
			sp("$i = C\\,dV/dt = %s\\text{ pF} \\cdot %s\\text{ V}/%s\\text{ ns} = %s$ mA (pF times V per ns is mA).", num(c, 0), num(dv, 1), num(dt, 1), num(i, 2)),
			facts("i_ma", num(i, 2)))
	case 2:
		l := float64(edakit.Int(1, 20))
		di := float64(5 * edakit.Int(1, 20))
		dt := rnd(0.2, 0.1, 28)
		v := l * di / dt
		return edakit.Numeric(
			sp("The current through a %s nH package inductance changes by %s mA in %s ns at a constant rate. What voltage appears across the inductance, in mV? Round to 1 decimal place.", num(l, 0), num(di, 0), num(dt, 1)),
			num(v, 1),
			sp("$v = L\\,di/dt = %s\\text{ nH} \\cdot %s\\text{ mA}/%s\\text{ ns} = %s$ mV.", num(l, 0), num(di, 0), num(dt, 1), num(v, 1)),
			facts("v_mv", num(v, 1)))
	default:
		if edakit.Int(0, 1) == 0 {
			c := float64(edakit.Int(1, 100))
			v := edakit.Pick(supplies)
			e := 0.5 * c * v * v
			return edakit.Numeric(
				sp("A %s pF capacitor is charged to %s V. How much energy does it store, in picojoules? Round to 3 decimal places.", num(c, 0), num(v, 1)),
				num(e, 3),
				sp("$E = \\tfrac12 C V^2 = 0.5 \\cdot %s \\cdot %s^2 = %s$ pJ.", num(c, 0), num(v, 1), num(e, 3)),
				facts("e_pj", num(e, 3)))
		}
		l := float64(edakit.Int(1, 20))
		i := float64(10 * edakit.Int(1, 20))
		e := 0.5 * l * i * i / 1000
		return edakit.Numeric(
			sp("A %s nH inductor carries %s mA. How much energy does its magnetic field store, in picojoules? Round to 3 decimal places.", num(l, 0), num(i, 0)),
			num(e, 3),
			sp("$E = \\tfrac12 L I^2 = 0.5 \\cdot %s \\times 10^{-9} \\cdot (%s \\times 10^{-3})^2$ J, which is %s pJ.", num(l, 0), num(i, 0), num(e, 3)),
			facts("e_pj", num(e, 3)))
	}
}

func rcTimeConstant(ctx generator.GeneratorContext) generator.Problem {
	r := pick(0.5, 1, 1.5, 2, 3, 4, 5, 8, 10, 12, 20)
	c := float64(5 * edakit.Int(1, 40))
	tau := r * c
	switch kp(ctx) {
	case 1:
		if edakit.Int(0, 1) == 0 {
			return edakit.Numeric(
				sp("A %s kohm driver resistance charges a %s fF load. What is the time constant, in picoseconds? Round to 1 decimal place.", num(r, 1), num(c, 0)),
				num(tau, 1),
				sp("$\\tau = RC = %s\\text{ kohm} \\cdot %s\\text{ fF} = %s$ ps (one kohm times one fF is one ps).", num(r, 1), num(c, 0), num(tau, 1)),
				facts("tau_ps", num(tau, 1)))
		}
		r2 := pick(1, 2, 4, 5, 10)
		rt := par(r, r2)
		return edakit.Numeric(
			sp("A source with %s kohm series resistance drives a node that has a %s kohm resistor and a %s fF capacitor to ground. What is the time constant of the node, in picoseconds? Round to 1 decimal place.", num(r, 1), num(r2, 1), num(c, 0)),
			num(rt*c, 1),
			sp("The capacitor sees the Thevenin resistance of the rest of the circuit: $%s \\parallel %s = %s$ kohm. So $\\tau = %s \\cdot %s = %s$ ps.", num(r, 1), num(r2, 1), num(rt, 3), num(rt, 3), num(c, 0), num(rt*c, 1)),
			facts("rth", num(rt, 3), "tau_ps", num(rt*c, 1)))
	case 2:
		v := edakit.Pick(supplies)
		t := edakit.Round(float64(edakit.Int(1, 40))*tau/10, 1)
		if edakit.Int(0, 1) == 0 {
			x := v * (1 - math.Exp(-t/tau))
			return edakit.Numeric(
				sp("A %s kohm resistor charges a %s fF capacitor from 0 V toward a %s V step. What is the capacitor voltage %s ps after the step, in volts? Round to 3 decimal places.", num(r, 1), num(c, 0), num(v, 1), num(t, 1)),
				num(x, 3),
				sp("$\\tau = %s$ ps, so $v = %s(1 - e^{-%s/%s}) = %s$ V.", num(tau, 1), num(v, 1), num(t, 1), num(tau, 1), num(x, 3)),
				facts("tau_ps", num(tau, 1)))
		}
		x := v * math.Exp(-t/tau)
		return edakit.Numeric(
			sp("A %s fF capacitor charged to %s V discharges through %s kohm. What is its voltage %s ps later, in volts? Round to 3 decimal places.", num(c, 0), num(v, 1), num(r, 1), num(t, 1)),
			num(x, 3),
			sp("$\\tau = %s$ ps, so $v = %s e^{-%s/%s} = %s$ V.", num(tau, 1), num(v, 1), num(t, 1), num(tau, 1), num(x, 3)),
			facts("tau_ps", num(tau, 1)))
	default:
		f := pick(0.5, 0.6, 0.7, 0.8, 0.9, 0.95, 0.99)
		t := tau * math.Log(1/(1-f))
		return edakit.Numeric(
			sp("A %s kohm resistor charges a %s fF capacitor from 0 V toward a step. How long until the capacitor reaches %s%% of the final voltage, in picoseconds? Round to 1 decimal place.", num(r, 1), num(c, 0), num(100*f, 0)),
			num(t, 1),
			sp("Solve $1 - e^{-t/\\tau} = %s$: $t = \\tau \\ln(1/(1 - %s)) = %s \\cdot %s = %s$ ps.", num(f, 2), num(f, 2), num(tau, 1), num(math.Log(1/(1-f)), 4), num(t, 1)),
			facts("tau_ps", num(tau, 1)))
	}
}

func riseTime(ctx generator.GeneratorContext) generator.Problem {
	r := pick(0.5, 1, 2, 2.5, 4, 5, 8, 10)
	c := float64(5 * edakit.Int(1, 40))
	tau := r * c
	switch kp(ctx) {
	case 1:
		t := tau * math.Ln2
		return edakit.Numeric(
			sp("A gate with an effective drive resistance of %s kohm switches a %s fF load, modeled as one RC stage driven by a step. What is the 50%% delay, in picoseconds? Round to 1 decimal place.", num(r, 1), num(c, 0)),
			num(t, 1),
			sp("$t_{50} = \\tau \\ln 2 = 0.693 \\cdot %s = %s$ ps.", num(tau, 1), num(t, 1)),
			facts("tau_ps", num(tau, 1)))
	case 2:
		t := tau * math.Log(9)
		return edakit.Numeric(
			sp("A single RC stage with $R = %s$ kohm and $C = %s$ fF is driven by a step. What is the 10%% to 90%% rise time of its output, in picoseconds? Round to 1 decimal place.", num(r, 1), num(c, 0)),
			num(t, 1),
			sp("$t_{90} - t_{10} = \\tau(\\ln 10 - \\ln(10/9)) = \\tau \\ln 9 = 2.197 \\cdot %s = %s$ ps.", num(tau, 1), num(t, 1)),
			facts("tau_ps", num(tau, 1)))
	default:
		tr := float64(5 * edakit.Int(2, 100))
		f := math.Log(9) / (2 * math.Pi * tr) * 1000
		return edakit.Numeric(
			sp("A single-pole RC response has a 10%% to 90%% rise time of %s ps. What is its 3 dB bandwidth, in GHz? Round to 3 decimal places.", num(tr, 0)),
			num(f, 3),
			sp("$f_{3dB} = 1/(2\\pi\\tau)$ and $t_r = \\tau \\ln 9$, so $f_{3dB} = \\ln 9/(2\\pi t_r) \\approx 0.35/t_r = %s$ GHz.", num(f, 3)),
			facts("f3db_ghz", num(f, 3)))
	}
}

func rlc(ctx generator.GeneratorContext) generator.Problem {
	l := float64(edakit.Int(1, 20))
	c := float64(edakit.Int(1, 20))
	z0 := math.Sqrt(1000 * l / c) // sqrt(L/C) in ohms for nH and pF
	stem := sp("A series RLC circuit has $L = %s$ nH and $C = %s$ pF", num(l, 0), num(c, 0))
	switch kp(ctx) {
	case 1:
		r := float64(edakit.Int(1, 60))
		if edakit.Int(0, 1) == 0 {
			z := r / (2 * z0)
			return edakit.Numeric(stem+sp(", with $R = %s$ ohm. What is its damping ratio $\\zeta$? Round to 3 decimal places.", num(r, 0)),
				num(z, 3),
				sp("$\\zeta = \\frac{R}{2}\\sqrt{C/L}$; $\\sqrt{L/C} = %s$ ohm, so $\\zeta = %s/(2 \\cdot %s) = %s$.", num(z0, 3), num(r, 0), num(z0, 3), num(z, 3)),
				facts("sqrt_l_over_c", num(z0, 3)))
		}
		f0 := 1 / (2 * math.Pi * math.Sqrt(l*1e-9*c*1e-12)) / 1e9
		return edakit.Numeric(stem+". What is its undamped natural frequency $f_0 = \\omega_0/(2\\pi)$, in GHz? Round to 3 decimal places.",
			num(f0, 3),
			sp("$\\omega_0 = 1/\\sqrt{LC}$, so $f_0 = 1/(2\\pi\\sqrt{%s \\times 10^{-9} \\cdot %s \\times 10^{-12}}) = %s$ GHz.", num(l, 0), num(c, 0), num(f0, 3)),
			facts("f0_ghz", num(f0, 3)))
	case 2:
		rc := 2 * z0
		return edakit.Numeric(stem+". What series resistance makes it critically damped, in ohms? Round to 2 decimal places.",
			num(rc, 2),
			sp("Critical damping is $\\zeta = 1$, which needs $R = 2\\sqrt{L/C} = 2 \\cdot %s = %s$ ohm. Less resistance rings; more is overdamped.", num(z0, 3), num(rc, 2)),
			facts("r_crit", num(rc, 2)))
	default:
		z := float64(edakit.Int(5, 80)) / 100
		r := edakit.Round(2*z*z0, 2)
		z = r / (2 * z0)
		os := 100 * math.Exp(-z*math.Pi/math.Sqrt(1-z*z))
		return edakit.Numeric(stem+sp(", with $R = %s$ ohm, and the output is taken across the capacitor. A step drives it. What is the percent overshoot of the output? Round to 1 decimal place.", num(r, 2)),
			num(os, 1),
			sp("$\\zeta = R/(2\\sqrt{L/C}) = %s/(2 \\cdot %s) = %s$. Overshoot is $100\\,e^{-\\zeta\\pi/\\sqrt{1-\\zeta^2}} = %s$ percent.", num(r, 2), num(z0, 3), num(z, 3), num(os, 1)),
			facts("zeta", num(z, 3)))
	}
}

func distributedRC(ctx generator.GeneratorContext) generator.Problem {
	r := pick(0.05, 0.08, 0.1, 0.2, 0.4, 0.5, 0.8, 1, 2)
	c := pick(0.1, 0.15, 0.2, 0.25, 0.3)
	l := float64(100 * edakit.Scaled(ctx, 4, 30))
	rt, ct := r*l, c*l
	rc := rt * ct / 1000 // ps
	stem := sp("A wire of length %s um has resistance %s ohm/um and capacitance %s fF/um to ground", num(l, 0), num(r, 2), num(c, 2))
	f := facts("R_ohm", num(rt, 1), "C_fF", num(ct, 1), "RC_ps", num(rc, 3))
	switch kp(ctx) {
	case 1:
		return edakit.Numeric(stem+". Modeling it as a distributed RC line driven by an ideal source, what is its Elmore delay at the far end, in picoseconds? Round to 2 decimal places.",
			num(rc/2, 2),
			sp("$R = %s$ ohm and $C = %s$ fF. Each element of capacitance sees the resistance up to it, so the Elmore sum becomes $\\int_0^L r x\\, c\\, dx = RC/2 = %s$ ps.", num(rt, 1), num(ct, 1), num(rc/2, 2)), f)
	case 2:
		return edakit.Numeric(stem+". Driven by an ideal step, what is its 50% delay at the far end, in picoseconds, using the distributed-line value $0.38RC$? Round to 2 decimal places.",
			num(0.38*rc, 2),
			sp("$RC = %s \\cdot %s$ ohm fF $= %s$ ps, so $t_{50} \\approx 0.38 RC = %s$ ps. A lumped single RC would claim $0.69RC = %s$ ps.", num(rt, 1), num(ct, 1), num(rc, 3), num(0.38*rc, 2), num(0.69*rc, 2)), f)
	default:
		l = float64(1000 * edakit.Int(1, 5))
		rt, ct = r*l, c*l
		rc = rt * ct / 1000
		stem = sp("A wire of length %s um has resistance %s ohm/um and capacitance %s fF/um to ground", num(l, 0), num(r, 2), num(c, 2))
		k := edakit.Int(2, 6)
		d := 0.38 * rc / float64(k)
		return edakit.Numeric(stem+sp(". It is split into %d equal segments by ideal repeaters with zero delay and zero output resistance. Using $0.38RC$ per segment, what is the total wire delay, in picoseconds? Round to 2 decimal places.", k),
			num(d, 2),
			sp("Each segment has $RC/%d^2$, so the total is $%d \\cdot 0.38 RC/%d^2 = 0.38 RC/%d = %s$ ps, against %s ps unbuffered. Wire delay is quadratic in length.", k, k, k, k, num(d, 2), num(0.38*rc, 2)), f)
	}
}
