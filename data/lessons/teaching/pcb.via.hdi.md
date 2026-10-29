# HDI build-ups and microvia stacking

High-density interconnect (HDI) boards add thin build-up layers to a conventional core, connected by laser-drilled microvias, so fine-pitch BGAs can escape without through vias. They are described by the notation $i + N + i$: an $N$-layer core with $i$ sequential build-up layers on each side. The notation fixes the layer count, the number of lamination and drilling cycles, and how many via structures a connection stacks, which together drive cost and reliability.

## The i plus N plus i notation

In $i + N + i$, the core is a conventional $N$-layer board, often with buried through-core vias, and $i$ build-up layers are added on each side. The total is $2i + N$ layers. A 1+6+1 board has 8 layers; a 2+4+2 board also has 8, but with two sequential build-ups per side. The top build-up layers are L1 to $Li$, the core spans $L(i+1)$ to $L(i+N)$, and the bottom build-up layers follow. The notation tells a router which layer spans are available: microvias between adjacent build-up layers and from the outer build-up to the core surface, buried vias through the core, and possibly through vias. "Any-layer" HDI makes every layer a build-up layer with microvias between any adjacent pair.

## Sequential lamination and drilling

Each build-up layer is added in its own cycle: laminate a dielectric and copper foil onto both sides, laser-drill microvias, plate, image and etch. With the core laminated once, an $i + N + i$ board needs $1 + i$ lamination cycles: a 2+4+2 board needs 3. Drilling follows the same structure: one mechanical drill of the core for buried vias, one laser-drill step per build-up layer per side, and one final mechanical drill for through holes, $2i + 2$ operations in all for a board that uses every kind. Each cycle adds press time, registration error between layers, yield loss and cost, which is why HDI price rises steeply with $i$ and why designers try a 1+N+1 build before a 2+N+2.

## Stacked connections and reliability

A connection from L1 into the core uses one microvia per build-up layer, stacked or staggered, and then the buried via. Example on a 2+6+2 board: L1 to L5 uses microvias L1 to L2 and L2 to L3 (the core's top), then the buried via from L3 through the core, which also connects L4 to L8, so 3 via structures. A connection from L1 to L10 uses 2 microvias, the buried via and 2 more microvias, 5 structures. Every interface between stacked microvias, and between a microvia and the buried via's cap, is a place where thermal-cycling stress concentrates; stacked microvias over buried vias have shown latent failures, so many designs stagger them. A layer-assignment engine can count structures per connection and penalize deep stacks as a reliability cost, not just a fabrication cost.

## Going deeper

IPC-2226, "Sectional Design Standard for High Density Interconnect Printed Boards". Holden, "The HDI Handbook" (2009).
