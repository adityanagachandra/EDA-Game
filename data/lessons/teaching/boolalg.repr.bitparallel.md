# Bit-parallel truth tables

For functions of up to about 16 inputs, the fastest representation is the explicit truth table packed into machine words. A 6-input function is exactly one 64-bit word, and one AND instruction computes the AND of two functions on all 64 rows at once. ABC's cut enumeration, rewriting and NPN canonicalization, the truth-table libraries in mockturtle, and SAT-sweeping simulators all depend on this trick. A tool designer must be able to read, build and transform these words by hand.

## A function as one machine word

Store $f$ at minterm $i$ in bit $i$ of the word, bit 0 being the least significant. With minterms numbered over $a, b, c$ with $a$ as the most significant index bit, $f = ab + \overline{c}$ has ON-set $m(0, 2, 4, 6, 7)$, so bits 7, 6, 4, 2 and 0 are set: binary 11010101, which is 0xD5.

Read the word in hex from the high end: bits 7 to 4 give the high digit (1101 = D) and bits 3 to 0 the low digit (0101 = 5). The convention matters. Many libraries number variables the other way, with the first variable as the least significant index bit, and reading a table in the wrong convention silently permutes the variables.

## Variable masks and word operations

The projection word of a variable is the table of the function that equals that variable. With the convention above and three variables, $c$ (the lowest index bit) alternates every row, 0xAA; $b$ alternates in pairs, 0xCC; $a$ in fours, 0xF0. For four variables the words are $a = $ 0xFF00, $b = $ 0xF0F0, $c = $ 0xCCCC and $d = $ 0xAAAA.

Any expression is then evaluated with word operations. $ab + \overline{c}$ is (0xF0 AND 0xCC) OR (NOT 0xAA) = 0xC0 OR 0x55 = 0xD5, matching the table built row by row. NOT must be masked to the table width, since a 64-bit complement of an 8-bit table sets 56 unrelated bits.

## Cofactors by shifting

A variable whose projection word has stride $s$ (1 for the lowest index bit, then 2, 4, 8) splits the word into interleaved halves. The positive cofactor, as a full-width word that no longer depends on the variable, is obtained by keeping the half where the variable is 1 and copying it over the other half: $t_x = (t \wedge m) \vee ((t \wedge m) \gg s)$, where $m$ is the projection word. The negative cofactor shifts the other way.

For $f$ = 0xD5 and variable $a$ (mask 0xF0, stride 4): 0xD5 AND 0xF0 = 0xD0, shifted right by 4 is 0x0D, and the OR is 0xDD. Its low nibble D is $f_a = b + \overline{c}$ over $b, c$. A variable is in the support exactly when $t_x \ne t_{\overline{x}}$, so support, symmetry and unateness checks are a few shifts and compares.

## Going deeper

Knuth, "The Art of Computer Programming", volume 4A, section 7.1.3, on bitwise tricks. Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting", DAC (2006).
