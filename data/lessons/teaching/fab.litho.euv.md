# Extreme ultraviolet lithography

Extreme ultraviolet (EUV) lithography uses 13.5 nm light, about fourteen times shorter than ArF, so one exposure prints pitches that 193 nm immersion needs three or four masks to reach. It needs reflective optics in vacuum and reflective masks, and its small number of high-energy photons makes randomness a first-order effect. For design, EUV changes mask counts, rule complexity and the variability models for the tightest layers.

## Resolution at 13.5 nm

Every optic for 13.5 nm is a multilayer mirror, because no material is transparent at this wavelength, and the whole path runs in vacuum. Production tools have NA 0.33; high-NA tools have NA 0.55. With $k_1 = 0.4$, NA 0.33 prints a half-pitch of $0.4 \times 13.5/0.33 = 16.4$ nm, against $0.4 \times 193/1.35 = 57.2$ nm for 193 nm immersion. At NA 0.55 the same $k_1$ gives 9.8 nm. EUV masks are reflective, illuminated at about six degrees off normal, which introduces shadowing and mask three-dimensional effects that OPC must model.

## What a single EUV exposure replaces

The value of EUV is counted in masks avoided. If 193 nm immersion runs at $k_1 = 0.3$, its single-exposure pitch is $2 \times 0.3 \times 193/1.35 = 85.8$ nm. A layer at 36 nm pitch would need $\lceil 85.8/36 \rceil = 3$ litho-etch exposures with 193 nm light, or SAQP, against one EUV exposure. Fewer masks mean fewer overlay errors between masks, shorter cycle time and simpler colour-aware rules. EUV exposures are expensive, so the comparison is a cost and yield calculation, and at the tightest pitches EUV itself may be double patterned.

## Photon shot noise

An EUV photon carries $hc/\lambda \approx 91.8$ eV, about 14 times an ArF photon, so a given dose delivers 14 times fewer photons. At 30 mJ/cm^2, photons per square nanometre are $0.03\ \mathrm{J/cm^2} / 10^{14}\ \mathrm{nm^2/cm^2} / 1.47 \times 10^{-17}\ \mathrm{J} = 20.4$. A 10 nm by 10 nm pixel receives about 2040 photons, and as a Poisson count its relative fluctuation is $1/\sqrt{2040} = 2.2\%$. These stochastics cause line-edge roughness, missing or merged contacts and random micro-bridges, failures that a simple critical-dimension rule cannot describe. Raising dose reduces them but lowers scanner throughput.

## Going deeper

Bakshi (ed.), "EUV Lithography" (2nd ed., 2018); De Bisschop, "Stochastic effects in EUV lithography", Journal of Micro/Nanolithography, MEMS, and MOEMS (2017).
