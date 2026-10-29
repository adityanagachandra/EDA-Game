# CCS and LVF models

Two Liberty extensions carry what NLDM drops. Composite current source (CCS) models store the current a driver delivers as a function of time, so a timer can reconstruct waveforms into any load. Liberty variation format (LVF) stores the statistical spread of every delay, transition and constraint, so a timer can apply on-chip variation per arc rather than with flat margins. Both are now standard at advanced nodes.

## CCS: the driver as a current waveform

A CCS driver model stores, for each input transition and output load on the table grid, the output current waveform over time, sampled as a sequence of points. Given a load, a timer integrates the current to get the output voltage: for a capacitive load, $V(t) = \frac{1}{C}\int_0^t i\,dt$. With current samples held over 5 ps intervals of 40, 60 and 30 microamperes into 10 fF, the charge is $(40 + 60 + 30) \times 5 = 650$ attocoulombs and the output has risen by 65 mV, since aC per fF is mV. Because the timer has the current, it can drive an RC network properly instead of guessing an effective capacitance. CCS receiver models add input capacitances that differ in the early and late parts of a transition. ECSM is a similar approach that stores voltage waveforms.

## LVF: a sigma for every table entry

LVF adds, alongside each NLDM table, tables of the standard deviation of that quantity due to local random variation, on the same grid: for example ocv_sigma_cell_rise beside cell_rise. A timer doing parametric on-chip variation (POCV) uses a late delay of $\mu + n\sigma$ and an early delay of $\mu - n\sigma$ for some number of sigmas $n$. Example: at one grid point cell_rise is 48 ps and its sigma 3 ps; a 3-sigma late delay is $48 + 9 = 57$ ps. Along a path, independent sigmas add in root-sum-square, so a path's spread grows more slowly than its mean, which is why POCV is less pessimistic than flat derating.

## Asymmetric sigmas and what LVF assumes

Delay distributions are skewed, especially at low voltage, so LVF allows separate early and late sigmas, marked with sigma_type early and late. Early delay is $\mu - n\sigma_{early}$ and late delay is $\mu + n\sigma_{late}$; with $\mu = 60$ ps, $\sigma_{early} = 2$ ps, $\sigma_{late} = 3.5$ ps and $n = 3$, the early and late delays are 54 and 70.5 ps, 16.5 ps apart. LVF still assumes that local variation can be summarized by a mean and spread per arc and per grid point, and that variations on different cells are independent, so it does not by itself model correlated global variation, which is handled by corners. It also multiplies characterization cost, because each sigma comes from Monte Carlo or sensitivity-based simulation.

## Going deeper

The Liberty Reference Manual, CCS timing and variation (LVF) sections; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter on advanced timing models.
