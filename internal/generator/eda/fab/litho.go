package fab

import (
	"math"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

type scanner struct {
	name string
	lam  float64
	na   float64
}

var scanners = []scanner{
	{"a KrF scanner (248 nm, NA 0.80)", 248, 0.80},
	{"a dry ArF scanner (193 nm, NA 0.93)", 193, 0.93},
	{"an ArF immersion scanner (193 nm, NA 1.35)", 193, 1.35},
	{"an EUV scanner (13.5 nm, NA 0.33)", 13.5, 0.33},
	{"a high-NA EUV scanner (13.5 nm, NA 0.55)", 13.5, 0.55},
}

// resolution: Rayleigh half-pitch, the k1 a pitch demands, and the NA a target demands.
func resolution(ctx generator.GeneratorContext) generator.Problem {
	s := edakit.Pick(scanners)
	switch kpOf(ctx) {
	case 1:
		k1 := step(0.28, 0.60, 0.01)
		hp := k1 * s.lam / s.na
		ans := edakit.Num(hp, 1)
		q := edakit.Sprintf("Using the Rayleigh criterion $R = k_1 \\lambda / NA$ with $k_1 = %s$, what minimum half-pitch can %s print, in nanometres? Round to 1 decimal place.", edakit.Num(k1, 2), s.name)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$R = %s \\times %s / %s = %s$ nm, so the densest printable line pitch is twice that.", edakit.Num(k1, 2), edakit.Num(s.lam, 1), edakit.Num(s.na, 2), ans),
			map[string]string{"half_pitch": ans})
	case 2:
		hp := s.lam / s.na * step(0.26, 0.6, 0.01)
		pitch := math.Round(2 * hp)
		k1 := (pitch / 2) * s.na / s.lam
		ans := edakit.Num(k1, 3)
		q := edakit.Sprintf("A layer needs a line pitch of %.0f nm printed in one exposure on %s. What $k_1$ factor does that demand? Round to 3 decimal places.", pitch, s.name)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Half-pitch is %s nm, and $k_1 = R \\cdot NA / \\lambda = %s \\times %s / %s = %s$. Values near the 0.25 single-exposure limit need aggressive resolution enhancement.", edakit.Num(pitch/2, 1), edakit.Num(pitch/2, 1), edakit.Num(s.na, 2), edakit.Num(s.lam, 1), ans),
			map[string]string{"k1": ans})
	default:
		lam := edakit.Pick([]float64{193, 13.5})
		k1 := step(0.30, 0.45, 0.01)
		var hp float64
		if lam == 193 {
			hp = float64(edakit.Int(40, 90))
		} else {
			hp = float64(edakit.Int(8, 16))
		}
		na := k1 * lam / hp
		ans := edakit.Num(na, 3)
		q := edakit.Sprintf("A process team can run at $k_1 = %s$ with %s nm light. What numerical aperture is needed to print a %.0f nm half-pitch? Round to 3 decimal places.", edakit.Num(k1, 2), edakit.Num(lam, 1), hp)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$NA = k_1 \\lambda / R = %s \\times %s / %.0f = %s$. Dry optics stop below 1 and water immersion near 1.35; EUV optics are at 0.33 and 0.55.", edakit.Num(k1, 2), edakit.Num(lam, 1), hp, ans),
			map[string]string{"NA": ans})
	}
}

