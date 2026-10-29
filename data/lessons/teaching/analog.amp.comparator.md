# Comparators: regeneration and offset

A comparator decides which of two voltages is larger and produces a logic level. It is the decision element of every ADC, and its speed, offset and metastability set converter performance. Modern comparators are clocked latches that use positive feedback; their regeneration time constant and input-referred offset are the two quantities a converter designer and a sizing tool must compute.

## Regeneration by positive feedback

A latch is two cross-coupled inverters reset to a balanced state and then released with a small imbalance $V_0$. Near balance each side is a transconductor driving a capacitance, and the cross-coupling gives $C\,dv/dt = g_m v$, so the imbalance grows exponentially: $v(t) = V_0 e^{t/\tau}$ with $\tau = C/g_m$. The time to reach a logic level $V_{out}$ is $t = \tau \ln(V_{out}/V_0)$. With $C = 20$ fF and $g_m = 400\ \mu$S, $\tau = 50$ ps (femtofarads over microsiemens is nanoseconds); a 1 mV imbalance reaches 0.8 V in $50 \ln 800 = 334$ ps. Each factor of $e$ in input costs only one $\tau$, which is why latches resolve tiny inputs quickly.

## Metastability

If the decision must be ready by time $T$, the latch fails when the initial imbalance is below $V_{min} = V_{out} e^{-T/\tau}$. With $\tau = 20$ ps, $T = 200$ ps and $V_{out} = 1$ V, $V_{min} = e^{-10} = 45\ \mu$V. Inputs closer to the threshold than this leave the output undecided, a metastable state. For a random input distributed over a range much wider than $V_{min}$, the probability of metastability scales as $e^{-T/\tau}$, so doubling the regeneration time squares down the error rate. In ADCs a metastable comparator can produce large code errors (sparkle codes in flash converters, wrong MSBs in SARs), so designers budget enough time constants per decision and sometimes add metastability detection.

## Offset and the preamplifier

Mismatch between the latch devices gives an input-referred offset, often tens of millivolts for small devices. A preamplifier of gain $A_p$ in front divides the latch's offset by $A_p$ when referred to the input, and contributes its own offset. Independent offsets add in variance: $\sigma_{in} = \sqrt{\sigma_{pre}^2 + (\sigma_{latch}/A_p)^2}$. With a 3 mV preamp offset, a 25 mV latch offset and a gain of 10, $\sigma_{in} = \sqrt{9 + 6.25} = 3.9$ mV. The preamp also isolates the input from latch kickback. The cost is power and delay, so converter designers trade preamp gain against offset calibration, which measures and cancels offset digitally instead of with area or power.

## Going deeper

Razavi, "The StrongARM latch" (IEEE Solid-State Circuits Magazine, 2015); Johns and Martin, "Analog Integrated Circuit Design", chapter on comparators.
