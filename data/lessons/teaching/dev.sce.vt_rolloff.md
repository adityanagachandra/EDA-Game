# Threshold roll-off and the natural length

Threshold roll-off is the drop in threshold as channel length shrinks. Its size is governed by an electrostatic natural length set by the device's geometry, and the threshold reduction decays roughly exponentially in $L/\lambda$. This scale length is the cleanest way to compare planar, SOI, double-gate and gate-all-around devices, and to see why device architecture, not only lithography, decides how short a gate can be.

## The exponential roll-off model

Solving Poisson's equation in the channel region shows that the potential disturbance from the source and drain decays into the channel with a characteristic length $\lambda$. The threshold reduction then behaves approximately as
$$\Delta V_t(L) \approx A\, e^{-L/(2\lambda)},$$
where $A$ depends on the junction potentials and bias, and is of order the built-in potential plus a share of $V_{DS}$.

Example: $A = 1.0$ V and $\lambda = 8$ nm. At $L = 40$ nm, $\Delta V_t = e^{-2.5} = 82$ mV. At $L = 30$ nm, $e^{-1.875} = 153$ mV.

The steep exponential is the key message. Each nanometre shaved off near the roll-off knee costs more threshold than the last. That is why gate-length variation, around a nanometre, becomes significant threshold variation in short devices, and why processes choose gate lengths several natural lengths long.

## The natural length of a device

The natural length depends on how many gates surround a thin body and how thick the body and dielectric are. For a single-gate fully depleted SOI device, a standard estimate is
$$\lambda = \sqrt{\frac{\varepsilon_{Si}}{\varepsilon_{ox}}t_{Si}t_{ox}}.$$
For a symmetric double gate it becomes $\sqrt{\varepsilon_{Si}t_{Si}t_{ox}/(2\varepsilon_{ox})}$, smaller by $\sqrt 2$, and gate-all-around geometries shrink it further.

Example: $\varepsilon_{Si}/\varepsilon_{ox} = 3$, $t_{Si} = 10$ nm and $t_{ox} = 1$ nm. Single gate: $\lambda = \sqrt{30} = 5.48$ nm. Double gate: $\lambda = \sqrt{15} = 3.87$ nm.

The lesson for device scaling is that a thinner body and more gates shrink $\lambda$ and so allow a shorter gate at the same roll-off. That is the electrostatic argument for FinFETs (thin fin, gate on three sides) and nanosheets (gate on all sides). Planar bulk devices have an effective body depth set by depletion, which doping can only reduce so far.

## Minimum gate length for a roll-off budget

Inverting the model gives the shortest gate that keeps roll-off within a budget $\Delta V_{max}$:
$$L_{min} = 2\lambda \ln\frac{A}{\Delta V_{max}}.$$

Example: $A = 1.0$ V, $\lambda = 5.5$ nm and a budget of 20 mV. $L_{min} = 11 \times \ln 50 = 43$ nm. Halving $\lambda$ halves $L_{min}$.

A common rule of thumb follows: keep $L$ at roughly five to ten natural lengths for acceptable short-channel behavior. Device architects use this budget to choose fin width, sheet thickness and EOT for a target gate pitch. For EDA, the consequence shows up in design rules and in the models: gate lengths in a library are fixed by the process, and length biasing (drawing slightly longer gates on non-critical cells) is used to trade a little speed for lower leakage, exploiting exactly the steep roll-off curve.

## Going deeper

Yan, Ourmazd and Lee, "Scaling the Si MOSFET: from bulk to SOI to bulk", IEEE Transactions on Electron Devices (1992); Colinge, "FinFETs and Other Multi-Gate Transistors" (2008), chapter 1.
