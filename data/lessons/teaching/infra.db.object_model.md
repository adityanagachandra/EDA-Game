# Blocks, instances, nets and terminals

Every netlist database names the same handful of objects: a block (the design), instances of master cells, nets, and terminals where nets meet pins. OpenDB calls them dbBlock, dbInst, dbMaster, dbNet, dbITerm, dbMTerm and dbBTerm. Getting the counts and links right is the first step in sizing a database, writing a parser that fills it, and choosing a graph model for an algorithm. This lesson counts terminals, connections and the edges a graph expansion creates.

## Instance terminals and block terminals

A master cell has pins, called master terminals. An instance is one use of a master in the block. Each pairing of an instance with one of its master's pins is an instance terminal (ITerm), the object a net actually connects to. The block's own ports are block terminals (BTerms). So the number of ITerms is the sum over masters of instances times signal pins. Example: 12 NAND2 cells with 3 signal pins, 8 inverters with 2 and 5 flip-flops with 3 give $12 \cdot 3 + 8 \cdot 2 + 5 \cdot 3 = 36 + 16 + 15 = 67$ ITerms. The block's 20 ports are 20 BTerms, a separate object kind with its own geometry on the die boundary. Power pins are often handled separately, as special connections, so state which pins are counted. ITerms usually outnumber instances by a factor of 3 to 4, which makes them the largest object population in the database.

## Nets as hyperedges of terminals

A net is a set of terminals: one driver and one or more sinks, possibly including BTerms. In graph terms the netlist is a hypergraph whose vertices are instances (or terminals) and whose hyperedges are nets. Each connected terminal belongs to exactly one net, so the total number of net-terminal connections equals the connected ITerms plus the BTerms. Example: 800 ITerms, of which 6 are spare outputs left unconnected, plus 24 BTerms give $800 - 6 + 24 = 818$ connections. Spread over 240 nets that is $818 / 240 = 3.41$ terminals per net, the average net degree. Most nets in real designs have two or three terminals; a few clock and reset nets have thousands. Averages hide those high-fanout nets, which dominate the cost of algorithms that are quadratic in net size.

## Expanding nets into graphs

Many algorithms, such as quadratic placement or spectral partitioning, need an ordinary graph rather than a hypergraph. The clique model replaces a net with $p$ terminals by an edge between every pair: $p(p-1)/2$ edges. A star model adds one new vertex per net and connects it to the $p$ terminals: $p$ edges. Example: nets with 2, 4 and 7 terminals give $1 + 6 + 21 = 28$ clique edges against $2 + 4 + 7 = 13$ terminal entries in the hyperedge form. A 1,000-terminal clock net alone would add 499,500 clique edges, which is why tools either use stars for large nets, drop or reweight them, or skip clock nets. The database therefore stores the hyperedge, the compact and exact form, and each engine builds the derived graph model that suits its algorithm.

## Going deeper

The OpenDB documentation and its header files for dbBlock, dbInst, dbNet and dbITerm; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 2.
