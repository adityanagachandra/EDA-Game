# Linear time-invariant systems and convolution

Linear time-invariant (LTI) systems are the reason circuit and interconnect analysis is tractable. Linearity lets responses superpose; time invariance lets one response be shifted to any time. Together they mean a system is fully described by its impulse response, and any output is a convolution. Reduced-order interconnect models, crosstalk superposition and waveform-based delay calculation all rest on this.

## Linearity and time invariance

A system is linear if scaling and adding inputs scales and adds outputs: input $a x_1 + b x_2$ gives output $a y_1 + b y_2$. It is time invariant if delaying the input by $T$ delays the output by $T$ and changes nothing else. RLC networks with fixed element values are LTI; transistors are not linear, which is why gate delays are tabulated rather than superposed.

Example: an interconnect model gives 0.4 V at $t = 30$ ps for input $x_1$ and 0.25 V at the same time for $x_2$. For input $3x_1 - 2x_2$ the output at 30 ps is $3 \cdot 0.4 - 2 \cdot 0.25 = 0.7$ V. Crosstalk analysis uses exactly this: the victim's waveform is the sum of its own response and the coupled responses from each aggressor, each computed separately.

## Discrete convolution

Any discrete input is a sum of scaled, shifted unit impulses, $x[n] = \sum_k x[k]\delta[n-k]$. By linearity and time invariance, the output is the same sum of scaled, shifted impulse responses:
$$y[n] = \sum_k h[k]\,x[n-k].$$
For finite sequences, only terms with both indices in range contribute.

Example: $h = [2, 1, -1]$ and $x = [1, 3, 2]$, both starting at $n = 0$. Then $y[0] = 2 \cdot 1 = 2$, $y[1] = 2 \cdot 3 + 1 \cdot 1 = 7$, $y[2] = 2 \cdot 2 + 1 \cdot 3 + (-1) \cdot 1 = 6$. The output length is $3 + 3 - 1 = 5$. Waveform engines that combine a sampled input with a sampled channel response compute exactly this sum, often with an FFT when the sequences are long, since direct convolution costs the product of the two lengths.

## Pulse response from the step response

A rectangular pulse of height $V$ from 0 to $T$ is a step of $V$ at $t = 0$ minus a step of $V$ at $t = T$. By superposition and time invariance, its response is the step response minus a delayed copy. For an RC stage with $s(t) = V(1 - e^{-t/\tau})$, after the pulse ends:
$$y(t) = V\left(e^{-(t-T)/\tau} - e^{-t/\tau}\right), \quad t > T.$$

Example: $\tau = 50$ ps, a 1 V pulse lasting 100 ps, observed at 150 ps: $y = e^{-1} - e^{-3} = 0.368 - 0.050 = 0.318$ V. During the pulse the output only reached $1 - e^{-2} = 0.865$ V. This decomposition is how glitch and noise propagation is analyzed: a narrow pulse may never reach a receiver's threshold, which is the basis of noise filtering in static noise analysis.

## Going deeper

Oppenheim and Willsky, "Signals and Systems" (2nd ed.), chapters 1 and 2; Lathi, "Linear Systems and Signals", chapter 2.
