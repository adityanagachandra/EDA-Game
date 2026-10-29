# Mask data preparation and fracturing

Between tape-out and the mask writer lies mask data preparation: deriving mask layers, applying OPC, and fracturing polygons into the simple figures a writer exposes. Data volume and write time grow with the number of figures, not with area. Post-OPC data is far larger than the drawn layout, which is why OASIS and efficient fracturing matter.

## Shot count of a rectangle

A variable-shaped-beam writer exposes one rectangle (or trapezoid) per shot, up to a maximum size $m$. A rectangle $W \times H$ split on a regular grid of maximal shots needs $\lceil W / m \rceil \cdot \lceil H / m \rceil$ shots. A 3200 by 1800 nm rectangle with $m = 1000$ nm needs $4 \cdot 2 = 8$ shots. Writing time is roughly the number of shots times the time per shot, plus settling, so fracturing tries to minimize the count and avoid slivers, very thin shots that print poorly.

## Choosing a fracture

A polygon can be cut into rectangles in different ways, and the shot counts differ. An L-shape, the union of $[0, X] \times [0, Y_1]$ and $[0, X_1] \times [Y_1, Y]$, can be cut horizontally at $y = Y_1$ or vertically at $x = X_1$. Each cut gives two rectangles, each split into shots, and the better cut is the one with fewer shots. The general problem of covering a rectilinear polygon with the fewest rectangles is solved by maximum matching on chords for partitions; with shot-size limits and sliver avoidance, fracturing tools use heuristics.

## Vertex growth after OPC

OPC splits each edge into fragments and moves each fragment independently to pre-compensate printing distortions. A rectangle whose sides of $W$ and $H$ are cut into $m_W = \lceil W/f \rceil$ and $m_H = \lceil H/f \rceil$ fragments of at most $f$ gains a jog of two vertices at every internal break when neighbouring offsets differ: $4 + 2(2(m_W - 1) + 2(m_H - 1))$ vertices. A 400 by 200 nm rectangle with $f = 50$ goes from 4 to $4 + 2(14 + 6) = 44$ vertices. This growth, repeated over billions of shapes, is why mask data sets reach terabytes.
