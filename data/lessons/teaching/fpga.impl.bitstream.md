# Configuration bits and bitstreams

The bitstream is the file that configures an FPGA: every LUT bit, every multiplexer select, every routing switch and every hard-block mode, organized into frames that the configuration logic writes. Bitstream generation is the last step of the flow and the only one whose output the hardware sees, so its size, structure and load time are part of the architecture.

## Counting configuration bits

Each tile contributes the bits of its logic and its routing. A tile with one cluster of $N$ BLEs of $k$-LUTs holds $N(2^k + 1)$ logic bits, and its share of connection and switch boxes adds one bit per switch. With $N = 8$, $k = 5$ and 910 routing bits, a tile holds $8 \times 33 + 910 = 1174$ bits, and a $35 \times 35$ array holds $1225 \times 1174 = 1{,}438{,}150$ bits, before I/O, clocking and hard blocks. Routing bits usually outnumber logic bits several times, which mirrors the area split: most of an FPGA is programmable interconnect. The bitstream is therefore mostly routing, and most of it is zeros for a lightly used device.

## Frames and load time

Configuration memory is written in frames, the smallest addressable unit, each a fixed number of bits that usually spans a column of tiles. A full configuration of $B$ bits in frames of $f$ bits needs $\lceil B/f \rceil$ frames. Load time is set by the configuration port: $B$ bits over a $w$-bit port at clock $F$ take $B / (wF)$ seconds. A 29-megabit device loaded over a 32-bit port at 100 MHz takes about 9 ms; over a 1-bit serial interface at 50 MHz, about 0.6 s. Frames also define what can be rewritten independently, which is the basis of partial reconfiguration and of readback for scrubbing radiation-induced upsets.

## Encoding multiplexer selects

Routing in modern FPGAs is built from multiplexers. A mux with $m$ inputs can store its select one-hot, one bit per input, or binary-encoded in $\lceil \log_2 m \rceil$ bits with a decoder. With 2050 muxes of 33 inputs, one-hot needs $2050 \times 33 = 67{,}650$ bits and encoded storage $2050 \times 6 = 12{,}300$. Fewer bits mean a smaller configuration memory and a shorter bitstream, while one-hot drives pass transistors directly with no decoder delay or area; real devices use mixtures, such as two-level muxes with one-hot selection at each level. The encoding is invisible to the router but defines how the bitstream writer translates a routed switch into bits, which is exactly what a bitstream database records.

## Going deeper

Kuon, Tessier and Rose (2008), sections on configuration; Lemieux and Lewis (2004) on multiplexer-based routing; the Project X-Ray documentation on frames and configuration.
