package hls

import (
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func partitioning(ctx generator.GeneratorContext) generator.Problem {
	switch whichKP(ctx) {
	case 1:
		B := edakit.Int(2, 8)
		N := B * edakit.Int(4, 32)
		i := edakit.Int(0, N-1)
		q := edakit.Sprintf("An array of %d elements, indexed from 0, is partitioned cyclically into %d banks numbered from 0. In which bank does element %d live?", N, B, i)
		return num(q, edakit.IntStr(i%B),
			edakit.Sprintf("Cyclic partitioning deals elements round-robin: bank = i mod B = %d mod %d = %d, at offset floor(%d / %d) = %d inside the bank.", i, B, i%B, i, B, i/B),
			map[string]string{"bank": edakit.IntStr(i % B), "offset": edakit.IntStr(i / B)})
	case 2:
		B := edakit.Int(2, 8)
		N := edakit.Int(3*B, 40*B)
		size := ceilDiv(N, B)
		i := edakit.Int(0, N-1)
		q := edakit.Sprintf("An array of %d elements, indexed from 0, is block-partitioned into %d banks numbered from 0; each bank holds a contiguous run of ceil(%d / %d) elements. In which bank does element %d live?", N, B, N, B, i)
		return num(q, edakit.IntStr(i/size),
			edakit.Sprintf("Each bank holds %d elements, so element %d is in bank floor(%d / %d) = %d.", size, i, i, size, i/size),
			map[string]string{"bank_size": edakit.IntStr(size)})
	default:
		U := edakit.Int(2, 16)
		p := edakit.Int(1, 2)
		B := ceilDiv(U, p)
		q := edakit.Sprintf("A loop is unrolled by %d so each cycle reads the %d consecutive elements a[%d·k] through a[%d·k + %d]. Each bank of a cyclically partitioned array has %d port(s). What is the smallest cyclic partition factor that lets all reads happen in one cycle?", U, U, U, U, U-1, p)
		return num(q, edakit.IntStr(B),
			edakit.Sprintf("Consecutive elements fall in different banks, so B banks serve up to %d × B reads per cycle when the %d reads spread evenly. B = ceil(%d / %d) = %d is the smallest that works.", p, U, U, p, B),
			map[string]string{"factor": edakit.IntStr(B)})
	}
}

func bankConflicts(ctx generator.GeneratorContext) generator.Problem {
	B := edakit.Int(2, 4)
	k := edakit.Scaled(ctx, 3, 6)
	set := map[int]bool{}
	for len(set) < k {
		set[edakit.Int(0, 8)] = true
	}
	var offs []int
	for o := range set {
		offs = append(offs, o)
	}
	sort.Ints(offs)
	var acc []string
	for _, o := range offs {
		if o == 0 {
			acc = append(acc, "a[i]")
		} else {
			acc = append(acc, edakit.Sprintf("a[i+%d]", o))
		}
	}
	i0 := edakit.Int(0, 20)
	cnt := map[int]int{}
	var banks []int
	m := 0
	for _, o := range offs {
		b := (i0 + o) % B
		banks = append(banks, b)
		cnt[b]++
		m = maxInt(m, cnt[b])
	}
	head := edakit.Sprintf("Each iteration of a loop reads %s. The array is partitioned cyclically into %d banks numbered from 0.", strings.Join(acc, ", "), B)
	switch whichKP(ctx) {
	case 1:
		j := edakit.Int(0, len(offs)-1)
		return num(head+edakit.Sprintf(" In the iteration with i = %d, which bank does the read %s use?", i0, acc[j]), edakit.IntStr(banks[j]),
			edakit.Sprintf("The read touches element %d + %d = %d, and %d mod %d = %d.", i0, offs[j], i0+offs[j], i0+offs[j], B, banks[j]),
			map[string]string{"bank": edakit.IntStr(banks[j])})
	case 2:
		return num(head+edakit.Sprintf(" In the iteration with i = %d, what is the largest number of reads that hit one bank?", i0), edakit.IntStr(m),
			edakit.Sprintf("The banks touched are %s, and the most popular bank receives %d reads. Changing i rotates the bank numbers but not the counts.", joinInts(banks), m),
			map[string]string{"max_per_bank": edakit.IntStr(m)})
	default:
		p := edakit.Int(1, 2)
		ii := ceilDiv(m, p)
		return num(head+edakit.Sprintf(" Each bank has %d read port(s). If memory is the only constraint, what is the minimum II of the pipelined loop?", p), edakit.IntStr(ii),
			edakit.Sprintf("In any iteration the busiest bank receives %d reads (banks %s at i = %d). With %d port(s) per bank that takes ceil(%d / %d) = %d cycles, so II ≥ %d.", m, joinInts(banks), i0, p, m, p, ii, ii),
			map[string]string{"max_per_bank": edakit.IntStr(m), "ii": edakit.IntStr(ii)})
	}
}

func dataflow(ctx generator.GeneratorContext) generator.Problem {
	k := edakit.Scaled(ctx, 3, 5)
	ls := make([]int, k)
	is := make([]int, k)
	var p []string
	sumL, mI, arg := 0, 0, 0
	for i := range ls {
		ls[i] = 10 * edakit.Int(5, 60)
		is[i] = ls[i]
		if edakit.Int(0, 1) == 1 {
			is[i] = 10 * edakit.Int(3, ls[i]/10)
		}
		sumL += ls[i]
		if is[i] > mI {
			mI, arg = is[i], i
		}
		p = append(p, edakit.Sprintf("T%d (latency %d, interval %d)", i+1, ls[i], is[i]))
	}
	N := edakit.Int(10, 200)
	head := "A dataflow region chains tasks through ping-pong buffers: " + strings.Join(p, ", ") +
		", in cycles per frame. A task's interval is how often it can accept a new frame. Model the region as starting a new frame every max-interval cycles, with the first frame taking the sum of the latencies."
	tot := sumL + (N-1)*mI
	switch whichKP(ctx) {
	case 1:
		return num(head+" How many cycles separate successive frames at the output in steady state?", edakit.IntStr(mI),
			edakit.Sprintf("The slowest stage sets the rate: T%d with interval %d cycles.", arg+1, mI),
			map[string]string{"interval": edakit.IntStr(mI)})
	case 2:
		return num(head+edakit.Sprintf(" How many cycles does it take to process %d frames?", N), edakit.IntStr(tot),
			edakit.Sprintf("First frame: %d cycles. Each of the other %d frames adds %d. Total %d.", sumL, N-1, mI, tot),
			map[string]string{"total": edakit.IntStr(tot)})
	default:
		seq := N * sumL
		sp := float64(seq) / float64(tot)
		return num(head+edakit.Sprintf(" Without dataflow the tasks run one after another for each frame. What is the speedup of the dataflow version for %d frames? Round to 2 decimal places.", N), edakit.Num(sp, 2),
			edakit.Sprintf("Sequential: %d × %d = %d cycles. Dataflow: %d + %d × %d = %d. Speedup %s; it tends to %d / %d for many frames.", N, sumL, seq, sumL, N-1, mI, tot, edakit.Num(sp, 2), sumL, mI),
			map[string]string{"sequential": edakit.IntStr(seq), "dataflow": edakit.IntStr(tot)})
	}
}

func lineBuffers(ctx generator.GeneratorContext) generator.Problem {
	K := edakit.Pick([]int{3, 3, 5, 7})
	W := edakit.Pick([]int{320, 640, 720, 1024, 1280, 1920, 2048, 3840})
	switch whichKP(ctx) {
	case 1:
		e := (K - 1) * W
		q := edakit.Sprintf("A %d×%d convolution streams an image of width %d pixels in raster order, one pixel per cycle, and reads each pixel from external memory only once. How many pixels must the line buffers hold?", K, K, W)
		return num(q, edakit.IntStr(e),
			edakit.Sprintf("To form a %d-row window the design keeps the previous %d rows: %d × %d = %d pixels. The newest row streams in directly.", K, K-1, K-1, W, e),
			map[string]string{"pixels": edakit.IntStr(e)})
	case 2:
		b := edakit.Pick([]int{8, 10, 12, 16, 24, 32})
		per := ceilDiv(W*b, 18432)
		tot := (K - 1) * per
		q := edakit.Sprintf("A %d×%d stencil on an image %d pixels wide with %d-bit pixels keeps each buffered row in its own block RAMs. Treat one block RAM as 18,432 bits of capacity. How many block RAMs do the line buffers need?", K, K, W, b)
		return num(q, edakit.IntStr(tot),
			edakit.Sprintf("One row is %d × %d = %d bits, which needs ceil(%d / 18432) = %d block RAM(s). With %d buffered rows that is %d.", W, b, W*b, W*b, per, K-1, tot),
			map[string]string{"per_row": edakit.IntStr(per)})
	default:
		c := (K-1)*W + K
		q := edakit.Sprintf("A %d×%d stencil streams an image of width %d in raster order, one pixel per cycle, with pixel (0, 0) arriving in cycle 1. The first valid output needs the full window whose bottom-right pixel is (row %d, column %d). In which cycle does that pixel arrive?", K, K, W, K-1, K-1)
		return num(q, edakit.IntStr(c),
			edakit.Sprintf("Pixel (r, c) arrives in cycle r × W + c + 1, so pixel (%d, %d) arrives in cycle %d × %d + %d + 1 = %d.", K-1, K-1, K-1, W, K-1, c),
			map[string]string{"cycle": edakit.IntStr(c)})
	}
}
