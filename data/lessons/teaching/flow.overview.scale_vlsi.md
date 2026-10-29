# Integrated circuits and VLSI scale

An integrated circuit puts transistors and their wiring on one piece of silicon. The names SSI, MSI, LSI, VLSI and ULSI (small, medium, large, very large and ultra large scale integration) count how many devices fit on a chip, and today's chips hold billions. Scale is the reason EDA exists: every tool must run in time roughly linear, or $n \log n$, in the size of the design, because anything quadratic is already too slow.

## Integration levels and exponential growth

The integration names grew from tens of gates (SSI) to thousands (LSI) to hundreds of thousands and beyond (VLSI). The cut-off numbers vary by author, so treat them as orders of magnitude, not definitions. Behind the growth is the observation usually called Moore's law: the number of transistors on a chip of given cost doubled roughly every two years for decades. If a chip holds $N_0$ transistors and the count doubles every $T$ years, then after $y$ years it holds $N_0 \cdot 2^{y/T}$. Starting from 2 billion with $T = 2$, eight years later the count is $2 \cdot 2^4 = 32$ billion.

## Why exhaustive methods fail

A block with $n$ input bits has $2^n$ input combinations. At 40 inputs that is $2^{40} \approx 1.1 \times 10^{12}$ vectors. A simulator that evaluates $10^9$ vectors per second needs about 1100 seconds. At 64 inputs the same simulator needs about 585 years. So EDA cannot enumerate. It uses structure instead: symbolic representations that share work (BDDs and AIGs), solvers that prune (SAT), and analyses that avoid enumerating input patterns altogether (static timing analysis computes worst cases over a graph rather than simulating vectors). The same blow-up applies to state: a design with 100 flip-flops has $2^{100}$ possible states.

## The design productivity gap

Transistor counts grew faster than the number of engineers. Suppose a skilled layout engineer draws 100 transistors per working day and a year has 250 working days. A 10-million-transistor design then needs $10^7 / (100 \cdot 250) = 400$ engineer-years of drawing alone, before any checking. Modern chips are hundreds of times larger. The gap between what the process can build and what a team can design by hand is closed only by raising the abstraction level (writing RTL instead of drawing transistors) and by reusing blocks. Every concept in this curriculum is a piece of that closing.

## Going deeper

Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits: A Design Perspective" (2nd ed., 2003), chapter 1.
