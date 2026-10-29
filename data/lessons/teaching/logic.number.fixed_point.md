# Fixed-point formats

Signal-processing and machine-learning accelerators rarely use floating point in their inner loops. They use fixed point: an integer datapath with an implied binary point. Choosing word widths is a numerical-analysis problem (range against resolution against area), and high-level synthesis tools, quantization flows and bit-accurate simulators all need to compute the exact consequences of a format choice.

## An integer with an implied binary point

A signed fixed-point format with $w$ total bits and $f$ fractional bits stores an integer $k$ in two's complement and means the value $x = k / 2^f$. The hardware never sees the binary point; it only shifts the interpretation.

Take $w = 8$, $f = 4$ and the pattern `11010110`. As a two's-complement integer it is -42. Its value is $-42 / 16 = -2.625$. Check with weights: the top bit is worth $-2^{3}$, then $2^{2}, 2^{1}, 2^{0}, 2^{-1}, \dots, 2^{-4}$: $-8 + 4 + 1 + 0.25 + 0.125 = -2.625$. Writers abbreviate formats as Q-notation, but conventions differ on whether the sign bit is counted in the integer part, so always state $w$ and $f$ explicitly.

## Range and resolution

The format's resolution (its step) is $2^{-f}$. Its range is $[-2^{w-1}, 2^{w-1}-1]$ scaled by $2^{-f}$, that is from $-2^{w-1-f}$ to $2^{w-1-f} - 2^{-f}$. With $w = 12$ and $f = 8$ the step is $1/256 = 0.00390625$, the most negative value is -8 and the largest is $8 - 1/256 = 7.99609375$.

Every bit moved from the integer side to the fractional side halves both the step and the range. Rounding a real number to the nearest representable value makes an error of at most half a step, $2^{-(f+1)}$; truncating (dropping bits) makes an error up to a full step and biases the result downward. Word-length optimization in HLS chooses $f$ per signal so the accumulated error stays inside a budget at minimum area.

## Formats grow through multiplication and addition

Multiplying a $(w_1, f_1)$ value by a $(w_2, f_2)$ value gives an exact product with $w_1 + w_2$ bits and $f_1 + f_2$ fractional bits, because the integers multiply and the scale factors $2^{-f_1}$ and $2^{-f_2}$ multiply. A 16-bit, 15-fraction coefficient times a 12-bit, 8-fraction sample is exact in 28 bits with 23 fractional bits.

Addition needs the binary points aligned first (shift the operand with fewer fractional bits left), and summing $N$ terms grows the integer part by $\lceil \log_2 N \rceil$ bits to stay overflow-free. A filter with 32 taps therefore needs 5 guard bits in its accumulator. After the accumulation the design rounds or truncates back down, and that is where the quantization noise enters. Bit-accurate C models exist precisely so that the RTL and the algorithm agree on these widths.

## Going deeper

Oppenheim and Schafer, "Discrete-Time Signal Processing" (3rd ed., 2010), chapter 6 on quantization effects; Padgett and Anderson, "Fixed-Point Signal Processing" (2009).
