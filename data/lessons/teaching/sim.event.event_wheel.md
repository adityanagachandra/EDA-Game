# The event wheel and time queue

An event-driven simulator spends much of its time inserting and removing events from a time-ordered queue. A general priority queue costs $O(\log n)$ per operation. Simulators exploit the fact that most delays are short and integer-valued: a **timing wheel** (a circular array of buckets indexed by time) makes scheduling and retrieval $O(1)$ for near-future events. Understanding the wheel explains simulator performance, the cost of very long delays, and why time advances in jumps.

## Slots indexed by time modulo the wheel size

A timing wheel is an array of $N$ buckets, each holding a list of events due at one time. An event due at absolute time $t$ goes into slot $t \bmod N$. The current time pointer walks around the wheel; when it reaches a slot, every event in that bucket is due now. Example: a wheel of 16 slots at current time 37. A gate with delay 5 schedules its output for $37 + 5 = 42$, and $42 \bmod 16 = 10$, so the event goes into slot 10. Insertion is a constant-time append to a list, with no comparison against other events. Within a slot, events are processed in arrival order, which is one reason the order of simultaneous events is a simulator choice rather than a language guarantee.

## The overflow list for far-future events

The wheel only represents times within one revolution, $[\text{base}, \text{base} + N)$. An event due later would collide with an earlier time that maps to the same slot. Such events go to an **overflow** structure, typically a sorted list or heap. Example: a 32-slot wheel covering times 64 through 95, with events scheduled at 70, 95, 96 and 120. The events at 70 and 95 go into the wheel; 96 and 120 go to overflow. As the base advances by a full revolution, overflow events that now fall inside the window move into the wheel. Designs with a few long delays (a clock period much longer than gate delays, a timeout) pay the slower overflow cost only for those events. Choosing $N$ trades memory against how often events spill.

## Advancing time in jumps

Simulated time does not tick through every value. After the current slot empties, the kernel looks for the next nonempty slot (or the earliest overflow event) and jumps there. If the queue holds events at times 5, 5, 15, 30 and 30 ns, the kernel visits exactly three times: 5, 15 and 30 ns, processing every event at each before moving on. Empty times cost nothing except the scan for the next occupied slot, which implementations accelerate with occupancy bitmaps. This is why event-driven simulation of a mostly idle design is fast, and why adding a fine-grained free-running signal (such as a very fast clock) slows everything down: it guarantees an occupied slot at every one of its edges.

## Going deeper

Brown, "Calendar queues: a fast O(1) priority queue implementation for the simulation event set problem", Communications of the ACM (1988); Varghese and Lauck, "Hashed and hierarchical timing wheels" (1987).