// depthOfFocus: k2 lambda / NA^2, its scaling with NA, and the focus budget.
func depthOfFocus(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		s := edakit.Pick(scanners)
		k2 := step(0.5, 1.0, 0.05)
		dof := k2 * s.lam / (s.na * s.na)
		ans := edakit.Num(dof, 1)
		q := edakit.Sprintf("With $DOF = k_2 \\lambda / NA^2$ and $k_2 = %s$, what is the depth of focus of %s, in nanometres? Round to 1 decimal place.", edakit.Num(k2, 2), s.name)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$DOF = %s \\times %s / %s^2 = %s$ nm.", edakit.Num(k2, 2), edakit.Num(s.lam, 1), edakit.Num(s.na, 2), ans),
			map[string]string{"dof": ans})
	case 2:
		na1 := edakit.Pick([]float64{0.33, 0.75, 0.85, 0.93, 1.2})
		na2 := na1 + step(0.1, 0.4, 0.02)
		r := (na2 * na2) / (na1 * na1)
		ans := edakit.Num(r, 2)
		q := edakit.Sprintf("At fixed wavelength and $k_2$, a scanner's numerical aperture rises from %s to %s. By what factor does the depth of focus shrink (old DOF divided by new DOF)? Round to 2 decimal places.", edakit.Num(na1, 2), edakit.Num(na2, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("DOF goes as $1/NA^2$, so the factor is $(%s/%s)^2 = %s$, while resolution improves only by $%s/%s = %s$.", edakit.Num(na2, 2), edakit.Num(na1, 2), ans, edakit.Num(na2, 2), edakit.Num(na1, 2), edakit.Num(na2/na1, 2)),
			map[string]string{"factor": ans})
	default:
		dof := float64(10 * edakit.Int(6, 20))
		flat := float64(edakit.Int(10, 30))
		focus := float64(edakit.Int(10, 30))
		ans := edakit.Num(dof-flat-focus, 0)
		if dof-flat-focus <= 0 {
			dof = flat + focus + float64(10*edakit.Int(2, 6))
			ans = edakit.Num(dof-flat-focus, 0)
		}
		q := edakit.Sprintf("A critical layer prints with %.0f nm of usable depth of focus. Wafer flatness inside the field uses %.0f nm and scanner focus control uses %.0f nm. How much topography may the underlying CMP leave across the field, in nanometres?", dof, flat, focus)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Focus errors stack linearly in this budget: %.0f - %.0f - %.0f = %s nm. That remainder becomes a planarity requirement, which reaches design as density and fill rules.", dof, flat, focus, ans),
			map[string]string{"budget": ans})
	}
}

// photomask: 4x reduction, the mask error enhancement factor and dies per reticle field.
func photomask(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		cd := float64(edakit.Int(10, 120))
		if edakit.Int(0, 1) == 0 {
			ans := edakit.Num(4*cd, 0)
			q := edakit.Sprintf("A scanner projects the mask at 4x reduction. A line must print %.0f nm wide on the wafer. How wide is it on the mask, in nanometres (ignoring OPC bias)?", cd)
			return edakit.Numeric(q, ans, edakit.Sprintf("Mask dimensions are 4 times wafer dimensions: %s nm.", ans), map[string]string{"mask_nm": ans})
		}
		ans := edakit.Num(cd, 0)
		q := edakit.Sprintf("A 4x reduction mask carries a %.0f nm wide line (ignoring OPC bias). What width does it print on the wafer, in nanometres? Round to the nearest nanometre.", 4*cd)
		return edakit.Numeric(q, ans, edakit.Sprintf("Divide by the reduction ratio: %.0f / 4 = %s nm.", 4*cd, ans), map[string]string{"wafer_nm": ans})
	case 2:
		meef := step(1.0, 4.0, 0.1)
		dm := step(0.4, 4.0, 0.2)
		dw := meef * dm / 4
		ans := edakit.Num(dw, 3)
		q := edakit.Sprintf("A mask line is %s nm too wide at mask scale. The mask error enhancement factor for this feature is %s and the reduction is 4x. How much too wide does it print on the wafer, in nanometres? Round to 3 decimal places.", edakit.Num(dm, 1), edakit.Num(meef, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\Delta CD_{wafer} = MEEF \\cdot \\Delta CD_{mask}/4 = %s \\times %s / 4 = %s$ nm. Without enhancement it would be only %s nm.", edakit.Num(meef, 1), edakit.Num(dm, 1), ans, edakit.Num(dm/4, 3)),
			map[string]string{"dcd": ans})
	default:
		half := edakit.Int(0, 1) == 1
		fh := 33.0
		name := "a full 26 mm by 33 mm scanner field"
		if half {
			fh = 16.5
			name = "a high-NA half field of 26 mm by 16.5 mm"
		}
		w := float64(edakit.Int(3, 26))
		h := step(3, fh, 0.5)
		n := int(math.Floor(26/w)) * int(math.Floor(fh/h))
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("A die measures %s mm (along the 26 mm side) by %s mm. How many whole dies fit in %s, without rotating the die?", edakit.Num(w, 1), edakit.Num(h, 1), name)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$\\lfloor 26/%s \\rfloor \\times \\lfloor %s/%s \\rfloor = %d \\times %d = %s$ dies per exposure.", edakit.Num(w, 1), edakit.Num(fh, 1), edakit.Num(h, 1), int(math.Floor(26/w)), int(math.Floor(fh/h)), ans),
			map[string]string{"per_field": ans})
	}
}

