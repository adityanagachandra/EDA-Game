# Database units and the manufacturing grid

Layout databases store integers, not lengths. The stream's UNITS record says how long one database unit (DBU) is, and every coordinate counts those units. Integer coordinates make Boolean operations and spacing checks exact and reproducible, which a signoff tool needs. Two limits follow: the finite integer width caps the coordinate range, and the foundry's manufacturing grid restricts which integers are legal vertex positions.

## Database units

GDSII's UNITS record holds two 8-byte reals: the size of a database unit in user units (often 0.001, meaning micrometres) and in metres (often $10^{-9}$, meaning 1 nm). A coordinate in micrometres becomes an integer by multiplying by DBUs per micrometre. With a 0.5 nm unit there are 2000 DBU per um, so $x = 213.43$ um is stored as 426860. Converting a design to a finer unit multiplies every coordinate exactly; converting to a coarser unit needs rounding and can move edges, which is why mixing IP with different units is a classic source of off-grid errors.

## The coordinate range

GDSII coordinates are signed 32-bit integers, so the largest value is $2^{31} - 1 = 2147483647$ DBU. With a 1 nm unit that is about 2147.48 mm, far larger than any die. With a 0.1 nm unit it is 214.75 mm, still enough for a wafer-scale layout but no longer generous. Inside a tool, products of coordinates (areas, cross products in orientation tests) need twice the bits, so engines either use 64-bit intermediates or restrict the range. Choosing the unit is a technology decision made once, because it fixes both precision and reach.

## The manufacturing grid

Foundries require every vertex to lie on a manufacturing grid, for example 5 nm, so that mask writers and OPC work on a common lattice. A vertex is legal only if both coordinates are multiples of the grid; being an integer DBU is not enough. With a 10 nm grid, $(3050, 190)$ passes, while $(163, 2120)$ and $(1310, 1053)$ fail. Off-grid vertices usually come from scaling, rotation by non-right angles, mirroring about a half-unit line, or Boolean results of diagonal edges. A grid check is cheap (a divisibility test per vertex) and runs early, because off-grid geometry corrupts the precise checks that follow.
