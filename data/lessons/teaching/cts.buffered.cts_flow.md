# The clock tree synthesis flow in OpenROAD

In an open-source flow, CTS is a sequence of steps around the clock_tree_synthesis command of OpenROAD's CTS module (TritonCTS). Knowing the sequence explains why timing changes after CTS, why hold is fixed then, and where clock nets get special routing. A tool designer should be able to say what each step consumes and produces, and what breaks if a step is skipped.

## Where CTS sits

CTS runs after global and detailed placement and pre-CTS optimisation, and before global and detailed routing. Before CTS the timer treats clocks as ideal: zero network latency, with set_clock_uncertainty standing in for the skew the tree will add. CTS needs the placed sink locations and a list of clock buffer cells it may use; in OpenROAD the user passes the buffer list and optionally a root buffer to clock_tree_synthesis. The output is a placed buffer tree for every clock net, which the router then treats as fixed, high-priority nets.

## Inside the CTS step

OpenROAD's CTS groups nearby sinks into clusters, builds a balanced, H-tree-style topology over the clusters, and inserts buffers chosen from the allowed list using characterised delay and slew. Clustering keeps the leaf nets short and the top-level problem small. After tree construction the flow repairs long clock wires, and legalises the inserted buffers, because CTS places them where the topology wants, possibly on occupied sites. Clock nets are often given non-default routing rules: wider wires and larger spacing, or shields, which lower resistance and crosstalk on nets that toggle every cycle.

## After CTS

Timing then switches to propagated clocks, so every register sees its real latency through the tree. Results change. If the built tree has 120 ps of skew between critical pairs where 50 ps of uncertainty was budgeted, paths that met timing before CTS can now fail. Hold is repaired now, by inserting delay on short paths, because hold depends on real skew. On-chip variation matters too: pairs whose clock paths split near the root have long non-common paths, which carry more uncorrelated variation and get less credit from CPPR. CTS quality is reported as latency, skew, buffer count and wirelength together.

## Going deeper

The OpenROAD documentation for the CTS module and the OpenROAD flow scripts; Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