// modelOPC: edge placement error and the fixed-point OPC iteration.
func modelOPC(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 3, 6)
		var parts []string
		sum := 0.0
		for i := 0; i < n; i++ {
			e := float64(edakit.Int(-8, 8))
			sum += e * e
			parts = append(parts, edakit.Num(e, 0))
		}
		rms := math.Sqrt(sum / float64(n))
		ans := edakit.Num(rms, 2)
		q := edakit.Sprintf("After one OPC pass, the simulated edge placement errors (printed minus target, in nm) at the %d fragments of a polygon are %s. What is the RMS edge placement error, in nanometres? Round to 2 decimal places.", n, strings.Join(parts, ", "))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Square, average and take the root: $\\sqrt{%s/%d} = %s$ nm. OPC tools report both RMS and maximum EPE because a single bad fragment can still cause a bridge.", edakit.Num(sum, 0), n, ans),
			map[string]string{"sum_sq": edakit.Num(sum, 0)})
	case 2:
		e0 := float64(edakit.Int(3, 15)) * edakit.Pick([]float64{-1, 1})
		g := step(0.6, 1.6, 0.1)
		f := step(0.3, 1.0, 0.1)
		e1 := e0 * (1 - f*g)
		ans := edakit.Num(e1, 2)
		q := edakit.Sprintf("Near its current position, moving a mask fragment by $b$ nm (wafer scale) moves the printed edge by $%s\\,b$ nm. The fragment's EPE is %s nm. OPC updates the fragment by $-f \\cdot EPE$ with feedback $f = %s$. What is the EPE after this update, in nanometres? Round to 2 decimal places.", edakit.Num(g, 1), edakit.Num(e0, 0), edakit.Num(f, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The move is $%s \\times %s = %s$ nm, which shifts the edge by %s nm. New EPE $= %s(1 - %s \\times %s) = %s$ nm.", edakit.Num(-f, 1), edakit.Num(e0, 0), edakit.Num(-f*e0, 2), edakit.Num(-f*e0*g, 2), edakit.Num(e0, 0), edakit.Num(f, 1), edakit.Num(g, 1), ans),
			map[string]string{"factor": edakit.Num(1-f*g, 2)})
	default:
		e0 := float64(edakit.Int(6, 20))
		var g, f, r float64
		for {
			g = step(0.5, 1.8, 0.1)
			f = step(0.2, 0.9, 0.1)
			r = math.Abs(1 - f*g)
			if r > 0.08 && r < 0.9 {
				break
			}
		}
		tol := step(0.2, 1.0, 0.1)
		n, e := 0, e0
		for e > tol {
			e *= r
			n++
		}
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("A fragment starts with an EPE of %.0f nm. The printed edge moves %s nm per nm of fragment move, and each OPC iteration moves the fragment by $-%s \\cdot EPE$. How many iterations until $|EPE| \\le %s$ nm?", e0, edakit.Num(g, 1), edakit.Num(f, 1), edakit.Num(tol, 1))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each iteration multiplies the EPE by $1 - %s \\times %s = %s$, so its magnitude shrinks by a factor of %s per step. Starting at %.0f nm it first reaches %s nm or less after %s iterations.", edakit.Num(f, 1), edakit.Num(g, 1), edakit.Num(1-f*g, 2), edakit.Num(r, 2), e0, edakit.Num(tol, 1), ans),
			map[string]string{"contraction": edakit.Num(r, 2)})
	}
}

