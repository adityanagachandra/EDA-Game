# Library cells as pattern graphs

A mapper cannot match a Liberty cell directly. It needs each cell's function as a small graph in the subject graph's base. Building these pattern graphs is a one-time step when the library is loaded, but it fixes what the mapper can find. Counting the NAND2 and INV nodes of a pattern tells you how much of a subject graph one match absorbs, which is what makes complex cells attractive.

## From cell function to pattern

Take the cell's Boolean function from the library, write it with two-input operators, and decompose it with the same rules used for the subject graph: AND(x, y) = INV(NAND(x, y)), OR(x, y) = NAND(INV x, INV y), NOT x = INV x, then cancel inverter pairs. The pattern's free leaves stand for the cell's input pins. AOI21, y = NOT((a AND b) OR c), gives INV(NAND(NAND(a, b), INV c)): two NAND2 and two INV. NAND2 is the single node NAND(x, y); INV is INV(x). Because the subject graph and the patterns share the base, the matcher only compares node types and connectivity.

## Counting pattern nodes

The number of NAND2 nodes in a pattern equals the number of binary operators in the function as written. The inverter count follows from the raw count, (AND) + 2 (OR) + (NOT), minus two per cancelled pair. For NOR3 written NOT((a OR b) OR c): 2 OR and 1 NOT give 5 raw inverters and no pair cancels, so the pattern has 2 NAND2 + 5 INV = 7 nodes. For OAI22, NOT((a OR b) AND (c OR d)): 3 NAND2, raw inverters 1 + 4 + 1 = 6, one pair cancels at the output, leaving 4. The total node count is how much subject graph one match of the cell covers.

## Multiple shapes and what patterns cannot express

A cell with associative inputs has several pattern shapes: NAND4 can be NAND(INV(NAND(a, b)), INV(NAND(c, d))) or a chain. A structural matcher finds only the shapes it is given, so libraries are often expanded with every shape up to a size limit, which grows the pattern set quickly. Cells whose function needs an input twice, such as XOR2 or a multiplexer, have no tree pattern with distinct leaves, so a tree matcher cannot use them without extensions. These limits are why modern mappers move to Boolean matching on cut functions, where a cell is recognized by its truth table regardless of shape.

## Going deeper

Keutzer, "DAGON", DAC 1987; Rudell, "Logic Synthesis for VLSI Design", PhD thesis, UC Berkeley, 1989, for pattern generation in the SIS-era mappers.
