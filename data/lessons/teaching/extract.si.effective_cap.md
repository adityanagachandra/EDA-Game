# Switching-dependent effective capacitance

Once each coupling has a Miller factor, a victim's effective load is its ground and pin capacitance plus every coupling scaled by its factor. That number goes straight into the delay calculation, and it changes with the aggressors' behaviour. Computing it for specific scenarios, and its worst and best cases, shows how much of a net's timing is decided by its neighbours rather than by its own driver and wire.

## A specific switching scenario

For a victim with ground capacitance $C_g$, pin capacitance $C_p$ and couplings $C_{c,i}$ with factors $k_i$, the effective load is $C_{eff} = C_g + C_p + \sum_i k_i C_{c,i}$. A victim with 3.7 fF of wire to ground and 4 fF of pins, with three aggressors of 0.8, 3.6 and 2.8 fF all switching against it, sees $3.7 + 4 + 2 \times 7.2 = 22.1$ fF. Had they all been quiet, it would see 14.9 fF. The coupling elements must be in the parasitic netlist with their aggressor identities for this to be possible; a grounded-coupling SPEF only supports the quiet case. In tools, the scenario comes from timing windows and switching directions computed by the timer.

## The worst case

For setup analysis without window information, the worst case sets every factor to 2: $C_{eff,max} = C_g + C_p + 2\sum C_{c,i}$. A victim with 10.8 fF to ground, 3.9 fF of pins and couplings of 0.6, 1.3 and 0.5 fF has $C_{eff,max} = 14.7 + 4.8 = 19.5$ fF. Note that this goes beyond the extracted total capacitance of 17.1 fF: the effective load is not bounded by the physical capacitance. Small couplings are often filtered out below a threshold, either grounded or ignored as aggressors, to keep analysis tractable; the filter must be applied carefully, because a hundred small couplings to one bus can add up to a large one.

## The spread between cases

The best case for delay, factor 0 on every coupling, leaves $C_g + C_p$. The ratio of worst to best is $1 + 2\sum C_c / (C_g + C_p)$. For 11.3 fF to ground, 3.3 fF of pins and 7 fF of coupling, it is $28.6/14.6 = 1.96$: the victim's delay could nearly double between scenarios. That spread is why crosstalk is a signoff check of its own. It also tells a router where to act: a net whose coupling is a large share of its capacitance gains most from extra spacing or shielding, and a net whose capacitance is mostly ground and pins gains little.

## Going deeper

Sapatnekar, "Timing" (2004), chapter on crosstalk; the signal-integrity sections of the OpenSTA and OpenROAD documentation.