// euv: EUV half-pitch, matching EUV with ArF multiple patterning, and photon shot noise.
func euv(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		na := edakit.Pick([]float64{0.33, 0.55})
		k1 := step(0.30, 0.50, 0.01)
		hp := k1 * 13.5 / na
		ans := edakit.Num(hp, 2)
		q := edakit.Sprintf("An EUV scanner uses 13.5 nm light at NA %s. At $k_1 = %s$, what half-pitch does one exposure print, in nanometres? Round to 2 decimal places.", edakit.Num(na, 2), edakit.Num(k1, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("$R = %s \\times 13.5 / %s = %s$ nm. The same $k_1$ on 193 nm immersion at NA 1.35 gives %s nm.", edakit.Num(k1, 2), edakit.Num(na, 2), ans, edakit.Num(k1*193/1.35, 1)),
			map[string]string{"hp": ans})
	case 2:
		k1 := step(0.28, 0.35, 0.01)
		pArF := 2 * k1 * 193 / 1.35
		pitch := float64(edakit.Int(20, 48))
		n := int(math.Ceil(pArF/pitch - 1e-9))
		if n < 2 {
			pitch = math.Floor(pArF / 2)
			n = int(math.Ceil(pArF/pitch - 1e-9))
		}
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("An EUV tool prints a %.0f nm line pitch in one exposure. On 193 nm immersion (NA 1.35) the team can run at $k_1 = %s$. How many litho-etch exposures (LE to the n) would 193 nm immersion need for the same pitch?", pitch, edakit.Num(k1, 2))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Single-exposure ArF pitch is $2 k_1 \\lambda/NA = %s$ nm. Splitting into n masks relaxes each to $n \\times %.0f$ nm, so $n = \\lceil %s / %.0f \\rceil = %s$.", edakit.Num(pArF, 1), pitch, edakit.Num(pArF, 1), pitch, ans),
			map[string]string{"pitch_arf": edakit.Num(pArF, 1)})
	default:
		dose := float64(5 * edakit.Int(4, 16))
		a := float64(edakit.Int(5, 20))
		eph := 1239.84 / 13.5 * 1.602177e-19 // J
		perNm2 := dose * 1e-3 / 1e14 / eph
		N := perNm2 * a * a
		noise := 100 / math.Sqrt(N)
		ans := edakit.Num(noise, 2)
		q := edakit.Sprintf("An EUV exposure delivers %.0f mJ/cm^2 at 13.5 nm (photon energy about 91.8 eV). Treat the photon count in a square pixel %.0f nm by %.0f nm as Poisson. What is the relative fluctuation $1/\\sqrt{N}$ of that count, in percent? Round to 2 decimal places.", dose, a, a)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("One photon carries %.4g J. The dose is %s photons per nm^2, so the pixel receives $N = %s$ photons and $1/\\sqrt{N} = %s$%%. Fewer, higher-energy photons are why EUV is shot-noise limited.", eph, edakit.Num(perNm2, 2), edakit.Num(N, 0), ans),
			map[string]string{"N": edakit.Num(N, 0)})
	}
}

