# The ideal diode equation

The ideal diode equation describes current through a p-n junction as an exponential of applied voltage. It is the template for subthreshold MOSFET current, junction leakage, ESD clamps, temperature sensors and the bandgap references that set on-chip voltages. Circuit simulators evaluate it inside every source and drain junction of every transistor, so its exponential behavior is also a classic source of Newton convergence trouble.

## The diode equation

Forward bias $V$ lowers the junction barrier by $qV$, so the minority density at each edge of the depletion region rises by $e^{qV/kT}$. These injected carriers diffuse away and recombine, and the diffusion current they carry gives the Shockley equation
$$I = I_S\left(e^{V/(n\,kT/q)} - 1\right).$$
$I_S$ is the saturation current, proportional to junction area and to $n_i^2$, and $n$ is the ideality factor: 1 for pure diffusion current, up to 2 when recombination in the depletion region dominates.

Example: $I_S = 10^{-14}$ A, $n = 1$, $V = 0.6$ V at 300 K. Then $V/(kT/q) = 23.21$ and $I = 10^{-14} \times e^{23.21} = 1.21 \times 10^{-4}$ A $= 121$ µA.

Under reverse bias the exponential vanishes and $I \to -I_S$, which is the ideal junction leakage. Real reverse leakage adds generation in the depletion region and, at high fields, band-to-band tunneling.

## Voltage for a given current

Inverting the equation gives the forward voltage needed for a current:
$$V = n\frac{kT}{q}\ln\left(\frac{I}{I_S} + 1\right).$$
For currents far above $I_S$ the $+1$ is negligible.

Example: $I_S = 10^{-14}$ A, $n = 1$, target 1 mA. $V = 0.02585 \times \ln(10^{11}) = 0.02585 \times 25.33 = 0.655$ V.

Because of the logarithm, the forward voltage is nearly constant over a wide range of current, which is why a silicon diode is often described as dropping about 0.6 to 0.7 V. It also means the voltage is a sensitive thermometer: at fixed current it falls by roughly 2 mV/K as $I_S$ rises with temperature. Bandgap references combine this negative temperature coefficient with the positive coefficient of a difference between two such voltages to produce a stable reference.

## Voltage per decade and the ideality factor

Far above $I_S$, a current ratio $R$ needs a voltage step
$$\Delta V = n\frac{kT}{q}\ln R.$$
For a factor of 10 and $n = 1$ this is 59.5 mV at 300 K. An ideality factor of 1.2 stretches it to 71.4 mV per decade.

Measuring the slope of $\log I$ against $V$ therefore extracts $n$ directly, and a slope that changes along the curve reveals which mechanism dominates: recombination ($n$ near 2) at low current, diffusion ($n = 1$) in the middle, and series resistance flattening the curve at high current.

For a circuit simulator the exponential is a hazard. A Newton step that overshoots by 1 V in a diode would ask for $e^{38}$ times the current. SPICE limits the junction voltage update between iterations (the classic pnjlim routine) to keep the iteration stable. Every compact model containing junctions inherits this limiting.

## Going deeper

Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), section 2.3; Nagel, "SPICE2: A Computer Program to Simulate Semiconductor Circuits" (1975), for junction voltage limiting.
