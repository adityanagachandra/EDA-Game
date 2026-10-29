# Process corners and what they bound

Process corners are the most used and most misused variation tool in circuit simulation. A corner is a shifted set of global process parameters that makes devices fast or slow. Corners are cheap and deterministic, but they bound only what they were built to bound. A tool designer must know what each corner models, which corner stresses which check, and which effects no global corner can capture.

## The named corners

Corner names list NMOS speed, then PMOS: TT is typical-typical, FF fast-fast, SS slow-slow, FS fast NMOS with slow PMOS, SF the reverse. Fast devices have lower threshold voltage and higher mobility, slow devices the opposite. Each corner is a model library section, selected with a line such as .lib models.lib ss, that shifts parameters for every device on the die together: it models global, die-to-die or wafer-to-wafer variation. The skewed corners FS and SF exist because some circuits fail when the two device types are unbalanced: inverter switching thresholds, ratioed logic and sense amplifiers depend on the NMOS-to-PMOS strength ratio.

## Which corner for which check

Each specification has a worst corner. Setup timing fails when paths are slow: SS devices at low supply and the slowest temperature. Hold timing and dynamic power are worst when devices are fast and supply is high: FF at high voltage. Leakage grows exponentially as thresholds drop and temperature rises, so it is worst at FF and high temperature. The slowest temperature is not always hot: in advanced low-voltage processes, cooling raises the threshold more than it improves mobility, so cold becomes slow, an effect called temperature inversion. A sign-off plan pairs each metric with the corners that bound it rather than running every corner for everything.

## What corners miss

Global corners shift matched devices together, so their difference stays zero: a comparator's offset passes every corner and still fails on silicon. Local mismatch needs Monte Carlo with per-device variation. Corners also bound parameters, not performance: a $\pm 3\sigma$ parameter corner pushes all parameters to extremes at once, which may be too pessimistic for one metric and not worst for another, because performance depends on many correlated parameters nonlinearly. Statistical or performance-specific corners, derived from a block's own Monte Carlo distribution, give sharper bounds. Reusing a digital library's SS and FF as the only checks for an analog reference misses its temperature and mismatch behavior entirely.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed., Addison-Wesley, 2010), chapter on variation; Pelgrom, Duinmaijer and Welbers, "Matching properties of MOS transistors" (IEEE JSSC, 1989).
