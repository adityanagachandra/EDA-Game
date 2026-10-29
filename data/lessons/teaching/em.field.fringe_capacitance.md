# Fringe capacitance of on-chip wires

On-chip wires are about as tall as they are wide and sit close to the layers below, so a large part of their capacitance comes from field lines leaving their sidewalls and top: fringe capacitance. Parasitic extractors model it with fitted formulas or pattern tables built from field-solver runs. Knowing how large it is, and how it scales, explains why area-only estimates fail in every modern process.

## Area capacitance alone

The simplest model treats a wire of width $W$ at height $H$ above a ground plane as a parallel plate, giving a capacitance per unit length $C/\ell = \varepsilon W/H$. On chip it is convenient to use $\varepsilon_0 = 8.854$ aF/um, so the answer comes out in attofarads per micrometre.

Example: $W = 0.1$ um and $H = 0.2$ um in a dielectric with $\varepsilon_r = 2.7$ give $8.854 \cdot 2.7 \cdot 0.5 = 11.95$ aF/um. A real wire of that shape has well over twice this, as the next section shows. The area term is accurate only for wide wires, $W \gg H$, such as power straps and wide buses, where the field between the bottom face and the ground plane dominates. For minimum-width signal wires it is the smaller part.

## An empirical fringe model

Fitted closed forms add the fringe. A widely used one, from Sakurai and Tamaru, models a single wire of thickness $T$ over a ground plane as
$$\frac{C}{\varepsilon\ell} = 1.15\frac{W}{H} + 2.80\left(\frac{T}{H}\right)^{0.222},$$
fitted to numerical field solutions over roughly $0.3 \le W/H, T/H \le 30$. The first term is a slightly enlarged area term; the second captures sidewall and top fringing and depends only weakly on $T/H$.

Example: $W = 0.1$, $T = 0.2$, $H = 0.2$ um, $\varepsilon_r = 2.7$: $W/H = 0.5$, $T/H = 1$, so $C/(\varepsilon\ell) = 0.575 + 2.80 = 3.375$, and $C/\ell = 23.906 \cdot 3.375 = 80.7$ aF/um. That is almost seven times the area-only value. Outside its fitted range the formula is an extrapolation, so a tool should flag it there.

## The fringe share as wires narrow

The fringe share is the fraction of the total that the parallel-plate term $\varepsilon W/H$ misses. Because the fringe term barely depends on $W$, it stays roughly constant while the area term shrinks with width, so the share grows as wires narrow.

Example with the fit, $T/H = 1$: at $W/H = 10$ the total is $11.5 + 2.8 = 14.3$ (in units of $\varepsilon$) and the area term is 10, a fringe share of 30%. At $W/H = 1$ the total is 3.95 and the share 75%. At $W/H = 0.5$ it is 85%. This is why wire capacitance per unit length has fallen far more slowly than wire width over process generations, and why it saturates near a few tenths of a femtofarad per micrometre. Neighbouring wires steal some of this fringe field as coupling capacitance, which a single-wire formula cannot show.

## Going deeper

Sakurai and Tamaru, "Simple formulas for two- and three-dimensional capacitances", IEEE Transactions on Electron Devices, 1983; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), chapter 4.
