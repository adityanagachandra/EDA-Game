# Clock waveforms and create_clock

Every check in a timer is measured against clock edges, and those edges come from the clock definitions in SDC. create_clock gives a clock a name, a source, a period and a waveform, the times of its rising and falling edges within one period. Reading a waveform correctly, and knowing where each edge falls in later cycles, is the first step in deriving any launch and capture relationship.

## Period and duty cycle

The command create_clock -name clk -period 10 -waveform {0 5} [get_ports clk] defines a clock on port clk with period 10 time units, rising at 0 and falling at 5. The waveform lists edge times in increasing order within the first period, first a rise and then a fall; without -waveform the default is a rise at 0 and a fall at half the period. The duty cycle is the fraction of the period the clock is high, $(t_{fall} - t_{rise}) / P$. For -period 13 -waveform {6 8} the clock is high for 2 of every 13 units, about 15.4 percent. Units are those of the library or the session, usually ns.

## Edges in later cycles

The waveform repeats every period: rising edges occur at $t_{rise} + kP$ and falling edges at $t_{fall} + kP$ for integers $k$. Counting the rise at $t_{rise}$ as the first, rising edge number $n$ is at $t_{rise} + (n-1)P$. For -period 17 -waveform {4 16}, the fourth rising edge is at $4 + 3 \times 17 = 55$. A timer enumerates edges this way when it pairs launch and capture edges, especially between clocks with different periods, where the pattern repeats only after the least common multiple of the periods. The clock's source latency and network latency then shift each edge as seen at a particular flop.

## Finding the next edge after a time

Many questions reduce to finding the first edge of a given kind after some time $t_0$: the capture edge for a launch, the closing edge of a latch, or the edge a generated clock follows. For falling edges, take the smallest $k$ with $t_{fall} + kP > t_0$. For -period 14 -waveform {0 2} and $t_0 = 19$, the falling edges are at 2, 16, 30, so the answer is 30. Whether the edge must be strictly later or may coincide matters: a setup check pairs a launch with a strictly later capture edge, while a hold check uses the edge that coincides with or precedes it. Writing the inequality explicitly avoids an off-by-one-cycle error, the most common bug in hand edge arithmetic.

## Going deeper

The Synopsys Design Constraints format specification, create_clock; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 7.
