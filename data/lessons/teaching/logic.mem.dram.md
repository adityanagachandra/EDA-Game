# DRAM refresh, bandwidth and latency

Dynamic RAM stores each bit as charge on a capacitor behind one transistor, which makes it far denser than SRAM but leaky and slow to access. DRAM lives off-chip (or stacked, as HBM), and a memory controller on the chip schedules its commands. The controller and PHY are standard blocks in a modern SoC, and their performance arithmetic (refresh overhead, peak bandwidth, row-buffer latency) drives system architecture.

## The 1T1C cell and refresh

A DRAM cell is one access transistor and one capacitor. Reading shares the cell's charge with a long bitline, and the sense amplifier detects a small voltage change, so a read is destructive and the row must be written back. The stored charge also leaks away, so every row must be refreshed within a retention window, conventionally 64 ms at normal temperature for DDR devices.

If the controller must issue $R$ refresh commands evenly over the window, the interval between commands is $t_{REFI} = 64\ \text{ms} / R$. With $R = 8192$, $t_{REFI} = 7.8125\ \mu\text{s}$. During each refresh the device is busy for $t_{RFC}$, so refresh costs a fraction $t_{RFC}/t_{REFI}$ of the time; with $t_{RFC} = 350$ ns that is about 4.5%. Denser devices have longer $t_{RFC}$, and refresh overhead grows with capacity.

## Peak bandwidth

A DDR interface transfers data on both clock edges. Its peak bandwidth is the bus width times the transfer rate:

$$BW = \frac{\text{width (bits)} \times \text{rate (MT/s)}}{8}\ \text{MB/s}.$$

A 64-bit DDR4-3200 channel (3200 mega-transfers per second) peaks at $64 \cdot 3200 / 8 = 25600$ MB/s, 25.6 GB/s. Sustained bandwidth is lower because of refresh, row changes, bus turnarounds between reads and writes, and bank conflicts; controller scheduling (for example open-row-first policies) exists to recover as much of the peak as possible.

## Row hits, misses and conflicts

Each bank has a row buffer holding one open row. An access to the open row is a row hit and costs the column latency $t_{CL}$. An access to a bank with no open row needs an activate first: $t_{RCD} + t_{CL}$. An access to a different row than the one open (a conflict) must precharge first: $t_{RP} + t_{RCD} + t_{CL}$.

With $t_{CL} = t_{RCD} = t_{RP} = 14$ ns, a hit costs 14 ns and a conflict 42 ns. If 70% of accesses hit and the rest conflict, the expected latency is $0.7 \cdot 14 + 0.3 \cdot 42 = 22.4$ ns. This expected-value calculation is the core of memory-controller policy: keeping rows open helps workloads with locality and hurts random ones, and address mapping (which address bits select bank, row and column) is chosen to spread conflicts across banks.

## Going deeper

Jacob, Ng and Wang, "Memory Systems: Cache, DRAM, Disk" (2008); JEDEC Standard JESD79-4 (DDR4 SDRAM).
