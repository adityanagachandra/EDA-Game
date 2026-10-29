# Input impedance of a terminated line

In sinusoidal steady state a terminated line transforms its load impedance into a different impedance at its input, depending on its electrical length. That transformation is how quarter-wave transformers match impedances, how stubs act as capacitors or inductors, and why a via stub or an unused branch can resonate. Frequency-domain channel and package tools compute it everywhere.

## The input impedance formula

A lossless line of characteristic impedance $Z_0$ and electrical length $\beta\ell$ terminated in $Z_L$ presents
$$Z_{in} = Z_0\frac{Z_L + jZ_0\tan\beta\ell}{Z_0 + jZ_L\tan\beta\ell}.$$
It repeats every half wavelength ($\beta\ell = 180^\circ$), and at any length a matched load gives $Z_{in} = Z_0$.

Example: a 50 ohm line of $30^\circ$ terminated in 100 ohm. $\tan 30^\circ = 0.5774$, so $Z_{in} = 50(100 + 28.87j)/(50 + 57.74j) = 57.14 - 37.12j$ ohm, magnitude 68.1 ohm. The purely resistive load has become partly capacitive. This is why a frequency-domain engine cannot treat a line as a delay alone when the ends are mismatched: the impedance the driver sees varies with frequency.

## The quarter-wave transformer

At $\beta\ell = 90^\circ$ the tangent diverges and the formula reduces to
$$Z_{in} = \frac{Z_0^2}{Z_L}.$$
A quarter-wave section therefore inverts impedances. To match a load $Z_2$ to a line $Z_1$ at one frequency, choose a section with $Z_T = \sqrt{Z_1Z_2}$.

Example: a 50 ohm quarter-wave section terminated in 100 ohm presents 25 ohm. To match a 100 ohm load to a 50 ohm system, use a section of $\sqrt{5000} = 70.7$ ohm. The match is exact only at the design frequency and its odd multiples; its bandwidth shrinks as the impedance ratio grows, so broadband matching uses several sections. The inversion also explains why an open stub a quarter-wave long looks like a short at its input, the basis of stub notches.

## Shorted and open stubs

Setting $Z_L = 0$ gives a shorted stub with $Z_{in} = jZ_0\tan\beta\ell$, purely reactive: inductive for $0 < \beta\ell < 90^\circ$ and capacitive between $90^\circ$ and $180^\circ$. Setting $Z_L \to \infty$ gives an open stub with $Z_{in} = -jZ_0\cot\beta\ell$: capacitive for short lengths.

Example: a 50 ohm shorted stub of $30^\circ$ has $X = 50\tan 30^\circ = 28.9$ ohm, inductive; an open stub of $30^\circ$ has $X = -50\cot 30^\circ = -86.6$ ohm, capacitive. Microwave designers build tuning elements this way. Signal-integrity engineers meet the unwanted version: a via barrel extending past the layer where the signal exits is an open stub; at the frequency where it is a quarter wave long it presents a short to the through path and carves a deep notch in $|S_{21}|$.

## Going deeper

Pozar, "Microwave Engineering" (4th ed.), chapters 2 and 5; Ramo, Whinnery and Van Duzer, "Fields and Waves in Communication Electronics" (3rd ed.), chapter 5.
