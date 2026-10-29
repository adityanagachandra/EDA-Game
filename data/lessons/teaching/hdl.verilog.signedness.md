# Signed and unsigned expression evaluation

Whether a Verilog expression is signed decides whether its operands are sign-extended or zero-extended, and the rule is global: a single unsigned operand makes the whole expression unsigned. That rule surprises even experienced designers and produces wrong arithmetic that simulates and synthesizes consistently, so only careful reading finds it. A front end must apply it exactly, together with the sizing rules.

## Signed operands are sign-extended

Declaring a variable signed, as in logic signed [3:0] a, makes its bit pattern a two's-complement number: 4'b1110 is -2. When every operand of an expression is signed, the expression is signed, and each operand is sign-extended to the context width before the operation.

With a = 4'b1110 (-2), c = 4'b1100 (-4) and an 8-bit target, assign y = a + c; extends both to 8 bits (11111110 and 11111100) and adds them: -6, which is 250 when y is read as unsigned. The result in the low 4 bits alone would look the same either way; the difference appears only in the extended bits, which is exactly where width-mixing bugs live.

## One unsigned operand makes everything unsigned

The signedness of an expression depends only on its operands: if any operand is unsigned, the whole expression is unsigned, and all operands are zero-extended, including the ones declared signed. With a signed a = 4'b1010 and an unsigned c = 4'b1000, a + c becomes $10 + 8 = 18$ in 8 bits, not $-6 + 8 = 2$.

The signedness of the left-hand side plays no role. Declaring y signed does not make a + c signed. To get signed arithmetic with mixed operands, convert explicitly with the signed system function or a cast, so that every operand is signed before the expression is evaluated.

## Part-selects and literals

Two rules catch most remaining bugs. A part-select is always unsigned, even when it covers the whole of a signed vector: a[3:0] + c with both declared signed is unsigned, so a = 4'b1111 and c = 4'b1011 give $15 + 11 = 26$ rather than $-1 + -5$.

Literals carry signedness too. A sized based literal such as 1'b1 or 4'd3 is unsigned, so a + 1'b1 is unsigned and zero-extends a. An unsized decimal such as 1 is a signed 32-bit integer, so a + 1 is signed, and its context width is 32 bits; a = 4'b1000 (-8) gives -7, which is 249 in an 8-bit target. Writing 4'sd1 gives a signed sized literal. Linters such as Verilator flag many of these mixes, and a robust front end reports them too.

## Going deeper

IEEE Std 1800-2017, clause 11.8 (expression evaluation rules); Sutherland and Mills, "Verilog and SystemVerilog Gotchas", gotchas on signed arithmetic.
