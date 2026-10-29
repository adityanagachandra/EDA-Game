# Via aspect ratios and drill limits

The aspect ratio of a plated hole is its depth divided by its drilled diameter. It limits how well plating solution reaches the middle of the barrel, so fabricators cap it, and the cap ties the board thickness to the smallest through via. Microvias have their own cap, near 1:1, which ties the build-up dielectric thickness to the laser drill size. These are simple divisions, but they set the density floor of the whole design.

## Through-hole aspect ratio

For a through via, aspect ratio is board thickness divided by drilled hole diameter, written $x$:1. Example: a 62 mil board with 10 mil drills has 6.2:1. Plating builds copper on the barrel wall by electrochemical deposition; in a deep, narrow hole the solution exchanges poorly and the current density falls in the middle, so the copper there thins. Thin barrel copper cracks under the expansion of the laminate in thermal cycling, one of the classic board reliability failures. That is why IPC performance classes specify minimum barrel plating thickness, and why fabricators publish a maximum aspect ratio for reliable plating. A drill-table check in DRC computes the ratio for every hole size against the stackup thickness.

## Minimum drill from a fab limit

Inverting the ratio gives the smallest allowed drill: $d_{min} = T / AR_{max}$, rounded up to the drill increment. Example: a fabricator that plates reliably to 10:1 and a 93 mil board give $d_{min} = 9.3$ mil, so 10 mil drills. Thicker boards, needed for high layer counts or connector press-fit, therefore force larger vias, larger pads and fewer routing channels between BGA balls. Designers respond by choosing a thinner board, by moving fanout to blind or microvias whose depth is much smaller than the board, or by paying for an advanced fabricator with a higher limit. Finished hole size is smaller than drill size by twice the plating thickness, so drawings must say which one a number means.

## Microvia aspect ratio

A microvia's depth is the thickness of the dielectric it spans, usually one build-up layer. Laser-drilled holes are cone-shaped and are plated or filled from one side, so the reliable aspect ratio is near 1:1 or lower; many fabricators cap it at 0.75:1 to 1:1. Example: a 4 mil laser drill and a 0.8:1 cap allow dielectrics up to 3.2 mil, so a 3 mil build-up qualifies and a 4 mil one does not. This couples stackup choice to via technology: the build-up dielectric thickness, which also sets the impedance of outer-layer traces, must satisfy the microvia rule. A stackup editor that knows the drill table can flag the conflict as soon as the dielectric is chosen.

## Going deeper

IPC-6012, "Qualification and Performance Specification for Rigid Printed Boards", on plating thickness by class. IPC-2226 on microvia structures.
