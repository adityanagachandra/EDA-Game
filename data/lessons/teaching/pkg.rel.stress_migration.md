# Stress migration and voiding

Stress migration, also called stress-induced voiding, moves metal atoms under a gradient of mechanical stress rather than electron current. Copper lines and vias embedded in dielectric are left in tension after processing, because copper shrinks more than its surroundings on cooling, and vacancies diffuse toward stress concentrations, often under vias, until a void opens. It is a failure mechanism with no current at all, and it has a peculiar temperature dependence.

## Mechanism and driving force

After copper is deposited and annealed, it cools inside a dielectric and barrier that expand much less. The copper ends in hydrostatic tension, and the tension is not uniform: it is highest where geometry concentrates it, such as under a via landing on a wide line. Atoms diffuse from regions of high tensile stress, which is equivalent to vacancies diffusing toward them, and vacancies collect into a void at the weakest interface. Wide lines provide a large reservoir of vacancies, so a via on a wide plate is more vulnerable than a via on a narrow line, the opposite of the intuition from electromigration.

## The temperature peak

Stress migration is slowest at both temperature extremes. Near the stress-free temperature, which is near the process anneal, the driving stress vanishes. At low temperature the stress is large but diffusion is frozen. The product of driving force, rising as temperature falls below the stress-free point, and diffusivity, falling as temperature falls, peaks in between, typically in a band around 150 to 250 °C for copper interconnect. Qualification therefore bakes parts at several temperatures to find the peak instead of assuming that hotter is always worse. This is a key difference from electromigration and dielectric breakdown, which accelerate monotonically with temperature.

## Design and verification response

Designers fight stress voiding with redundant vias, which give a second path when one opens, and with layout rules that limit the width of lines landing on single vias or add slots to wide plates. Verification tools check these rules: a via-doubling check flags single vias on wide lines, and some flows compute a vulnerability score from the plate area feeding each via. In packages and 3D stacks the stress field also includes TSV and bump stress, so chip-package interaction analysis maps package stress into the back-end stack. Thermal cycling then combines with the residual tension, which is why stress migration is often qualified alongside thermal-cycle tests.

## Going deeper

Ogawa, McPherson, Rosal and others, "Stress-induced voiding under vias connected to wide Cu metal leads", IEEE IRPS (2002); Tu, "Electronic Thin-Film Reliability" (2010).
