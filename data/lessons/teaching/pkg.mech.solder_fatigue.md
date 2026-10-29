# Solder joint fatigue

Solder joints between a die and a substrate, or a package and a board, are sheared every time the temperature cycles, because the two sides expand by different amounts. Repeated plastic shear grows cracks until the joint fails. The joint strain follows from geometry, and the Coffin-Manson relation turns strain into cycles to failure, so a designer can predict how package size, standoff height and underfill trade against lifetime.

## Shear strain in a joint

A joint at distance $DNP$ from the neutral point, the centre that does not move, sees the two sides slide past each other by $\Delta\alpha\Delta T \cdot DNP$. Spread over the joint height $h$, the shear strain is $\gamma = \Delta\alpha\Delta T \cdot DNP/h$. A silicon die on a 17 ppm/K substrate ($\Delta\alpha = 14.4$ ppm/K), a 100 K swing, a corner joint 10 mm from centre and a 100 um standoff give $\gamma = 14.4 \times 10^{-6} \times 100 \times 10/0.1 = 0.144$, a huge 14.4%. Such strain would fail an unprotected joint within a few cycles, which is why flip chip on organic substrates uses underfill: the epoxy couples die and substrate so that the joints no longer carry the full mismatch.

## Coffin-Manson scaling

The Coffin-Manson relation says plastic strain range and cycles to failure follow a power law, $N_f \propto \gamma^{-n}$, with $n$ often near 2 for solders in simple models. So the ratio of lifetimes at two strains is $N_2/N_1 = (\gamma_1/\gamma_2)^n$. Raising strain by 1.5 times at $n = 2$ cuts life to $1/2.25 = 0.444$ of its previous value. The exponent is fitted from accelerated thermal cycling, and models such as Engelmaier's add dependence on cycle frequency and mean temperature. The power law makes the corner joints, with the largest DNP, fail first, so they are often made non-critical or depopulated.

## Geometry levers for life

Because $\gamma \propto DNP/h$ and $N_f \propto \gamma^{-n}$, life scales as $(h/DNP)^n$. A joint that lasts 2000 cycles at 10 mm DNP lasts $2000 \times (10/14)^2 = 1020$ cycles at 14 mm on a larger package, other things equal. Taller joints help, as do smaller packages, smaller temperature swings and better CTE matching. Accelerated test cycles, for example between about -40 and 125 °C, are translated to field conditions with an acceleration factor from the same model. In practice the simple shear formula overstates strain when the package and board bend, so finite-element models calibrate the strain while Coffin-Manson remains the life law.

## Going deeper

Engelmaier, "Fatigue life of leadless chip carrier solder joints during power cycling", IEEE Transactions on Components, Hybrids and Manufacturing Technology (1983); Lau (ed.), "Solder Joint Reliability" (1991).
