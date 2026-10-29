# Half and full adders

Binary addition is column addition with carries, and the full adder is one column. Every adder architecture, from the ripple-carry adder to the Kogge-Stone tree, is a different way of computing the same carries, and every multiplier is a pile of full adders. The generate and propagate signals introduced here are the vocabulary of all fast adders.

## Carries column by column

Adding $a$ and $b$ bit by bit, column $i$ receives $a_i$, $b_i$ and a carry-in $c_i$ from column $i - 1$ (with $c_0$ the adder's carry-in), and produces a sum bit and a carry-out $c_{i+1}$. The carry-out is 1 when at least two of the three inputs are 1.

Add $a = 01101011$ and $b = 00110110$ with $c_0 = 0$. Column 0: $1 + 0 = 1$, no carry. Column 1: $1 + 1 = 10$, carry into column 2. Column 2: $0 + 1 + 1$ carries again, and so does column 3 ($1 + 0 + 1$). Column 4: $0 + 1 + 1$ carries; column 5: $1 + 1 + 1 = 11$ carries; column 6: $1 + 0 + 1$ carries; column 7: $0 + 0 + 1$ does not. Six columns produce a carry, and the sum is 10100001 ($107 + 54 = 161$).

## Generate, propagate and the full adder equations

A half adder adds two bits: $s = a \oplus b$, $c = ab$. A full adder adds three:

$$s_i = a_i \oplus b_i \oplus c_i, \qquad c_{i+1} = a_i b_i + c_i (a_i \oplus b_i).$$

Name the two pieces: generate $g_i = a_i b_i$ (the column makes a carry by itself) and propagate $p_i = a_i \oplus b_i$ (the column passes an incoming carry through). Then $c_{i+1} = g_i + p_i c_i$ and $s_i = p_i \oplus c_i$. In the example above, columns 1 and 5 generate ($a_i = b_i = 1$), and columns 0, 2, 3, 4 and 6 propagate. A full adder is two half adders plus an OR gate: the first half adder produces $g_i$ and $p_i$, the second adds $c_i$.

## Gate counts and the cost of a column

Built from two half adders and an OR, a full adder uses 2 XOR, 2 AND and 1 OR gate: 5 gates. An $n$-bit adder made of $n$ full adders therefore uses $5n$ gates; if the carry-in is known to be 0, the least significant column can be a half adder (2 gates), giving $5(n - 1) + 2$.

Real standard-cell libraries have dedicated full-adder cells (often a mirror adder in CMOS, about 28 transistors) that are smaller and faster than five separate gates, and synthesis tools map arithmetic onto them when it pays. More important for architecture is the structure of the carry path: in each column $c_{i+1}$ depends on $c_i$ through one AND-OR (or a majority gate), and how that dependence is chained across columns is what distinguishes the adders that follow.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 5.2; Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), section 11.2.
