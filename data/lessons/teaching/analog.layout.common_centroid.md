# Common-centroid layout

Random mismatch shrinks with area, but systematic mismatch from gradients across the die does not. Oxide thickness, doping, stress and temperature vary smoothly with position, so two devices placed side by side see slightly different parameters. Common-centroid layout splits each device into unit cells arranged so that every device has the same centroid, which cancels any linear gradient. Checking and generating such patterns is a core task of analog layout tools.

## Centroids of split devices

Split device A into $n$ equal unit cells at positions $(x_i, y_i)$. Under a parameter that varies linearly with position, $p(x, y) = p_0 + g_x x + g_y y$, the device's average parameter is $p_0 + g_x\bar{x} + g_y\bar{y}$, where $(\bar{x}, \bar{y})$ is the mean position of its cells: its centroid. In a two-row array with row 1 "ABBA" and row 2 "BAAB", A's cells sit at columns 1, 4 (row 1) and 2, 3 (row 2), so its column centroid is $(1+4+2+3)/4 = 2.5$; its row centroid is 1.5. Computing a centroid is averaging coordinates, and it is the quantity that decides how a device responds to any linear gradient.

## Checking a pattern

A pattern is common-centroid when all matched devices share one centroid. Then any linear gradient shifts them equally and the mismatch cancels. For two devices compute both centroids and their distance; zero means common-centroid. "ABBA / BAAB" is: both centroids are (2.5, 1.5). "ABAB / ABAB" is not: A sits at columns 1 and 3, B at 2 and 4, a separation of one pitch, so a gradient of $g$ per pitch creates a mismatch of $g$. A sufficient construction is point symmetry: if rotating the array 180 degrees about its centre maps every A cell to an A cell, A's centroid is the centre, and likewise for B. Layout generators build patterns this way and verifiers check them by computing centroids, which is linear in the number of cells.

## What else matching needs

Common centroid cancels linear gradients only; it is one of several matching principles. Dispersion spreads each device's cells across the array so that higher-order gradients average out. Equal surroundings matter: a cell at the array edge sees different etching and stress than an interior cell, so dummy cells are added around the array. Devices should share orientation, since mobility depends on current direction relative to the crystal, and the routing to each device should be symmetric so that parasitic resistance does not reintroduce mismatch. Each principle costs area and wiring. A layout tool must trade them against each other, which is why matching constraints are expressed explicitly, for example "these four devices common-centroid with dummies", rather than inferred from wirelength.

## Going deeper

Hastings, "The Art of Analog Layout" (2nd ed., 2005), chapters on device matching.
