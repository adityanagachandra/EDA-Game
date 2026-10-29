# Sampling noise and aperture jitter

Before quantization, a converter samples its input onto a capacitor at a clock edge. Two physical limits act here, independent of the quantizer: thermal noise frozen onto the capacitor, kT/C, and timing error of the clock edge, jitter. Either can cap a converter's SNR below its ideal $6.02N + 1.76$ dB, and both set concrete sizes and specs that a sizing tool must respect.

## kT/C noise

A switch with on-resistance $R$ charging a capacitor $C$ forms an RC low-pass filter. The resistor's thermal noise density $4kTR$ integrated over the filter's noise bandwidth $1/(4RC)$ gives a total noise power of $kT/C$, independent of $R$: more resistance means more noise density but less bandwidth. When the switch opens, this noise is frozen onto the capacitor. At 300 K, 1 pF holds $\sqrt{kT/C} = \sqrt{4.14 \times 10^{-21}/10^{-12}} = 64.4\ \mu$V RMS. Halving the noise needs four times the capacitance. Because the result depends only on temperature and capacitance, kT/C sets a hard lower bound on the sampling capacitor for a given noise target, and with it on the power needed to drive that capacitor.

## Sizing the capacitor against quantization

A common design rule makes the kT/C noise no larger than the quantization noise: $kT/C \le \Delta^2/12$, so $C \ge 12kT/\Delta^2$. For a 10-bit converter over 1 V, $\Delta = 977\ \mu$V and $C \ge 12 \cdot 4.14 \times 10^{-21}/(9.54 \times 10^{-7}) = 52$ fF. At 12 bits it is 834 fF, and at 14 bits over 13 pF. Each bit quadruples the capacitor, and the input driver and reference buffer must charge it within a fraction of a sampling period, so their power scales with it. This is why high-resolution converters become noise-limited and why the energy per conversion of precise converters rises about fourfold per bit, a trend visible in converter surveys.

## Aperture jitter

A sampling instant that wanders by $\delta t$ samples a sine $A\sin(2\pi f t)$ with an error of about $2\pi f A\cos(2\pi f t)\,\delta t$, largest at the zero crossings where the slope is steepest. For RMS jitter $\sigma_t$ the jitter-limited SNR is $-20\log_{10}(2\pi f_{in}\sigma_t)$. A 100 MHz input with 1 ps of jitter gives $-20\log_{10}(6.28 \times 10^{-4}) = 64.0$ dB, about 10.3 bits whatever the quantizer's resolution. The bound falls 6 dB per doubling of input frequency. That is why high-speed converters specify clock jitter in tens of femtoseconds, and why the PLL generating the clock becomes part of the converter's error budget. Total noise adds as powers: quantization, kT/C, jitter and circuit noise together set the measured SNDR.

## Going deeper

Pelgrom, "Analog-to-Digital Conversion" (3rd ed., 2017); Murmann, "ADC Performance Survey" (online, updated annually).
