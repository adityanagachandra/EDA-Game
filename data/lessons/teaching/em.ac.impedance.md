# Impedance and admittance

Impedance generalizes resistance to sinusoidal steady state, and admittance is its reciprocal. With them, series and parallel rules, dividers, nodal analysis and Thevenin equivalents all carry over to AC with complex arithmetic. Frequency-domain simulators assemble a complex admittance matrix exactly the way DC analysis assembles a conductance matrix, so the same structural intuition applies.

## Series impedance magnitude

Impedances in series add. A series RLC has $Z = R + j(\omega L - 1/(\omega C))$, so
$$|Z| = \sqrt{R^2 + \left(\omega L - \frac{1}{\omega C}\right)^2}.$$
The inductive and capacitive reactances partly cancel, and at one frequency they cancel completely.

Example: $R = 20$ ohm, $L = 50$ nH and $C = 20$ pF at 200 MHz. $X_L = 2\pi \cdot 2 \times 10^8 \cdot 5 \times 10^{-8} = 62.83$ ohm and $X_C = 1/(2\pi \cdot 2 \times 10^8 \cdot 2 \times 10^{-11}) = 39.79$ ohm, so the net reactance is 23.04 ohm and $|Z| = \sqrt{400 + 530.9} = 30.51$ ohm. Below the cancellation frequency the circuit looks capacitive, above it inductive. A package lead plus decap is exactly this circuit, and its minimum impedance, equal to $R$, is what the decap buys.

## Admittance and parallel combinations

For parallel elements it is easier to add admittances, $Y = 1/Z$: a conductance $G = 1/R$, a capacitor $j\omega C$, an inductor $1/(j\omega L)$. Then $Z = 1/Y$. For a resistor in parallel with a capacitor, $Y = 1/R + j\omega C$ and
$$|Z| = \frac{R}{\sqrt{1 + (\omega RC)^2}}.$$

Example: 1 kohm in parallel with 10 pF at 100 MHz. $\omega RC = 2\pi \cdot 10^8 \cdot 10^3 \cdot 10^{-11} = 6.283$, so $|Z| = 1000/\sqrt{1 + 39.48} = 157.2$ ohm, close to the capacitor's own 159 ohm because it dominates. Frequency-domain MNA stamps $G + j\omega C$ into the same positions as the DC conductance, so an AC solve is a complex sparse solve per frequency point, with the same sparsity pattern at every frequency.

## Phase angle of an impedance

The angle of $Z = R + jX$ is $\theta = \arctan(X/R)$, positive for inductive and negative for capacitive impedance. It tells how far the current is shifted from the voltage, and it decides how much of the apparent power $|V||I|$ is real power, $|V||I|\cos\theta$.

Example: $R = 50$ ohm, $L = 10$ nH and $C = 100$ pF at 50 MHz. $X_L = 3.14$ ohm and $X_C = 31.83$ ohm, so $X = -28.69$ ohm and $\theta = \arctan(-28.69/50) = -29.8^\circ$: capacitive. The sign flips as frequency passes the point where the reactances cancel. A frequency sweep that plots $|Z|$ and $\theta$ together is how a power-integrity engineer reads a power-delivery network: magnitude says how much noise per ampere, phase says which element dominates there.

## Going deeper

Nilsson and Riedel, "Electric Circuits", chapters 9 and 10; Swaminathan and Engin, "Power Integrity Modeling and Design for Semiconductors and Systems", chapter 1.
