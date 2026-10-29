# Lithographic resolution

Lithography projects a mask pattern onto photoresist through a lens, and the finest pattern it can print sets the pitch of every critical layer. The Rayleigh relation ties that limit to wavelength, numerical aperture and a process factor $k_1$. A tool designer uses it to see why a layer must be split across masks, why EUV was worth its cost, and why design rules restrict pitches and orientations.

## Diffraction and the Rayleigh relation

A periodic mask of pitch $p$ diffracts light into orders at angles with $\sin\theta = m\lambda/p$. The lens collects angles up to its numerical aperture $NA = n \sin\theta_{max}$. An image forms only if at least two orders pass, so the lens acts as a low-pass filter on the mask's spatial frequencies. The smallest printable half-pitch is written $R = k_1 \lambda / NA$, where $k_1$ absorbs illumination, resist and process quality. For a 193 nm immersion scanner at NA 1.35 and $k_1 = 0.3$: $R = 0.3 \times 193/1.35 = 42.9$ nm, so the densest single-exposure line pitch is about 86 nm.

## The k1 a pitch demands

Turning the relation around tells how hard a pattern is: $k_1 = R \cdot NA / \lambda$. A 76 nm pitch on 193 nm immersion is a half-pitch of 38 nm, so $k_1 = 38 \times 1.35/193 = 0.266$. Production comfort lives around 0.35 to 0.4; below about 0.3 a layer needs strong resolution enhancement, restricted layouts and careful OPC. The hard floor is 0.25: with two-beam imaging the zeroth and one first order can at best sit at opposite edges of the pupil, giving a minimum pitch of $\lambda/(2NA)$. No mask trick prints a denser grating in one exposure, which is why multiple patterning exists.

## Wavelength and numerical aperture as levers

The other two levers are hardware. Wavelength fell from 365 nm (i-line) to 248 nm (KrF) to 193 nm (ArF), and then to 13.5 nm for EUV. NA rose until dry lenses neared 0.93; immersion in water, with index about 1.44, pushed it to 1.35. To print a 40 nm half-pitch at $k_1 = 0.35$ with 193 nm light needs $NA = 0.35 \times 193/40 = 1.69$, beyond any water-immersion lens; at 13.5 nm it needs only 0.12. Each lever has a price: shorter wavelengths need new sources, masks and resists, and higher NA costs depth of focus, which shrinks as $1/NA^2$.

## Going deeper

Mack, "Fundamental Principles of Optical Lithography: The Science of Microfabrication" (2007), chapters 2 and 3; Wong, "Resolution Enhancement Techniques in Optical Lithography".
