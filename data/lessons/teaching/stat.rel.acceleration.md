# Accelerated life testing

Nobody can wait ten years to qualify a chip. Life tests run at elevated temperature, voltage or current so that failures appear in hundreds of hours, and a physical model of each mechanism translates stress time into field time. The acceleration factor is the hinge of that translation, and a wrong activation energy or exponent changes a lifetime claim by orders of magnitude.

## The Arrhenius acceleration factor

Thermally activated mechanisms proceed at a rate proportional to $e^{-E_a/kT}$, where $E_a$ is the activation energy, $k = 8.617 \times 10^{-5}$ eV/K Boltzmann's constant and $T$ the absolute temperature. The ratio of rates between stress temperature $T_s$ and use temperature $T_u$ is the acceleration factor
$$AF = \exp\!\left[\frac{E_a}{k}\left(\frac{1}{T_u} - \frac{1}{T_s}\right)\right].$$
For $E_a = 0.7$ eV, $T_u = 55$ C (328.15 K) and $T_s = 125$ C (398.15 K): $E_a/k = 8123$ K, the bracket is $0.0030474 - 0.0025116 = 0.0005358$, and $AF = e^{4.353} = 77.7$. Temperatures must be in kelvin; using Celsius is the classic error.

## Equivalent field time

Each stress hour counts as $AF$ field hours. A 1000-hour test at the conditions above represents about $77{,}700$ field hours, about 8.9 years at 55 C. That conversion is how a high-temperature operating life test of roughly a thousand hours supports a ten-year claim. The claim is only as good as the assumption that the stress accelerates the same mechanism that dominates in the field: too much stress can activate a mechanism, such as a package material degradation, that never occurs at use conditions, and its failures then mislead the extrapolation. Different mechanisms have different activation energies, so one test temperature accelerates them by different factors.

## Current and voltage acceleration

Some mechanisms depend on more than temperature. Black's equation for electromigration writes $MTTF = A J^{-n} e^{E_a/kT}$, with current density $J$ and an exponent $n$ commonly near 1 to 2. The ratio of use to stress lifetime is the product of a current factor and a thermal factor:
$$\frac{MTTF_u}{MTTF_s} = \left(\frac{J_s}{J_u}\right)^n \exp\!\left[\frac{E_a}{k}\left(\frac{1}{T_u} - \frac{1}{T_s}\right)\right].$$
Doubling the current with $n = 2$ alone cuts lifetime by four. Voltage-driven mechanisms such as time-dependent dielectric breakdown use exponential or power-law voltage models. Electromigration signoff tools apply exactly these models per wire segment to turn a current density into a lifetime check against a target.

## Going deeper

Black, "Electromigration: a brief survey and some recent results," IEEE Transactions on Electron Devices, 1969; JEDEC JEP122, "Failure Mechanisms and Models for Semiconductor Devices."
