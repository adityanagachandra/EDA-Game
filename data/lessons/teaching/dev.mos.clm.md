# Channel-length modulation and output resistance

In saturation the drain current is not perfectly flat: it rises slowly with drain voltage because the pinch-off point moves toward the source. This channel-length modulation sets the output resistance $r_o$ of every transistor, which caps the gain of every amplifier stage and makes the drive current of digital gates depend on the output voltage. Compact models and current-source cell models both have to capture it.

## The lambda model

Above pinch-off, extra drain voltage widens the depleted region near the drain, shortening the effective channel by $\Delta L$. A shorter channel carries more current. The first-order model scales the saturation current linearly:
$$I_D = I_{D,sat}(1 + \lambda V_{DS}).$$
The parameter $\lambda$, in V$^{-1}$, grows as channel length shrinks, roughly as $1/L$.

Example: $I_{D,sat} = 200$ µA, $\lambda = 0.1$ V$^{-1}$, $V_{DS} = 0.9$ V. $I_D = 200 \times 1.09 = 218$ µA.

The model is crude but convenient. Real devices also show DIBL and, at high drain voltage, substrate-current effects, all of which raise output conductance. Compact models such as BSIM model these separately, with distinct parameters, because their bias and length dependences differ.

## Output resistance

The small-signal output resistance is the inverse slope of the output characteristic:
$$r_o = \left(\frac{\partial I_D}{\partial V_{DS}}\right)^{-1} \approx \frac{1}{\lambda I_D}.$$

Example: $\lambda = 0.2$ V$^{-1}$ and $I_D = 100$ µA give $r_o = 1/(0.2 \times 10^{-4}) = 50$ kilohms.

Output resistance falls as current rises, and falls as channel length shrinks. For digital timing it means the current a gate delivers depends on how far its output has swung, which the current-source models in modern timing libraries (CCS, ECSM) capture by tabulating current versus both input and output voltage. For analog it means a minimum-length device is a poor current source, so analog designers routinely draw longer channels.

## Intrinsic gain

The largest voltage gain a single transistor can provide, with an ideal current-source load, is its intrinsic gain:
$$A_i = g_m r_o = \frac{2I_D}{V_{ov}} \cdot \frac{1}{\lambda I_D} = \frac{2}{\lambda V_{ov}}.$$
The current cancels.

Example: $\lambda = 0.25$ V$^{-1}$ and $V_{ov} = 0.2$ V give $A_i = 2/0.05 = 40$.

Scaling has steadily reduced intrinsic gain, because short channels raise $\lambda$ and low supply voltages limit overdrive choices; at minimum length in late planar nodes it fell to the order of ten. FinFETs recovered some of it through better electrostatics, which lowers DIBL and output conductance. For a tool designer this explains why analog circuits are sensitive to how well a compact model fits output conductance, and why fitting only $I_D$ while ignoring $g_{ds}$ produces models that are useless for amplifiers.

## Going deeper

Gray, Hurst, Lewis and Meyer, "Analysis and Design of Analog Integrated Circuits" (5th ed.), chapter 1; Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed.), chapter 2.
