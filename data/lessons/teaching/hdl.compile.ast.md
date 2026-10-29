# Abstract syntax trees for HDLs

An abstract syntax tree (AST) keeps the structure of a program and drops the punctuation that only served to fix that structure. It is the first real data structure of a hardware compiler: Yosys's Verilog front end builds one, slang builds a richer one, and elaboration, type checking and lowering to netlist IRs all walk it. Counting its nodes, its height and its shared subtrees exercises what a compiler stores.

## Nodes of an expression tree

In an expression AST every operator is an internal node whose children are its operands, and every operand occurrence is a leaf. Parentheses produce no node: their only job was to choose the shape, which the tree records directly. A unary operator has one child, a binary operator two, and the conditional ?: is one node with three children.

So (t ? (2 & c) : a) ^ ~(a - d) has the operators ^, ?:, &, ~ and -, five internal nodes, and the leaf occurrences t, 2, c, a, a and d, six leaves, for 11 nodes in all. Note that a appears twice as a leaf: an AST is a tree, so a repeated operand is a repeated node. A concrete syntax tree, by contrast, would keep every token, parentheses included.

## Height and traversal

The height of a tree is the number of edges on the longest path from the root to a leaf. The root is the operator applied last, the loosest-binding one at the top level. For (1 ^ b) ^ (1 - c), the root ^ has two children of height 1, so the tree has height 2. For an expression tree the height is the depth of nesting, which for a combinational circuit built directly from it is the number of logic levels on the longest path.

Compilers walk ASTs in two orders. A post-order traversal visits children before their parent, which is how an evaluator or a code generator processes operands before the operator. A pre-order traversal visits the parent first, which is how attributes such as expected widths are pushed down, exactly the two passes Verilog's sizing rules require.

## Sharing subtrees in a DAG

Identical subtrees compute identical values, so a compiler can store each one once. Hash-consing does this: before creating a node, look up its operator and the identities of its children in a table, and reuse an existing node if one matches. The tree becomes a directed acyclic graph.

In ((2 - d) * (b ^ 1)) | (c + (2 - d)) the tree has 13 nodes. The subtree 2 - d appears twice and becomes one node, and the leaves 2 and d are shared with it, as is every repeated operand. The DAG has 10 nodes. This is common subexpression elimination performed at construction time, and it is the representation behind and-inverter graphs in ABC and the structural hashing done by every logic synthesizer. Two subtrees match only if they are written identically, operand order included, unless the compiler first puts commutative operators into a canonical order.

## Going deeper

Appel, "Modern Compiler Implementation in ML", chapter 4; Mishchenko, Chatterjee and Brayton, "DAG-Aware AIG Rewriting" (DAC 2006).
