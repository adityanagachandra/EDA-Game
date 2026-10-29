# Liberty structure: cells, pins and arcs

Liberty is the text format in which a cell library tells synthesis, timing and power tools what each cell does. Its structure, a library header, cells, pins, and timing groups that describe arcs between pins, defines the timing graph a static timer builds. A tool designer writing a Liberty reader or a timer must know exactly what each group and attribute means.

## Cells, pins and timing arcs

A Liberty file opens with a library group whose header sets units (time_unit, capacitive_load_unit, voltage and power units), operating conditions, measurement thresholds and lookup-table templates. Each cell group holds the cell's area, leakage and a pin group per pin. An input pin has a direction and a capacitance. An output pin has a direction, a function attribute such as function : "!(A & B)", and one timing group per arc that ends on it. Each timing group names its start pin with related_pin, so the arc runs from related_pin to the pin that contains it, and carries lookup tables for delay and output transition. A timer turns every timing group into an edge of the timing graph and every pin into a node.

## Timing sense and timing type

Each timing group says how input edges map to output edges with timing_sense. A positive_unate arc can only move the output in the same direction as the input (AND, OR, buffer, a mux data input). A negative_unate arc moves it in the opposite direction (inverter, NAND, NOR, every AOI or OAI input). A non_unate arc can do either, depending on other inputs (XOR, a mux select pin, a full adder's sum). The timing_type attribute says what kind of arc it is: combinational for ordinary delay, rising_edge or falling_edge for clock-to-output delay in sequential cells, and setup_rising, hold_rising and similar for constraint checks. A timer uses timing_sense to decide which output edge, and so which table, follows each input edge; for a non_unate arc it propagates both.

## What the structure hides and assumes

Liberty describes a cell only by its pin-to-pin arcs, so a missing arc is a missing graph edge, and paths through it go silently untimed. Pin capacitance is a single value (or a rise and a fall value), summed into the load of the driving net, so it must match the real layout. State-dependent arcs, written with a when condition such as when : "B", let one input pin carry different tables for different side-input states; a timer either selects by state or takes the worst. Because table entries are bare numbers, a reader that ignores the header's units or its threshold definitions computes wrong delays from correct data. Writing a robust Liberty reader is largely about honouring these conventions consistently.

## Going deeper

The Liberty Reference Manual (Liberty User Guides and Reference Manual Suite); the OpenSTA source, Liberty reader.
