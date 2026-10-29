# Registers and register-transfer traces

A register is a row of flip-flops sharing one clock, holding a multi-bit value. Register-transfer level (RTL) design describes a circuit as registers plus the combinational functions that compute their next values every cycle. Tracing an RTL description cycle by cycle is what a cycle-based simulator such as Verilator does, and getting the simultaneous-update semantics right is the first thing an RTL author or a tool writer must internalize.

## Shift registers

A shift register passes its bits one position along on every clock edge. In a left-shifting 8-bit register with serial input at the least significant bit, each edge computes $R^+ = (2R + s_{in}) \bmod 2^8$: the top bit falls off and the serial bit enters at the bottom.

Start at $R = 00001101$ (13) and shift in the bits 1, 0, 1. After the first edge $R = 00011011$ (27), after the second 00110110 (54), after the third 01101101 (109). Shift registers implement serial-to-parallel conversion, delay lines, and, with XOR feedback, linear-feedback shift registers used for pseudo-random test patterns and CRCs.

## Simultaneous update

In RTL every register's next value is computed from the values *before* the edge, then all registers update together. With nonblocking assignments

`A <= A + B; B <= A;`

starting from $A = 1, B = 1$, the first edge computes $A^+ = 2$ and $B^+ = 1$ (the old $A$). The second gives $A = 3, B = 2$, the third $A = 5, B = 3$: Fibonacci numbers. If the statements were read sequentially, as blocking assignments in a software mindset, $B$ would receive the *new* $A$ and the trace would be different. Writing `A <= B; B <= A;` swaps two registers every cycle, which a sequential reading would get wrong. Mismatches between the RTL and its intended semantics here are among the most common simulation-versus-synthesis bugs.

## Enables and synchronous clear

Real registers load only when told. A register with enable `en` and synchronous clear `clr` behaves as: if `clr` then $R^+ = 0$; else if `en` then $R^+ = $ new value; else $R^+ = R$. The hold case is a 2-to-1 mux feeding $R$ back to its own D input, which is how synthesis implements an enable unless it chooses a clock-gating cell instead.

Trace a 4-bit counter $R^+ = (R + 1) \bmod 16$ starting at 14, with (en, clr) over five edges $= (1,0), (1,0), (0,0), (1,1), (1,0)$: $R$ goes 15, 0 (wrap), 0 (hold), 0 (clear), 1. Priority matters: here clear beats enable, and swapping that priority in a specification change is a typical small RTL edit an equivalence checker must catch.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 4 on HDL sequential logic; Cummings, "Nonblocking Assignments in Verilog Synthesis, Coding Styles That Kill!" (SNUG 2000).
