# Package-board co-design

Package-board co-design treats die, package and board as one system. The ball map, the power delivery path and the high-speed channels cross all three, and a choice made in one level, such as a ball assignment, fixes constraints for the others. Three computations anchor the practice: how many balls power needs, what impedance the power network must meet, and how many parallel paths achieve a target inductance.

## Power ball count

Each ball carries a maximum current, set by electromigration and joule heating at the operating temperature, often a few hundred milliamps to an ampere depending on ball size. A rail drawing $I$ amps needs at least $\lceil I/I_{ball} \rceil$ supply balls and as many return balls. A 120 A core rail at 0.6 A per ball needs 200 power and 200 ground balls. On a package of a few thousand balls, power can take most of the map, which is why high-power packages assign the central region under the die to power and ground, keep signals on the outer rings for escape, and interleave power and ground to reduce loop inductance.

## Target impedance

A power delivery network keeps the supply within a ripple tolerance despite current transients. The target impedance is $Z_t = V \times r / \Delta I$, where $r$ is the allowed ripple fraction and $\Delta I$ the transient current step. For a 0.75 V rail with 3% ripple and a 50 A step, $Z_t = 0.75 \times 0.03/50 = 0.45$ milliohm. The network must stay below this impedance from DC to the frequencies where on-die capacitance takes over. The board and its voltage regulator cover low frequencies, package capacitors and planes cover the middle, and on-die decoupling covers the highest. Each level is designed against the same target, so the package designer must know the board's capability and the die's decap.

## Parallel paths and inductance

Above a few megahertz the network is limited by inductance, of the balls, vias and planes. If one ball-and-via path has inductance $L_1$, $N$ paths in parallel give roughly $L_1/N$, ignoring mutual inductance between them. Reaching 5 pH with 300 pH per path needs $\lceil 300/5 \rceil = 60$ paths. Mutual coupling between same-direction currents raises the effective value, while interleaving power and ground lowers it because opposing currents cancel each other's field, so real designs extract the network rather than divide. This is a co-design result: the package fixes the ball pattern, the board fixes the via pattern underneath, and only together do they meet the target.

## Going deeper

Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems" (2007); Bogatin, "Signal and Power Integrity, Simplified" (3rd ed., 2018).
