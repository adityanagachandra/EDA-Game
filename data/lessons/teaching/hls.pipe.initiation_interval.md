# Initiation interval

The initiation interval (II) of a pipelined loop or function is the number of cycles between the starts of successive iterations. It sets throughput directly, independent of how long one iteration takes. Every pipelining directive in an HLS tool is a request for a target II, and every report states the II achieved and what prevented a lower one, so II is the number an HLS user watches most.

## II sets the throughput

In a pipelined loop a new iteration starts every II cycles while earlier iterations are still in progress. At clock frequency $f$, one iteration per sample, the steady-state throughput is $f / \mathrm{II}$ samples per second. A loop at 100 MHz with II = 8 delivers 12.5 million samples per second; at II = 1 it would deliver 100 million. The iteration latency, the cycles from an iteration's start to its end, does not appear: a 30-cycle iteration started every cycle still yields one result per cycle once the pipeline is full. That is the essential difference between latency and throughput. II = 1 is the ideal for a single loop, and a lower II needs either more hardware or a shorter recurrence; directives that ask for an II that the loop cannot reach produce a report explaining the bound.

## The II a rate requires

A throughput requirement turns into a maximum II. A filter that must process $r$ samples per second at clock $f$ may spend at most $f / r$ cycles per sample, and since II is a whole number, $\mathrm{II} \le \lfloor f / r \rfloor$. At 100 MHz and 30 million samples per second, the bound is $\lfloor 3.33 \rfloor = 3$: II = 3 gives 33.3 million samples per second and meets the rate, II = 4 gives 25 million and fails. The largest acceptable II is the cheapest design that meets the requirement, because a larger II lets each unit be reused by more operations of the same iteration. So a designer should not ask for II = 1 by default when the rate allows II = 3; the tool can then share units three ways and save area.

## Forbidden latencies on a shared resource

When one iteration uses a single resource several times, overlap with later iterations can collide. If an iteration uses the multiplier at offsets 0, 3 and 7 cycles from its start, iteration $k$ uses it at $k\,\mathrm{II} + \{0, 3, 7\}$. Two uses collide when two offsets are congruent modulo II, so II must give all offsets distinct residues. II = 3 fails ($0 \equiv 3$), II = 4 fails ($3 \equiv 7$), and II = 5 gives residues 0, 3, 2, which works. This is the reservation-table view from classical pipeline theory: differences between offsets are forbidden latencies, and an II is legal when no forbidden latency is a multiple of it. Modulo schedulers enforce the same rule with a modulo reservation table, and they can also move the uses themselves to make a lower II legal.

## Going deeper

Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapter on pipelining; Kogge, "The Architecture of Pipelined Computers" (1981), on reservation tables and forbidden latencies.
