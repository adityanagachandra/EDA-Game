# Floating-point encoding in hardware

Floating point trades a fixed step for a fixed relative precision. IEEE 754 defines the bit layouts and rounding that CPUs, GPUs and accelerators implement, and it is a large share of datapath area in those chips. A tool designer meets it in floating-point unit verification, which is a hard formal problem, and in the reduced formats (half precision and narrower) that machine-learning accelerators now generate in quantity.

## Sign, biased exponent and fraction

An IEEE 754 binary format stores a sign bit $s$, a biased exponent field $E$ and a fraction field $F$. Half precision (binary16) has 5 exponent bits with bias 15 and 10 fraction bits. For a normal number ($1 \le E \le 30$) the value is

$$x = (-1)^s \cdot 2^{E - 15} \cdot \left(1 + \frac{F}{2^{10}}\right).$$

The leading 1 is implicit: it is not stored, which buys one bit of precision. Take $s = 1$, $E = 17$, $F = 512$: the value is $-2^{2} \cdot (1 + 0.5) = -6$. Single precision uses the same scheme with 8 exponent bits, bias 127 and 23 fraction bits; double uses 11, 1023 and 52.

## Encoding a value

To encode a nonzero value, normalize it as $m \cdot 2^{k}$ with $1 \le m < 2$. The exponent field is $E = k + \text{bias}$ and the fraction field is $F = (m - 1) \cdot 2^{10}$ in half precision, rounded to nearest-even if it is not an integer.

Encode 12.5 in half precision. $12.5 = 1.5625 \cdot 2^{3}$, so $E = 3 + 15 = 18$ and $F = 0.5625 \cdot 1024 = 576$. The bias exists so that the exponent field is an unsigned integer and the whole word, for positive numbers, orders the same way as the values do: a hardware comparator for floats is mostly an integer comparator. The largest finite half-precision value is $(2 - 2^{-10}) \cdot 2^{15} = 65504$, which is why half-precision training needs loss scaling.

## Spacing, subnormals and special values

Within one binade (fixed $E$) the representable numbers are evenly spaced. The spacing, one unit in the last place, is $\text{ulp} = 2^{E - 15 - 10}$ in half precision. Between 8 and 16 ($E = 18$) the spacing is $2^{-7}$; between 1024 and 2048 it is 1.

Exponent field $E = 0$ encodes zero and the subnormals, $x = (-1)^s 2^{-14} (F / 2^{10})$, with no implicit 1. They fill the gap below the smallest normal $2^{-14}$ in steps of $2^{-24}$, giving gradual underflow. $E = 31$ encodes infinity ($F = 0$) and NaN ($F \ne 0$). Subnormals cost hardware (a normalizing shifter on input and output), so some accelerators flush them to zero, and an equivalence checker comparing such a unit to an IEEE reference model must be told.

## Going deeper

IEEE Std 754-2019; Goldberg, "What Every Computer Scientist Should Know About Floating-Point Arithmetic" (ACM Computing Surveys, 1991); Muller et al., "Handbook of Floating-Point Arithmetic" (2nd ed., 2018).
