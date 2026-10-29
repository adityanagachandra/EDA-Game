package synth

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Tree covering over a NAND2/INV subject tree. Pattern nodes: 'N' nand2, 'I' inverter,
// 'V' a pattern variable that binds any subject node.
type tnode struct {
	op   byte // 'N', 'I', 'L'
	name string
	kids []*tnode
}

type pattern struct {
	name        string
	area, delay int
	tree        *tnode
}

func pN(a, b *tnode) *tnode { return &tnode{op: 'N', kids: []*tnode{a, b}} }
func pI(a *tnode) *tnode    { return &tnode{op: 'I', kids: []*tnode{a}} }
func pV() *tnode            { return &tnode{op: 'V'} }

var fullLibrary = []pattern{
	{"INV", 2, 1, pI(pV())},
	{"NAND2", 3, 2, pN(pV(), pV())},
	{"NOR2", 3, 2, pI(pN(pI(pV()), pI(pV())))},
	{"AND2", 4, 3, pI(pN(pV(), pV()))},
	{"OR2", 4, 3, pN(pI(pV()), pI(pV()))},
	{"NAND3", 4, 3, pN(pV(), pI(pN(pV(), pV())))},
	{"AOI21", 4, 3, pI(pN(pN(pV(), pV()), pI(pV())))},
	{"OAI21", 4, 3, pN(pN(pI(pV()), pI(pV())), pV())},
}

// match returns every way the pattern can bind its variables to subject nodes.
func match(p, s *tnode) [][]*tnode {
	if p.op == 'V' {
		return [][]*tnode{{s}}
	}
	if p.op != s.op {
		return nil
	}
	if p.op == 'I' {
		return match(p.kids[0], s.kids[0])
	}
	var out [][]*tnode
	for _, order := range [][2]int{{0, 1}, {1, 0}} {
		l := match(p.kids[0], s.kids[order[0]])
		r := match(p.kids[1], s.kids[order[1]])
		for _, x := range l {
			for _, y := range r {
				out = append(out, append(append([]*tnode(nil), x...), y...))
			}
		}
	}
	return out
}

// cover returns the minimum cost of covering s: area sums, delay takes the max over leaves.
func cover(s *tnode, lib []pattern, delay bool, memo map[*tnode]int) int {
	if s.op == 'L' {
		return 0
	}
	if v, ok := memo[s]; ok {
		return v
	}
	best := 1 << 30
	for _, p := range lib {
		for _, b := range match(p.tree, s) {
			c := 0
			for _, leaf := range b {
				lc := cover(leaf, lib, delay, memo)
				if delay {
					if lc > c {
						c = lc
					}
				} else {
					c += lc
				}
			}
			if delay {
				c += p.delay
			} else {
				c += p.area
			}
			if c < best {
				best = c
			}
		}
	}
	memo[s] = best
	return best
}

func (s *tnode) str() string {
	switch s.op {
	case 'L':
		return s.name
	case 'I':
		return "INV(" + s.kids[0].str() + ")"
	}
	return "NAND(" + s.kids[0].str() + ", " + s.kids[1].str() + ")"
}

func (s *tnode) count() (inv, nand int) {
	switch s.op {
	case 'L':
		return 0, 0
	case 'I':
		i, n := s.kids[0].count()
		return i + 1, n
	}
	i1, n1 := s.kids[0].count()
	i2, n2 := s.kids[1].count()
	return i1 + i2, n1 + n2 + 1
}

func randSubject(nand int, leaves *int, parentInv bool) *tnode {
	if nand == 0 {
		l := &tnode{op: 'L', name: varNames[*leaves]}
		*leaves++
		if !parentInv && edakit.Int(0, 2) == 0 {
			return pI(l)
		}
		return l
	}
	k := edakit.Int(0, nand-1)
	n := pN(randSubject(k, leaves, false), randSubject(nand-1-k, leaves, false))
	if !parentInv && edakit.Int(0, 1) == 0 {
		return pI(n)
	}
	return n
}

// techMapping: area and delay tree covering.
func techMapping(ctx generator.GeneratorContext) generator.Problem {
	which := kp(ctx)
	leaves := 0
	s := randSubject(edakit.Scaled(ctx, 3, 5), &leaves, false)
	inv, nand := s.count()
	stem := "The subject tree is " + s.str() + ", with a to " + varNames[leaves-1] + " primary inputs."
	switch which {
	case 1:
		lib := fullLibrary[:2]
		a := cover(s, lib, false, map[*tnode]int{})
		return edakit.Numeric(stem+" The library holds only INV (area 2) and NAND2 (area 3). What is the area of the cover?",
			edakit.IntStr(a), edakit.Sprintf("With only the base gates, each subject node is its own match: %d inverters and %d NAND2 gates, area %d.", inv, nand, a),
			map[string]string{"area": edakit.IntStr(a)})
	case 2:
		a := cover(s, fullLibrary, false, map[*tnode]int{})
		return edakit.Numeric(stem+" The library, with each gate's NAND2/INV pattern, is INV 2, NAND2 3, NOR2 3, AND2 4, OR2 4, NAND3 4, AOI21 4, OAI21 4 (areas). Find the minimum-area cover by dynamic programming over the tree. What is its area?",
			edakit.IntStr(a), edakit.Sprintf("Visit nodes bottom up; at each node try every pattern that matches there and add the best costs of the subtrees at its inputs. The root's best cost is %d.", a),
			map[string]string{"area": edakit.IntStr(a)})
	}
	d := cover(s, fullLibrary, true, map[*tnode]int{})
	return edakit.Numeric(stem+" Gate delays are INV 1, NAND2 2, NOR2 2, and 3 for AND2, OR2, NAND3, AOI21 and OAI21, all inputs arriving at time 0, with no load dependence. What is the minimum arrival time at the root over all covers?",
		edakit.IntStr(d), edakit.Sprintf("The same dynamic program with cost = gate delay + the latest input arrival gives %d at the root.", d),
		map[string]string{"delay": edakit.IntStr(d)})
}
