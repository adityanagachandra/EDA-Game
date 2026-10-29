# Layout XOR between versions

An XOR compares two versions of a layout: it reports the area covered by exactly one of them. An empty XOR proves two layouts are geometrically identical; a nonzero XOR shows exactly what changed. Tape-out flows XOR the final stream against the signed-off database, and engine developers XOR old and new outputs as a differential test.

## XOR of two shapes

For shapes $A$ and $B$, the XOR area is $|A| + |B| - 2|A \cap B|$. Two rectangles of 120 and 90 square um overlapping in 60 square um give $120 + 90 - 120 = 90$ square um. XOR is a Boolean operation, computed by the same sweep that computes AND and OR. On a layer that should be unchanged, any nonzero XOR stops the flow until someone explains it; on a layer that was edited, the XOR should match the intended edit.

## XOR of two versions of a layer

Real comparisons involve many shapes per version, possibly overlapping within a version. Each version is first merged into its union; the XOR is the area covered by exactly one union. A shape moved by a few units shows up as two thin slivers, not as a whole shape, and a shape added in one version appears in full. XOR engines tile large layouts and run in parallel, and hierarchical XOR compares cells once when both versions share them.

## The signature of a shift

When a $w \times h$ rectangle shifts by $(d_x, d_y)$ with $|d_x| < w$ and $|d_y| < h$, the two copies overlap in a $(w - |d_x|) \times (h - |d_y|)$ box, so the XOR is $2(wh - (w - |d_x|)(h - |d_y|))$. A 100 by 50 rectangle shifted by (3, 0) gives $2(5000 - 4850) = 300$. A whole cell shifted by the same vector produces this signature on every shape, which a reviewer recognizes as a placement change rather than an edit. Tolerances can suppress differences thinner than a threshold, at the risk of hiding real changes.
