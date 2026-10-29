# Metal fill and density

Chemical-mechanical polishing removes material at a rate that depends on the local pattern density: sparse regions dish, dense regions erode. Foundries therefore require the metal density of every window to lie between a minimum and a maximum. Density is a union area clipped to a window, computed over many windows, so a density engine is a Boolean and area engine organized around windows.

## Density of one window

The density of a window is the area of metal inside it divided by the window area. Shapes must be clipped to the window (AND), and overlapping shapes merged (OR), or overlaps are counted twice. In a $70 \times 70$ um window, a shape $[62, 80] \times [-8, 16]$ contributes only its clipped part $[62, 70] \times [0, 16]$, 128 square um. If the merged, clipped shapes total 1276 square um, the density is $1276 / 4900 = 26.0$ percent. A union-area sweep computes the merged area in $O(n \log n)$ for $n$ shapes.

## Checking a tiling of windows

A density rule is checked over many windows. With a $2 \times 2$ tiling of an $80 \times 80$ um region into $40 \times 40$ windows and a 25 percent minimum, each window's density is computed from its own clipped shapes, and a shape that spans two windows contributes to both. Densities of 30.6, 6.9, 12.5 and 40.1 percent mean two windows fail. Shapes near window boundaries are where off-by-one clipping errors appear, so engines clip with half-open intervals and test against a brute-force reference.

## Density from routing pitch

Regular structures give density directly. Parallel wires of width $w$ at pitch $p$ fill a fraction $w / p$ of the area they cover. A $170 \times 170$ um window with a $60 \times 20$ um blockage and the rest routed at 0.28 um width and 0.84 um pitch has density $(28900 - 1200)(1/3) / 28900 = 31.9$ percent. Such estimates let a floorplanner predict density problems before routing: large blockages, analog areas and pad rings are the usual low-density regions, and wide power straps the usual high-density ones.

## Going deeper

Kahng, Robins, Singh and Zelikovsky, "Filling algorithms and analyses for layout density control" (IEEE TCAD, 1999).
