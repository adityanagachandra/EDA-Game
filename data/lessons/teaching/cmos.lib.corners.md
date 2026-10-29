# PVT corners and library sets

A library is characterized at one operating point: a process condition, a supply voltage and a temperature (PVT). A design is signed off across a set of such corners, chosen so that every check is bounded by some corner. Knowing which corner stresses which check, and the exceptions such as temperature inversion, is what lets a signoff plan be both safe and affordable.

## Process, voltage and temperature

Process corners capture manufacturing spread of device speed: slow-slow (SS), typical (TT), fast-fast (FF), and skewed corners such as fast nMOS with slow pMOS (FS) and the reverse (SF). Voltage corners cover the supply's tolerance, for example nominal plus and minus some percentage. Temperature corners cover the operating range, such as -40 and 125 degrees Celsius for some products. Each combination is a separate characterization run producing a separate Liberty file: 3 process, 3 voltage and 3 temperature points give 27 libraries before pruning. Interconnect has its own extraction corners (Cworst, Cbest, RCworst, RCbest), which come from parasitic extraction and combine with the device corners.

## Which corner stresses which check

Setup fails on slow paths, so the traditional setup corner is slow process, low voltage and high temperature, where devices are weakest. Hold fails on fast paths, so the traditional hold corner is fast process, high voltage and low temperature. Leakage is worst at fast process, high voltage and high temperature, since leakage grows steeply with temperature. At low supply voltages in advanced nodes, temperature inversion breaks the setup rule: colder devices have higher thresholds, and when gate overdrive is small that effect outweighs the mobility gain, so cells can be slower when cold. Signoff plans therefore check setup at both temperature extremes at low voltage rather than assuming hot is slowest.

## Choosing an affordable corner set

The full product of device corners, extraction corners and design modes is large, and each combination costs a full timing run. Teams choose dominant combinations for each check, guided by experience and by sensitivity studies, and add on-chip variation (OCV derates or LVF) to cover local variation within a corner. Fixes interact across corners: a hold buffer inserted for the fast corner adds delay at the slow corner and eats setup slack, so multi-corner multi-mode optimization evaluates every change at all active corners together. Voltage-scaling library sets, characterized at several supplies, let a timer handle dynamic voltage and frequency scaling designs between characterized points. Every extra corner multiplies characterization cost as well.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter on operating conditions; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 7.2 on variation.
