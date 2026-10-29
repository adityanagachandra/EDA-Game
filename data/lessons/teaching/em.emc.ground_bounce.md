# Ground bounce and simultaneous switching noise

When many outputs switch together, their combined current change flows through the package's ground and power inductance and shifts the on-die reference: ground bounce, or simultaneous switching noise. It corrupts quiet outputs and input thresholds and limits how many I/Os can switch at once. Decoupling capacitance supplies charge locally so the package does not have to. Power-integrity tools estimate both effects.

## L di/dt noise

Current returning through a ground pin and bond wire with inductance $L$ develops a voltage $V = L\,di/dt$ across it, so the die's ground moves relative to the board's ground while the current changes. Nanohenries times milliamperes per nanosecond give millivolts.

Example: a driver whose ground current rises by 20 mA in 0.25 ns through 2 nH produces $2 \cdot 20/0.25 = 160$ mV of bounce. That shift appears on every signal referenced to the die's ground, including quiet outputs held low, which glitch by the same amount, and input receivers, whose thresholds move. Faster edges make it worse: halving the edge time doubles $di/dt$ and the bounce, which is why output drivers include slew-rate control.

## Simultaneous switching through shared pins

With $N$ drivers switching together, each changing current by $\Delta I$ in $\Delta t$, through $P$ ground pins in parallel each of inductance $L$, the effective inductance is about $L/P$ (ignoring mutual inductance between pins) and
$$V \approx \frac{L}{P}\cdot\frac{N\,\Delta I}{\Delta t}.$$

Example: 32 drivers, each changing by 10 mA in 0.5 ns, sharing 4 ground pins of 2 nH: $L/P = 0.5$ nH and total $di/dt = 640$ mA/ns, so $V = 320$ mV. Mutual inductance between adjacent pins carrying current in the same direction raises the effective inductance above $L/P$, so this estimate is optimistic. I/O planning tools limit the number of simultaneously switching outputs per ground pin, interleave power and ground pins, and use this kind of estimate to place them.

## Decoupling charge

On-die decoupling capacitance supplies current for the first few nanoseconds of a load step, before the package and board can respond. To keep the droop below $\Delta V$ while delivering an extra current $I$ for a time $\Delta t$, the capacitor must hold the charge $Q = I\Delta t$ with $\Delta V = Q/C$, so
$$C \ge \frac{I\,\Delta t}{\Delta V}.$$
Milliamperes times nanoseconds over millivolts give nanofarads.

Example: an extra 200 mA for 2 ns with at most 50 mV droop needs $200 \cdot 2/50 = 8$ nF. On-die decap density is limited, so this calculation drives how much die area is filled with decap cells. It is a first-order estimate: the decap's own series resistance and the inductance between decap and load reduce its effectiveness at the highest frequencies, and the package resonance with the decap creates the anti-resonance peak that power-integrity analysis must also check.

## Going deeper

Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems", chapter 1; Johnson and Graham, "High-Speed Digital Design: A Handbook of Black Magic", chapter 7; Bakoglu, "Circuits, Interconnections, and Packaging for VLSI", chapter 7.
