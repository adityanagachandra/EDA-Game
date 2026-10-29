# Real-number modeling

Real-number modeling represents analog signals as real values passed between processes in an event-driven simulator, without an analog solver. A filter, an ADC or a PLL becomes a few lines of SystemVerilog or Verilog-AMS that update a real variable at events. Models run at digital-simulation speed, which makes them the standard way to verify the digital control around analog blocks. Their accuracy rests on how continuous dynamics are discretized.

## Discretizing a continuous block

An RC low-pass filter obeys $\tau\dot{y} = x - y$. A real-number model updates $y$ only at events, for example every $T$ nanoseconds. If $x$ is held constant over the step, the ODE solves exactly: $y_{k+1} = y_k + \alpha(x_k - y_k)$ with $\alpha = 1 - e^{-T/\tau}$. For $\tau = 50$ ns and $T = 5$ ns, $\alpha = 1 - e^{-0.1} = 0.0952$. This exact discretization is stable for any step and matches the continuous response at the sample points, so the model's error comes only from the assumption that the input is constant between events. Higher-order blocks discretize the same way, through the matrix exponential of their state-space form, or by bilinear transformation of their transfer function.

## Stepping a model through time

Exact steps compose: after $k$ steps of a constant input $x$ from initial value $y_0$, $y_k = x + (y_0 - x)e^{-kT/\tau}$. With $\tau = 50$ ns, a step from 0 to 1 V reaches $1 - e^{-1} = 0.632$ V after 10 steps of 5 ns, exactly the continuous value at 50 ns. Event-driven models can also update only when the input changes and compute the output analytically at any later query time, which avoids periodic events altogether. Either way, the model's time resolution is chosen by the modeller, not by an error-controlled solver, so the modeller must pick it relative to the fastest time constant and the signals the digital logic samples.

## Discretization error and validation

A naive model uses forward Euler, $\alpha = T/\tau$, which is accurate only for $T \ll \tau$ and becomes unstable for $T > 2\tau$. With $T/\tau = 0.25$, Euler's per-step factor 0.25 exceeds the exact 0.221, so after 10 steps of a 0.8 V step input Euler gives $0.8(1 - 0.75^{10}) = 0.755$ V against the exact 0.734 V, an error of 21 mV. Errors like this can change when a comparator in the model trips and therefore which branch the digital logic takes. Real-number models are therefore validated against transistor-level or Verilog-A simulation on the same stimulus, with tolerances on the metrics the digital side depends on. Real-number models also cannot represent loading between blocks unless the model includes it explicitly, since there is no shared nodal solve.

## Going deeper

Kundert and Zinke, "The Designer's Guide to Verilog-AMS" (2004); IEEE 1800 SystemVerilog standard, sections on real types and user-defined nettypes.
