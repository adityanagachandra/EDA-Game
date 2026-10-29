# SAR ADC bit cycles

The successive-approximation ADC finds the code by binary search: one comparator, a capacitive DAC and a small digital state machine resolve one bit per cycle. Its simplicity and energy efficiency make it the dominant architecture at medium resolution and speed, and its operation is the bisection method run in hardware. Tracing a conversion cycle by cycle is the foundation for understanding its errors, redundancy and timing.

## The DAC trial voltage each cycle

The SAR register starts at zero. In cycle 1 it sets the MSB, so the DAC outputs $V_{ref}/2$, and the comparator checks whether $V_{in} \ge V_{DAC}$. If yes the bit stays 1, otherwise it is cleared. Cycle 2 sets the next bit, giving a trial of $V_{ref}/4$ or $3V_{ref}/4$, and so on: each trial is the midpoint of the remaining interval. For a 6-bit converter with $V_{ref} = 1.6$ V and $V_{in} = 0.34$ V: cycle 1 tries 800 mV (bit 0), cycle 2 tries 400 mV (bit 0), cycle 3 tries 200 mV (bit 1), cycle 4 tries 300 mV (bit 1), cycle 5 tries 350 mV (bit 0). Each trial halves the uncertainty, exactly as bisection halves a bracket around a root.

## Reading out the code

After $N$ cycles the register holds the code. Continuing the example, cycle 6 tries 325 mV (bit 1), so the bits are 001101, code 13. An ideal SAR gives the same code as an ideal truncating converter, $\lfloor V_{in}/\mathrm{LSB} \rfloor$ with $\mathrm{LSB} = 25$ mV: $0.34/0.025 = 13.6$, code 13. The value of the trace is in understanding failures. A wrong early decision, from comparator noise, DAC settling or kickback, cannot be undone by later cycles in a plain binary search, so the code is off by up to half the range. Redundant SARs use non-binary steps (more than $N$ cycles, or radix less than 2) so that later cycles can correct an early error, a standard technique for speed.

## Conversion time and the capacitive DAC

A synchronous SAR spends some cycles sampling and then one cycle per bit. With 2 sampling cycles and 10 bits, a conversion takes 12 cycles, and a 240 MHz clock gives 20 MS/s. Asynchronous SARs trigger each decision when the comparator finishes, saving the margin of a fixed clock. The DAC is usually a binary-weighted capacitor array that also samples the input: $2^N$ unit capacitors for an $N$-bit design, so mismatch and kT/C noise set the unit size, and the array dominates area and switching energy. Switching schemes that avoid charging large capacitors unnecessarily reduce DAC energy substantially. For a tool designer, a SAR is a natural mixed-signal verification case: analog comparator and DAC models, digital control, and timing that couples them.

## Going deeper

Pelgrom, "Analog-to-Digital Conversion" (3rd ed., 2017), chapter on successive approximation; McCreary and Gray, "All-MOS charge redistribution analog-to-digital conversion techniques" (IEEE JSSC, 1975).
