# Partial and cumulative antenna ratios

Antenna damage happens layer by layer. When metal $i$ is etched, the layers above it do not yet exist, so only metal connected to the gate through layers 1 to $i$ is electrically attached. Antenna checks therefore compute a partial ratio for each layer and, in many decks, a cumulative ratio that adds them. Getting the process-step connectivity right is the core of an antenna engine.

## Partial ratios

The partial ratio of layer $i$ counts only the layer-$i$ metal already connected to the gate when layer $i$ is etched: connected through layers 1 to $i$. Consider a gate of 0.36 square um with 58 square um of metal 2 reaching it through metal 1 and metal 2, and another 69 square um of metal 2 that reaches it only through metal 4. When metal 2 is etched, the second piece is floating, so the metal 2 partial ratio is $58 / 0.36 = 161.1$, not $127 / 0.36$. The engine computes, per layer, the connectivity of the stack truncated at that layer.

## Cumulative ratios

Charge damage accumulates over steps, so many decks also bound a cumulative ratio: at layer $i$, the sum of the partial ratios of layers 1 to $i$. With partial areas 17, 3, 36 and 41 square um on metals 1 to 4 and a gate of 0.07 square um, the cumulative ratio at metal 4 is $97 / 0.07 = 1385.7$. Metal that joins the gate only through a higher layer never counts in its own layer's partial ratio, so it never enters the sum. This is why breaking a long lower-layer wire with an upper-layer jumper works.

## The worst layer

Each layer has its own limit, often different for thin and thick metals, and the deck checks every one. The worst layer is the one with the largest partial ratio relative to its limit. A net with partial areas 32, 56, 8 and 36 square um on a 0.27 square um gate peaks at metal 2 with 207.4. A repair should target that layer: moving metal 2 length to a higher layer, or adding a diode that protects all layers at once. An antenna engine reports the violating layer so the router can choose the fix.

## Going deeper

The LEF/DEF Language Reference, sections on ANTENNAPARTIALMETALAREA and the cumulative antenna ratios.
