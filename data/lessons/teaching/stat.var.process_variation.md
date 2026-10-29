# Process variation

No two transistors come out of a fab identical. Gate length, oxide thickness, dopant count and line-edge shape all scatter, and so do the delays, leakages and margins built from them. Every variation-aware engine, from a statistical timer to a Monte Carlo SPICE run to a yield optimizer, starts from a model of that scatter. The model's structure, not just its size, decides which analyses are pessimistic, which are optimistic, and which are simply wrong.

## Global, local, systematic and random

Variation is classified by what it is shared with. **Global** (die-to-die, wafer-to-wafer, lot-to-lot) variation moves every device on a die together: one wafer is a little slow. **Local** (within-die) variation makes devices on the same die differ. Within-die variation splits again. **Systematic** variation is predictable from layout context, such as gate lengths printing shorter in dense arrays, and belongs in the nominal model. **Random** variation, such as random dopant fluctuation and line-edge roughness, is independent from device to device and can only be described statistically. Supply voltage and temperature add **environmental** variation in the field, with their own statistics. A model writes a parameter as $p = p_0 + G + S + R$: nominal, global, spatial and random terms.

## Splitting the variance

Independent components add in variance, not in sigma. If a device's threshold voltage has a global part with $\sigma_g = 12$ mV and a local random part with $\sigma_l = 9$ mV, its total spread across many dies is $\sqrt{12^2 + 9^2} = 15$ mV, not 21 mV. The split also sets correlation. Two devices on one die share $G$, so their covariance is $\sigma_g^2$ and their correlation is $\sigma_g^2/(\sigma_g^2 + \sigma_l^2) = 144/225 = 0.64$. Their difference cancels $G$ and keeps two independent random parts, so the mismatch sigma is $\sqrt{2}\,\sigma_l \approx 12.7$ mV. One total sigma cannot answer any of these questions: the engine must carry the components separately.

## Why the structure matters downstream

Shared and independent components behave oppositely along a path. Twenty gates with a shared global sigma of 2 ps add coherently to 40 ps, while twenty independent local sigmas of 2 ps add to only $2\sqrt{20} \approx 8.9$ ps. A timer that lumps them gets one of the two badly wrong. Matched analog pairs and SRAM cells care about the opposite part: a differential pair's offset and a bitcell's read margin come from mismatch, because the global shift cancels. Corners model global variation well and mismatch not at all. A Monte Carlo deck that folds systematic shifts into a random sigma gets both the mean and the correlation wrong. Choosing the decomposition is a modeling decision with consequences in every later stage.

## Going deeper

Orshansky, Nassif and Boning, "Design for Manufacturability and Statistical Design: A Constructive Approach" (2008).
