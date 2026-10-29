# FPGA routing architecture

Most of an FPGA's area and delay is in its programmable routing. In the island-style architecture, logic blocks sit in a grid surrounded by channels of $W$ tracks; connection boxes join block pins to tracks, and switch boxes join tracks where channels meet. The parameters $W$, $F_c$ and $F_s$ define the routing-resource graph that every FPGA router searches.

## Connection boxes and Fc

A connection box links a logic-block pin to some of the $W$ tracks of the adjacent channel. The flexibility $F_c$ is the fraction of tracks each pin can reach, given separately for inputs ($F_{c,in}$) and outputs ($F_{c,out}$). Each connection is one programmable switch. A block with $I$ inputs and $O$ outputs therefore uses $I \cdot F_{c,in} W + O \cdot F_{c,out} W$ connection-box switches. With $I = 34$, $O = 7$, $W = 80$, $F_{c,in} = 0.25$ and $F_{c,out} = 0.125$, that is $34 \times 20 + 7 \times 10 = 750$ switches for one block. A higher $F_c$ makes routing easier but costs area and adds capacitive load to every track, so architects choose the smallest $F_c$ that keeps the design routable at the target $W$; logically equivalent cluster inputs help here, because the router may enter through any of them.

## Switch boxes and Fs

Where a horizontal and a vertical channel cross, a switch box joins track ends. The flexibility $F_s$ is the number of other track ends each incoming end can connect to; $F_s = 3$ (one on each other side) is the usual choice. An interior switch box has $4W$ track ends, each with $F_s$ connections, and one switch serves two ends, so it contains $4W F_s / 2 = 2 F_s W$ switches: 48 for $W = 8$ and $F_s = 3$. A fabric with an $m \times m$ grid of such boxes has $2F_s W m^2$ of them. The switch pattern matters as much as the count. Disjoint (subset), Wilton and universal switch boxes all have $F_s = 3$ but connect different track indices, which changes how many routes exist between two pins.

## Track domains in disjoint switch boxes

In a disjoint switch box, track $i$ connects only to track $i$ on the other sides. A signal that enters the routing on track 3 can therefore only ever travel on track 3, and the fabric splits into $W$ independent domains. A route from an output pin to an input pin exists only on the track indices that both pins' connection boxes reach. If the source reaches tracks $\{2, 6, 7, 14\}$ and the sink $\{1, 2, 3, 7, 8, 11, 13\}$, only tracks 2 and 7 can carry the connection, whatever the rest of the fabric looks like. The Wilton switch box rotates track indices at each turn so that a signal can migrate between domains, which raises routability for the same $W$. This is why the switch pattern, not only $W$, sets the minimum channel width a router achieves.

## Going deeper

Betz, Rose and Marquardt (1999), chapter 4; Wilton, "Architectures and Algorithms for Field-Programmable Gate Arrays with Embedded Memories", PhD thesis, University of Toronto (1997); Lemieux and Lewis, "Design of Interconnection Networks for Programmable Logic" (2004).
