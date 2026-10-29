# Silicon wafers and substrates

Every chip starts as a slice of a single silicon crystal. Wafer diameter, orientation and substrate type are fixed before any design exists, yet they decide how many dies a fab can make per wafer, how much area is lost at the edge, and which transistor structures are available. For a tool designer, the wafer is the domain over which yield and cost are computed, and its geometry is the first input to every dies-per-wafer and cost model.

## From crystal to wafer

Electronic-grade silicon is melted and pulled as a single crystal ingot from a seed by the Czochralski process. The ingot is ground to a precise diameter, sliced, lapped, etched and polished to a mirror finish. Logic wafers today are 300 mm in diameter and about 775 um thick; 200 mm and smaller wafers remain common for analog, power and older nodes. The surface is cut along a crystal plane, usually (100) for CMOS, because orientation sets oxide interface quality and carrier mobility. Wafer area is $\pi r^2$: a 300 mm wafer has radius 15 cm and area $\pi \times 15^2 \approx 706.9$ square cm, so a lot of 25 wafers carries about 17,671 square cm of silicon.

## Edge exclusion and usable area

The outer few millimetres of a wafer are unreliable: films are thinner, the edge is bevelled, and handling and clamping damage it. Fabs define an edge exclusion inside which no product die may sit. The usable disc has radius $r - e$, so the fraction lost is $1 - ((r - e)/r)^2$. For a 300 mm wafer with a 3 mm exclusion, $1 - (147/150)^2 = 3.96\%$ of the area is gone before a single die is placed. Because the ring has fixed width, it costs a small wafer proportionally more: the same 3 mm takes $1 - (97/100)^2 = 5.91\%$ of a 200 mm wafer.

## Substrate choices and wafer size

Beyond bulk silicon, a fab can start from an epitaxial wafer, a lightly doped crystalline layer grown on a heavily doped substrate to suppress latch-up, or a silicon-on-insulator (SOI) wafer, a thin silicon film on a buried oxide. Fully depleted SOI gives good electrostatics and back-gate biasing without fins, which changes the device models and some layout rules in the PDK. Moving to a larger wafer raises usable area by the square of the radius ratio: going from 200 mm to 300 mm with a 3 mm exclusion multiplies usable area by $(147/97)^2 = 2.297$, while processing cost per wafer grows less, which is why each transition cut cost per die.

## Going deeper

Plummer, Deal and Griffin, "Silicon VLSI Technology", chapter 3 on crystal growth and wafer preparation.
