# Interface synthesis: handshakes, streams and buses

A synthesized kernel must exchange data with the rest of the system, and C function arguments say nothing about how. Interface synthesis maps each argument to a hardware protocol: a plain wire, a valid/ready handshake, a FIFO stream, a RAM port, or a memory-mapped bus master. The choice fixes bandwidth, sets II bounds and determines which verification problems can arise.

## Port-level protocols

The simplest interface is a plain port with no handshake: the block samples it whenever its schedule needs the value, so the environment must hold it stable. A valid signal adds timing: the producer raises valid when data is present. A full valid/ready handshake adds back-pressure: data transfers only on a clock edge where both valid and ready are high, and once a producer raises valid, it must keep the data and valid steady until the transfer happens. If the consumer drops ready for three cycles, the producer waits three cycles. Block-level control wraps the whole function call with start, done and idle signals, so a host can launch the accelerator and learn when it finishes. HLS tools generate these protocols from directives, and the protocol must match what the surrounding system actually does.

## Streams, memories and buses

An argument read once in order maps naturally to a stream, such as AXI4-Stream backed by a FIFO: no addresses, a handshake per item, and producer and consumer can overlap. An argument accessed randomly maps to a RAM interface with address, data and enable signals, limited to the RAM's ports. A large array in DRAM maps to a memory-mapped bus master, such as an AXI4 master port, which issues read and write transactions. Each transaction pays address and latency overhead, so the tool tries to infer bursts: one address phase followed by many data beats for consecutive addresses. Scalars and control usually go through a simple register interface such as AXI4-Lite, which a host processor writes and reads. The choice is made per argument and recorded in the generated RTL's ports.

## How interfaces bound performance

Interfaces are resources, and they bound the II. Two reads per iteration from a single-port RAM interface force II = 2. A DRAM read with 60 cycles of latency, issued one at a time with no bursts and no outstanding requests, limits a loop to roughly one element per 60 cycles, whatever its datapath II; bursts and multiple outstanding transactions hide that latency. Widening a bus helps only if each beat carries useful data: a 512-bit bus delivering 32-bit elements one per beat gives no gain, while packing 16 elements per beat gives up to 16 times the bandwidth. Streams add a different hazard: blocking FIFOs make depth part of correctness, since a producer that must write four items to a depth-2 FIFO before the consumer starts reading it will deadlock.

## Going deeper

The interface chapter of the AMD Vitis HLS user guide (UG1399); the AMBA AXI protocol specification from Arm; Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018).
