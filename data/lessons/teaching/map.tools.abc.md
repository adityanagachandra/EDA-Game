# ABC map, amap and if

ABC, from UC Berkeley, is the open-source engine behind most academic mapping and the mapping step of Yosys. Its mappers implement the ideas of this domain directly: AIGs as subject graphs, structural choices, priority cuts, area flow and exact area, Boolean matching with supergates. Reading an ABC script as a sequence of these algorithms is the fastest way to connect theory to a working tool.

## The script skeleton

A mapping script reads a design, converts it to an AIG with strash (structural hashing), optimizes it, optionally computes structural choices, maps, and reports. For LUTs: strash, optimization commands, then if -K 6 to map into 6-input LUTs, then print_stats for LUT count and depth. For standard cells: read_lib to load a Liberty file, strash, optimization, then map or amap, then stime for timing with the library data. Every mapper works on the current AIG, so what runs before it shapes the subject graph and therefore the result.

## Which mapper does what

if is the cut-based FPGA mapper built on priority cuts: depth-oriented first, then area recovery with area flow and exact area, with -K setting the LUT size. map is the standard-cell mapper based on Boolean matching against supergates generated from the library; it is delay-oriented and recovers area afterwards. amap is an area-oriented standard-cell mapper. dch, run before a mapper, computes structural choices, so the mapper can find matches in several equivalent structures at once. After mapping, mfs resynthesizes nodes using don't-cares and lutpack repacks LUTs, both working on the mapped network.

## Checking and finishing the result

Always verify: cec checks combinational equivalence between the original and mapped networks with simulation and SAT. For standard cells, stime reports the critical path under the library's timing data, and upsize and dnsize resize gates on and off critical paths. Results depend on the script, the structure of the input AIG and the library: two runs that differ only in whether dch ran can produce the same depth with different LUT counts, because choices change which cuts exist. When comparing scripts, keep everything else fixed and compare QoR from an independent timer, not only the mapper's own numbers.

## Going deeper

Brayton and Mishchenko, "ABC: an academic industrial-strength verification tool", CAV 2010; the ABC source code and its built-in command help; Mishchenko, Cho, Chatterjee and Brayton, "Combinational and sequential mapping with priority cuts", ICCAD 2007.
