# Synchronizer MTBF

A synchronizer is a chain of flip-flops that gives a possibly metastable sample time to resolve before anyone uses it. Its quality is measured as mean time between failures (MTBF). The MTBF formula combines the event rate and the exponential resolution law, and it is the calculation behind every clock-domain-crossing sign-off: a design with hundreds of synchronizers must still meet a system-level reliability target.

## The MTBF formula

A failure is a metastable event that is still unresolved after the available resolution time $t_r$. Multiplying the event rate by the survival probability gives the failure rate, and its inverse is

$$\text{MTBF} = \frac{e^{t_r/\tau}}{T_w\, f_c\, f_d}.$$

With $\tau = 10$ ps, $t_r = 400$ ps, $T_w = 20$ ps, $f_c = 1$ GHz and $f_d = 100$ MHz: the denominator is $2\times10^{-11} \cdot 10^9 \cdot 10^8 = 2\times10^{6}$ per second, the numerator $e^{40} \approx 2.35\times10^{17}$, so MTBF $\approx 1.18\times10^{11}$ s, about 3700 years. Such numbers span dozens of orders of magnitude, so they are best compared as $\log_{10}$ of the MTBF in seconds: here $\log_{10} \text{MTBF} = 40/\ln 10 - \log_{10}(2\times10^{6}) \approx 17.37 - 6.30 = 11.07$.

## Resolution time in a two-flop synchronizer

In the standard two-flip-flop synchronizer the first flip-flop samples the asynchronous signal and may go metastable; the second samples the first one clock period later. The first flop's output must be resolved by the time the second flop needs it, so

$$t_r = T_c - t_{cq} - t_{su}$$

(any logic between the two flops would subtract further, which is why none is allowed). With $T_c = 1000$ ps, $t_{cq} = 50$ ps and $t_{su} = 30$ ps, $t_r = 920$ ps. Because $t_r$ sits in an exponent, the clock period matters enormously: halving the period in this example cuts $t_r$ from 920 ps to 420 ps and divides the MTBF by $e^{50}$ when $\tau = 10$ ps.

## Adding stages to meet a target

Each additional flip-flop in the chain adds about one more $T_c - t_{cq} - t_{su}$ of resolution time, so an $N$-flop synchronizer has $t_r \approx (N - 1)(T_c - t_{cq} - t_{su})$ and its $\log_{10}$ MTBF grows linearly in $N$. Designers pick the smallest $N \ge 2$ that meets the target, typically stated per synchronizer from a system budget: a chip with $10^{3}$ synchronizers and a 1000-year system target needs each one to beat $10^{6}$ years.

The price of extra stages is latency, one cycle each, which is why fast clock domains with weak (slow) flip-flops sometimes need three stages. Clock-domain-crossing tools check structure (that a synchronizer exists, with no logic inside it) rather than computing MTBF; the MTBF analysis is done once per synchronizer cell and clock combination.

## Going deeper

Ginosar, "Metastability and Synchronizers: A Tutorial" (IEEE Design and Test of Computers, 2011); Dally and Poulton, "Digital Systems Engineering" (1998), chapter 10.
