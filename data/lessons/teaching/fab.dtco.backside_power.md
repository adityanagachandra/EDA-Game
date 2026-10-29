# Backside power delivery

Backside power delivery (BSPDN) builds the power distribution network on the back of a thinned wafer and feeds the transistors from below. The front side then carries only signals. It attacks two problems at once: IR drop through tall front-side via stacks and the routing tracks lost to rails and power straps. It also changes the thermal path and the models of several EDA engines.

## Power from the back of the wafer

After the front-side devices and wiring are built, the wafer is bonded face-down to a carrier wafer and the original substrate is ground and etched until the device layer is only a short distance from the new surface. Backside vias, sometimes called nano-through-silicon vias, then connect a backside metal network either to buried power rails or directly to the transistors' source and drain. Thick, low-resistance backside metal layers form the global grid and connect to the package. The front side keeps its signal stack with the power rails and straps removed.

## IR drop with and without backside power

On the front side, current from the package bumps descends the whole metal stack to reach the cells. A path through ten via levels, each with four vias of 5 ohm in parallel, has $10 \times 5/4 = 12.5$ ohm. At 2 mA that is 25 mV of IR drop before the rails. A backside path of 2 ohm drops 4 mV, saving 21 mV. Lower IR drop lets a design run at lower supply margin or higher frequency, and power-grid analysis must be rebuilt for the new stack. The comparison also ignores the rails themselves, which are shorter and wider on the back.

## Trade-offs and tool consequences

Backside power frees front-side tracks and removes power straps, which reduces congestion and can let cells shrink. It costs wafer bonding, thinning and extra metal levels, and it changes the heat path: the bulk silicon that spread heat to the sink is gone, replaced by bonding layers and a carrier, so thermal analysis matters more. Engines that need new models include power-grid and IR analysis, extraction for both sides, placement legalization near backside taps, and thermal simulation. A DTCO comparison of buried-rail front-side delivery against backside delivery must therefore look at routed block PPA, IR drop, temperature and wafer cost together.

## Going deeper

Beyne et al., imec papers on backside power delivery networks in IEDM proceedings; the IRDS "More Moore" chapter on backside power.
