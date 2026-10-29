# Timestep control

Transient analysis spends most of its effort deciding how long the next step should be. Too long and the local truncation error exceeds tolerance; too short and the run wastes factorizations. Because a method of order $k$ has local error proportional to $h^{k+1}$, an estimate at one step size predicts the error at another, which gives a simple and powerful rule for the next step.

## The next step for trapezoidal

The trapezoidal rule is order 2, so its local error scales as $h^3$. If the last step $h$ gave error estimate $E$ and the tolerance is $\varepsilon$, the step that would just meet tolerance satisfies $E (h_{new}/h)^3 = \varepsilon$:
$$h_{new} = h \left(\frac{\varepsilon}{E}\right)^{1/3}.$$
For $h = 0.5$ ns, $E = 0.2$ mV and $\varepsilon = 1$ mV, $h_{new} = 0.5 \times 5^{1/3} = 0.855$ ns. The cube root makes control gentle: an error five times below tolerance only lets the step grow by 71%. Simulators also cap growth per step and the maximum step, so a quiet waveform does not let the step jump past an approaching edge.

## The next step for backward Euler

Backward Euler is order 1, so its error scales as $h^2$ and the rule uses a square root: $h_{new} = h(\varepsilon/E)^{1/2}$. For $h = 1$ ns, $E = 0.32$ mV and $\varepsilon = 2$ mV, the ratio is 6.25 and $h_{new} = 2.5$ ns. Lower-order methods react more strongly to the same error ratio. The general rule is $h_{new} = h(\varepsilon/E)^{1/(k+1)}$ for order $k$, which variable-order BDF codes use to compare orders too: they estimate the error at several orders and pick the order that allows the largest next step.

## Rejecting a step

If the error estimate exceeds the tolerance, the step is rejected: its solution is discarded and the simulator retries from $t_n$ with a smaller step from the same rule, usually with a safety factor below 1 so the retry is likely to pass. For a trapezoidal step of 0.5 ns with $E = 4$ mV against $\varepsilon = 1$ mV and safety factor 0.9, $h_{retry} = 0.9 \times 0.5 \times 4^{-1/3} = 0.284$ ns. A rejection costs a full Newton solve with nothing to show for it, so controllers aim slightly below tolerance. A Newton failure also triggers rejection, typically with a fixed cut such as a factor of 8, and repeated failures end in timestep-too-small.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Hairer, Norsett and Wanner, "Solving Ordinary Differential Equations I" (Springer, 1993), on step-size control.
