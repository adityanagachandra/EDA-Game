# Error-correcting codes for memories

Some memory errors appear after test: soft errors from particle strikes, retention loss, and wear. Error-correcting codes add check bits to each word so that single-bit errors are corrected and double-bit errors detected at every read. ECC complements test and repair: test removes hard defects at manufacture, ECC tolerates what happens in the field, and the ECC logic itself must be tested.

## Check bits for SEC-DED

A Hamming code with $r$ check bits protects $m$ data bits against any single error if $2^r \ge m + r + 1$: the syndrome must name each of the $m + r$ bit positions or "no error". For $m = 64$, $r = 7$ suffices ($2^7 = 128 \ge 72$). Adding one overall parity bit gives **SEC-DED** (single error correcting, double error detecting): 8 check bits for 64 data bits, the familiar (72, 64) code; (39, 32) for 32-bit words. The overhead fraction falls as words get wider, which is one reason memories use wide ECC words.

## Decoding the syndrome

In a Hamming(7, 4) code, check bits sit at positions 1, 2 and 4, and each covers the positions whose binary index contains it. With even parity, the **syndrome** of a received word is the XOR of the indices of all positions holding 1; for a valid codeword it is 0, and after a single flip it equals the flipped position. A received word 1000010 has 1s at positions 1 and 6, so the syndrome is $1 \oplus 6 = 7$: bit 7 flipped, and flipping it back corrects the word. In SEC-DED, the overall parity distinguishes a single error (parity wrong, syndrome nonzero) from a double error (parity right, syndrome nonzero).

## Uncorrectable words

If each bit of an $n$-bit codeword flips independently with probability $p$, the word is uncorrectable when two or more bits flip:
$$P_{\ge 2} = 1 - (1-p)^n - np(1-p)^{n-1}.$$
For $n = 72$ and $p = 6 \times 10^{-5}$, that is about 9.18 per million words. Errors accumulate in rarely read words, so controllers **scrub** memory periodically, reading and correcting words before a second error arrives. ECC also masks hard defects, so test must check memories with ECC bypassed or count corrected errors, or a repaired-looking memory may ship with no margin left.

## Going deeper

Lin and Costello, "Error Control Coding" (2nd ed., 2004); Hsiao, "A class of optimal minimum odd-weight-column SEC-DED codes," IBM Journal, 1970.
