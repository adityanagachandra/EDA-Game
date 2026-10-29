# Noise margins

Noise margins quantify how much a logic level can be disturbed, by crosstalk, supply bounce or leakage, before the receiving gate misreads it. They connect the inverter's transfer curve to the noise and signal-integrity checks that signoff tools run on every net. A designer of those tools needs the exact definitions, and needs to know that a margin belongs to a driver-receiver pair rather than to one gate.

## Defining the four levels

Four voltages define the margins. $V_{OH}$ and $V_{OL}$ are the output levels a driver guarantees for a 1 and a 0, close to $V_{DD}$ and 0 for static CMOS. $V_{IH}$ and $V_{IL}$ are input thresholds of the receiver: any input above $V_{IH}$ is read as 1 and any below $V_{IL}$ as 0. The standard choice puts $V_{IL}$ and $V_{IH}$ at the two points where the VTC slope is -1, the unity-gain points; between them the gate amplifies noise. The margins are
$$NM_H = V_{OH} - V_{IH}, \qquad NM_L = V_{IL} - V_{OL}.$$
With $V_{OH} = 1.19$ V, $V_{OL} = 0.01$ V, $V_{IL} = 0.48$ V and $V_{IH} = 0.70$ V, $NM_H = 490$ mV and $NM_L = 470$ mV.

## Margins from a piecewise-linear VTC

A quick estimate replaces the transition by a straight line through the switching point $(V_M, V_M)$ with slope $-g$, where $g$ is the gain, and treats the corners of that line as $V_{IL}$ and $V_{IH}$. The line meets $V_{DD}$ at $V_{IL} = V_M - (V_{DD} - V_M)/g$ and meets 0 at $V_{IH} = V_M + V_M/g$. With $V_{OH} = V_{DD}$ and $V_{OL} = 0$, $NM_L = V_{IL}$ and $NM_H = V_{DD} - V_{IH}$. Example: $V_{DD} = 1$ V, $V_M = 0.5$ V and $g = 10$ give $V_{IL} = 0.45$ V and $V_{IH} = 0.55$ V, so both margins are 450 mV. Higher gain pushes both margins toward $V_M$ and $V_{DD} - V_M$, which is why a centred threshold and high gain together maximize the smaller margin.

## Margins belong to connections

A margin compares one gate's outputs with another gate's inputs, so it is a property of a connection, not of a single cell. A driver with $V_{OL} = 50$ mV feeding a receiver with $V_{IL} = 380$ mV has $NM_L = 330$ mV, whatever the driver's own thresholds are. This matters when cells from different families, voltages or Vt flavours meet: a weak driver, a skewed receiver, or a signal crossing supply domains can shrink a margin to nothing, which is why level shifters exist. Signoff noise analysis goes further: it estimates the glitch that crosstalk injects on a quiet victim net, given the victim driver's holding resistance, and checks whether the glitch, after propagating through the receiver, exceeds what the receiver tolerates. The static margin is the first-order version of that check.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), sections 1.3 and 5.3; Weste and Harris, "CMOS VLSI Design" (4th ed.), section 2.5.
