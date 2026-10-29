# Layer jumpers and antenna repair

A layer jumper breaks a long lower-layer route near the gate and carries it over a higher layer. When the lower layer is etched, only the short stub is connected to the gate, so the partial ratio collapses. Jumpers cost two vias and a little upper-layer routing; diodes cost area and capacitance. Antenna repair in a router, such as OpenROAD's repair step, chooses between them per net.

## How a jump cuts the antenna

A metal 1 route 661 um long and 0.14 um wide on a 0.15 square um gate has a ratio of 616.9. If the router jumps to metal 2 at 28 um from the gate and returns to metal 1 only at the far end, then when metal 1 is etched the far piece is not yet connected (it joins only through metal 2), and the gate sees $0.14 \cdot 28 / 0.15 = 26.1$. The metal 2 jumper is short and adds little to the metal 2 partial ratio. The fix uses exactly the process-step connectivity that partial ratios model.

## Where to jump

The jump must happen close enough to the gate that the stub, plus any pin shape, stays within the limit: $(A_{pin} + w s) / G \le R$, so $s \le (R G - A_{pin}) / w$. With $G = 0.065$ square um, a 0.6 square um pin, $w = 0.28$ um and $R = 600$, $s \le 137.1$ um, so the router may run 137 um on metal 1 before jumping. Routers compute this distance during route planning and insert the via there, rather than discovering the violation later.

## Jump or diode

A jump moves the violating length up one layer, so it is impossible on the top layer, and it works only if the upper layer's own partial ratio stays within its limit after the addition. A net violating on metal 3 whose metal 4 partial ratio is 108 can take a jump adding 250 under a limit of 350, while a net violating on the top metal needs a diode. Repair engines try the cheaper fix first, check the result by recomputing ratios, and fall back to diodes. The decision is a small rule system, which makes it a natural Tier 1 policy for an agent to tune.

## Going deeper

The OpenROAD documentation for repair_antennas and the global router's antenna avoidance.
