# Four-state logic and X propagation

Real wires are 0 or 1, but a simulator must also represent what it does not know. Verilog uses four values: 0, 1, X (unknown) and Z (high impedance, undriven). X appears for uninitialized registers, conflicting drivers and undefined operations; Z appears on tri-state buses. Every operator needs a four-state truth table, and those tables decide whether an unknown is reported or silently lost. Tool designers care because four-state evaluation is slower than two-state and its rules shape every X-propagation debate.

## Four values and what they mean

The value **0** and **1** are known logic levels. **X** means the simulator cannot determine the value: it might be 0, it might be 1, or the net might be contended. **Z** means no driver is driving the net, as on a released tri-state bus. A register declared as `reg` or `logic` starts at X until it is assigned, which is how simulation reveals missing resets. A `bit` variable is two-state and starts at 0. Implementations store each four-state bit as two machine bits (for example a value bit and an unknown bit), so a 64-bit four-state vector takes two words and every operation needs extra instructions to propagate the unknown bits. Two-state simulators such as Verilator avoid that cost by mapping X and Z to concrete values.

## Truth tables with X and Z

At a logic operator, a Z input is treated as X: an undriven input to a gate has an unknown level. AND has the **controlling value** 0: if either input is 0, the result is 0 regardless of the other, so $0 \,\&\, X = 0$, but $1 \,\&\, X = X$ and $1 \,\&\, Z = X$. OR has controlling value 1: $1 \mid X = 1$, $0 \mid X = X$. XOR has no controlling value: any X or Z input makes the result X, so $1 \oplus X = X$ and even $X \oplus X = X$, although a real signal XORed with itself is 0. NOT maps X and Z to X. Evaluate expressions inside out: with $a = 1$, $b = Z$, $c = 0$, the expression $\lnot(a \,\&\, b) \mid c$ gives $1 \,\&\, Z = X$, then $\lnot X = X$, then $X \mid 0 = X$.

## Expressions and driver resolution

When several continuous assignments or gates drive one plain `wire`, the simulator resolves their values. A driver at Z yields to the others: Z with 1 gives 1. Equal values agree: 0 with 0 gives 0. A 0 against a 1 is contention and gives X, as does anything against X (other than Z). Example: `assign w = en ? a : 1'bz; assign w = c;` with $en = 1$, $a = 0$, $c = 1$ gives 0 against 1, so $w = X$; with $en = 0$ the first driver is Z and $w = c$. Verilog also defines `wand`, `wor` and strength levels for more detailed resolution, and VHDL's `std_logic` uses a nine-valued resolved type, but the four-value plain-wire table covers most RTL debugging.

## Going deeper

IEEE Std 1800-2017, clause 6 (data types) and clause 11 (operators, four-state truth tables); IEEE Std 1164 for VHDL's nine-valued logic.
