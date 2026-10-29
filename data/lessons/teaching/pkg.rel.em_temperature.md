# EM temperature acceleration and Blech length

Electromigration accelerates strongly with temperature through the Arrhenius factor in Black's equation. Reliability engineers exploit this to test in days what happens over years, and designers must account for it, because the local wire temperature, including Joule heating, sets the life. Short lines can also be immune: the Blech effect makes lines below a critical length-current product stop failing.

## Arrhenius acceleration

At fixed current density, the ratio of lifetimes at two absolute temperatures is $AF = \exp\left(\frac{E_a}{k}\left(\frac{1}{T_{use}} - \frac{1}{T_{stress}}\right)\right)$, with $k = 8.617 \times 10^{-5}$ eV/K. For copper with $E_a = 0.85$ eV, a 105 °C use condition (378.15 K) and a 300 °C stress (573.15 K) give $AF = \exp(9864 \times 0.0009) \approx 7150$. Temperatures must be in kelvin. The large factor is what makes accelerated tests possible, but it also means a 10 K error in the use temperature changes predicted life by tens of percent, so Joule heating of the wire itself must be included.

## Extrapolating a stress test

Combining current and temperature, the use life is $MTTF_{use} = MTTF_{stress} (J_s/J_u)^n \cdot AF$. A test at 3 MA/cm^2 and 300 °C fails at a median of 200 hours. At 1.5 MA/cm^2 and 105 °C with $n = 2$ and the factor 7150, the median use life is $200 \times 4 \times 7150 = 5.72$ million hours, about 653 years. That seems huge, but the design target is the time to a small failure fraction across millions of lines, not the median of one, which brings the margin down by orders of magnitude. The extrapolation also assumes that the same failure mechanism operates at test and use conditions.

## The Blech length

In a line blocked at its ends by diffusion barriers, electromigration builds a back-stress gradient that pushes atoms the other way. If the line is short enough, the back stress balances the electron wind and net transport stops. The criterion is a product: lines with $jL < (jL)_c$ are immortal, with a threshold of the order of a few thousand A/cm for copper, depending on the dielectric and cap. At $j = 2$ MA/cm^2 and $(jL)_c = 3000$ A/cm, the critical length is $3000/(2 \times 10^6) = 1.5 \times 10^{-3}$ cm, or 15 um. EM checkers use the Blech effect to waive short segments, and designers deliberately break long power straps into short segments between vias.

## Going deeper

Blech, "Electromigration in thin aluminum films on titanium nitride", Journal of Applied Physics (1976); Lienig and Thiele, "Fundamentals of Electromigration-Aware Integrated Circuit Design" (2018).
