# Partial reconfiguration

Partial reconfiguration rewrites the configuration of one region while the rest of the device keeps running. A static region holds the interfaces; reconfigurable regions swap between modules. The flow must floorplan regions on frame boundaries, freeze or decouple their interfaces during a swap, and generate partial bitstreams, and the size of those bitstreams sets how fast the device can switch functions.

## Partial bitstream size

Because frames are the unit of configuration, a reconfigurable region is aligned to frame boundaries, typically whole columns within a clock region. The partial bitstream contains every frame of the region: if the region spans $c$ columns over $r$ rows, each column-row needs $F$ frames, and a frame has $w$ 32-bit words, the payload is $c \cdot r \cdot F \cdot w \cdot 4$ bytes. With $c = 5$, $r = 2$, $F = 28$ and $w = 101$, that is 113,120 bytes, about 110.5 KB. The partial bitstream scales with the region, not with the device, which is the point: a small region swaps quickly. Headers and commands add a little; frames that hold hard blocks or memory contents can add much more.

## Reconfiguration time

An internal configuration port lets the design itself load a partial bitstream, typically one 32-bit word per cycle. At 125 MHz the port moves $4 \times 125 \times 10^6 = 500$ MB/s, so a 428,544-byte partial bitstream takes $428{,}544 / 500 = 857.1$ microseconds. Real throughput can be lower when the bitstream is fetched from external memory or when the port runs below its maximum clock. During the swap, the region's outputs are undefined, so the static region must decouple them, holding handshakes idle or isolating the interface, and must reset the new module once it is loaded. These are the steps that vendor flows such as AMD's dynamic function exchange automate around the bitstream itself.

## When it pays

Partial reconfiguration time-multiplexes hardware. If each swap costs $t_r$ and the loaded accelerator then runs for $t_t$, the region computes for a fraction $t_t / (t_t + t_r)$ of the time. With $t_r = 8.5$ ms and $t_t = 30$ ms, that is 77.9%; with $t_t = 1$ ms it would be 10.5%. Swapping pays when tasks run long compared with reconfiguration, or when the alternative is a larger, more expensive device. It costs floorplanning freedom: the region's resources are reserved, the static logic cannot use them, and the placer and router must keep static nets out of the region except for its declared partition pins. Every module must fit the region's resource mix, so regions are sized for the largest module.

## Going deeper

Vipin and Fahmy, "FPGA Dynamic and Partial Reconfiguration: A Survey of Architectures, Methods, and Applications", ACM Computing Surveys (2018); AMD's Vivado user guide on dynamic function exchange (UG909).
