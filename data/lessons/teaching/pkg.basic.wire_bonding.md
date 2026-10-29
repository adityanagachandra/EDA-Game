# Wire bonding

Wire bonding connects die pads to package leads with thin gold, copper or aluminium wires, bonded by heat, force and ultrasonic energy. It is the oldest and still the cheapest chip-to-package interconnect, used on most low pin-count parts. Its limits explain the move to flip chip: pads must sit on the die perimeter, so I/O count grows only linearly with die size, and each wire adds roughly a nanohenry of inductance per millimetre to every power and signal path.

## Perimeter-limited pad count

Wire-bond pads sit along the die edge, so the pad count is set by perimeter and pad pitch, not by area. A square die of side $S$ with one row of pads at pitch $p$ holds about $\lfloor S/p \rfloor$ pads per side and $4\lfloor S/p \rfloor$ in total. A 5 mm die at 60 um pitch holds $\lfloor 5000/60 \rfloor = 83$ per side, 332 in all. Two staggered rows can roughly double this, but the wires then cross at two loop heights and the bonder needs more clearance. Doubling the die side only doubles the pads while the logic grows fourfold, so large dies become pad-limited: the die must grow to fit its I/O, and silicon area is wasted. That scaling argument is the main reason high pin-count parts use area-array flip chip.

## Bond wire resistance

A bond wire is a cylinder of length $L$ and diameter $d$, so its resistance is $R = \rho L / (\pi d^2/4)$. Gold has $\rho \approx 2.44$ micro-ohm cm and copper about 1.7. A 2 mm gold wire of 25 um diameter has area $\pi (25 \times 10^{-6})^2/4 = 4.91 \times 10^{-10}$ m^2, so $R = 2.44 \times 10^{-8} \times 2 \times 10^{-3} / 4.91 \times 10^{-10} = 0.099$ ohm, about 99 milliohm. That matters for power: a 1 A supply through one such wire drops about 0.1 V. Packages therefore use many parallel power and ground wires, and copper wire is popular partly for its lower resistivity, at the cost of a harder bonding process.

## Inductance and ground bounce

The bigger electrical problem is inductance. A bond wire has a self-inductance of roughly 1 nH per millimetre, and a current change $\Delta I$ over time $t_r$ through inductance $L$ induces $V = L\,\Delta I/t_r$. If 16 outputs switch together and draw 0.4 A of extra ground current in 1 ns through a 3 nH ground wire, the ground bounces by $3 \times 0.4 / 1 = 1.2$ V, enough to corrupt logic levels. Putting $N$ ground wires in parallel divides the inductance by about $N$ when mutual coupling is ignored, so eight wires cut the bounce to 0.15 V. Real wires couple, so the improvement is smaller, which is why simultaneous switching noise is checked with extracted package models rather than with this estimate alone.

## Going deeper

Harman, "Wire Bonding in Microelectronics" (3rd ed., 2010); Bogatin, "Signal and Power Integrity, Simplified" (3rd ed., 2018), chapters on inductance.
