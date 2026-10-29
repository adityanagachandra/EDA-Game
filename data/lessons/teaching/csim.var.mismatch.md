# Device mismatch and the Pelgrom model

Two transistors drawn identically, side by side, are never identical. Random dopant fluctuation, line-edge roughness and other local effects give each device its own threshold and current factor. This local mismatch sets the offset of comparators, the accuracy of current mirrors and the stability of SRAM cells. The Pelgrom model captures its scaling with one law: the standard deviation falls with the square root of gate area.

## Mismatch shrinks with area

Pelgrom's law states that the threshold difference between two identical devices has standard deviation
$$\sigma(\Delta V_T) = \frac{A_{VT}}{\sqrt{WL}},$$
where $A_{VT}$ is a process coefficient, typically a few mV um. Averaging local fluctuations over a larger gate reduces their effect, and since the fluctuations are independent per unit area, the variance falls as one over area. With $A_{VT} = 3$ mV um, $W = 2$ um and $L = 0.5$ um, $\sqrt{WL} = 1$ um and $\sigma(\Delta V_T) = 3$ mV. A similar law with a coefficient $A_\beta$ describes the relative mismatch of the current factor. Distance between devices adds a further term in the original model, which careful layout minimizes.

## Sizing for a mismatch target

Inverting the law gives the area required for a target: $WL \ge (A_{VT}/\sigma)^2$. With $A_{VT} = 3$ mV um and a target of 1 mV, $WL \ge 9$ um$^2$; at $W = 4$ um that means $L \ge 2.25$ um. Halving the mismatch quadruples the area, and with it the capacitance and power at a given speed. This square-law trade-off is why precision analog circuits are large and why techniques such as offset calibration, chopping and auto-zeroing exist. A sizing tool that ignores mismatch will shrink matched pairs until Monte Carlo reveals unacceptable offset, so mismatch belongs in the optimization from the start.

## Per-device sigma in a Monte Carlo model

Pelgrom's coefficient describes the difference of a pair. A Monte Carlo simulator instead draws an independent threshold shift for each device. The difference of two independent draws with standard deviation $\sigma_{dev}$ has variance $2\sigma_{dev}^2$, so the model must use
$$\sigma_{dev} = \frac{\sigma(\Delta V_T)}{\sqrt{2}} = \frac{A_{VT}}{\sqrt{2WL}}.$$
For a pair with $\sigma(\Delta V_T) = 3$ mV, each device gets $\sigma_{dev} = 2.12$ mV. Putting the pair value on each device overstates mismatch by $\sqrt{2}$, and placing one random value on a shared model card makes all devices move together, which models global variation and gives zero mismatch. Foundry statistical models encode the correct per-instance draws.

## Going deeper

Pelgrom, Duinmaijer and Welbers, "Matching properties of MOS transistors" (IEEE JSSC, 1989); Kinget, "Device mismatch and tradeoffs in the design of analog circuits" (IEEE JSSC, 2005).
