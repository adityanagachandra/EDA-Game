# Hierarchy flattening and uniquification

After elaboration a design is a tree of instances. Flattening replaces each instance with a copy of its module's contents, so the whole design becomes one module; uniquification gives each distinct specialization its own definition. Both change what later optimizations can see and how large the netlist becomes, so tool builders must be able to predict their effect exactly.

## Counting copies after flattening

Flattening copies a module once for every path from the top to an instance of it. If top contains 4 alu instances and 3 ctrl instances, each alu contains 2 regfile instances and each ctrl contains 3, the flat netlist contains $4 \times 2 + 3 \times 3 = 17$ copies of regfile's logic. The count is the number of root-to-leaf paths in the instance tree that end at regfile, found by multiplying instance counts along each path and summing over the paths.

Hierarchical names survive as part of the flattened names, typically joined with dots, such as alu_2.rf_1.mem, so reports can still be traced to the original hierarchy. Yosys's flatten pass performs this inlining, leaving black boxes and modules marked keep_hierarchy in place.

## Counting cells in the flat module

The size of the flat netlist follows from a bottom-up sum. Each module's total is its own primitive cells plus the totals of its children times their instance counts. With top 12 cells, alu 34 cells and 3 regfiles, ctrl 8 cells and 3 regfiles, regfile 18 cells, 2 alus and 1 ctrl: alu totals $34 + 3 \times 18 = 88$, ctrl totals $8 + 54 = 62$, and the flat module has $12 + 2 \times 88 + 62 = 250$ cells.

This is a post-order traversal of the module DAG, computed once per definition rather than once per instance, which matters for large designs: the flat netlist can be far larger than the source suggests. That growth is the cost of flattening; the benefit is that constants and duplicate logic across former boundaries become visible to optimization.

## Uniquification by parameter set

Before flattening, elaboration already creates one module definition per distinct parameter assignment. Instances that agree on every parameter value share one specialization: two fifo instances with default parameters share the original fifo, while fifo #(.W(32), .DEPTH(8)) and fifo #(.W(16)) each need their own. A design whose top contains ctrl twice, fifo twice with defaults and one fifo with W = 32 and DEPTH = 8 has four module definitions: top, ctrl, fifo and the specialized fifo.

Uniquification goes one step further when instances of the same specialization are later optimized differently, for example because their inputs are tied to different constants: each must become a separate definition first. Flattening makes the question moot by copying everything, which is one reason synthesis scripts flatten before the main optimization passes.

## Going deeper

The Yosys manual entries for hierarchy, flatten and uniquify; IEEE Std 1800-2017, clause 23.
