# Hierarchy and instancing

Designs are written hierarchically: a top module instantiates cores, which instantiate ALUs, which hold leaf cells. A database can keep that hierarchy, storing each module definition once and referring to it from each instance, or flatten it into one level of leaf cells. The choice decides memory, how edits work and how names are formed. Physical design engines usually work flat, while hierarchical storage pays off for repeated blocks, layout (GDS) and hierarchical verification.

## Flattening a module hierarchy

Flattening replaces every module instance by a copy of its contents, recursively, until only leaf cells remain. Count bottom up. Suppose alu holds 400 leaf cells, core holds 150 leaf cells and 4 alu instances, periph holds 90 leaf cells, and top holds 20 leaf cells, 3 cores and 2 periph instances. One core flattens to $150 + 4 \cdot 400 = 1750$ leaves. Top flattens to $20 + 3 \cdot 1750 + 2 \cdot 90 = 5450$ leaves. This is a post-order traversal of the module tree: each module's flattened count is its own leaves plus the sum of child counts times multiplicities. Flat names concatenate the path, such as core_1/alu_2/U7, so they grow with depth, and a flat database must store every one of those 5450 instances with its own placement.

## Definitions versus occurrences in memory

A hierarchical database stores each definition once: one object for each leaf cell and each child-module instance inside each module. In the example, top has $20 + 3 + 2 = 25$ objects, core $150 + 4 = 154$, alu 400 and periph 90, for 669 objects in all, against 5450 in the flat form, a saving of 4781. The saving grows with replication: a chip with 64 identical cores stores one core definition. Layout formats exploit the same idea; a GDS structure is referenced many times. The price is that per-occurrence data, such as a location, a routed wire or a timing value, cannot live in the shared definition, because each occurrence has its own. Tools that keep hierarchy must keep occurrence data in separate tables keyed by hierarchical path.

## Occurrences and uniquification

An occurrence is one path from the top to an instance of a module. With 3 cores each holding 4 alus there are $3 \cdot 4 = 12$ occurrences of alu, all sharing one definition. If an ECO resizes a cell in core_1/alu_2 only, editing the definition would change all 12. The tool must first uniquify: copy the definition into a new module, say alu_u1, point that one occurrence at the copy and edit it. Synthesis does the same when one instance of a module is optimized differently from the others. Uniquification trades memory for independence, and repeated uniquification drifts toward a flat design. That is one reason physical design engines such as those in OpenROAD operate on a flattened netlist and keep the hierarchy only in names or in optional module objects.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 1; the Yosys manual on the hierarchy and flatten passes.
