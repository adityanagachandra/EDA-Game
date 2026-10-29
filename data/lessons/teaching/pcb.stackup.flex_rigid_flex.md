# Flex and rigid-flex boards

Flexible circuits put copper on thin polyimide film so the board can bend, either once at installation or repeatedly in service. Rigid-flex boards laminate flexible layers through rigid sections and leave them bare in between, replacing connectors and cables. The design rules change: materials, coverlay instead of mask, bend-radius limits, routing across bends and plane hatching. Board tools model flex regions as zones of the stackup with their own layer sets and rules.

## Flex materials and construction

The base dielectric of a flex circuit is **polyimide** film, which survives soldering temperatures and bends without cracking. Copper is often **rolled-annealed**, whose elongated grains resist fatigue better than electrodeposited copper, which matters for hinges that flex thousands of times. Instead of liquid solder mask, which cracks, flex uses a **coverlay**: polyimide film with adhesive, laminated over the circuit with openings at pads. In **rigid-flex**, the flexible layers continue through the rigid sections, where FR-4 layers are laminated onto them, and are left exposed in the bend region. The board tool therefore needs a stackup per region: the rigid stack, the flex stack, and the transitions where coverlay and rigid material end.

## Bend radius rules

Bending strains the outer copper in proportion to its distance from the neutral axis over the bend radius, so the minimum bend radius scales with the flex thickness. Fabricators publish rules of the form radius at least $k$ times thickness, with $k$ growing with the number of copper layers and much larger for **dynamic** flexing (repeated bending in service) than for **static** bends made once at installation. Example with a rule of 12 times thickness for two copper layers: a 0.2 mm two-layer flex needs at least 2.4 mm; if the same section must flex continuously under a 100 times rule, it needs 20 mm. IPC-2223 gives guidance; the fabricator's own numbers govern. A DRC can check bend lines drawn in the mechanical layer against these limits.

## Routing and features in flex regions

Traces should cross a bend **perpendicular** to the bend line and be staggered between layers, not stacked, because stacked traces act like an I-beam that stiffens the section and concentrates strain. Corners are rounded, traces get teardrops at pads, and **vias and pads stay out of bend areas**, since a rigid plated barrel is a stress riser. Planes in flex are often **cross-hatched** to keep the section flexible, which changes the impedance of traces referenced to them and must be modeled. **Stiffeners** (FR-4 or polyimide bonded locally) support connector and component areas so joints are not flexed. Rigid-flex costs more to fabricate, but it removes connectors and cables, which are a common source of field failures, and saves volume.

## Going deeper

IPC-2223, "Sectional Design Standard for Flexible/Rigid-Flexible Printed Boards". Fjelstad, "Flexible Circuit Technology" (4th ed., 2011).
