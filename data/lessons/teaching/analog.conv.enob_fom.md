# ENOB and converter figures of merit

A real converter falls short of its ideal SQNR because noise and distortion add to quantization. The effective number of bits translates measured SNDR back into bits, and figures of merit normalize power by speed and resolution so that converters of different architectures can be compared. Both appear in every converter paper and in the objective functions of converter design-space exploration.

## Effective number of bits

SNDR (also called SINAD) is the ratio of signal power to all noise and distortion, measured with a full-scale sine. Inverting the ideal formula gives the effective number of bits: $\mathrm{ENOB} = (\mathrm{SNDR} - 1.76)/6.02$. A 12-bit converter measuring 68 dB has $\mathrm{ENOB} = 10.3$ bits; the lost bits went to thermal noise, jitter and distortion. ENOB is frequency-dependent, usually falling as the input frequency rises because jitter and settling errors grow, so a datasheet quotes it at a stated input frequency. A system designer budgets in ENOB because it is the resolution the application actually receives.

## The Walden figure of merit

The Walden figure of merit divides power by the number of effective conversion steps per second: $\mathrm{FOM}_W = P/(2^{\mathrm{ENOB}} f_s)$, in joules per conversion step. A converter at 10 mW, 100 MS/s and 9 effective bits has $\mathrm{FOM}_W = 10^{-2}/(512 \cdot 10^8) = 195$ fJ per step. Lower is better. The measure assumes power doubles per added bit, which matches converters limited by matching and speed, typically at low and medium resolution. It is a convenient single number for comparing designs, and a converter design-space explorer often uses it as the scalar objective when sweeping architecture parameters.

## The Schreier figure of merit

For thermal-noise-limited converters, halving noise needs four times the power, not two, so the Walden measure penalizes high-resolution designs unfairly. The Schreier figure of merit uses $\mathrm{FOM}_S = \mathrm{SNDR} + 10\log_{10}(\mathrm{BW}/P)$ in dB, with bandwidth BW in hertz (half the sample rate for a Nyquist converter) and power in watts. A converter with SNDR 74 dB, BW 10 MHz and 5 mW has $\mathrm{FOM}_S = 74 + 10\log_{10}(2 \times 10^9) = 167.0$ dB. Higher is better. Each extra bit adds 6 dB of SNDR and costs four times the power, a 6 dB penalty, so noise-limited designs keep the same score. Which figure of merit is appropriate depends on the regime, and comparing designs with the wrong one misleads an exploration loop.

## Going deeper

Walden, "Analog-to-digital converter survey and analysis" (IEEE JSAC, 1999); Murmann, "ADC Performance Survey" (online, updated annually).
