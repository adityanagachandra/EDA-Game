# Mask data preparation

Between tapeout and a finished mask, the layout is transformed several times: layers are derived, corrected, verified, broken into writer primitives and written. This mask data preparation (MDP) chain decides whether the mask can be made, how long it takes to write, and whether a layout's hierarchy survives. Each step is a geometric engine with its own data structures and cost.

## The data path from tapeout to mask

Tapeout delivers drawn layers. Boolean and sizing operations derive the mask layers: a cut mask from line ends, implant masks from device layers, separate colours for multiple patterning. RET and OPC then correct each critical layer, and lithography verification simulates the result. Mask rule check (MRC) enforces the mask shop's own limits on minimum width, space and jog size, so OPC does not create slivers the mask writer cannot reproduce. Mask process correction (MPC) pre-compensates the mask writer's and mask etch's distortions. Finally the data is fractured into the writer's format and written, then inspected and repaired.

## Fracturing and shot count

A variable-shaped-beam (VSB) electron-beam writer exposes one rectangle or trapezoid per shot, up to a maximum size. Fracturing decomposes each post-OPC polygon into such primitives, a rectilinear polygon decomposition, and the write time grows with the number of shots. A rectangle $w$ by $h$ with maximum shot size $S$ needs $\lceil w/S \rceil \lceil h/S \rceil$ shots. An L-shaped polygon fractured into 1200 by 400 nm and 300 by 900 nm pieces with $S = 500$ nm needs $3 \times 1 + 1 \times 2 = 5$ shots. OPC jogs multiply rectangle counts, and curves multiply them again, which is why curvilinear masks were impractical on VSB writers.

## Data volume, hierarchy and multi-beam writing

OPC breaks hierarchy: identical cells in different contexts receive different corrections, so instances diverge and data grows. MDP tools exploit hierarchy where context repeats and flatten where it does not. Post-OPC files for a critical layer are far larger than the drawn layout. Multi-beam mask writers expose the mask with many thousands of beamlets in a raster, so write time depends on area and pixel size rather than on shape count, which finally made curvilinear ILT masks practical at full-chip scale. For a tool designer, MDP is a lesson in where geometry engines meet physical throughput limits.

## Going deeper

Rizvi (ed.), "Handbook of Photomask Manufacturing Technology"; the SEMI P39 OASIS specification for post-OPC data formats.
