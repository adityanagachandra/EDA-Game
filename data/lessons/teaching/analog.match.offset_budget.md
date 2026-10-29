# Offset and current-mismatch budgets

Pelgrom's law gives the mismatch of one pair. A circuit's offset or current error combines the mismatch of several pairs, each weighted by how strongly it affects the output. Building that budget tells a designer which devices to enlarge and tells a sizing tool which constraint to impose on which device. The weighting factors turn out to be ratios of transconductances, so gm/Id choices directly shape mismatch.

## Current mismatch in a mirror

A threshold difference $\Delta V_T$ between mirror devices acts like an error on the shared gate voltage. The resulting relative current error is $\Delta I/I = (g_m/I_D)\Delta V_T$. A mirror at $g_m/I_D = 10\ \mathrm{V}^{-1}$ with $\sigma(\Delta V_T) = 2$ mV has $\sigma(\Delta I/I) = 2\%$; at $g_m/I_D = 5$ it would be 1%. So current mirrors and current sources prefer low $g_m/I_D$ (strong inversion): their threshold mismatch converts into less current error. Input pairs prefer the opposite, high $g_m/I_D$, because for them what matters is the input-referred voltage, not the current. One device-level choice therefore has opposite best values depending on the device's role.

## Adding independent contributions

Threshold and current-factor mismatch are independent, so their contributions add in variance: $\sigma^2(\Delta I/I) = (g_m/I_D)^2\sigma^2(\Delta V_T) + \sigma^2(\Delta\beta/\beta)$. With a 2% threshold term and a 0.5% current-factor term, the total is $\sqrt{4 + 0.25} = 2.06\%$. In strong inversion with large overdrive the current-factor term can dominate; in weak inversion the threshold term always does. The same root-sum-square rule applies to any set of independent error sources, and it is why a budget is written as a table of variances: the largest term is where area or redesign pays off, while shrinking a small term changes almost nothing.

## Input-referred offset of an OTA

In a five-transistor OTA, input-pair mismatch appears directly at the input. Mismatch in the mirror load produces a current error $g_{m3}\Delta V_{T3}$, which the input pair must cancel with a voltage $\Delta V_{T3}\,g_{m3}/g_{m1}$. So $\sigma^2(V_{os}) = \sigma^2(\Delta V_{T1}) + (g_{m3}/g_{m1})^2\sigma^2(\Delta V_{T3})$. With 1.5 mV for each pair and $g_{m3}/g_{m1} = 0.5$, $\sigma(V_{os}) = \sqrt{2.25 + 0.5625} = 1.68$ mV. The design lesson: make the load's $g_m$ small relative to the input pair's, by running the load at lower $g_m/I_D$. An equation-based sizer can impose this budget as a constraint and size each pair's area from Pelgrom's law.

## Going deeper

Kinget, "Device mismatch and tradeoffs in the design of analog circuits" (IEEE JSSC, 2005); Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapter 19.
