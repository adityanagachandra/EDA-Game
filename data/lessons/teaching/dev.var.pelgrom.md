# Pelgrom mismatch scaling

Pelgrom's law says that random threshold mismatch between two identical, adjacent transistors shrinks as the square root of their gate area. It is the single most used variability formula in circuit design and the basis of the local-variation parameters in statistical compact models. Analog sizing, SRAM design and the area cost of precision all follow from it.

## Mismatch scales with area

Random local variation averages over the gate area. If each small patch of channel contributes an independent random threshold deviation, the average over area $WL$ has a standard deviation proportional to $1/\sqrt{WL}$. For the threshold difference between two adjacent identical devices,
$$\sigma_{\Delta V_t} = \frac{A_{V_t}}{\sqrt{WL}},$$
where the Pelgrom coefficient $A_{V_t}$ is in mV µm and is a property of the process.

Example: $A_{V_t} = 2$ mV µm, $W = 0.5$ µm, $L = 0.1$ µm. $WL = 0.05$ µm$^2$, so $\sigma_{\Delta V_t} = 2/0.2236 = 8.94$ mV.

The coefficient has fallen with oxide thickness over the generations, since thinner oxide reduces the threshold impact of a given charge fluctuation, and FinFETs reduced it further by removing channel dopants. A distance term for gradients exists in Pelgrom's original paper, but for neighbors the area term dominates.

## Area needed for a matching target

Inverting the law gives the gate area required for a mismatch target:
$$WL = \left(\frac{A_{V_t}}{\sigma_{target}}\right)^2.$$

Example: a comparator needs $\sigma_{\Delta V_t} \le 1$ mV with $A_{V_t} = 2.5$ mV µm. $WL = 6.25$ µm$^2$, more than a hundred times a minimum logic device.

The quadratic cost is the precision-area trade-off of analog design: halving mismatch costs four times the area, and four times the capacitance, which slows the circuit and costs power. This is why offset calibration, chopping and digital correction exist. For an EDA flow, it means device sizes in analog blocks are set by statistics rather than by drive current, and a layout tool must preserve matching by symmetric, interdigitated or common-centroid placement, because gradients add systematic mismatch that Pelgrom's random term does not include.

## From pair mismatch to single-device sigma

Pelgrom's coefficient describes the difference between two devices. If each device's local deviation is independent with sigma $\sigma_{V_t}$, the difference has variance $2\sigma_{V_t}^2$, so
$$\sigma_{V_t} = \frac{\sigma_{\Delta V_t}}{\sqrt 2}.$$
Statistical compact models specify the single-device local sigma, so this conversion matters when comparing a datasheet's mismatch number with a model card.

Example: $A_{V_t} = 1.5$ mV µm and a minimum device of $0.05 \times 0.02$ µm$^2$. $\sigma_{\Delta V_t} = 1.5/\sqrt{0.001} = 47.4$ mV, so $\sigma_{V_t} = 33.5$ mV and $3\sigma = 101$ mV, a large fraction of a 300 mV threshold.

For small devices this is enormous, and for an SRAM array with millions of cells the relevant tail is 5 to 6 sigma. That is why SRAM bitcells, built from near-minimum devices, are the most variability-limited circuits on a chip, and why SRAM yield analysis needs rare-event statistics rather than plain Monte Carlo.

## Going deeper

Pelgrom, Duinmaijer and Welbers, "Matching properties of MOS transistors", IEEE Journal of Solid-State Circuits (1989); Kinget, "Device mismatch and tradeoffs in the design of analog circuits", IEEE Journal of Solid-State Circuits (2005).
