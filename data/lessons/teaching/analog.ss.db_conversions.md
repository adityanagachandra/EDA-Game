# Decibels in analog and RF design

Analog and RF specifications are written in decibels: gain in dB, noise figure in dB, phase noise in dBc/Hz, power in dBm, converter SNR in dB. Logarithmic units turn cascades of multiplications into sums and compress ratios spanning many orders of magnitude. Getting the factor of 10 versus 20 wrong is one of the most common errors in hand analysis and in tool scripts, so the conversions must be automatic.

## Power ratios and voltage ratios

A decibel measures a power ratio: $G_{dB} = 10\log_{10}(P_2/P_1)$. Power is proportional to voltage squared across the same impedance, so a voltage ratio uses $20\log_{10}(V_2/V_1)$. Both give the same number when the impedances are equal. A voltage gain of 1000 V/V is $20\log_{10}(1000) = 60$ dB; a power gain of 1000 is 30 dB. Useful anchors: a factor of 2 in power is 3.01 dB, a factor of 2 in voltage is 6.02 dB, and a factor of 10 in voltage is 20 dB. Op-amp gain, filter attenuation and S-parameter magnitudes are voltage-like quantities and use 20; noise power, signal power and noise factor use 10.

## Absolute levels in dBm

Appending a reference turns a ratio into an absolute level. dBm is power relative to 1 mW: $P_{dBm} = 10\log_{10}(P/1\ \mathrm{mW})$. So 1 mW is 0 dBm, 100 mW is 20 dBm and 1 microwatt is -30 dBm. Converting back, $P = 10^{P_{dBm}/10}$ mW: -17 dBm is $10^{-1.7} = 0.020$ mW. Related units appear constantly: dBc is power relative to the carrier (used for spurs and phase noise), dBV is voltage relative to 1 V, and dBFS is relative to a converter's full scale. A gain in dB added to a level in dBm gives a level in dBm; adding two levels in dBm is meaningless, because the powers must be summed in milliwatts first.

## Cascades add in dB

Because logarithms turn products into sums, the gain of a cascade is the sum of the stage gains in dB. A chain of a 12 dB stage, a stage of 7 V/V and a 3 dB attenuator has total gain $12 + 20\log_{10}7 - 3 = 12 + 16.90 - 3 = 25.90$ dB. Losses are negative gains. This is why RF link budgets, ADC signal chains and photonic loss budgets are written as columns of dB values. The subtlety is what does not add in dB: noise powers from independent sources, intercept points and noise factors in a cascade must be combined in linear units (the Friis and IIP3 cascade formulas), then converted back. A tool that stores everything in dB must convert at exactly those points.

## Going deeper

Lee, "The Design of CMOS Radio-Frequency Integrated Circuits" (2nd ed., 2004), chapter 1; Pozar, "Microwave Engineering" (4th ed., 2012).
