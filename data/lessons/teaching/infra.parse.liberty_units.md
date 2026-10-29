# Liberty units and scaling

Liberty tables are just numbers; their meaning comes from unit attributes in the library header. A delay of 0.05 is 50 ps in one library and 0.05 ps in another. Every reader must scale each value by its library's units into one internal unit system, and every consumer must trust that it did. Unit bugs are silent: the tool runs, and every slack is wrong by a factor of a thousand.

## Time units scale delay tables

The header attribute `time_unit` takes one of the values "1ps", "10ps", "100ps" or "1ns", and every delay and transition value in the library is a multiple of it. With `time_unit : "1ns" ;`, a `cell_rise` entry of 0.0532 is $0.0532 \times 1000 = 53.2$ ps. With "10ps", an entry of 5.32 is also 53.2 ps. The input-transition axis of a table (index_1 in a typical delay template) is in the same time unit. A reader converts at parse time: it computes the factor once from the header, multiplies each value as it stores it, and keeps everything internally in one unit, often seconds as a double, so that engines never see library units at all.

## Capacitive load units scale pins and indices

Capacitance uses the complex attribute `capacitive_load_unit (1,pf) ;` or `(1,ff)`. Pin capacitances and the output-load axis of delay tables (often index_2) are multiples of it. A pin with `capacitance : 0.0016 ;` in a pF library is 1.6 fF. The conversion runs both ways. When a timer computes a net load of 4.5 fF and must look it up in a pF table, it needs the index value $4.5 / 1000 = 0.0045$. Doing the lookup with 4.5 instead would land far beyond the table's last index and extrapolate wildly. Libraries also declare `voltage_unit`, `current_unit`, `pulling_resistance_unit` and `leakage_power_unit`, each handled the same way.

## Units belong to each library

A session often loads several libraries: standard cells, IO cells, memories from a generator. Nothing forces them to share units. Suppose library A uses "1ns" and library B uses "1ps", and a path has arc delays 0.12 (A), 35 (B) and 0.08 (A). The path delay is $120 + 35 + 80 = 235$ ps. Adding the raw numbers gives 35.2, which means nothing. This is why scaling belongs in the reader, per library, rather than in the timer: the timer adds delays from many cells and cannot know which library each number came from unless the reader already normalized it. A related trap is SDC: its numbers are unitless and take their meaning from `set_units` or from the library units, so changing libraries can silently change what a constraint means.

## Going deeper

The Liberty Reference Manual, library-level attributes; the OpenSTA source for unit handling in its Liberty reader.
