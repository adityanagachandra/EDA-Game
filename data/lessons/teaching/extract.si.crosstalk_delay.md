# Crosstalk-induced delay

When an aggressor switches during a victim's transition, the victim's delay changes: slower if the aggressor opposes it, faster if it helps. The change is crosstalk-induced delay, or delta delay. Setup analysis must use the slow-down and hold analysis the speed-up, so crosstalk can break both checks. Computing delta delay with a simple driver model connects the Miller factor to slack and shows why SI analysis is part of every signoff timing run.

## Delay with a quiet aggressor

Model the victim's driver as a resistance $R_d$ charging a lumped load, with 50% delay $t = 0.69 R_d C$; with kohm and fF the result is in ps. With a quiet aggressor, coupling counts at face value, so $C = C_g + C_c$. A 3.2 kohm driver with 4.5 fF to ground (including pins) and 2.7 fF of coupling has $t = 0.69 \times 3.2 \times 7.2 = 15.9$ ps. This quiet case is what a timer without SI reports when it grounds the couplings. Real delay calculators use effective-capacitance iteration and the reduced RC network instead of a lumped load, but the dependence of delay on the coupling's effective value is the same, and the lumped model keeps the arithmetic transparent.

## The delta delay

When the aggressor opposes the victim, the factor rises from 1 to 2 and the effective load grows by $C_c$; the delay grows by $\Delta t = 0.69 R_d C_c$. When it switches the same way, the factor falls to 0, and the delay falls by the same amount. A 2.4 kohm driver with 2.6 fF of coupling has $\Delta t = 0.69 \times 2.4 \times 2.6 = 4.31$ ps in either direction. The delta is proportional to the driver's resistance: weak drivers on long, coupled nets are the usual victims. Upsizing the victim's driver shrinks the delta; so do spacing, shielding and moving the net to another layer. Delta delay also changes the victim's output slew, which propagates to downstream stages.

## Effect on setup and hold

Setup slack uses the latest arrival, so it takes the slow-down; hold slack uses the earliest arrival, so it takes the speed-up. A hold check with 31 ps of slack through a victim with a 1.7 kohm driver and 3.1 fF of coupling loses $0.69 \times 1.7 \times 3.1 = 3.64$ ps when the aggressor can switch with it, leaving 27.36 ps. A setup path with 5 ps of slack and a 6 ps slow-down fails. Signoff timing therefore runs SI-aware analysis for both checks, with delta delays computed per arc. Because delta delay depends on when aggressors switch, and their switching depends on the delays, the analysis iterates timing windows until they converge, and every ECO that moves wires reopens it.

## Going deeper

Sapatnekar, "Timing" (2004), the chapter on crosstalk-aware timing; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure.
