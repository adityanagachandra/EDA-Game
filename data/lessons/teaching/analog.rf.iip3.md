# Linearity: IIP3 and intermodulation

A weakly nonlinear amplifier fed two nearby tones produces third-order intermodulation products close to the tones, inside the band, where no filter can remove them. The third-order intercept point summarizes this nonlinearity in one number. It sets a receiver's tolerance to interferers, is cascaded like noise figure but in reverse, and is computed by two-tone harmonic-balance simulation.

## Two-tone test and the intercept

Model the nonlinearity as $y = a_1x + a_2x^2 + a_3x^3$. Two tones at $f_1$ and $f_2$ produce, among other products, terms at $2f_1 - f_2$ and $2f_2 - f_1$ whose amplitude grows as the cube of the input. On a dB scale the fundamental output rises 1 dB per dB of input and the IM3 product 3 dB per dB, so the lines meet at a fictitious input power, the input-referred intercept IIP3. From one measurement: if the output fundamentals sit $\Delta P$ dB above the IM3 products, $\mathrm{IIP3} = P_{in} + \Delta P/2$. With -30 dBm per tone and IM3 50 dB below the fundamentals, $\mathrm{IIP3} = -30 + 25 = -5$ dBm.

## Predicting IM3 levels

Reversing the relation predicts distortion for any input: the input-referred IM3 level is $3P_{in} - 2\,\mathrm{IIP3}$, and adding the gain gives the output level. An amplifier with IIP3 = 0 dBm and 15 dB gain, driven by two -25 dBm tones, has input-referred IM3 at -75 dBm and output IM3 at -60 dBm. The output intercept is OIP3 = IIP3 + G. Because IM3 falls 3 dB for every dB of input reduction, backing off the input improves the signal-to-IM3 ratio by 2 dB per dB. The same third-order coefficient also sets the 1 dB compression point, roughly 9.6 dB below IIP3 for a memoryless cubic nonlinearity. The extrapolation is only valid in the weakly nonlinear region; near compression, measured slopes deviate from 3.

## Cascading intercepts

For a cascade, with IM3 products assumed to add in phase (a worst case), $1/\mathrm{IIP3} = 1/\mathrm{IIP3}_1 + G_1/\mathrm{IIP3}_2 + G_1G_2/\mathrm{IIP3}_3 + \dots$, in linear units such as milliwatts. Here gain in front makes a stage's linearity matter more, the opposite of the noise cascade. An LNA with IIP3 = 0 dBm (1 mW) and 15 dB gain (31.6) before a mixer with IIP3 = 10 dBm (10 mW): $1/\mathrm{IIP3} = 1 + 3.16 = 4.16$, so IIP3 = 0.24 mW, -6.2 dBm. The mixer, though more linear, dominates because it sees amplified signals. Receiver planning sets the gain distribution so that noise figure and IIP3 targets are both met, and the dynamic range between noise floor and intermodulation is the figure that results.

## Going deeper

Razavi, "RF Microelectronics" (2nd ed., 2011), chapter 2; Lee, "The Design of CMOS Radio-Frequency Integrated Circuits" (2nd ed., 2004).
