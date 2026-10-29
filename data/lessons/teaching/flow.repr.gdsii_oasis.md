# GDSII and OASIS mask layout

GDSII and OASIS are the formats in which final layout is handed to the foundry. They store polygons on numbered layers, organized into a hierarchy of cells, and nothing else: no nets, no timing, no logic. Physical verification, mask data preparation and layout viewers read them. GDSII dates from the 1970s Calma systems and is still universal; OASIS (SEMI standard P39) is its compact successor for very large layouts.

## GDSII structures and units

A GDSII file is a binary stream of records. It holds a library of structures (cells). Each structure contains elements: BOUNDARY (a polygon), PATH (a wire with a width), SREF (a placed instance of another structure), AREF (a regular array of instances) and TEXT (labels). Coordinates are 4-byte signed integers in database units, and a UNITS record gives the size of one database unit in metres, commonly $10^{-9}$ (1 nm). A coordinate of 152400 at 1 nm per unit is 152.4 $\mu m$. Layers are identified only by number and datatype; their meaning comes from the foundry's layer map.

## Polygons and file size

A BOUNDARY stores its vertices explicitly, and the polygon is closed by repeating the first point, so a rectangle takes 5 points. Each point is two 4-byte integers, 8 bytes. A layout of 10 million rectangles therefore needs at least $10^7 \cdot 5 \cdot 8 = 400$ million bytes of coordinate data alone, before record headers. Hierarchy is the main compression: a cell drawn once and referenced a million times costs one SREF per use, or a single AREF for a regular array, instead of a million copies of its polygons.

## OASIS and repetition

OASIS stores the same information far more compactly. It uses variable-length integers, relative coordinates, modal variables that let a record omit values repeated from the previous one, and explicit repetitions that place one element many times. An array of $r \times c$ identical placements can be one record with a repetition instead of $rc$ separate records, saving $rc - 1$ placements. Files are commonly several times to an order of magnitude smaller than GDSII. Both formats lose connectivity, which is why LVS must re-extract the circuit from geometry to check it.

## Going deeper

SEMI P39, the OASIS specification; the KLayout documentation on GDS2 and OASIS readers.
