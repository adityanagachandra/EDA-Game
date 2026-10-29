package sim

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Four-state values.
const (
	v0 = iota
	v1
	vX
	vZ
)

var fsName = []string{"0", "1", "X", "Z"}

func fsAnd(a, b int) int {
	if a == v0 || b == v0 {
		return v0
	}
	if a == v1 && b == v1 {
		return v1
	}
	return vX
}

func fsOr(a, b int) int {
	if a == v1 || b == v1 {
		return v1
	}
	if a == v0 && b == v0 {
		return v0
	}
	return vX
}

func fsXor(a, b int) int {
	if a > v1 || b > v1 {
		return vX
	}
	return a ^ b
}

func fsNot(a int) int {
	if a > v1 {
		return vX
	}
	return 1 - a
}

// fsResolve is the resolution of two drivers on a plain wire.
func fsResolve(a, b int) int {
	if a == vZ {
		return b
	}
	if b == vZ {
		return a
	}
	if a == b {
		return a
	}
	return vX
}

func fsApply(op string, a, b int) int {
	switch op {
	case "&":
		return fsAnd(a, b)
	case "|":
		return fsOr(a, b)
	default:
		return fsXor(a, b)
	}
}

func fsChoice(q string, ans int, why string) generator.Problem {
	var d []string
	for i, n := range fsName {
		if i != ans {
			d = append(d, n)
		}
	}
	return edakit.Choice(q, fsName[ans], why, d...)
}

func fourState(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	ops := []string{"&", "|", "^"}
	switch kp {
	case 1:
		op := edakit.Pick(ops)
		a := edakit.Int(0, 3)
		b := edakit.Pick([]int{vX, vZ, edakit.Int(0, 3)})
		r := fsApply(op, a, b)
		q := edakit.Sprintf("In a four-state Verilog simulation, a = %s and b = %s. What is the value of a %s b?", fsName[a], fsName[b], op)
		why := "A Z input to a logic operator is treated as unknown, X. "
		switch op {
		case "&":
			why += "For AND, a 0 on either input forces 0; otherwise any X or Z input gives X."
		case "|":
			why += "For OR, a 1 on either input forces 1; otherwise any X or Z input gives X."
		default:
			why += "XOR has no controlling value, so any X or Z input gives X."
		}
		return fsChoice(q, r, why)
	case 2:
		a, b, c := edakit.Int(0, 3), edakit.Int(0, 3), edakit.Int(0, 3)
		op1, op2 := edakit.Pick(ops), edakit.Pick(ops)
		neg := edakit.Int(0, 1) == 1
		inner := fsApply(op1, a, b)
		if neg {
			inner = fsNot(inner)
		}
		r := fsApply(op2, inner, c)
		expr := edakit.Sprintf("(a %s b)", op1)
		if neg {
			expr = "~" + expr
		}
		q := edakit.Sprintf("With a = %s, b = %s and c = %s in a four-state simulator, what is %s %s c?", fsName[a], fsName[b], fsName[c], expr, op2)
		why := edakit.Sprintf("Evaluate inside out: %s is %s, and combining it with c = %s under %s gives %s. Z behaves as X at a logic operator, 0 controls AND and 1 controls OR.", expr, fsName[inner], fsName[c], op2, fsName[r])
		return fsChoice(q, r, why)
	default:
		// Two continuous assignments drive one wire.
		en := edakit.Int(0, 1)
		a := edakit.Int(0, 1)
		c := edakit.Int(0, 3)
		d1 := vZ
		if en == 1 {
			d1 = a
		}
		r := fsResolve(d1, c)
		q := edakit.Sprintf("A plain wire w has two drivers: assign w = en ? a : 1'bz; and assign w = c;. With en = %d, a = %d and c = %s, what value does w resolve to?", en, a, fsName[c])
		why := edakit.Sprintf("The first driver gives %s. Wire resolution lets a Z driver yield to the other one, keeps agreeing values, and turns a 0 against a 1 (or anything against X) into X, so w = %s.", fsName[d1], fsName[r])
		return fsChoice(q, r, why)
	}
}

// xPessimism compares how different coding styles of a 2:1 mux behave when the select is X.
func xPessimism(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	a, b := edakit.Int(0, 1), edakit.Int(0, 1)
	prev := edakit.Int(0, 1)
	type style struct {
		code string
		res  int
		why  string
	}
	gateRes := vX
	if a == 0 && b == 0 {
		gateRes = v0
	}
	ternRes := vX
	if a == b {
		ternRes = a
	}
	styles := []style{
		{"assign y = (s & a) | (~s & b);", gateRes, "With s = X, s & a and ~s & b are each 0 only when their data input is 0, otherwise X. OR of two X terms is X, so the gate-level mux is X unless both data inputs are 0."},
		{"assign y = s ? a : b;", ternRes, "The conditional operator with an unknown select combines both data inputs bit by bit: equal bits pass through and differing bits become X."},
		{"always @* if (s) y = a; else y = b;", b, "An if statement treats an X condition as false, so the else branch runs and y = b."},
		{"always @* case (s) 1'b1: y = a; default: y = b; endcase", b, "Case items compare with ===, so X matches neither 1'b1 and the default branch assigns b."},
		{edakit.Sprintf("always @* case (s) 1'b0: y = b; 1'b1: y = a; endcase (y held %d before)", prev), prev, "X matches no case item and there is no default, so nothing is assigned and y keeps its previous value."},
	}
	var st style
	if kp == 1 {
		st = styles[edakit.Int(0, 1)]
	} else if kp == 2 {
		st = styles[edakit.Int(2, 4)]
	} else {
		st = styles[edakit.Int(0, 3)]
	}
	if kp != 3 {
		q := edakit.Sprintf("A 2:1 mux is coded as %s The select s is X, a = %d and b = %d. What does a four-state simulator produce for y?", st.code, a, b)
		return fsChoice(q, st.res, st.why)
	}
	// Classify against the hardware: the real output is a or b.
	var verdict, why string
	switch {
	case st.res == vX && a == b:
		verdict, why = "X-pessimistic", edakit.Sprintf("Real hardware outputs %d whichever way s resolves, because both data inputs are %d. The simulator reports X, which is pessimistic.", a, a)
	case st.res != vX && a != b:
		verdict, why = "X-optimistic", edakit.Sprintf("Real hardware could output 0 or 1, depending on how s resolves. The simulator reports a definite %s, hiding the uncertainty, which is optimistic.", fsName[st.res])
	default:
		verdict, why = "Accurate", "The simulated value matches what the hardware can actually produce: a known value when both data inputs agree, X when they differ."
	}
	q := edakit.Sprintf("A 2:1 mux is coded as %s The select s is X, a = %d and b = %d, and the simulator produces y = %s. Compared with real hardware, where s is some unknown 0 or 1, how is this result best described?", st.code, a, b, fsName[st.res])
	return edakit.Choice(q, verdict, why+" "+st.why, otherVerdicts(verdict)...)
}

func otherVerdicts(v string) []string {
	var out []string
	for _, s := range []string{"X-pessimistic", "X-optimistic", "Accurate"} {
		if s != v {
			out = append(out, s)
		}
	}
	return out
}