// multiPatterning: per-mask pitch, exposures needed and mask counts.
func multiPatterning(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		p := float64(edakit.Int(18, 64))
		n := edakit.Pick([]int{2, 3, 4})
		name := map[int]string{2: "LELE (two exposures)", 3: "LELELE (three exposures)", 4: "four litho-etch exposures"}[n]
		ans := edakit.Num(p*float64(n), 0)
		q := edakit.Sprintf("A metal layer has a %.0f nm line pitch and is split by %s, with lines assigned to masks in rotation. What line pitch must each individual mask print, in nanometres?", p, name)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Each mask takes one line in every %d, so its pitch is $%d \\times %.0f = %s$ nm. The combined pattern interleaves them.", n, n, p, ans),
			map[string]string{"per_mask": ans})
	case 2:
		k1 := step(0.27, 0.35, 0.01)
		pmin := 2 * k1 * 193 / 1.35
		p := float64(edakit.Int(30, 70))
		n := int(math.Ceil(pmin/p - 1e-9))
		if n < 1 {
			n = 1
		}
		ans := edakit.IntStr(n)
		q := edakit.Sprintf("On 193 nm immersion (NA 1.35) a single exposure resolves at $k_1 = %s$. A layer needs a %.0f nm pitch. What is the fewest litho-etch exposures that can print it by pitch splitting?", edakit.Num(k1, 2), p)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Minimum single-exposure pitch is $2 k_1 \\lambda/NA = %s$ nm. Each of $n$ masks prints pitch $n p$, so need $n p \\ge %s$: $n = %s$.", edakit.Num(pmin, 1), edakit.Num(pmin, 1), ans),
			map[string]string{"pmin": edakit.Num(pmin, 1)})
	default:
		type sch struct {
			name  string
			masks int
			cut   bool
		}
		schemes := []sch{{"single exposure", 1, false}, {"LELE", 2, false}, {"LELELE", 3, false}, {"SADP", 1, true}, {"SAQP", 1, true}}
		nl := edakit.Scaled(ctx, 3, 7)
		total := 0
		var parts []string
		for i := 0; i < nl; i++ {
			s := edakit.Pick(schemes)
			m := s.masks
			desc := edakit.Sprintf("layer %d by %s", i+1, s.name)
			if s.cut {
				c := edakit.Int(1, 2)
				m += c
				desc += edakit.Sprintf(" with %d cut or block mask", c)
				if c > 1 {
					desc += "s"
				}
			}
			total += m
			parts = append(parts, desc)
		}
		ans := edakit.IntStr(total)
		q := "Count the photomasks for this interconnect module: " + strings.Join(parts, "; ") +
			". Each litho-etch pass needs one mask, SADP and SAQP need one mandrel mask plus their cut or block masks. How many masks in total?"
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Add one mask per exposure, and mandrel plus cuts for each self-aligned layer: %s masks. Mask count drives cost and cycle time, which is what a DTCO study trades against pitch.", ans),
			map[string]string{"masks": ans})
	}
}

// chromatic returns the chromatic number of a small graph by backtracking.
func chromatic(n int, adj [][]bool) int {
	col := make([]int, n)
	var try func(v, k int) bool
	try = func(v, k int) bool {
		if v == n {
			return true
		}
		for c := 1; c <= k; c++ {
			ok := true
			for u := 0; u < v; u++ {
				if adj[u][v] && col[u] == c {
					ok = false
					break
				}
			}
			if ok {
				col[v] = c
				if try(v+1, k) {
					return true
				}
			}
		}
		return false
	}
	for k := 1; k <= n; k++ {
		if try(0, k) {
			return k
		}
	}
	return n
}

// minMono is the fewest same-mask conflicts over all two-mask assignments.
func minMono(n int, edges [][2]int) int {
	best := len(edges)
	for m := 0; m < 1<<n; m++ {
		c := 0
		for _, e := range edges {
			if (m>>e[0])&1 == (m>>e[1])&1 {
				c++
			}
		}
		if c < best {
			best = c
		}
	}
	return best
}

func randomGraph(n int, p float64) ([][2]int, [][]bool) {
	adj := make([][]bool, n)
	for i := range adj {
		adj[i] = make([]bool, n)
	}
	var edges [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if float64(edakit.Int(0, 99)) < p*100 {
				edges = append(edges, [2]int{i, j})
				adj[i][j], adj[j][i] = true, true
			}
		}
	}
	if len(edges) == 0 {
		edges = append(edges, [2]int{0, 1})
		adj[0][1], adj[1][0] = true, true
	}
	return edges, adj
}

