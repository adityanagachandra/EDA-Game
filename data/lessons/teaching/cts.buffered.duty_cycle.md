# Duty-cycle distortion in clock buffer chains

A clock buffer rarely has equal rising and falling delays. Each stage shifts one edge more than the other, so the high pulse grows or shrinks a little, and a deep tree accumulates the error. Designs that use both clock edges, latches or half-cycle paths, or that specify minimum pulse widths, are sensitive to this duty-cycle distortion. It is one reason clock trees are often built from inverters rather than buffers.

## One stage at a time

Let a buffer delay its rising output by $t_{LH}$ and its falling output by $t_{HL}$. An input high pulse of width $W$ rises at 0 and falls at $W$; at the output it rises at $t_{LH}$ and falls at $W + t_{HL}$. The output high pulse is $W + (t_{HL} - t_{LH})$. With $t_{LH} = 58$ and $t_{HL} = 66$ ps, each buffer widens the high pulse by 8 ps, and seven buffers turn a 925 ps pulse into 981 ps. The error is systematic, so it adds linearly with depth rather than averaging out.

## Duty cycle at the output

The duty cycle is the high time over the period. A 2000 ps clock at 50% has a 1000 ps high pulse; seven buffers with $t_{HL} - t_{LH} = 7$ ps make it 1049 ps, a duty cycle of 52.45%. Half-cycle paths, which launch on one edge and capture on the other, lose exactly the corresponding time on one of the two half cycles. Static timing sees this when clocks are propagated through characterised rise and fall arcs, and minimum-pulse-width checks catch pulses that become too narrow.

## Why inverters cancel

An inverter turns a high pulse into a low pulse: input rising at 0 gives output falling at $t_{HL}$, input falling at $W$ gives output rising at $W + t_{LH}$, a low pulse of $W + t_{LH} - t_{HL}$. The next inverter turns it back into a high pulse of $W + t_{LH} - t_{HL} + t_{HL} - t_{LH} = W$. Each pair of identical inverters cancels exactly, so an even inverter chain preserves the duty cycle and an odd one leaves one stage's error. Six buffers with a 3 ps mismatch followed by nine inverters with $t_{LH} - t_{HL} = -4$ ps give $625 + 18 - 4 = 639$ ps. Mismatch between nominally identical inverters still remains in silicon.

## Going deeper

Friedman, "Clock distribution networks in synchronous digital integrated circuits", Proceedings of the IEEE (2001); Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003).
