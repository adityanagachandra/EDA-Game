# MOSFET threshold voltage

The threshold voltage $V_t$ is the gate voltage at which a MOSFET's channel inverts and conduction turns on. It is the single most consequential device parameter for a digital designer: it sets drive current, leakage, delay and its variation. Every compact model has a threshold equation at its core, and every threshold-voltage flavor in a library is a different setting of the terms below.

## The inversion condition

Bias the gate of an NMOS on a p-type body. As $V_{GS}$ rises, holes are pushed away (depletion) and then electrons gather at the surface. The conventional strong-inversion condition is that the surface electron density equals the bulk hole density, which happens when the surface potential $\psi_s$ bends the bands by twice the Fermi potential:
$$\psi_s = 2\phi_F, \qquad \phi_F = \frac{kT}{q}\ln\frac{N_A}{n_i}.$$

Example: $N_A = 10^{17}$ cm$^{-3}$ at 300 K gives $\phi_F = 0.02585 \times \ln(10^7) = 0.4167$ V, so $2\phi_F = 0.833$ V.

The definition is a convention. Inversion is a smooth process, and conduction rises exponentially well before $2\phi_F$, which is why measured thresholds are defined by an extraction method (constant current, linear extrapolation) rather than by this condition directly.

## The depletion-charge term

At threshold the gate must support the depletion charge under the channel. The depletion depth at $\psi_s = 2\phi_F$ is $W_{dep} = \sqrt{2\varepsilon_{Si}(2\phi_F)/(qN_A)}$, and the charge per area is
$$Q_{dep} = qN_A W_{dep} = \sqrt{2\varepsilon_{Si}qN_A(2\phi_F)}.$$
The oxide drops a voltage $Q_{dep}/C_{ox}$ to hold this charge.

Example: $N_A = 10^{17}$ cm$^{-3}$, $2\phi_F = 0.833$ V. $Q_{dep} = \sqrt{2 \times 1.036 \times 10^{-12} \times 1.602 \times 10^{-19} \times 10^{17} \times 0.833} = 1.66 \times 10^{-7}$ C/cm$^2$. With $t_{ox} = 3$ nm ($C_{ox} = 1.15 \times 10^{-6}$ F/cm$^2$) the term is 0.144 V.

This term grows as $\sqrt{N_A}$ and shrinks with thinner oxide. It is why raising channel doping raises threshold, and why oxide scaling weakened doping's leverage on $V_t$.

## The full threshold expression

Adding the voltage needed to reach flat bands gives the long-channel threshold of a bulk NMOS:
$$V_t = V_{FB} + 2\phi_F + \frac{\sqrt{2\varepsilon_{Si}qN_A(2\phi_F)}}{C_{ox}}.$$
The flat-band voltage $V_{FB}$ is the work-function difference between gate and body, minus any fixed oxide charge over $C_{ox}$. An n+ polysilicon gate on a p-type body gives $V_{FB}$ near -0.9 to -1.0 V.

Example: $N_A = 10^{17}$ cm$^{-3}$, $t_{ox} = 3$ nm, $V_{FB} = -0.95$ V. $V_t = -0.95 + 0.833 + 0.144 = 0.027$ V, too low for logic, which is why real processes add threshold-adjust implants or choose gate metals with suitable work functions.

In high-k metal-gate processes, threshold flavors (ULVT, LVT, SVT, HVT) are largely set by work-function metal stacks rather than by channel doping. Compact models replace this textbook expression with fitted forms that add short-channel, narrow-width and body-bias terms on top.

## Going deeper

Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 3.1; Tsividis and McAndrew, "Operation and Modeling of the MOS Transistor" (3rd ed.), chapter 4.