func edgeList(edges [][2]int) string {
	var s []string
	for _, e := range edges {
		s = append(s, edakit.Sprintf("%c%c", 'A'+e[0], 'A'+e[1]))
	}
	return strings.Join(s, ", ")
}

// decomposition: conflict graphs for multiple patterning.
func decomposition(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		n := edakit.Scaled(ctx, 4, 7)
		d := float64(10 * edakit.Int(5, 9))
		xs := make([]float64, n)
		x := 0.0
		var parts []string
		for i := 0; i < n; i++ {
			xs[i] = x
			parts = append(parts, edakit.Sprintf("%c at %.0f", 'A'+i, x))
			x += float64(10 * edakit.Int(2, 8))
		}
		cnt := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if xs[j]-xs[i] < d {
					cnt++
				}
			}
		}
		ans := edakit.IntStr(cnt)
		q := edakit.Sprintf("Small via shapes sit on one line at x positions (nm): %s. Two shapes on the same mask must be at least %.0f nm apart, centre to centre. How many conflict edges does the decomposition graph have?", strings.Join(parts, ", "), d)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Every pair closer than %.0f nm must go on different masks, so each such pair is an edge, including non-neighbours. There are %s.", d, ans),
			map[string]string{"edges": ans})
	case 2:
		n := edakit.Scaled(ctx, 4, 7)
		edges, adj := randomGraph(n, step(0.35, 0.7, 0.05))
		k := chromatic(n, adj)
		ans := edakit.IntStr(k)
		q := edakit.Sprintf("A layout's conflict graph has features A to %c and conflict edges %s. Each edge joins two features that cannot share a mask. What is the fewest masks (litho-etch exposures) that decompose this layout without stitching?", 'A'+n-1, edgeList(edges))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("The fewest masks is the chromatic number of the conflict graph, here %s. An odd cycle forces at least 3; a 4-clique forces 4.", ans),
			map[string]string{"chi": ans})
	default:
		n := edakit.Scaled(ctx, 4, 7)
		edges, _ := randomGraph(n, step(0.4, 0.75, 0.05))
		m := minMono(n, edges)
		ans := edakit.IntStr(m)
		q := edakit.Sprintf("A layer must be split into two masks (LELE). Its conflict graph has features A to %c and edges %s. In the best two-mask assignment, how many conflict edges remain with both ends on the same mask and must be fixed by a stitch or a layout change?", 'A'+n-1, edgeList(edges))
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Trying all two-colourings, the fewest monochromatic edges is %s. Zero means the graph is bipartite (no odd cycle); each odd cycle needs at least one fix.", ans),
			map[string]string{"min_conflicts": ans})
	}
}

