# Pattern matching for DFM

Pattern matching scans a layout for every occurrence of a known problematic pattern. Patterns come from foundry libraries, from earlier hotspots, or from failure analysis. On a rasterized clip, matching is a sliding-window comparison; production engines hash windows and use topological signatures so that a full chip is searched in one pass. Rotations, mirrors and don't-care regions make one pattern cover a family of layouts.

## Exact matching of a raster pattern

Rasterize the layer into pixels (1 for metal) and slide a $p \times q$ pattern over every position of an $R \times C$ grid: there are $(R - p + 1)(C - q + 1)$ positions, and a match requires every pixel to agree. A $2 \times 2$ pattern over a $5 \times 5$ grid has 16 positions. Production matchers avoid pixel-by-pixel comparison by hashing each window (rolling hashes update in constant time as the window slides), or by encoding a clip's edge topology and dimensions so that matches survive small shifts.

## Orientation-invariant matching

Lithography treats a pattern and its rotations alike (for symmetric illumination), so a library entry should match in all orientations. The matcher generates the rotated versions (90, 180, 270 degrees; mirrors add four more) and counts positions where at least one matches. A position is counted once even if several orientations match it. A $2 \times 2$ pattern 10, 10 matches twice unrotated in a grid but eight times with rotations. Hash-based matchers store the canonical form of each pattern, the smallest hash over its orientations, so one lookup covers all of them.

## Don't-care pixels

Fuzzy patterns mark some pixels as don't-care, matching either value, so one entry covers every layout that agrees on the constrained pixels. A $3 \times 3$ pattern with three don't-cares compares six pixels per position. More don't-cares mean more matches, both real hotspots and false alarms. Library designers trade coverage against precision; learned detectors can be seen as a soft version of the same idea, where the model decides which pixels matter.

## Going deeper

The ICCAD 2012 CAD contest on fuzzy pattern matching for physical verification.
