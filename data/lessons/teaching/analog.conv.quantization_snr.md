# Quantization noise and the 6.02N + 1.76 rule

Quantization replaces each sample by the nearest code, adding an error that behaves like noise for busy signals. Its power fixes the best signal-to-noise ratio any $N$-bit converter can achieve, the famous $6.02N + 1.76$ dB. Every converter spec, ENOB and figure of merit is measured against this ideal, so a tool designer needs its derivation and its limits.

## The ideal SQNR

For a full-scale sine wave of amplitude $V_{FS}/2$, the signal power is $(V_{FS}/2)^2/2 = V_{FS}^2/8$. The quantization noise power is $\Delta^2/12$ with $\Delta = V_{FS}/2^N$. The ratio is $\frac{V_{FS}^2/8}{V_{FS}^2/(12 \cdot 2^{2N})} = 1.5 \cdot 2^{2N}$. In decibels, $10\log_{10}(1.5) + 2N \cdot 10\log_{10}2 = 1.76 + 6.02N$. A 10-bit converter has an ideal SQNR of 61.96 dB; a 16-bit one 98.08 dB. The 6.02 dB per bit comes from halving the step, which quarters the noise power; the 1.76 dB comes from the sine's peak-to-RMS ratio. A different signal, such as a full-scale triangle or noise, gives a different constant.

## The uniform-noise model

Model the quantization error $e$ as uniform on $[-\Delta/2, \Delta/2]$. Its variance is $\int_{-\Delta/2}^{\Delta/2} e^2\,de/\Delta = \Delta^2/12$, so its RMS is $\Delta/\sqrt{12} \approx 0.289\,\Delta$. For a 12-bit converter over 2 V, $\Delta = 488.3\ \mu$V and the RMS error is 141 microvolts. The model is good when the signal crosses many codes between samples and is uncorrelated with the sampling clock. It fails for slow or periodic inputs synchronized with the clock, where the error becomes deterministic and shows up as harmonic spurs rather than a flat floor. Dither, a small added noise, restores the model. A mixed-signal model that adds uniform noise in place of real quantization is only valid in this busy-signal regime.

## Oversampling

Quantization noise of power $\Delta^2/12$ spreads uniformly from 0 to $f_s/2$. If the signal occupies only a band $f_B$, digital filtering removes the noise outside it, and the oversampling ratio $\mathrm{OSR} = f_s/(2f_B)$ divides the in-band noise. The SQNR becomes $6.02N + 1.76 + 10\log_{10}\mathrm{OSR}$ dB: 3 dB, or half a bit, per doubling. A 4-bit quantizer at OSR 256 reaches $25.84 + 24.08 = 49.9$ dB. That is a slow way to buy resolution. Delta-sigma modulators add noise shaping, a feedback loop that pushes quantization noise out of band, so an order-$L$ loop gains about $(2L + 1) \cdot 3$ dB per doubling of OSR. This is why high-resolution audio and sensor converters are delta-sigma.

## Going deeper

Pelgrom, "Analog-to-Digital Conversion" (3rd ed., 2017), chapter on quantization; Schreier and Temes, "Understanding Delta-Sigma Data Converters" (2005).
