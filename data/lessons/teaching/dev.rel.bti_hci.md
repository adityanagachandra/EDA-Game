# Bias temperature instability and hot carriers

Transistors age. Under bias and temperature, traps form at the gate dielectric interface and in the dielectric, and the threshold drifts upward. Bias temperature instability (BTI) happens under static gate bias; hot-carrier injection (HCI) happens while current flows at high drain voltage. Both slow circuits over years, and timing signoff must cover the end of life, not only time zero.

## Bias temperature instability

Negative bias temperature instability (NBTI) affects PMOS devices whose gate is held negative relative to the source and body, at elevated temperature. No drain current is needed: the vertical field across the dielectric and heat create interface traps and trap charge in the dielectric. The threshold magnitude rises and the drive current falls. A PMOS whose input sits low, such as a pull-up in an idle gate, is under NBTI stress.

Positive BTI (PBTI) is the NMOS counterpart under positive gate bias. It was minor with SiO$_2$ dielectrics but became significant with high-k stacks, where electrons trap in the high-k layer.

BTI partially recovers when stress is removed, as some trapped charge escapes. Measured degradation therefore depends on the delay between stress and measurement, and on the duty cycle of the stress. Aging models express degradation as a function of stress time, voltage, temperature and duty cycle, often with separate recoverable and permanent components.

## Power-law time dependence

BTI shift grows sublinearly with stress time, commonly modeled as a power law:
$$\Delta V_t = A\,t^{n},$$
with $n$ typically between about 0.1 and 0.3, and $A$ depending on voltage and temperature.

Example: a stress test shows 10 mV after 1000 hours with $n = 0.2$. At 10 years (87,600 hours) the ratio of times is 87.6, so $\Delta V_t = 10 \times 87.6^{0.2} = 24.5$ mV. Linear extrapolation in time would predict 876 mV, absurdly wrong.

A power law with a small exponent means most degradation happens early, and that extrapolation over many decades of time magnifies any error in the exponent. Changing $n$ from 0.2 to 0.25 changes the 10-year prediction to 30.6 mV. Reliability engineers fit $n$ from stress data measured over several decades of time, at accelerated voltage and temperature, and then extrapolate both in time and in stress conditions.

## Hot-carrier injection

Hot-carrier injection needs current and a high lateral field together. Carriers accelerated through the high field near the drain gain enough energy to create interface traps or to be injected into the gate dielectric. Damage concentrates near the drain end of the channel, where the field peaks, which makes it asymmetric: a device degraded in forward operation behaves differently when source and drain are swapped.

HCI damage accumulates per switching event, so it is driven by activity: clock buffers and output drivers that switch often at full supply age most. Idle gates suffer BTI but little HCI.

Historically HCI was worse at low temperature, because reduced phonon scattering lets carriers reach higher energies, while BTI is thermally activated and worse when hot. In advanced FinFET and nanosheet devices the picture is more complex, with self-heating raising channel temperature and multiple-carrier mechanisms at low voltage. Aging-aware flows therefore need both activity (for HCI) and static signal probabilities (for BTI) per device.

## Going deeper

Schroder and Babcock, "Negative bias temperature instability: road to cross in deep submicron silicon semiconductor manufacturing", Journal of Applied Physics (2003); Grasser (ed.), "Bias Temperature Instability for Devices and Circuits" (2014).
