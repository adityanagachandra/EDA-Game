# Cache locality and data-oriented layouts

Placement, timing and routing engines spend most of their time in loops over millions of objects, and those loops are usually limited by memory traffic rather than arithmetic. Memory moves in cache lines, typically 64 bytes, so what matters is how many lines a loop touches. Data-oriented layout arranges data by how loops read it: arrays of structures, structures of arrays, or hot and cold splits. It is one of the cheapest large speedups available to a tool designer.

## Cache lines touched by a loop

A cache line is the unit the hardware fetches: read one byte and the whole 64-byte line comes in. Suppose 65,536 cells are stored as an array of 96-byte structs, and a density loop reads only each cell's 4-byte x coordinate. The structs are contiguous, so the loop sweeps the whole array: $65536 \times 96 = 6{,}291{,}456$ bytes, or 98,304 lines. Only 4 of every 96 bytes brought in are used, about 4%. If the structs were instead scattered across the heap, each accessed through its own pointer, every cell would cost at least one separate line, and the hardware prefetcher could not predict the next address, so each access might wait for main memory. Contiguity and density are the two levers: walk memory in order, and make the bytes you walk be the bytes you need.

## Structure of arrays

A structure of arrays stores each field in its own contiguous array: all x values together, all y values together, and so on. The same loop over x now reads $65536 \times 4 = 262{,}144$ bytes, 4,096 lines, instead of 98,304. The reduction factor is the struct size over the field size, $96 / 4 = 24$. Dense arrays of one type also let the compiler vectorize: one SIMD instruction can add a displacement to 8 x values at once. Electrostatic global placers keep cell positions, sizes and gradients in exactly this form, often as arrays of floats ready for vector units or GPUs. The cost is that code reading many fields of one object touches several arrays, and inserting or deleting an object must update every array, so the layout suits the hot inner loops more than the editing interface.

## Hot and cold splitting

A full structure of arrays is not always worth the disruption. A middle ground splits each object into a hot part, the few fields the inner loops read, and a cold part, everything else: names, property lists, rarely used flags. If the hot part is 16 bytes (x, y, width, a net index), a pass over 65,536 cells touches $65536 \times 16 / 64 = 16{,}384$ lines instead of 98,304, most of the gain of a full split while related fields stay together. The cold part is reached by the same index when needed. The same reasoning drives other choices: storing net terminals in one flat array with per-net offsets rather than a linked list per net, sorting cells by location so neighbours share lines, and keeping coordinates as 32-bit integers rather than 64-bit doubles. Measure first: a profiler that counts cache misses shows which loops are memory bound.

## Going deeper

Ulrich Drepper, "What Every Programmer Should Know About Memory" (2007); Richard Fabian, "Data-Oriented Design" (2018).
