# Timing sense and unateness of arcs

A timer does not propagate one arrival time per pin but two, one for a rising and one for a falling transition. Which input transition produces which output transition depends on the cell's logic function, captured in Liberty as each arc's timing sense. Getting this mapping right is what lets the timer pair the correct delay table with the correct transition, and it explains why some cells, such as XOR, make timing more pessimistic than others.

## Positive, negative and non-unate arcs

An arc is positive unate when a rising input can only cause a rising output and a falling input only a falling output. Buffers, AND and OR gates are positive unate from each input. An arc is negative unate when the output moves the opposite way: inverters, NAND, NOR, AOI and OAI gates. Static CMOS gates are naturally negative unate, because a rising input turns on an NMOS pull-down. An arc is non-unate when the output direction depends on the other inputs. In an XOR, a rising input A makes the output rise if B is 0 and fall if B is 1. Liberty records this in the arc's timing_sense attribute as positive_unate, negative_unate or non_unate, and the timer reads it when it builds the graph.

## Following a transition through a chain

To know which transition reaches the end of a path, walk the arcs in order. A positive-unate arc keeps the direction, a negative-unate arc flips it. A falling input through an OR, a buffer and a NOR becomes falling, falling, then rising: the path ends with a rising transition. A timer does this for both transitions at once. It keeps a rise arrival and a fall arrival at each pin, and a negative-unate arc maps the input rise arrival plus the cell_fall delay to the output fall arrival, and the input fall arrival plus the cell_rise delay to the output rise arrival. Path reports show the direction at each pin with a rise or fall marker, which is how an engineer checks that the right delay tables were used.

## Non-unate arcs propagate both edges

Through a non-unate arc either output transition is possible, so the timer must propagate both: the output rise arrival is the worse of the input rise and fall arrivals plus the rise delay, and likewise for fall. Once both directions exist they continue downstream, even through later unate arcs, because a buffer passes a rise as a rise and a fall as a fall. For example, a falling edge into an XOR, then a NAND, then an AND yields both a rise and a fall at the end. This is a source of pessimism: the timer combines the worst of two transitions that may not both be reachable from one real input vector. Mode-specific constants (set_case_analysis) can resolve a non-unate arc to a unate one when a side input is known.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 3 on timing arcs and unateness; the Liberty reference manual on timing_sense.
