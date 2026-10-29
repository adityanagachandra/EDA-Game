# Noise factor and noise figure

Every RF block adds noise. The noise factor measures how much a block degrades the signal-to-noise ratio of what passes through it, relative to the thermal noise of a standard 290 K source. It is the key sensitivity specification of receivers, the quantity cascaded by the Friis formula, and a standard output of RF noise simulation.

## Definition from SNR

The noise factor is $F = \mathrm{SNR}_{in}/\mathrm{SNR}_{out}$, measured with the input noise set by a source at the standard temperature $T_0 = 290$ K. In decibels it is the noise figure, $\mathrm{NF} = 10\log_{10}F = \mathrm{SNR}_{in,dB} - \mathrm{SNR}_{out,dB}$. A block that turns a 40 dB input SNR into 37 dB has NF = 3 dB and $F = 2$: it doubles the noise power relative to the signal. A noiseless block has $F = 1$, NF = 0 dB; no physical block does better. Because noise powers of independent sources add, $F = 1 + N_{added}/N_{source}$, where both noises are referred to the input.

## Noise temperature

An equivalent description replaces the block's added noise by a noise temperature $T_e$: the extra source temperature that would produce the same output noise from a noiseless block. Since source noise power is $kT_0B$, $F = 1 + T_e/T_0$, so $T_e = (F - 1)\cdot 290$ K. An amplifier with NF = 1 dB has $F = 1.259$ and $T_e = 75.1$ K; NF = 3 dB gives 289 K. Noise temperature is preferred for very low-noise systems, such as radio astronomy or satellite receivers, where noise figures are fractions of a dB and temperatures give better resolution. It also makes the physics explicit: added noise power is $kT_eB$.

## Sensitivity of a receiver

The thermal noise available from a matched source is $kT_0 = -174$ dBm per hertz at 290 K. Over a bandwidth $B$ the floor is $-174 + 10\log_{10}B$ dBm, and the receiver adds its noise figure. To demodulate with a required SNR, the minimum input power is $P_{min} = -174 + 10\log_{10}B + \mathrm{NF} + \mathrm{SNR}_{req}$. For a 1 MHz channel, NF = 6 dB and 10 dB required SNR: $-174 + 60 + 6 + 10 = -98$ dBm. Every dB of noise figure is a dB of sensitivity, which is why the first stages of a receiver are designed for low noise. RF simulators compute NF by noise analysis around the operating point, or by periodic noise analysis for mixers, whose frequency translation folds noise from several bands.

## Going deeper

Pozar, "Microwave Engineering" (4th ed., 2012), chapter 10; Razavi, "RF Microelectronics" (2nd ed., 2011), chapter 2.
