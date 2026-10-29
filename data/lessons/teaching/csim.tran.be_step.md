# One backward-Euler transient step

One transient step is one companion-model solve. For a linear circuit there is no Newton iteration: stamp the companions, solve once, advance. Taking backward-Euler steps by hand on an RC circuit shows the whole machinery at work, and comparing the result with the exact exponential shows what first-order accuracy means in volts.

## Assembling the step equation

A source steps to $V_s$ at $t = 0$ and drives resistor $R$ into a node with capacitor $C$ to ground. With the backward-Euler companion, KCL at the node is
$$\left(\frac{1}{R} + \frac{C}{h}\right) v_{n+1} = \frac{V_s}{R} + \frac{C}{h} v_n.$$
With $R$ in kilohms, $C$ in picofarads and $h$ in nanoseconds, every term is in millisiemens or milliamperes. For $V_s = 2.5$ V, $R = 1$ k$\Omega$, $C = 5$ pF, $h = 5$ ns and $v_n = 0.9$ V: $(1 + 1) v_{n+1} = 2.5 + 0.9$, so $v_{n+1} = 1.7$ V. The new voltage is a weighted average of the source and the old voltage, with weights $1/R$ and $C/h$.

## Marching two steps

Each accepted point becomes the history for the next step. Starting from 0 V with $h/\tau = 1$, where $\tau = RC$, the weights are equal, so $v_1 = V_s/2$ and $v_2 = (V_s + v_1)/2 = 0.75 V_s$. The exact values are $V_s(1 - e^{-1}) = 0.632 V_s$ and $V_s(1 - e^{-2}) = 0.865 V_s$. In general backward Euler multiplies the distance to the final value by $1/(1 + h/\tau)$ per step, against the exact $e^{-h/\tau}$. The numerical solution always lags the true one for a charging RC, but it can never overshoot or oscillate, whatever the step: backward Euler is L-stable.

## Error against the exact exponential

The local error of one step measures accuracy. From 0 V, one step gives $V_s \tfrac{h/\tau}{1 + h/\tau}$, while the exact value is $V_s(1 - e^{-h/\tau})$. For $V_s = 3.5$ V, $\tau = 5$ ns and $h = 1.25$ ns, that is $3.5 \times 0.2 = 0.700$ V against $3.5 \times 0.2212 = 0.774$ V, an undershoot of 74 mV. Expanding both in $h$ shows the difference is about $\tfrac{h^2}{2} |v''|$: local error of order $h^2$, global error of order $h$. Halving the step halves the accumulated error, which is why accuracy-hungry simulations prefer second-order methods.

## Going deeper

Najm, "Circuit Simulation" (Wiley, 2010), chapter 5; Hairer and Wanner, "Solving Ordinary Differential Equations II: Stiff and Differential-Algebraic Problems" (Springer, 1996).
