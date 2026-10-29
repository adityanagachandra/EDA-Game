# Intrinsic gain gm times ro

The intrinsic gain $g_m r_o$ is the largest voltage gain a single transistor can deliver, reached when it drives an ideal current-source load. It sets how many stages or cascodes an amplifier needs, and it has collapsed with scaling: from several hundred in long-channel devices to tens or less in nanometre processes. Sizing tools treat it as a device figure of merit read from lookup tables.

## The best single-transistor gain

A common-source transistor with an ideal current-source load has output resistance $r_o$, so its gain is $-g_m r_o$. Any real load is in parallel with $r_o$ and can only lower the gain. That makes $A_i = g_m r_o$ an upper bound on what one device contributes. With $g_m = 1.5$ mS and $r_o = 80$ kilohms, $A_i = 120$ V/V, about 41.6 dB. Millisiemens times kilohms is dimensionless, a convenient unit pair for hand work. In a differential pair with a mirror load, or in the second stage of an op-amp, each stage's gain is a fraction of this bound, so an 80 dB op-amp needs roughly two stages of 40 dB each, or cascodes.

## The current cancels

Substituting square-law expressions, $g_m = 2I_D/V_{OV}$ and $r_o = 1/(\lambda I_D)$, gives $g_m r_o = 2/(\lambda V_{OV})$. The drain current cancels: intrinsic gain depends on the bias point (overdrive) and on channel length through $\lambda$, not on how much current flows. For $\lambda = 0.1\ \mathrm{V}^{-1}$ and $V_{OV} = 0.2$ V the gain is 100. In gm/Id terms, $g_m r_o = (g_m/I_D)(I_D r_o) = (g_m/I_D)V_A$, where $V_A = 1/\lambda$ is the Early voltage. A device at $g_m/I_D = 15\ \mathrm{V}^{-1}$ with $V_A = 6$ V has intrinsic gain 90, or 39.1 dB. Higher $g_m/I_D$ (lower overdrive) and longer channels both raise gain.

## Why scaling hurts and what designers do

Short-channel effects shrink the Early voltage: drain-induced barrier lowering lets the drain modulate the channel, so $r_o$ falls as $L$ shrinks. In advanced planar nodes a minimum-length device may have intrinsic gain of 10 or less; FinFETs recover some through better electrostatics, but long channels are still used for gain. Designers respond with longer $L$ for gain-critical devices (at the cost of capacitance and speed), cascodes that multiply output resistance by another $g_m r_o$, and gain-boosting loops. For a tool designer, intrinsic gain is a function of $L$ and $g_m/I_D$ that must come from simulated tables, because $\lambda$ is not a constant: it is bias- and length-dependent, and a single number in an optimizer's model will mislead it.

## Going deeper

Jespers and Murmann, "Systematic Design of Analog CMOS Circuits Using Pre-Computed Lookup Tables" (2017), chapter 2.
