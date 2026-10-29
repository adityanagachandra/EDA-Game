# Stackup build: thickness, layer count and symmetry

A stackup is the ordered list of copper layers and dielectrics that make up the board: cores, prepregs, copper weights and the finished thickness. It sets every impedance, the via aspect ratio, the layer count available to routing and the board's mechanical flatness. The layout tool stores it as a table that every SI calculation reads, and the fabricator turns it into a lamination recipe. Three computations recur: finished thickness, the number of cores and prepregs, and a symmetry check.

## Cores, prepreg and copper

A **core** is a cured glass-epoxy laminate with copper on both sides. **Prepreg** is glass cloth with partly cured resin that flows and cures in the lamination press, bonding cores and foils. Copper weight is quoted in ounces per square foot: 1 oz is about 1.37 mil (35 micrometres) thick. Finished thickness is the sum of the dielectrics plus the copper. Example: a 4-layer stack of 1 oz, 8 mil prepreg, 1 oz, 40 mil core, 1 oz, 8 mil prepreg, 1 oz has $56 + 4 \cdot 1.37 = 61.48$ mil, close to the common 62 mil board. Real fabricators adjust for resin flow into inner copper, plating on outer layers and solder mask, which is why the fabricator's stackup proposal, not the designer's estimate, is the reference.

## Layer count and construction

Most boards use an even number of copper layers, because cores carry copper in pairs. In a **foil build** of an $L$-layer board, the inner layers sit on $(L - 2)/2$ cores, and the two outer layers are copper foils pressed onto prepreg. The board has $L - 1$ dielectric layers, so $L - 1 - (L - 2)/2 = L/2$ of them are prepreg. An 8-layer foil build has 3 cores and 4 prepreg layers. Layer count is chosen from routing demand (how many signal layers BGA escape needs), the number of reference planes high-speed signals require, and power distribution. Each signal layer should sit next to a solid reference plane, which is why a typical 6-layer board has 4 signal layers or fewer.

## Symmetry and why boards warp

A stackup should be **symmetric** about its center: layer $k$ and layer $L + 1 - k$ with the same copper weight, and each dielectric matching its mirror. Copper and glass-epoxy expand at different rates, so an unbalanced stack bends during lamination cooling and again in every reflow, producing bow and twist that break BGA joints and jam assembly fixtures. A symmetry check compares mirrored pairs: in a 6-layer stack there are 3 copper pairs and 2 dielectric pairs, since the center dielectric mirrors onto itself. Copper distribution matters too: a solid plane mirrored by a sparse signal layer is unbalanced even with equal weights, so designers add copper thieving or balance planes. A stackup editor can run this check automatically on every edit.

## Going deeper

Ritchey, "Right the First Time: A Practical Handbook on High Speed PCB and System Design", volume 1, on stackup planning. IPC-2221, section on board construction.
