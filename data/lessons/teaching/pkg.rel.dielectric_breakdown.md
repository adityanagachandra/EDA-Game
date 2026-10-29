# Time-dependent dielectric breakdown

Time-dependent dielectric breakdown (TDDB) is the gradual degradation of an insulator under an electric field until a conductive path forms and it shorts. It limits gate oxides, the low-k dielectric between closely spaced copper lines, the liners of TSVs and the dielectrics of hybrid bonds. Lifetime models in field and Weibull statistics in area let engineers extrapolate from accelerated tests to use conditions and to the full chip.

## The electric field

The field across a dielectric of thickness $t$ at voltage $V$ is $E = V/t$. A useful conversion: 1 V across 1 nm is $10^7$ V/cm, or 10 MV/cm. So 0.8 V across 2 nm of gate dielectric is 4 MV/cm, and 0.75 V across 15 nm of low-k between lines is 0.5 MV/cm. Low-k dielectrics break down at far lower fields than thermal oxide, and they are porous and damaged by etch, so line-to-line TDDB becomes a limit as spacing shrinks. In hybrid bonding, pad misalignment shrinks the dielectric gap between a pad and its neighbour, which raises the field and makes overlay a reliability parameter as well as a yield one.

## Field acceleration models

Lifetime falls steeply with field. The E-model takes time to failure as $t_f \propto e^{-\gamma E}$, so the ratio between a use field $E_u$ and a stress field $E_s$ is $e^{\gamma(E_s - E_u)}$. With $\gamma = 2$ cm/MV, stressing at 8 MV/cm and using at 5 MV/cm gives a factor $e^{6} = 403$. Other models are used too: a $1/E$ model, power-law voltage models for thin gate oxides, and the square-root-of-E model proposed for low-k. They agree at stress conditions but diverge by orders of magnitude when extrapolated to low use fields, so the model choice is itself a reliability decision and must be justified by data across several stress fields.

## Weibull statistics and area scaling

Breakdown times follow a Weibull distribution with shape $\beta$ and characteristic life $\eta$. Because breakdown is a weakest-link event, a larger area fails sooner: $\eta_2 = \eta_1 (A_1/A_2)^{1/\beta}$. A test structure of 1 unit area with $\eta = 10$ years scaled to a chip with 100 times the area and $\beta = 1.5$ gives $10 \times 0.01^{1/1.5} = 0.464$ years for the chip's characteristic life. That is why a chip with many metres of minimum-spaced lines needs a much larger per-line margin than a test structure suggests. Low $\beta$, typical of thin and defective dielectrics, makes area scaling harsher.

## Going deeper

McPherson, "Reliability Physics and Engineering" (3rd ed., 2019).
