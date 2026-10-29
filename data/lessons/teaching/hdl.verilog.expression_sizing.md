# Expression bit widths under Verilog sizing rules

Every Verilog expression has a bit width determined by fixed rules, and many expressions are evaluated at a width set by their context rather than by their operands. These rules decide whether a carry is kept or lost and whether a comparison is done in 4 or 32 bits. A front end must implement them exactly, and a designer who misreads them writes adders that silently drop carries.

## Self-determined widths

The rules form a small table. Arithmetic and bitwise binary operators (+, -, *, /, %, &, |, ^) take the larger operand width: if a is 8 bits and b is 4, a + b is 8 bits. Comparisons and logical operators give 1 bit. A concatenation {a, b} sums its parts, and a replication {3{b}} multiplies. Unary minus and bitwise negation keep the operand width.

Apply the rules bottom up on the expression tree. With a 3 bits, c 3 bits and d 8 bits, (d | a) & (a - c) is max(max(8, 3), max(3, 3)) = 8 bits, and {d, a + c} is $8 + 3 = 11$ bits. Notice that + does not add a carry bit: the width of a + b is the larger operand width, not one more.

## Comparisons, reductions, shifts and conditionals

Four more rules complete the table. A comparison a < b or a == b returns 1 bit, although its operands are first extended to the larger of their two widths. A reduction &a, |a or ^a returns 1 bit. A shift a << n keeps the width of a; the shift amount n is self-determined and does not affect the result width. A conditional s ? x : y takes the larger of the widths of x and y; the condition s does not contribute.

So (a < b) + c has width max(1, width of c), and {2{a < b}} is 2 bits. A useful mental check: only the operands that flow into the result value can widen it. Selectors, shift amounts and conditions are evaluated on their own.

## Context-determined evaluation

Most operands are context-determined: before evaluation they are extended to the width of the whole expression, including the assignment target. In assign y = (a + b) >> 1; with a and b 8 bits and y 9 bits, the sum is computed at 9 bits, so the carry survives and the shift keeps it. With a = 176 and b = 220, a + b = 396 fits in 9 bits, and y = 198. If y were 8 bits, the sum would be computed at 8 bits, become $396 - 256 = 140$, and y would be 70.

Concatenation stops context propagation: its operands are self-determined. So assign y = {a + b} >> 1; computes the sum at 8 bits even when y is 9 bits, losing the carry. Front ends implement this as a two-pass algorithm: compute self-determined widths bottom up, then propagate the final width top down into context-determined operands.

## Going deeper

IEEE Std 1800-2017, clause 11.6 (expression bit lengths) and table 11-21; Sutherland and Mills, "Verilog and SystemVerilog Gotchas" (Springer, 2007).
