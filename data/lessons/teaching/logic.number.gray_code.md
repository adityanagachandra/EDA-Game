# Gray code

A Gray code orders all $2^n$ bit patterns so that consecutive ones differ in exactly one bit. That single property makes it the standard encoding for values that cross between clock domains (FIFO pointers), for mechanical and optical encoders, and for low-power state machines and counters, where every bit that toggles costs switching energy. Converting between binary and Gray is a few XOR gates, and knowing which way the XOR chain runs matters.

## Binary to Gray

The reflected binary Gray code of a binary number $b$ is

$$g = b \oplus (b \gg 1),$$

so each Gray bit is $g_i = b_i \oplus b_{i+1}$ and the top bit is copied. This is one layer of XOR gates with no carry chain, so it is fast.

Convert $b = 1011_2 = 11$: $b \gg 1 = 0101$, and $1011 \oplus 0101 = 1110$. So the Gray code word for 11 is $1110_2$, which read as an unsigned integer is 14. The Gray word for 12 is $1100 \oplus 0110 = 1010$, which differs from 1110 in one bit, as promised.

## Gray back to binary

Going back needs a prefix XOR from the top: $b_{n-1} = g_{n-1}$, then $b_i = b_{i+1} \oplus g_i$. Equivalently $b_i$ is the XOR of all Gray bits at positions $\ge i$.

Decode $g = 1110$: $b_3 = 1$, $b_2 = 1 \oplus 1 = 0$, $b_1 = 0 \oplus 1 = 1$, $b_0 = 1 \oplus 0 = 1$, so $b = 1011 = 11$. Unlike the forward direction, this is a serial chain of depth $n - 1$ (or a prefix tree of depth $\lceil \log_2 n \rceil$), which is why an asynchronous FIFO keeps its pointers in Gray form for the crossing and converts back to binary only where it must compute a difference.

## Why one bit at a time matters

Counting from $a$ to $b$ in plain binary toggles $\sum_{i=a}^{b-1} \text{popcount}(i \oplus (i+1))$ bits. Step 7 to 8 in 4 bits toggles all four. In Gray code every step toggles exactly one bit, so the same count toggles $b - a$ bits. Over a full $n$-bit binary count the average is almost 2 toggles per step, against 1 for Gray.

The deeper reason is correctness, not power. If a multi-bit value is sampled by another clock while it changes, each bit may be caught old or new independently. With binary 0111 to 1000, the receiver could see any of 16 patterns. With Gray, only one bit is in flight, so the receiver sees either the old value or the new one, never a third value. That guarantee is what makes Gray-coded pointers safe to synchronize bit by bit.

## Going deeper

Knuth, "The Art of Computer Programming", volume 4A, section 7.2.1.1; Cummings, "Simulation and Synthesis Techniques for Asynchronous FIFO Design" (SNUG 2002).
