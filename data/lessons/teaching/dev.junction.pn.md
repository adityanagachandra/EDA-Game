# The p-n junction

A p-n junction forms wherever p-type and n-type silicon meet: every MOSFET source and drain, every well boundary, every ESD diode. Electrons and holes diffuse across the boundary, leave behind ionized dopants and build a potential barrier that stops further net flow. That barrier, the built-in potential, sets junction capacitance, leakage and the threshold of every device built on top of it.

## The built-in potential

Bring p-type silicon ($N_A$) and n-type silicon ($N_D$) into contact. Electrons diffuse from the n side into the p side and holes the other way. They leave ionized donors and acceptors uncovered near the boundary, and that space charge creates a field opposing further diffusion. Equilibrium is reached when the Fermi level is flat across the junction. The bands then bend by $qV_{bi}$, the difference between the two sides' Fermi offsets:
$$V_{bi} = \frac{kT}{q}\ln\frac{N_A N_D}{n_i^2}.$$

Example: $N_A = 10^{17}$, $N_D = 10^{19}$ cm$^{-3}$, $n_i = 10^{10}$ cm$^{-3}$ at 300 K. Then $N_A N_D/n_i^2 = 10^{16}$ and $V_{bi} = 0.02585 \times 36.84 = 0.952$ V.

The built-in potential cannot be measured with a voltmeter, because contact potentials at the metal terminals cancel it. It appears in capacitance-voltage data instead.

## How the built-in potential depends on doping

Because $V_{bi}$ depends on the logarithm of $N_A N_D$, it is a weak function of doping. Multiplying either doping by a factor $F$ adds exactly $(kT/q)\ln F$, whatever the starting values. A factor of 10 adds 59.5 mV at 300 K; a factor of 2 adds 17.9 mV.

This makes $V_{bi}$ a robust quantity. Silicon junctions in practice have $V_{bi}$ between about 0.6 V and 1.0 V, and it cannot exceed roughly the band gap, since the Fermi levels on the two sides are bounded by the band edges.

For a one-sided junction such as an $n^+$-p source on a p-well, the heavy side sits near its band edge and the light side's doping decides the small remaining change. Process tuning of a well implant therefore moves $V_{bi}$ by tens of millivolts, while it moves depletion width and capacitance much more, since those depend on the doping directly.

## Built-in potential at temperature

Temperature enters twice: through $kT/q$ in front and through $n_i^2$ inside the logarithm. Since $n_i$ grows very quickly with temperature, the logarithm shrinks faster than $kT$ grows, and $V_{bi}$ falls as the device heats.

Example: $N_A = 10^{17}$ and $N_D = 10^{19}$ cm$^{-3}$ at 400 K, with $kT/q = 0.03447$ V and $n_i \approx 4.7 \times 10^{12}$ cm$^{-3}$ (a value that must be supplied or computed for the stated gap). Then $\ln(10^{36}/2.2 \times 10^{25}) = 24.5$, so $V_{bi} = 0.845$ V, about 0.1 V below the 300 K value.

Compact models encode this through junction temperature coefficients. A lower $V_{bi}$ at high temperature increases junction capacitance slightly and reduces the barrier that holds off minority carriers. Together with the steep growth of $n_i^2$, this is why junction and diode leakage rise so strongly at hot corners.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), chapter 2; Pierret, "Semiconductor Device Fundamentals", chapter 5.
