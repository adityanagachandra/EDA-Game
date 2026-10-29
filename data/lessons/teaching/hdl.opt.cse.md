# Common subexpression elimination

Common subexpression elimination (CSE) finds computations that produce a value already computed and reuses the earlier result. In software it saves instructions; in hardware it saves gates, because two identical cells with identical inputs can be merged into one. Value numbering is the standard way to find them, and Yosys's opt_merge and the structural hashing in AIG packages are the same idea on netlists.

## Local value numbering

Value numbering assigns a number to every distinct value. Inputs get fresh numbers. For each operation, form a key from its operator and its operands' value numbers, and look it up in a table. If the key is new, the result gets a fresh number and the key is recorded. If the key is present, the operation is redundant: the result is a copy of the earlier one and receives its number.

In t1 = b * c; t2 = c & t1; t3 = b - t2; t4 = t2 & b; t5 = b - t2;, the key for t5 (minus, number of b, number of t2) matches t3's, so t5 is replaced by a copy of t3: one redundant operation. Because a copy carries the original's value number, any later operation using t5 matches operations using t3, so savings cascade.

## Commutativity and canonical keys

a + b and b + a compute the same value, but their keys differ unless the compiler puts commutative operators into a canonical order, for example by sorting operand value numbers. With canonical keys, t7 = t1 & c matches an earlier t4 = c & t1. Subtraction is not commutative, so c - a and a - c keep different keys.

This is also where hardware compilers differ in strength. Structural hashing in an AIG canonicalizes the inputs of every AND node, so commuted duplicates merge automatically. Yosys's opt_merge compares cells by type, parameters and connected input signals, and handles commutative ports of cells such as add and and by normalizing their order. Deeper equivalences (a + b + c versus c + (a + b)) need algebraic rewriting or SAT-based sweeping, beyond what value numbering sees.

## Redefinition kills availability

Outside SSA, a variable can be reassigned, and then an expression that looks identical may not be. In t1 = a + b; a = t2; t3 = a + b;, the copy a = t2 gives a a new value number (that of t2), so t3's key differs from t1's and t3 is not redundant. Value numbering handles this automatically by updating the variable's number on every assignment; a text-matching approach would wrongly merge them.

In SSA this problem disappears, because a reassignment creates a new name, which is why compilers convert to SSA before CSE. In a netlist it never arises, since every net has one driver. The savings count is the number of operations whose result is replaced by a copy; in hardware each one is a cell removed, provided the merged cell's fan-out does not hurt timing, which is a trade-off a timing-driven flow may decide to reverse later by duplication.

## Going deeper

Cooper and Torczon, "Engineering a Compiler" (3rd ed.), chapter 8 (local value numbering); the Yosys manual, opt_merge.
