# Four-state logic values

Verilog and SystemVerilog simulate with four values per bit: 0, 1, x (unknown) and z (high impedance). The extra two values let a simulator represent uninitialized registers, conflicting drivers and undriven nets, and they propagate through logic by precise rules. A tool builder must know those rules exactly, because a simulator, a linter and an equivalence checker that disagree about x produce contradictory verdicts on the same design.

## Bitwise operators on x and z

At the input of a logic operator, z behaves exactly like x: an undriven value is unknown to the gate reading it. The rules then follow from asking whether the unknown input could change the output. For AND, a 0 on either input forces 0, so $0 \wedge x = 0$, while $1 \wedge x = x$. For OR, a 1 forces 1, so $1 \vee x = 1$ and $0 \vee x = x$.

Work lane by lane on vectors. With a = 4'b10xz and b = 4'b1101, a & b is 4'b100x: lane 3 is 1 & 1 = 1, lane 2 is 0 & 1 = 0, lane 1 is x & 0 = 0, and lane 0 is z & 1 = x. Three classic mistakes give different answers: letting z pass through unchanged, reading unknowns as 0, and making every lane with an unknown input x even when the other input controls it.

## XOR, XNOR and inversion

XOR has no controlling value: flipping either input always flips the output. So any x or z input makes the XOR lane x, and the same holds for XNOR (written ~^ or ^~). Inversion maps 0 to 1 and 1 to 0, and maps both x and z to x.

Compound expressions are evaluated operator by operator. For a & ~b with a = 4'b1x01 and b = 4'b0z11, first ~b = 4'b1x00, then AND with a gives 4'b1x00: lane 2 is x & x = x, and lanes 1 and 0 are forced to 0 by the zeros of ~b. The structure of the rules mirrors the Boolean algebra of cofactors: an unknown input matters exactly when the output differs between its two possible values.

## Reductions and the two equalities

Reduction operators fold a vector to one bit with the same rules: &a is 0 if any bit is 0, otherwise x if any bit is unknown, otherwise 1. |a is 1 if any bit is 1, otherwise x or 0. ^a is x as soon as any bit is unknown.

Equality comes in two forms with different purposes. Logical equality a == b returns 0 if some lane holds two known, different bits; otherwise it returns x if any lane holds x or z; otherwise 1. Case equality a === b compares the four-state values literally, so x matches only x and z only z, and the result is always 0 or 1. That makes === useful in testbenches (checking that a signal really is x) and meaningless in hardware, where no gate can detect x. Synthesis therefore reduces === to == or rejects it, and a tool must never let x-semantics leak into a netlist.

## Going deeper

IEEE Std 1800-2017, clause 11.4 (operators) and table 11-6; Sutherland, "I'm Still in Love with My X!" (DVCon 2013).
