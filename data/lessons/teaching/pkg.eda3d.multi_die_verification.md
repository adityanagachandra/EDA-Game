# Multi-die assembly verification

Single-die DRC and LVS prove one die correct against its own netlist. A multi-die package adds new ways to fail: bumps that do not line up between dies, a die flipped or rotated the wrong way, nets connected to the wrong neighbour pin, or package wiring that does not match the system netlist. Multi-die verification extends physical verification to the assembly, using a description of how dies are placed and connected.

## The assembly description

Assembly verification needs a machine-readable description of the stack: which dies are present, their placement, rotation and flip in a common coordinate system, which interfaces connect, and a system-level netlist that says which bump of one die connects to which bump of another or to the package. Without it, a checker sees only separate layouts. Face-to-face bonding mirrors one die, so coordinates must be transformed consistently, and a single transform error shifts every bump. Industry formats and open descriptions for 3D assemblies exist for exactly this reason. The verification engine treats the description as the source of truth, so it must itself be checked against the design intent.

## Interface DRC

Interface DRC checks the geometry where dies meet. Bump or pad centres on both sides must align within tolerance, pad sizes and pitches must match the bonding rules, and keep-out regions such as TSV zones must not collide with active circuits on the other die. It also checks alignment marks, dummy pads for planarity in hybrid bonding, and minimum distances to die edges. Because each die was checked separately against its own rules, the interface checker works on a merged view, often only the top layers and bump layers of each die, transformed into the shared frame. A misalignment of one pitch is a classic error that single-die DRC cannot find.

## Assembly LVS

Assembly LVS proves that the physical connections between dies and the package implement the intended system netlist. It extracts connectivity through each die's bump layer, the interposer or substrate wiring, and the bonds between them, then compares the resulting graph with the system netlist, the same graph-matching idea as die LVS at a coarser level. Typical findings are swapped bump assignments, power nets shorted to ground through a package plane, and an unconnected redundant lane. Since dies can come from different teams or vendors, assembly LVS often uses abstract views of each die, only its bump-level pins, so the IP of each die stays private while the connections are verified.

## Going deeper

Lim, "Design for High Performance, Low Power, and Reliable 3D Integrated Circuits" (2013); the IEEE Heterogeneous Integration Roadmap chapter on co-design and simulation tools.