// sadp: spacer patterning pitch arithmetic.
func sadp(ctx generator.GeneratorContext) generator.Problem {
	switch kpOf(ctx) {
	case 1:
		quad := edakit.Int(0, 1) == 1
		P := float64(8 * edakit.Int(8, 20))
		if quad {
			ans := edakit.Num(P/4, 0)
			q := edakit.Sprintf("SAQP starts from a mandrel pitch of %.0f nm and applies the spacer process twice. What final line pitch results, in nanometres?", P)
			return edakit.Numeric(q, ans, edakit.Sprintf("Each spacer pass halves the pitch: %.0f / 4 = %s nm.", P, ans), map[string]string{"pitch": ans})
		}
		ans := edakit.Num(P/2, 0)
		q := edakit.Sprintf("SADP starts from printed mandrels at a pitch of %.0f nm. Spacers form on both sidewalls and the mandrels are removed. What final line pitch results, in nanometres?", P)
		return edakit.Numeric(q, ans, edakit.Sprintf("Each mandrel yields two spacer lines per mandrel pitch: %.0f / 2 = %s nm.", P, ans), map[string]string{"pitch": ans})
	case 2:
		P := float64(4 * edakit.Int(16, 40))
		s := float64(edakit.Int(8, int(P/5)))
		m := float64(edakit.Int(int(s), int(P-3*s)))
		gap := P - m - 2*s
		if gap <= 0 {
			m = P/2 - s
			gap = P - m - 2*s
		}
		ans := edakit.Num(gap, 0)
		q := edakit.Sprintf("In spacer-is-line SADP, mandrels of width %.0f nm sit at a pitch of %.0f nm and the spacer film is %.0f nm thick. After mandrel removal, spaces alternate between the gap where a mandrel stood and the gap between spacers of neighbouring mandrels. How wide is the second gap, in nanometres?", m, P, s)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Between two mandrels the space is $P - m = %.0f$ nm, and two spacers of %.0f nm line it, leaving $%.0f - %.0f - 2 \\times %.0f = %s$ nm. The other gap equals the mandrel width, %.0f nm.", P-m, s, P, m, s, ans, m),
			map[string]string{"gap_a": edakit.Num(m, 0)})
	default:
		P := float64(4 * edakit.Int(16, 40))
		s := float64(edakit.Int(8, int(P/5)))
		ideal := P/2 - s
		m := ideal + float64(edakit.Int(-6, 6))
		if m < 4 {
			m = ideal + 2
		}
		if edakit.Int(0, 2) == 0 {
			ans := edakit.Num(ideal, 1)
			q := edakit.Sprintf("Spacer-is-line SADP with a mandrel pitch of %.0f nm and a %.0f nm spacer. What mandrel width makes the final lines equally spaced, in nanometres? Round to 1 decimal place.", P, s)
			return edakit.Numeric(q, ans,
				edakit.Sprintf("Equal spaces need the mandrel gap $m$ to equal the inter-mandrel gap $P - m - 2s$, so $m = P/2 - s = %s$ nm.", ans),
				map[string]string{"ideal": ans})
		}
		walk := math.Abs(2*m + 2*s - P)
		ans := edakit.Num(walk, 0)
		q := edakit.Sprintf("Spacer-is-line SADP with mandrel pitch %.0f nm, spacer thickness %.0f nm and printed mandrel width %.0f nm. The line-to-line pitches alternate. What is the pitch walk, the difference between the two alternating pitches, in nanometres?", P, s, m)
		return edakit.Numeric(q, ans,
			edakit.Sprintf("Across a mandrel the line pitch is $m + s = %.0f$ nm; between mandrels it is $P - m - s = %.0f$ nm. The difference is $|2m + 2s - P| = %s$ nm, so any mandrel CD error shows up as pitch walk.", m+s, P-m-s, ans),
			map[string]string{"pitch_a": edakit.Num(m+s, 0)})
	}
}

// mdpCompute: variable-shaped-beam shot count after fracturing.
func mdpCompute(ctx generator.GeneratorContext) generator.Problem {
	S := float64(edakit.Pick([]int{250, 500, 1000}))
	n := edakit.Int(2, 4)
	var parts []string
	total, naive, perRect := 0, 0.0, n
	for i := 0; i < n; i++ {
		w := float64(50 * edakit.Int(2, 50))
		h := float64(50 * edakit.Int(2, 50))
		total += int(math.Ceil(w/S-1e-9)) * int(math.Ceil(h/S-1e-9))
		naive += w * h / (S * S)
		parts = append(parts, edakit.Sprintf("%.0f by %.0f", w, h))
	}
	q := edakit.Sprintf("Fracturing turns a mask polygon into rectangles of %s nm (mask scale). A variable-shaped-beam writer exposes rectangles up to %.0f nm by %.0f nm per shot. How many shots does the polygon need?", strings.Join(parts, ", "), S, S)
	return edakit.Choice(q, edakit.IntStr(total),
		edakit.Sprintf("Each rectangle needs $\\lceil w/S \\rceil \\times \\lceil h/S \\rceil$ shots; summing gives %d. Shot count drives mask write time, which is why OPC jogs and curvilinear shapes are expensive on VSB writers.", total),
		edakit.Distinct(total, 3, int(math.Ceil(naive)), perRect, total+n)...)
}
