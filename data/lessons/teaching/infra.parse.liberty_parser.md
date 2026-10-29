# Parsing Liberty

Liberty (.lib) files carry the timing, power and function of every library cell, and every synthesis, timing and power tool must read them. The format is a simple nested grammar of groups and attributes, but real libraries are huge and full of vendor extensions. A tool designer must parse the grammar generically, interpret the timing semantics correctly, and do both fast enough for files of hundreds of megabytes.

## Groups, simple and complex attributes

Liberty has three statement shapes. A group has a name, a parenthesized argument and a braced body: `cell (NAND2_X1) { ... }`, `pin (A1) { ... }`, `timing () { ... }`. A simple attribute is `name : value ;`, for example `direction : input ;` or `time_unit : "1ns" ;`. A complex attribute is `name (v1, v2, ...) ;`, for example `capacitive_load_unit (1,pf) ;` or `index_1 ("0.01, 0.04, 0.16") ;`. Because every statement has one of these three shapes whatever its name, one small recursive-descent parser can build a generic tree of groups and attributes for any Liberty file, including vendor attributes it has never seen. A second pass interprets the tree into cells, pins and arcs. This split keeps the grammar code tiny and puts all the semantics in one place.

## What the parser must keep for timing

Timing lives in `timing ()` groups inside the pin where the arc ends. For the arc from A1 to ZN of a NAND2, the group sits inside `pin (ZN)` with `related_pin : "A1"` and `timing_sense : negative_unate`. Inside it, tables such as `cell_rise (delay_template_7x7)` and `rise_transition (...)` hold `index_1`, `index_2` and `values`. The template name points to an `lu_table_template` declared earlier at library level, whose `variable_1` and `variable_2` say which axis is input transition and which is output load. The parser must keep the template binding, both index vectors, the value matrix (checking it has as many rows as index_1 entries and columns as index_2 entries), the related pin, the timing sense and any `when` condition. The output pin's `function` string, such as `"!(A1 & A2)"`, gives the logic that synthesis needs.

## Scale, streaming and the header

A library for one process corner can run to hundreds of megabytes, almost all of it tables. A placer needs only areas and pin capacitances; a timer needs the tables but maybe not the power data. A streaming parser that interprets the groups it needs and skips the rest, while still tracking brace nesting, saves both time and memory. Order matters: units such as `time_unit` and `capacitive_load_unit` and the table templates are declared in the library header before the cells, so a parser can convert values as it reads them. A table that names an undeclared template is an error, not something to guess at. Several libraries in one session can use different units, so conversion into the tool's internal units must happen per library at read time. Error messages need file and line numbers, since users fix libraries by hand.

## Going deeper

The Liberty Reference Manual (the open Liberty specification); the OpenSTA source for its Liberty reader; Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter on cell libraries.
