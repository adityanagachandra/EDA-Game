# DIBL threshold shift

Drain-induced barrier lowering is the short-channel effect that makes a transistor's threshold depend on its own drain voltage. It is quoted for every process in millivolts per volt, it is one of the first parameters a compact-model extraction fits, and it directly multiplies off-current at the supply voltage. Computing its size and its leakage cost turns the qualitative short-channel story into numbers a power tool uses.

## Threshold at high drain bias

In a short channel, the drain's field reaches the source end of the channel and lowers the potential barrier that electrons must cross. The gate then needs less voltage to turn the device on. To first order the threshold falls linearly with drain voltage:
$$V_t(V_{DS}) = V_{t,lin} - \eta\,(V_{DS} - V_{DS,lin}),$$
where $\eta$ is the DIBL coefficient in mV/V and $V_{t,lin}$ is measured at a small drain voltage, typically 50 mV.

Example: $V_{t,lin} = 0.38$ V at 0.05 V, $\eta = 80$ mV/V. At $V_{DS} = 0.85$ V, $V_t = 0.38 - 0.080 \times 0.80 = 0.316$ V.

Well-controlled FinFETs and nanosheets reach DIBL of a few tens of mV/V; late planar bulk devices at minimum length were often near or above 100 mV/V. Because the saturation threshold is what governs both off-current and drive current at the supply, DIBL feeds both speed and leakage.

## Extracting the DIBL coefficient

Extraction measures threshold twice, once in the linear regime at small $V_{DS}$ and once in saturation at $V_{DS} = V_{DD}$, using the same threshold definition (often constant current). Then
$$\eta = \frac{V_{t,lin} - V_{t,sat}}{V_{DD} - V_{DS,lin}}.$$

Example: $V_{t,lin} = 0.350$ V at 0.05 V and $V_{t,sat} = 0.290$ V at 0.75 V. $\eta = 60\ \text{mV}/0.70\ \text{V} = 85.7$ mV/V.

The value depends on the threshold definition and on length: it grows steeply as $L$ approaches the natural length. A compact model must reproduce DIBL across the length range with parameters that scale smoothly, and an extraction agent that fits only one length will produce a model that predicts the wrong leakage elsewhere. In BSIM-family models, DIBL parameters are among those fitted against short-channel devices at high drain bias.

## The leakage cost of DIBL

Off-current is exponential in threshold, so the DIBL shift multiplies leakage by
$$\frac{I_{off}(V_{DD})}{I_{off}\text{ at the low-}V_{DS}\text{ threshold}} = 10^{\eta\,\Delta V_{DS}/S}.$$

Example: $\eta = 80$ mV/V, $\Delta V_{DS} = 0.8$ V and $S = 75$ mV/decade. The threshold drop is 64 mV, which is 0.853 decades: a factor of 7.1 in off-current.

This is why DIBL appears on every process comparison table next to subthreshold swing: together they set how much threshold a device needs for a given leakage budget. It also explains why DIBL matters to timing: drive current rises with $V_{DS}$ through the same barrier lowering, so the current delivered during a transition depends on the output voltage. A model that gets DIBL right only at $V_{DS} = V_{DD}$ will mispredict the shape of the output transition.

## Going deeper

Taur and Ning, "Fundamentals of Modern VLSI Devices" (2nd ed.), section 3.2; the BSIM4 and BSIM-CMG technical manuals, sections on short-channel and DIBL parameters.
