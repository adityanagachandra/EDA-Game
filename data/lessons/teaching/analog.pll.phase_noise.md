# Phase noise and jitter

A real oscillator's zero crossings wander in time. In the frequency domain that wandering appears as phase noise, power spread into sidebands around the carrier; in the time domain it is jitter. Phase noise limits radio selectivity, and jitter limits data-converter SNR and serial-link bit-error rates. Converting between the two and scaling phase noise with offset and division ratio are routine calculations.

## Phase noise and its slope

Phase noise $\mathcal{L}(\Delta f)$ is the single-sideband noise power in a 1 Hz band at offset $\Delta f$ from the carrier, relative to the carrier power, in dBc/Hz. In the region dominated by white noise sources, the oscillator integrates frequency fluctuations into phase, so phase noise falls as $1/\Delta f^2$: 20 dB per decade. A VCO at -100 dBc/Hz at 100 kHz offset is -120 dBc/Hz at 1 MHz and -114 dBc/Hz at 500 kHz, since $20\log_{10}5 = 14$. Closer to the carrier, flicker noise steepens the slope to 30 dB per decade, and far out the noise floor flattens. Leeson's model describes this shape with the tank's quality factor, which is why LC oscillators with high-Q tanks are quieter.

## Multiplication by N

Multiplying frequency by $N$ multiplies phase deviations by $N$, which multiplies phase-noise power by $N^2$: the in-band phase noise of a PLL output is the reference's plus $20\log_{10}N$. A reference at -150 dBc/Hz multiplied by 100 appears at -110 dBc/Hz at the output. The same rule runs in reverse: dividing a clock by 4 lowers its phase noise by 12 dB. This sets a floor for synthesizers with large division ratios: the reference and phase detector noise, multiplied by $N^2$, dominate inside the loop bandwidth. It is one reason fractional-N designs, which allow a higher reference and therefore a smaller $N$, can achieve lower in-band noise.

## From phase noise to jitter

Integrating phase noise over offset frequency gives the phase variance: $\sigma_\phi^2 = 2\int \mathcal{L}(f)\,df$, the factor 2 accounting for both sidebands. RMS jitter is $\sigma_t = \sigma_\phi/(2\pi f_0)$. For a flat -120 dBc/Hz out to 10 MHz on a 2 GHz clock, $\sigma_\phi = \sqrt{2 \cdot 10^{-12} \cdot 10^7} = 4.47$ mrad and $\sigma_t = 356$ fs. Real profiles are integrated piecewise over the band of interest, which is set by the application: a serial link integrates over the band its clock recovery cannot track, an ADC over the full band that aliases into its signal. Through the jitter-limited SNR formula, this number connects the PLL's design to the converter's ENOB.

## Going deeper

Leeson, "A simple model of feedback oscillator noise spectrum" (Proceedings of the IEEE, 1966); Hajimiri and Lee, "A general theory of phase noise in electrical oscillators" (IEEE JSSC, 1998).
