# Junction voltage limiting

An exponential device turns a modest Newton overshoot into a catastrophe: evaluating a diode at 5 V gives a current of order $e^{200}$. SPICE prevents this with junction limiting, a rule that replaces a large forward step by a logarithmically compressed one before the device is evaluated. The limiter, called pnjlim in Berkeley SPICE, is a few lines of code that make most diode and bipolar circuits converge.

## The critical voltage

Limiting is applied only above a critical voltage where the exponential becomes steep:
$$v_{crit} = V_T \ln\!\left(\frac{V_T}{\sqrt{2}\, I_S}\right),$$
the point of minimum radius of curvature of the diode's current-voltage curve. Below it the curve is gentle and plain Newton is safe. For $I_S = 10$ fA and $V_T = 25$ mV, $V_T/(\sqrt{2} I_S) = 1.77 \times 10^{12}$, its log is 28.2, and $v_{crit} = 0.705$ V. Larger saturation currents lower $v_{crit}$: a 100 fA diode has $v_{crit} = 0.648$ V. The value is computed once per device from its model parameters.

## Limiting a forward step

If the proposed voltage exceeds $v_{crit}$ and differs from the previous junction voltage by more than $2V_T$, and the previous voltage is positive, SPICE replaces it by
$$v = v_{old} + V_T \ln\!\left(1 + \frac{v_{new} - v_{old}}{V_T}\right).$$
This matches the current the linearized model predicted: the step in voltage is chosen so that the exponential grows by about the factor the tangent line promised. With $v_{old} = 0.6$ V and $v_{new} = 3$ V, the argument is $1 + 2.4/0.025 = 97$, so $v = 0.6 + 0.025 \times 4.57 = 0.714$ V. A 2.4 V jump becomes a 114 mV step, and Newton proceeds from a sane point.

## When limiting does nothing

The limiter leaves the proposal alone in two cases: when it is at or below $v_{crit}$, where the exponential is mild, and when the step is at most $2V_T$, where the companion model is already accurate. A third case applies when the previous voltage was zero or reverse-biased: then the logarithmic step from $v_{old}$ is not meaningful, and SPICE uses $v = V_T \ln(v_{new}/V_T)$ instead. For $v_{old} = -0.5$ V and $v_{new} = 2$ V, that gives $0.025 \ln 80 = 0.110$ V. After limiting, the device is evaluated at the limited voltage, and the iteration is flagged as not converged, so Newton cannot stop on a limited step.

## Going deeper

Nagel, "SPICE2" (Berkeley ERL Memo M520, 1975), the discussion of junction voltage limiting; the pnjlim routine in the ngspice source.
