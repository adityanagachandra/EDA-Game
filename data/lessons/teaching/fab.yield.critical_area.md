# Critical area analysis

Not every defect that lands on a die kills it. A small extra particle between widely spaced wires is harmless; the same particle between two minimum-spaced wires is a short. Critical area is the region of the layout where the centre of a defect of a given size must fall to cause a failure. It turns a layout into a yield number and tells a DFM tool where spreading or widening wires pays off.

## Shorts between two parallel wires

Model a defect as a conducting disc of diameter $x$. Two parallel wires a distance $s$ apart, running together for length $L$, are shorted if the disc touches both. Along the direction across the gap, the disc centre must lie in a band of width $x - s$, and nothing happens if $x \le s$. Ignoring end effects, the critical area for shorts is $A_c(x) = L \max(0, x - s)$. For $L = 500$ um, $s = 30$ nm and $x = 70$ nm: $A_c = 500$ um $\times 0.040$ um $= 20$ um^2. A defect of 25 nm has zero critical area here, however many land.

## Opens in a single wire

A missing-material defect opens a wire of width $w$ only if it spans the whole width. Its centre must lie in a band of width $x - w$ along the wire, so $A_c^{open}(x) = L \max(0, x - w)$. For $w = 24$ nm, $L = 500$ um and $x = 70$ nm, that is 23 um^2. The two mechanisms pull in opposite directions: widening a wire reduces open critical area but, at fixed pitch, narrows the space and increases short critical area. A DFM optimizer balances them using the relative densities of extra-material and missing-material defects on each layer.

## Critical area of a layout

For a bus of $n$ parallel wires there are $n - 1$ gaps and $n$ wires, so the total critical area at size $x$ is $(n-1) L \max(0, x - s) + n L \max(0, x - w)$, ignoring defects that span three wires. For a real layout, tools compute critical area by geometry: grow every shape by $x/2$ and intersect the grown shapes of different nets; the overlap area is where a defect centre causes a short. Repeating for many defect sizes, with plane sweeps or spatial indices for speed, gives $A_c(x)$, which is then weighted by the defect size distribution to get the expected fault count.

## Going deeper

Stapper, "Modeling of integrated circuit defect sensitivities", IBM Journal of Research and Development (1983); Maly, "Computer-aided design for VLSI circuit manufacturability", Proceedings of the IEEE (1990).
