# Partition refinement in netlist comparison

Partition refinement is the algorithm at the heart of LVS tools such as Netgen and Gemini. It colours every device and net, then repeatedly recolours each vertex by its own colour and the colours of its neighbours, until no class splits. Vertices left alone in a class are matched; classes that differ between the two netlists localize errors. Running it by hand on a small netlist shows exactly how information flows.

## One round of refinement

Start with devices coloured by type, VDD and GND in classes of their own (matched by name), and other nets coloured by degree. In a NOR2 (PMOS A from VDD to n1, PMOS B from n1 to Y, NMOS A and B from GND to Y), the initial classes are PMOS, NMOS, VDD, GND, the degree-2 nets {A, B, n1} and the degree-3 net Y: 6 classes. One round recolours each vertex by (its colour, the multiset of (terminal kind, neighbour colour)). The two PMOS split (one touches VDD, one touches Y), and n1 splits from A and B (it touches source/drain, they touch gates): 8 classes.

## Rounds until stable

Refinement repeats until a round leaves the number of classes unchanged. In the NOR2, round two splits A from B, because A gates the PMOS next to VDD and B the PMOS next to Y. Round three splits the two NMOS, which differ only in their gate nets. Round four changes nothing, so the class counts are 6, 8, 9, 10 and there are three splitting rounds. Each round costs time linear in the number of edges, using hashing of signatures, and the number of rounds is bounded by the number of vertices but is small for real circuits because information spreads one hop per round.

## The stable partition

The stable partition is the coarsest partition in which vertices of the same class have the same number of neighbours of each class along each terminal kind. In the NOR2 every vertex ends in its own class: 10 singletons, so every device and net is matched. In general, singleton classes in both netlists match outright, equal-size larger classes are ambiguities resolved by tie-breaking, and unequal classes are errors. The procedure is colour refinement, the one-dimensional Weisfeiler-Leman algorithm, run on both graphs with shared colour names.

## Going deeper

Ebeling, "GeminiII: A Second Generation Layout Validation Program" (ICCAD, 1988); the Netgen documentation.
