# Sheet resistance and counting squares

Every resistance an extractor reports starts from the same idea: a metal layer of fixed thickness has a sheet resistance, and a wire's resistance is that number times how many squares of its own width fit along it. The extractor never needs absolute dimensions for this, only ratios, which is why process kits publish sheet resistance per layer. Getting the square count right on straight runs, width changes and corners is the first thing a resistance extractor does after it has fractured a net into rectangles.

## Resistance per square

A wire of resistivity $\rho$, thickness $T$, width $W$ and length $L$ has $R = \rho L / (W T)$. Group the terms as $R = (\rho / T)(L / W)$. The first factor, $R_s = \rho/T$, depends only on the layer and is called the sheet resistance, in ohms per square. The second, $L/W$, is the number of squares: how many $W \times W$ tiles fit along the wire. A square of any size has the same resistance, because doubling its side doubles the path length and the cross-section together. A wire 2.4 um long and 0.02 um wide holds $2.4/0.02 = 120$ squares. Counting squares is pure geometry, so an extractor can do it on the fractured rectangles before it knows anything about temperature or corner.

## From resistivity to ohms

When a process gives resistivity and thickness instead of sheet resistance, form $R_s = \rho/T$ first. With an effective resistivity of 0.04 ohm-um and a thickness of 0.08 um, $R_s = 0.5$ ohm/sq. A wire 32 um long and 0.2 um wide has 160 squares, so $R = 0.5 \times 160 = 80$ ohm. Series segments add: a net that runs 100 squares on one layer and 40 on another of different sheet resistance sums the two products. A width change along the wire simply splits it into two rectangles with their own square counts. The word effective matters: at advanced nodes the resistivity a kit quotes already includes the size effects that raise copper resistivity above its bulk value, so it is a per-layer, sometimes per-width, number rather than a constant of nature.

## Corners and current crowding

At a right-angle bend the current does not use the corner square uniformly. It crowds around the inner corner and leaves the outer corner nearly idle, so the corner square conducts better than one full square. The standard approximation counts a corner square of an L-shaped wire as about 0.56 squares. For an L of width 0.02 um whose arms hold 190 and 20 squares outside the corner, with $R_s = 0.5$ ohm/sq, the resistance is $0.5 \times (190 + 20 + 0.56) = 105.28$ ohm. Treating the corner as a full square, or as zero, misses by a few tenths of a square per bend, which is small on long wires but noticeable on short, jogged local nets with many bends. Extractors handle bends, T-junctions and contacts with precomputed correction tables for exactly this reason.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed., 2011), the chapter on interconnect; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed., 2003), chapter 4.
