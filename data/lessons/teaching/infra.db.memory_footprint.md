# Memory footprint of a netlist

A modern block can have tens of millions of instances and over a hundred million terminals, so bytes per object decide whether a design fits in memory and how fast engines sweep it. Estimating a footprint is simple arithmetic over object counts and record sizes, but the record sizes depend on choices a tool designer makes: pointer width, padding, and what is stored inline. This lesson sizes a netlist, compares pointers with indices and computes padded record sizes.

## Counting bytes per object

Footprint is the sum over object kinds of count times record size. Take 1,000,000 instances of 64 bytes, 1,100,000 nets of 48 bytes, and 3.2 ITerms per instance at 24 bytes each. Instances take $64 \times 10^6$ bytes, nets $52.8 \times 10^6$ bytes, ITerms $3.2 \times 10^6 \times 24 = 76.8 \times 10^6$ bytes, for 193.6 MB. ITerms dominate because they are the most numerous objects, so shaving 4 bytes from the ITerm record saves more than shaving 16 from the instance. The estimate is linear in design size, $O(n)$, but the constant matters: a 20-fold larger design at the same density needs about 3.9 GB for these objects alone, before names, wires, timing graphs and engine-private data, which together often exceed the netlist itself.

## Indices instead of pointers

On a 64-bit machine a pointer takes 8 bytes. A reference stored as a 32-bit index into a per-type table takes 4, and 32 bits can address over 4 billion objects of one kind, more than any block needs. An ITerm typically references its instance, its master terminal, its net and neighbours in per-net or per-instance lists. With 4 references per ITerm and 10 million ITerms, switching from pointers to indices saves $10^7 \times 4 \times 4 = 160$ MB. Indices bring further benefits: the table can grow and move in memory without breaking references, the tables can be written to disk and read back without pointer fix-up, and a side array indexed by the same number can attach engine data to every object. The cost is one extra addition per access, base plus index times size, which is almost always free next to the cache misses it prevents.

## Alignment and padding

A C++ compiler aligns each field to a multiple of its size and rounds the record up to a multiple of its largest field, so field order changes the size. Take fields flags (1 byte), master (8), x (4), orient (1), net (8). In that order, flags sits at 0, master is padded to 8, x at 16, orient at 20, net is padded to 24 and ends at 32: the record is 32 bytes, of which 22 are payload and 10 are padding. Sorted from largest to smallest, master, net, x, flags, orient end at 22 and round up to 24 bytes. Saving 8 bytes on a record repeated 10 million times is 80 MB. Database authors also pack flags into bit fields, store coordinates as 32-bit integers in database units rather than doubles, and keep variable-size data such as names out of line.

## Going deeper

Ulrich Drepper, "What Every Programmer Should Know About Memory" (2007); the OpenDB source in the OpenROAD repository for its object tables.
