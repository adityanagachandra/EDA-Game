# Routing-resource graphs

The routing-resource (RR) graph is the router's model of the FPGA: one node per wire segment and per pin, one directed edge per programmable switch. Every architecture detail, from segment lengths to switch patterns, is compiled into it, and the router never sees the architecture otherwise. Its size determines router memory and run time, so counting its nodes and edges is a basic design check.

## Wire nodes

In VPR's island-style model, each wire segment is one node of type CHANX (horizontal) or CHANY (vertical). For an $n \times n$ array of logic blocks with a channel on every side of every block, there are $n + 1$ horizontal channels and $n + 1$ vertical channels, each $n$ tiles long. With $W$ tracks of length-$L$ wires per channel, each track holds $n/L$ wires, so the graph has $2(n + 1) W n / L$ wire nodes. For $n = 10$, $W = 4$ and $L = 2$: $2 \times 11 \times 4 \times 5 = 440$. Each node carries a capacity (1 for a wire), a base cost and the electrical data the timing-driven router needs, such as resistance and capacitance.

## Pin, source and sink nodes

Each logic block adds an IPIN node per input pin, an OPIN node per output pin, a SOURCE node from which its output nets start, and a SINK node where its input connections end. Logically equivalent pins, such as the inputs of a cluster whose local crossbar can deliver any input to any LUT, share one SINK with capacity equal to the number of pins, so the router may end a connection at any of them. With 18 inputs and 9 outputs, each block contributes $18 + 9 + 2 = 29$ nodes; 64 blocks add 1856 to 1008 wire nodes, 2864 in all. Equivalence matters: it gives the router a choice of pins, which helps routability, and it is part of the architecture description, not something the router can infer.

## Edges

Every programmable switch becomes a directed edge. Connection boxes give each IPIN an edge from each of the $F_{c,in} W$ wires it can reach, and each OPIN an edge to each of its $F_{c,out} W$ wires; switch boxes give wire-to-wire edges, two per bidirectional switch. For a $10 \times 10$ array of blocks with 12 inputs and 10 outputs at $W = 24$, $F_{c,in} = 0.5$ and $F_{c,out} = 0.25$, the pin edges alone number $100 \times (12 \times 12 + 10 \times 6) = 20{,}400$. Edges usually outnumber nodes several times, so RR graphs for large devices take gigabytes; VPR and nextpnr reduce memory by sharing edge data per switch type and by building the graph lazily or in compact arrays.

## Going deeper

Betz, Rose and Marquardt (1999), chapter 6; the VTR documentation for the RR graph file format; Murray et al., "VTR 8: High-Performance CAD and Customizable FPGA Architecture Modelling", ACM TRETS (2020).
