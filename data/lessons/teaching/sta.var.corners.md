# PVT corners

A chip must work across manufacturing variation, supply voltage and temperature. Libraries are characterized at a set of process, voltage and temperature corners, and the timer is run at each one. Which corner limits which check is not fixed: setup is usually worst at slow corners and hold at fast ones, but advanced nodes add exceptions such as temperature inversion, so signoff checks every relevant corner rather than guessing.

## Slack at one corner

A corner name such as ss_0p72v_125c encodes the process (slow NMOS and slow PMOS), the supply (0.72 V) and the junction temperature (125 C). At each corner the library gives different cell delays, setup and hold times, and the extractor gives different wire parasitics. A path's setup slack at a corner is computed with that corner's numbers: $S = T - t_{su} - d_{max}$ with ideal clocks, where $d_{max}$ includes clock-to-Q. If at ss_0p72v_m40c the path takes 799 ps with a 39 ps setup against a 750 ps period, its slack is -88 ps. The same path at a fast corner might have 250 ps to spare.

## The worst setup corner

Delay grows as the process gets slower and the voltage lower, so setup is normally limited at a slow-process, low-voltage corner. Temperature is less simple. At higher supply voltages, delay rises with temperature because carrier mobility falls. At low voltages in advanced nodes, the threshold-voltage drop with temperature can dominate, so delay falls with temperature: temperature inversion, which makes the cold slow corner the worst for setup. The timer does not need to know why: it computes slack at every corner and reports the minimum. For corners at -76 and -88 ps among others that pass, the worst setup slack is -88 ps, at the cold corner.

## The worst hold corner

Hold fails when data is fast, so it is normally limited at fast-process, high-voltage corners, where path delays are shortest, and it must also be checked at the slow corners, because clock skew and hold times scale differently from data delays. Hold slack at a corner is $d_{min} - t_h$ with ideal clocks, using that corner's minimum path delay and hold time; the worst hold slack is the minimum over corners. A fix for one corner must be checked at the others: a hold buffer sized at the fast corner adds more delay at the slow corner and eats setup margin there. Signoff therefore runs many corners together, which is the multi-corner part of MCMM.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 10; Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), the section on temperature dependence of delay.
