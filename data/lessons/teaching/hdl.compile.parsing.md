# Parsing and operator precedence

The parser turns the token stream into a tree according to the language grammar. For expressions the grammar encodes operator precedence and associativity, which decide what an unparenthesized expression means. Verilog has a long precedence table with several surprises, and every front end, from Icarus to slang to Yosys, must apply it identically or the same source will mean different hardware in different tools.

## Precedence levels

Precedence decides which operator binds tighter. In 4 + 2 * 3 the multiplication binds first, giving $4 + 6 = 10$, not 18. Verilog's binary operators, from tightest to loosest, are: ** ; then * / % ; then + - ; then the shifts << >> <<< >>> ; then the relational < <= > >= ; then the equalities == != === !== ; then & ; then ^ and ~^ ; then | ; then && ; then || ; and finally the conditional ?:. Unary operators bind tighter than all of them.

Fully parenthesizing an expression makes its tree explicit. 2 + 3 << 1 is (2 + 3) << 1 = 10, because + binds tighter than the shift, which surprises programmers used to other conventions.

## Bitwise operators below comparisons

The most dangerous surprise in Verilog, as in C, is that the bitwise operators bind more loosely than the comparisons. a & b == c parses as a & (b == c), a bitwise AND with a 1-bit comparison result, not as (a & b) == c. Likewise 7 & 7 == 2 * 1 is 7 & (7 == 2) = 7 & 0 = 0. And 2 < 9 << 2 == 1 parses as ((2 < (9 << 2)) == 1), which is 1.

Binary operators of equal precedence associate to the left, so 8 - 3 - 2 is (8 - 3) - 2 = 3. A precedence-climbing or Pratt parser implements all of this with one table of levels: it parses an operand, then repeatedly absorbs operators whose level is at least the current minimum, parsing each right operand at one level tighter.

## The conditional and logical operators

The conditional operator ?: has the lowest precedence and associates to the right. In 1 == 3 ? 7 & 7 : 8 | 7, the condition is (1 == 3), which is 0, so the value is the false arm (8 | 7) = 15; a ? b : c ? d : e means a ? b : (c ? d : e). The logical && binds tighter than ||, so a || b && c is a || (b && c).

Parse trees are also where front ends attach source locations, the line and column of each node, which later stages use for error messages and for the src attributes that link netlist cells back to RTL. A front end that loses them makes every downstream report harder to use.

## Going deeper

IEEE Std 1800-2017, clause 11.3.2 and table 11-2 (operator precedence); Aho, Lam, Sethi and Ullman, "Compilers", chapter 4.
