# Analog-to-digital converters

An ADC samples a continuous voltage and assigns it one of $2^N$ codes. Every mixed-signal system has one at its boundary, and its architecture (flash, SAR, pipeline, delta-sigma) is chosen by the trade-off between speed, resolution and power. For a tool designer, converters are the canonical mixed-signal blocks: analog comparators and DACs driven by digital logic, verified with both kinds of simulation.

## Resolution and the LSB

An $N$-bit converter over a full-scale range $V_{FS}$ divides it into $2^N$ steps of one least significant bit, $\mathrm{LSB} = V_{FS}/2^N$. A 10-bit converter over 1 V has an LSB of $1/1024 = 0.977$ mV; at 16 bits it is 15.3 microvolts. The LSB sets the scale of every error budget in the converter: comparator offset, sampling noise and DAC mismatch must all be small fractions of it. Each additional bit halves the LSB, so the analog precision required doubles per bit, and since noise power and mismatch area scale as the square of precision, the cost of a high-resolution converter grows roughly fourfold per bit in the noise-limited regime.

## From input voltage to code

An ideal truncating converter outputs code $k = \lfloor V_{in}/\mathrm{LSB} \rfloor$, clipped to $[0, 2^N - 1]$; a rounding converter shifts its thresholds by half an LSB. For a 6-bit converter over 2 V, $\mathrm{LSB} = 31.25$ mV, and $V_{in} = 0.71$ V gives $0.71/0.03125 = 22.72$, so code 22, binary 010110. The transfer characteristic is a staircase; the difference between the input and the reconstructed level $k \cdot \mathrm{LSB}$ is the quantization error, bounded by one LSB for truncation and half an LSB for rounding. Real converters deviate from the ideal staircase through offset, gain error and nonlinearity, measured as INL and DNL in LSB.

## Architectures by speed and resolution

A flash ADC compares the input against all $2^N - 1$ thresholds at once: one clock per conversion, but 255 comparators at 8 bits and 1023 at 10, so it stops near 6 to 8 bits. A two-step or subranging converter resolves coarse bits with a small flash, subtracts, and resolves fine bits with another: a 10-bit two-step with 5 and 5 bits needs $2 \cdot 31 = 62$ comparators. Pipelines extend this into many stages for high speed at medium resolution. A SAR converter uses one comparator and a binary search, $N$ cycles per conversion, and is the most energy-efficient at medium speed. Delta-sigma converters oversample with a coarse quantizer and noise shaping for high resolution at low bandwidth. Each architecture is a different split of the work between analog precision and digital processing.

## Going deeper

Maloberti, "Data Converters" (2007); Pelgrom, "Analog-to-Digital Conversion" (3rd ed., 2017).
