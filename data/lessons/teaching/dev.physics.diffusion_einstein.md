# Diffusion current and the Einstein relation

Carriers also move without any field, from where they are dense to where they are sparse. This diffusion current dominates in the quasi-neutral regions of a forward-biased junction and in a MOSFET below threshold, so it sets diode currents and subthreshold leakage. The Einstein relation links diffusion to mobility, which means one transport parameter describes both mechanisms in device models.

## The Einstein relation

Random thermal motion moves carriers down a concentration gradient. The flux is proportional to the gradient, with the diffusion coefficient $D$ in cm$^2$/s as the constant. In equilibrium the drift and diffusion currents must cancel exactly at every point, and using the Boltzmann form of the carrier density forces
$$D = \mu \frac{kT}{q}.$$
This is the Einstein relation. At 300 K, $kT/q = 25.85$ mV, so an electron mobility of 1000 cm$^2$/(V s) gives $D_n = 25.85$ cm$^2$/s.

The relation holds for non-degenerate carriers near equilibrium, which covers most of a device. Because $D$ and $\mu$ are tied together, device simulators specify mobility models only and derive diffusion from them. Note that $D$ grows with temperature through $kT$ while mobility falls, so the two effects partly cancel.

## Diffusion current density

The electron diffusion current density is
$$J_{n,diff} = q D_n \frac{dn}{dx},$$
and for holes $J_{p,diff} = -q D_p\, dp/dx$; the sign difference comes from the opposite charge. Holes diffuse in the direction of falling hole density, and their current points the same way.

Example: the electron density falls linearly by $10^{16}$ cm$^{-3}$ over 1 µm, with $D_n = 30$ cm$^2$/s. The gradient is $10^{16}/10^{-4} = 10^{20}$ cm$^{-4}$, so $|J| = 1.602 \times 10^{-19} \times 30 \times 10^{20} = 481$ A/cm$^2$.

Large current densities need only modest concentration differences over short distances. In a MOSFET below threshold the channel carrier density falls from source to drain, and diffusion carries the subthreshold current. That is why subthreshold current is exponential in gate voltage: the gate sets the source-end density exponentially through the barrier height.

## Diffusion length

Injected minority carriers diffuse while they recombine. If the minority lifetime is $\tau$, the excess density decays as $e^{-x/L}$ with diffusion length
$$L = \sqrt{D\tau}.$$

Example: $D = 25$ cm$^2$/s and $\tau = 1$ µs give $L = \sqrt{25 \times 10^{-6}}$ cm $= 5 \times 10^{-3}$ cm $= 50$ µm. With $\tau = 10$ ns, $L = 5$ µm.

Diffusion length decides whether a diode is long or short. If the quasi-neutral region is much wider than $L$, the minority profile decays to zero inside it and the current scales as $D/L$. If the region is much thinner, as in heavily doped shallow source and drain regions, the profile is linear and the current scales as $D/W$. Compact junction models fold these cases into the saturation current density parameter, but process changes that alter lifetime, such as contamination or radiation damage, move junction leakage through $L$.

## Going deeper

Pierret, "Semiconductor Device Fundamentals", chapter 3; Sze and Ng, "Physics of Semiconductor Devices" (3rd ed.), section 1.7.
