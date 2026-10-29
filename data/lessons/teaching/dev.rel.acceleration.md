# Temperature and voltage acceleration

A ten-year lifetime cannot be measured before a product ships, so reliability is tested under accelerated stress: higher temperature and higher voltage compress years into weeks. Acceleration models translate stress results back to use conditions. They are the bridge between reliability experiments and the aging models and voltage limits that EDA flows enforce, and their extrapolation is a central source of risk.

## Thermal acceleration

Thermally activated mechanisms follow the Arrhenius model: the rate is proportional to $e^{-E_a/kT}$, where $E_a$ is the activation energy. The acceleration factor between a stress temperature and a use temperature (both in kelvin) is
$$AF_T = \exp\left[\frac{E_a}{k}\left(\frac{1}{T_{use}} - \frac{1}{T_{stress}}\right)\right].$$

Example: $E_a = 0.7$ eV, use at 85 °C (358 K), stress at 150 °C (423 K). $E_a/k = 8124$ K and $1/358 - 1/423 = 4.29 \times 10^{-4}$, so $AF_T = e^{3.49} = 32.7$.

Activation energies differ by mechanism, roughly a few tenths of an eV to about 1 eV, so a single stress test cannot be translated without knowing which mechanism dominates. Using the wrong $E_a$ misestimates lifetime exponentially. Temperature must also be the device's actual channel temperature, which self-heating can raise well above ambient in advanced devices.

## Voltage acceleration

Higher voltage raises the field across the dielectric and the energy of carriers, accelerating wear-out. Several empirical forms are used: exponential in voltage or field (the E model), exponential in inverse field (1/E model) and power law in voltage. In the exponential form,
$$t \propto e^{-\gamma V}, \qquad AF_V = e^{\gamma(V_{stress} - V_{use})}.$$

Example: $\gamma = 8$ V$^{-1}$, stress at 1.2 V and use at 0.8 V. $AF_V = e^{3.2} = 24.5$.

The models agree inside the measured range and diverge wildly when extrapolated to use voltage, which is the core difficulty: extrapolating from 1.2 V to 0.8 V with the wrong model can be off by orders of magnitude. Reliability teams choose models with physical backing for each mechanism (for thin-oxide TDDB, power-law voltage dependence is widely used) and stress at several voltages to constrain the fit.

## Lifetime at use conditions

When temperature and voltage effects are treated as independent, their factors multiply:
$$t_{use} = t_{stress}\cdot AF_T \cdot AF_V.$$

Example: a failure criterion is reached after 500 hours at 150 °C and 1.2 V. With $AF_T = 32.7$ and $AF_V = 24.5$, $t_{use} = 500 \times 32.7 \times 24.5 = 400{,}600$ hours, about 45.7 years.

Independence is itself an assumption; for some mechanisms the activation energy depends on voltage. And the result applies to the stated failure criterion and failure fraction, typically through the Weibull or log-normal statistics of the stressed population. For EDA, these extrapolated lifetimes become the operating limits encoded in reliability rules (maximum voltage per device type, allowed overdrive time, current-density limits) and the stress parameters of aging models that produce aged timing libraries.

## Going deeper

JEDEC JEP122, "Failure Mechanisms and Models for Semiconductor Devices"; McPherson, "Reliability Physics and Engineering: Time-to-Failure Modeling" (3rd ed.).
