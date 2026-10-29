# Wire segment lengths

FPGA tracks are built from wire segments that each span $L$ logic blocks before reaching a programmable switch. Short segments waste delay in switches on long connections; long segments waste track on short connections and pay quadratic RC delay. Choosing the segment-length mix is a core routing-architecture decision, and the router must know the segments to cost its search correctly.

## Counting segments

A track of a channel $n$ tiles long, built from length-$L$ segments laid end to end, holds $n/L$ segments when $n$ is a multiple of $L$. A channel of $W$ such tracks therefore holds $nW/L$ segments. With $n = 16$, $W = 24$ and $L = 4$, that is $16 \times 24 / 4 = 96$ segments. Architectures stagger the segment start points, so that in every tile about $W/L$ tracks begin a new segment; this gives every logic block access to fresh segment starts without changing the count. The number of segments matters because every segment is one node of the routing-resource graph and every segment end is a switch site: longer segments shrink the graph and the number of switches, which saves area and router memory.

## Delay of a connection on buffered segments

A connection of $d$ tiles on length-$L$ segments that starts at a segment boundary uses $\lceil d/L \rceil$ segments. Each buffered segment costs a switch-buffer delay $t_b$ plus the distributed RC of its wire, $\frac{1}{2} r c L^2$ for resistance $r$ and capacitance $c$ per tile (the Elmore delay of a uniform line). Total delay is $\lceil d/L \rceil \left(t_b + \frac{1}{2} r c L^2\right)$. With $d = 22$, $L = 4$, $t_b = 50$ ps, $r = 60$ ohm and $c = 30$ fF per tile, each segment costs $50 + 0.5 \times 60 \times 30 \times 16 \times 10^{-3} = 64.4$ ps, and six segments cost 386.4 ps. Buffers decouple segments, so delay adds linearly across them; an unbuffered pass-transistor chain would instead grow quadratically in the number of segments.

## Choosing the segment length

For a given distance, the best length balances the two terms. Short segments pay $t_b$ many times; long segments pay $\frac{1}{2}rcL^2$, which grows quadratically, and waste the unused part of the last segment when $d$ is not a multiple of $L$. For $d = 22$, $t_b = 70$ ps, $r = 210$ ohm and $c = 50$ fF, the delays are 1655.5 ps for $L = 1$, 1001 ps for $L = 2$, 924 ps for $L = 4$, 1218 ps for $L = 8$ and 2828 ps for $L = 16$, so $L = 4$ wins. Betz and Rose found medium lengths, around 4 to 8 tiles, to work best in their experiments, and commercial devices mix several lengths. The router then needs base costs per segment type, so that long wires are not used up by short connections.

## Going deeper

Betz and Rose, "FPGA Routing Architecture: Segmentation and Buffering to Optimize Speed and Density", FPGA (1999); Betz, Rose and Marquardt (1999), chapter 7.
