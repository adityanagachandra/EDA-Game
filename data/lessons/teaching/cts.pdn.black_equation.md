# Black's equation and EM lifetime

Electromigration limits come from a lifetime model. Black's equation relates the median time to failure of a wire to its current density and temperature. Foundries fit its parameters to stress tests, then set current-density limits so that wires last the product's life at the signoff temperature. Knowing the equation lets a tool designer convert an ECO's change in current or temperature into a change in lifetime, and adjust limits for hot spots.

## The form of the model

Black's equation is
$$\mathrm{MTTF} = A\,J^{-n}\,e^{E_a/kT},$$
with $J$ the current density, $T$ the absolute temperature, $E_a$ an activation energy (often quoted around 0.7 to 1 eV for copper interconnect, depending on the failure mechanism), $k = 8.617 \times 10^{-5}$ eV/K Boltzmann's constant, $n$ a current exponent (commonly 1 to 2) and $A$ a constant that depends on geometry and process. At fixed temperature, raising $J$ by a factor of 1.8 with $n = 1$ shortens lifetime to $1/1.8 = 0.556$ of its value; with $n = 2$ it falls to $1/3.24 = 0.309$.

## Temperature dependence

At fixed current density, the lifetime ratio between temperatures $T_1$ and $T_2$ (in kelvin) is $\exp\big(\tfrac{E_a}{k}(\tfrac{1}{T_2} - \tfrac{1}{T_1})\big)$. Raising a strap from 90 to 95 degrees C, that is from 363.15 to 368.15 K, with $E_a = 0.85$ eV gives $\exp(-0.369) = 0.691$: five degrees cost about 30% of the lifetime. This sensitivity is why EM analysis uses the local temperature, including Joule self-heating of the wire, and why hot spots near high-power blocks need tighter limits.

## Derating the limit

To keep the same lifetime at a higher temperature, current density must fall. Setting the two lifetimes equal gives $J_2/J_1 = \exp\big(\tfrac{E_a}{nk}(\tfrac{1}{T_2} - \tfrac{1}{T_1})\big)$. With $n = 2$, $E_a = 1$ eV and a rise from 85 to 90 degrees C, $J_2/J_1 = 0.800$: the limit drops to 80% of its original value. Foundry EM rules are tabulated this way per temperature. Black's equation describes a median; signoff limits target a low failure fraction over the product lifetime, so they include statistical margin beyond the median model.

## Going deeper

Black, "Electromigration: a brief survey and some recent results", IEEE Transactions on Electron Devices (1969); Lienig and Thiele, "Fundamentals of Electromigration-Aware Integrated Circuit Design" (2018).
