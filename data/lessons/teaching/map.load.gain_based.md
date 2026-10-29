# Gain-based mapping with fixed electrical effort

Gain-based mapping removes the unknown load from the mapping problem. Every gate is assumed to run at a fixed electrical effort, so its delay becomes a constant set by its logical effort and parasitic delay. The mapper then compares structures, not sizes, and a later sizing step chooses each gate's input capacitance to realize the assumed gain.

## Delay at fixed gain

In the logical-effort model a gate's delay is $d = g h + p$ in units of $\tau$, with logical effort g, electrical effort (gain) $h = C_{out}/C_{in}$ and parasitic delay p. Fix h for every gate, say h = 5. Then each gate has a load-independent delay: INV (g = 1, p = 1) gives 6, NAND2 (g = 4/3, p = 2) gives 8.67, NAND3 (g = 5/3, p = 3) gives 11.33. A path's delay is the sum, so NAND3 followed by INV takes 17.33 $\tau$. The mapper can run the exact load-independent DP on these numbers.

## Sizing realizes the gain

Fixing h is a promise about sizes. After mapping, each gate's output load is known: the input capacitances of its fanouts plus wire. The sizer sets $C_{in} = C_{out} / h$. A gate driving 10 fF of pins and 2 fF of wire at h = 3 gets $C_{in}$ = 4 fF, which becomes part of its driver's load, so sizing proceeds from the outputs backward. Because each gate keeps the same gain, the delays the mapper assumed are exactly realized, up to discreteness of the library sizes and limits on the largest and smallest cells.

## Comparing structures, not loads

Under fixed gain, choosing between implementations depends only on g and p. A 4-input AND as NAND4 then INV costs (2h + 4) + (h + 1); as two NAND2 feeding a NOR2 it costs (4/3 h + 2) + (5/3 h + 2) on its critical path. At h = 5 these are 20 and 19 $\tau$, so the two-level NAND-NOR form wins. The method avoids the crossover problem entirely, but it assumes every gate can be sized freely and that wire load is small or included in $C_{out}$. Large fanouts and long wires break the fixed-gain assumption and need buffering, which is why gain-based mapping is followed by buffer insertion.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort: Designing Fast CMOS Circuits" (1999), for the delay model and the sizing rule that realizes a chosen gain.
