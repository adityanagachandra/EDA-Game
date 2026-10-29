# Clock gating

Clock gating stops the clock to registers that do not need to change, removing the largest term in dynamic power for idle logic. Synthesis inserts integrated clock-gating cells (ICGs) where a register bank holds its value under an enable; CTS then builds the tree through them; and STA checks a special timing arc on their enable. Gating is simple to state and easy to get wrong, because a gated clock that glitches corrupts state silently.

## The integrated clock gate

For positive-edge flops, an ICG contains a latch that is transparent while the clock is low, followed by an AND of the latched enable with the clock. While the clock is low the latch follows the enable; when the clock rises the latch closes and the AND passes either a full pulse or none. A test-enable input forces the gated clock on during scan shift. For flops that capture on the falling edge the polarity flips: a latch transparent while the clock is high and an OR, so the disabled clock rests high. Libraries supply the ICG as one characterised cell, so tools recognise it as a clock element with a known internal timing relationship.

## Why a bare AND gate glitches

Gate a clock with a plain AND and the output follows the enable while the clock is high. For an 800 ps clock high from 0 to 400 ps, an enable that rises at 20 ps produces a 380 ps pulse that starts late; an enable that falls at 150 ps produces a 150 ps runt pulse. A runt pulse may clock some flops and not others, or violate their minimum pulse width, corrupting state. The latch in the ICG removes the problem by sampling the enable only while the clock is low, so the enable may change at any time during the high phase without effect.

## Gating in the flow

Synthesis finds register banks whose next-state logic is a feedback multiplexer, load the input when the enable is true and keep the old value otherwise, and replaces the multiplexer with an ICG on the bank's clock. CTS treats the ICG as a node of the tree, balancing through it so gated and ungated sinks still meet the skew target. STA performs a clock-gating setup and hold check of the enable against the ICG's clock pin. Placement of the gate is a trade-off: near the root it stops more of the tree but sees an earlier clock and tightens enable timing; near the flops it eases timing but saves less. Gating lowers average power, not the worst cycle: many gates opening together still create a large current step.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), the chapter on power; Keating, Flynn, Aitken, Gibbons and Shi, "Low Power Methodology Manual" (2007).
