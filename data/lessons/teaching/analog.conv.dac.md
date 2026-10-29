# DAC architectures and segmentation

A digital-to-analog converter turns a code into a voltage or current by switching weighted elements: currents, capacitors or resistors. DACs drive outputs and also sit inside SAR and pipeline ADCs. Their linearity is set by element mismatch and by how the code maps onto elements, so the choice between binary, thermometer and segmented decoding is a calculation in counting and in variance.

## Binary, thermometer and segmented

A binary-weighted DAC has one element per bit, bit $k$ weighing $2^k$ units: $N$ switched elements. A thermometer-coded DAC has $2^N - 1$ identical unit elements and turns on $k$ of them for code $k$. A segmented DAC splits the code: the top $M$ bits drive $2^M - 1$ thermometer elements, each worth $2^{N-M}$ units, and the remaining $N - M$ bits drive binary elements. A 12-bit DAC segmented 6 plus 6 has $63 + 6 = 69$ switched sources, compared with 12 for pure binary and 4095 for pure thermometer. The thermometer decoder grows with the number of elements, so segmentation balances decoder size and wiring against linearity.

## DNL and INL from measured levels

Differential nonlinearity measures each step against one LSB: $\mathrm{DNL}(k) = (V_k - V_{k-1})/\mathrm{LSB} - 1$. Integral nonlinearity measures each level against the ideal straight line: $\mathrm{INL}(k) = (V_k - V_0)/\mathrm{LSB} - k$. With the endpoint method, the LSB is the full span divided by $2^N - 1$. For a 3-bit DAC spanning 0 to 7.7 mV, the LSB is 1.1 mV; a step of 0.88 mV from code 3 to 4 has DNL $= 0.88/1.1 - 1 = -0.2$ LSB. A DNL of -1 means a step of zero, a missing code in an ADC that uses the DAC. These definitions are what a converter testbench computes from a ramp or histogram test, and what a behavioral model must reproduce to be useful.

## Mismatch and the major-carry transition

Let each unit element have independent relative error $\sigma_u$. In a thermometer DAC every step switches on exactly one new unit, so each step's error has sigma $\sigma_u$ LSB and the DAC is monotonic by construction. In a binary DAC the mid-scale transition from $011\ldots1$ to $100\ldots0$ turns off $2^{N-1} - 1$ units and turns on $2^{N-1}$ different ones, so $2^N - 1$ independent errors add: $\sigma_{DNL} = \sigma_u\sqrt{2^N - 1}$. For 10 bits and $\sigma_u = 0.5\%$ that is 0.16 LSB, against 0.005 for thermometer. Segmentation puts the large, error-prone transitions into the thermometer part. Combined with Pelgrom's law, these formulas give the unit-element area for a linearity target, a direct sizing equation.

## Going deeper

Pelgrom, "Analog-to-Digital Conversion" (3rd ed., 2017), chapter on DACs; Van den Bosch et al., "A 10-bit 1-GSample/s Nyquist current-steering CMOS D/A converter" (IEEE JSSC, 2001).
