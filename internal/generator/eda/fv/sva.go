package fv

import (
	"math/rand"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func randomWave(t int, pct int) []bool {
	w := make([]bool, t)
	for i := range w {
		w[i] = rand.Intn(100) < pct
	}
	return w
}

func waveText(name string, w []bool) string {
	var b strings.Builder
	for _, x := range w {
		if x {
			b.WriteByte('1')
		} else {
			b.WriteByte('0')
		}
	}
	return name + " = " + b.String()
}

// svaFailures counts failing attempts of a |-> ##[m:n] b; attempts whose window runs past the
// end of the trace without a match are pending and not counted.
func svaFailures(a, b []bool, m, n int) (fails, attempts, pending int) {
	T := len(a)
	for t := 0; t < T; t++ {
		if !a[t] {
			continue
		}
		attempts++
		ok := false
		for u := t + m; u <= t+n && u < T; u++ {
			if b[u] {
				ok = true
				break
			}
		}
		if ok {
			continue
		}
		if t+n >= T {
			pending++
			continue
		}
		fails++
	}
	return
}

func svaAssertions(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	for {
		T := edakit.Scaled(ctx, 8, 14)
		a, b := randomWave(T, 45), randomWave(T, 50)
		var prop string
		var m, n int
		switch kp {
		case 1:
			if rand.Intn(2) == 0 {
				prop, m, n = "a |-> b", 0, 0
			} else {
				prop, m, n = "a |=> b", 1, 1
			}
		case 2:
			d := edakit.Int(2, 3)
			prop, m, n = edakit.Sprintf("a |-> ##%d b", d), d, d
		default:
			m = edakit.Int(1, 2)
			n = m + edakit.Int(1, 2)
			prop = edakit.Sprintf("a |-> ##[%d:%d] b", m, n)
		}
		fails, attempts, pending := svaFailures(a, b, m, n)
		if attempts < 2 || (fails == 0 && rand.Intn(3) > 0) {
			continue
		}
		q := edakit.Sprintf("Sampled values on cycles 0 to %d (leftmost is cycle 0): %s and %s. The assertion is assert property (@(posedge clk) %s). An attempt starts at every cycle where a is 1; attempts still waiting when the trace ends are pending and not counted. How many attempts fail?",
			T-1, waveText("a", a), waveText("b", b), prop)
		return edakit.Numeric(q, edakit.IntStr(fails),
			edakit.Sprintf("a is 1 on %d cycles, so %d attempts start. Checking b in each attempt's window gives %d failures, and %d attempt(s) are still pending at the end of the trace.", attempts, attempts, fails, pending),
			map[string]string{"attempts": edakit.IntStr(attempts), "fails": edakit.IntStr(fails), "pending": edakit.IntStr(pending)})
	}
}

type seqElem struct {
	sig      []bool
	min, max int
}

// matchFrom reports whether elems[i:] matches starting at cycle t, with ##1 between elements.
func matchFrom(elems []seqElem, i, t int) bool {
	if i == len(elems) {
		return true
	}
	e := elems[i]
	for r := 1; r <= e.max; r++ {
		u := t + r - 1
		if u >= len(e.sig) || !e.sig[u] {
			return false
		}
		if r >= e.min && matchFrom(elems, i+1, u+1) {
			return true
		}
	}
	return false
}

func svaSequences(ctx generator.GeneratorContext) generator.Problem {
	kp := ctx.KPOrdinal
	if kp == 0 {
		kp = edakit.Int(1, 3)
	}
	T := edakit.Scaled(ctx, 8, 14)
	if kp == 1 {
		a := randomWave(T, 50)
		fn := edakit.Pick([]string{"$rose(a)", "$fell(a)", "$stable(a)"})
		cnt := 0
		prev := false
		for _, x := range a {
			switch fn {
			case "$rose(a)":
				if x && !prev {
					cnt++
				}
			case "$fell(a)":
				if !x && prev {
					cnt++
				}
			default:
				if x == prev {
					cnt++
				}
			}
			prev = x
		}
		name := strings.TrimPrefix(fn, "$")
		q := edakit.Sprintf("Sampled values of a on cycles 0 to %d (leftmost is cycle 0): %s. Assume a was 0 before cycle 0. On how many cycles is the SVA system function %s (spelled with a leading dollar sign in source) true?", T-1, waveText("a", a), name)
		return edakit.Numeric(q, edakit.IntStr(cnt),
			edakit.Sprintf("The function compares the value sampled on each cycle with the value sampled on the previous one. Scanning the trace pair by pair gives %d cycles.", cnt),
			map[string]string{"count": edakit.IntStr(cnt)})
	}
	for {
		a, b, c := randomWave(T, 50), randomWave(T, 55), randomWave(T, 50)
		var elems []seqElem
		var seq string
		waves := waveText("a", a) + ", " + waveText("b", b)
		if kp == 2 {
			switch rand.Intn(3) {
			case 0:
				seq, elems = "a ##1 b", []seqElem{{a, 1, 1}, {b, 1, 1}}
			case 1:
				seq, elems = "a ##1 b ##1 c", []seqElem{{a, 1, 1}, {b, 1, 1}, {c, 1, 1}}
				waves += ", " + waveText("c", c)
			default:
				// a ##2 b: a, then any cycle, then b
				any := make([]bool, T)
				for i := range any {
					any[i] = true
				}
				seq, elems = "a ##2 b", []seqElem{{a, 1, 1}, {any, 1, 1}, {b, 1, 1}}
			}
		} else {
			waves += ", " + waveText("c", c)
			if rand.Intn(2) == 0 {
				seq, elems = "a ##1 b[*2] ##1 c", []seqElem{{a, 1, 1}, {b, 2, 2}, {c, 1, 1}}
			} else {
				seq, elems = "a ##1 b[*1:2] ##1 c", []seqElem{{a, 1, 1}, {b, 1, 2}, {c, 1, 1}}
			}
		}
		cnt := 0
		var starts []string
		for t := 0; t < T; t++ {
			if matchFrom(elems, 0, t) {
				cnt++
				starts = append(starts, edakit.IntStr(t))
			}
		}
		if cnt == 0 && rand.Intn(3) > 0 {
			continue
		}
		q := edakit.Sprintf("Sampled values on cycles 0 to %d (leftmost is cycle 0): %s. At how many start cycles does the sequence %s have at least one match that ends within the trace?", T-1, waves, seq)
		why := "No start cycle matches."
		if cnt > 0 {
			why = "Matches start at cycles " + strings.Join(starts, ", ") + ". Each ##1 moves one cycle forward, and b[*m:n] needs b on m to n consecutive cycles."
		}
		return edakit.Numeric(q, edakit.IntStr(cnt), why, map[string]string{"matches": edakit.IntStr(cnt)})
	}
}
