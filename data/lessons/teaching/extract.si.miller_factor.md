# The Miller factor of a coupling capacitor

A coupling capacitor's effect on its victim depends on what the other end does. If the neighbour is quiet, the capacitor acts like one to ground. If the neighbour switches the opposite way, the voltage across the capacitor changes twice as much and the victim's driver must supply twice the charge. If it switches the same way, the capacitor carries no charge. The Miller factor, 0, 1 or 2, captures this in a single multiplier, and it is the basis of every switching-aware delay calculation.

## Effective coupling

The current through a capacitor is $C\, d(V_v - V_a)/dt$, where $V_v$ and $V_a$ are the victim and aggressor voltages. When the victim swings by $\Delta V$ and the aggressor is quiet, the change across the capacitor is $\Delta V$: it behaves like $C_c$ to ground, Miller factor 1. If the aggressor swings by $-\Delta V$ at the same time, the change is $2\Delta V$, so the capacitor looks like $2C_c$ to ground, Miller factor 2. If the aggressor swings by $+\Delta V$ simultaneously, the change is zero and the capacitor looks like nothing, Miller factor 0. A 2.2 fF coupling therefore contributes 0, 2.2 or 4.4 fF of effective load depending on the neighbour. Timers replace each coupling with a grounded capacitor of $k C_c$ for the chosen factor $k$.

## Charge through the coupling

The factor is a statement about charge. During the transition the victim's driver pushes $Q = C_c \Delta(V_v - V_a)$ through the coupling capacitor. For a victim rising 0.8 V with an aggressor falling 0.8 V through 5.3 fF, the change across the capacitor is 1.6 V and the charge is 8.48 fC, the same as a 10.6 fF grounded capacitor would draw. With the aggressor quiet the charge is 4.24 fC, and with it rising alongside, zero. The values 0 and 2 assume simultaneous transitions of equal slew; real aggressors switch with different slews and offsets, so the effective factor can lie between the extremes, and some analyses use factors outside 0 to 2 for very different slews. Signoff tools compute the coupling current from waveforms rather than trusting a fixed factor.

## Bounding with unknown aggressors

Without timing information about the aggressors, a timer must bound. Max-delay (setup) analysis assumes every aggressor switches against the victim, factor 2. Min-delay (hold) analysis assumes every aggressor switches with it, factor 0. The difference in effective coupling between the two analyses is $2\sum C_c$; for couplings of 2.5 and 3.5 fF that is 12 fF. That spread is the pessimism of not knowing the aggressors' timing, and on a net whose capacitance is mostly coupling it can double the delay range. Timing windows, logical correlation (two nets that cannot switch in opposite directions) and aggressor filtering by size all exist to shrink it.

## Going deeper

Kahng, Muddu and Sarto, "On switch factor based analysis of coupled RC interconnects", Design Automation Conference (2000); Cheng, Lillis, Lin and Chang, "Interconnect Analysis and Synthesis" (2000).
