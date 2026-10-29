# Multi-bit crossings, Gray codes and async FIFOs

A two-flop synchronizer makes one bit safe. A multi-bit value is not made safe by synchronizing every bit, because the bits resolve independently and can be captured from different cycles. Gray codes, handshakes and asynchronous FIFOs are the standard answers. Their arithmetic (how many bits change, what code to use, how deep a FIFO must be) is a routine part of CDC design and verification.

## Why multi-bit binary crossings fail

Suppose a 4-bit binary counter crosses domains with a synchronizer on each bit, and it increments from 7 (0111) to 8 (1000). All four bits change at once. The destination samples each bit near the transition, and each synchronizer may capture the old or the new value independently, so the destination can see any of 16 values, such as 1111 or 0000, neither of which the source ever held. In general, going from $k$ to $k + 1$ changes the trailing ones of $k$ plus one more bit: $\text{popcount}(k \oplus (k+1))$ bits. Only when exactly one bit changes per transition is every possible sample either the old value or the new one. That is the property a Gray code provides.

## Gray codes

The reflected binary Gray code of $b$ is $g = b \oplus (b \gg 1)$. Consecutive integers have Gray codes that differ in exactly one bit, including the wrap from the maximum back to 0. Example: $b = 13 = 1101_2$, $b \gg 1 = 0110_2$, so $g = 1011_2 = 11$. And $14 = 1110_2$ gives $1110 \oplus 0111 = 1001_2 = 9$; 1011 and 1001 differ in one bit. Decoding back to binary takes a prefix XOR: each binary bit is the XOR of all Gray bits at or above it. Gray coding only works for values that change by one step at a time, such as counters and FIFO pointers; arbitrary data words need a handshake that holds the data stable while a single synchronized control bit announces it.

## Asynchronous FIFOs and their depth

An asynchronous FIFO moves a stream of data words between domains. Its write and read pointers are kept in binary for addressing, converted to Gray code, and synchronized into the opposite domain, where they are compared to produce full and empty flags. The data itself sits in a dual-port memory and needs no synchronization, because the flags guarantee it is stable when read. Sizing: a burst of $B$ words written at $f_w$ while the reader drains at $f_r < f_w$ leaves a backlog of $B(1 - f_r/f_w)$ words at the end of the burst. For $B = 80$, $f_w = 200$ MHz and $f_r = 125$ MHz, that is $80 \times 0.375 = 30$ words, so at least 30 entries, plus margin for the few cycles of synchronizer latency on the pointers, which makes the flags conservative.

## Going deeper

Cummings, "Simulation and Synthesis Techniques for Asynchronous FIFO Design" (SNUG 2002); Ginosar, "Metastability and synchronizers: a tutorial", IEEE Design and Test of Computers (2011).
