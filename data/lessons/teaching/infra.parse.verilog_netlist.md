# Parsing structural Verilog netlists

After synthesis, a design travels as a structural Verilog netlist: a module with ports, wire declarations and instances of library cells. Every back-end tool starts by parsing this text into its object model, turning ports into block terminals, instances into instance objects, and each connection into an instance terminal on a net. A tool designer must know exactly what the text implies, because a mis-parsed connection silently becomes a wrong circuit.

## Instances, nets and terminals in the text

Read `module top (a, b, y0); input a, b; output y0; wire n1; NAND2_X1 u1 (.A1(a), .A2(b), .ZN(n1)); INV_X1 u2 (.A(n1), .ZN(y0)); endmodule`. The parser creates two instances, u1 and u2. Each port becomes a block terminal (BTerm) and also names a net, so the nets are a, b, y0 and the wire n1: four nets. Each instance pin becomes an instance terminal (ITerm): u1 has three pins and u2 two, so five ITerms. A useful invariant for a netlist where every cell has one output: each instance drives exactly one net, so the net count is the number of primary inputs plus the number of instances. If the count from the parser disagrees, something was dropped or duplicated.

## Fanout from named connections

A named connection `.A2(n1)` binds pin A2 of the instance to net n1, whatever order the connections are written in. The parser looks up A2 on the cell master, finds its direction from the library, and attaches an ITerm to n1. The fanout of a net is the number of input pins it drives. In `NAND2_X1 u3 (.A1(n1), .A2(n1), .ZN(n2)); INV_X1 u4 (.A(n1), .ZN(y0));`, net n1 drives u3/A1, u3/A2 and u4/A: fanout 3, even though only two instances read it, because each pin is a separate sink. Output pins such as ZN are drivers and are not counted. Fanout feeds everything downstream: load capacitance for timing, net weights for placement, and buffering decisions.

## Positional connections need the cell's port order

Verilog also allows positional connections: `NAND2_X1 u5 (n1, b, n3);` binds signals to ports in the order the cell's module declares them. The netlist alone cannot tell you which signal is the output: the parser needs the port order, from a Verilog cell library or from the order the tool assumes for the master. If NAND2_X1 is declared as (A1, A2, ZN), then n1 is A1, b is A2 and n3 is the output ZN. Get the order wrong and the parser silently turns an output into an input, creating a net with two drivers and a dangling sink. Robust readers therefore check that every net has exactly one driver after reading, and treat positional connections to cells without a declared port order as an error. Escaped identifiers such as `\bus[3] ` (a backslash name ended by whitespace) are another trap for hand-written lexers.

## Going deeper

IEEE Std 1364-2005, the Verilog standard, sections on module instantiation; the Yosys manual on reading and writing netlists; the OpenROAD documentation for read_verilog and link_design.
