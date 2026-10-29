# Latches and flip-flops

Sequential circuits remember. The two basic storage elements are the latch, which is transparent while its enable is active, and the edge-triggered flip-flop, which samples its input only at a clock edge. Every register, state machine and pipeline is built from them, and every timing check in static timing analysis is a statement about when a flip-flop samples. A flip-flop's behavior is a recurrence: the next state is a function of the present state and the inputs.

## Characteristic equations

Write $Q$ for the present state and $Q^+$ for the state after the next active clock edge. The common flip-flop types are defined by their characteristic equations:

- D: $Q^+ = D$.
- T (toggle): $Q^+ = T \oplus Q$.
- JK: $Q^+ = J\bar{Q} + \bar{K}Q$. So $J K = 00$ holds, 10 sets, 01 resets and 11 toggles.
- SR: $Q^+ = S + \bar{R}Q$, with $S = R = 1$ forbidden.

A JK flip-flop at $Q = 1$ with $J = 1, K = 1$ goes to $Q^+ = 1 \cdot 0 + 0 \cdot 1 = 0$: it toggles. In practice ASIC libraries contain D flip-flops almost exclusively; the other types are built by putting their characteristic equation in front of a D input.

## Iterating over clock edges

Because $Q^+ = f(Q, \text{inputs})$, a flip-flop's history over $k$ edges is $k$ iterations of a recurrence. Start a JK flip-flop at $Q = 0$ and apply $(J, K) = (1, 0), (1, 1), (0, 0), (1, 1), (0, 1)$ on five successive edges. The states after each edge are 1 (set), 0 (toggle), 0 (hold), 1 (toggle), 0 (reset). The final state is 0, and $Q$ was 1 after two of the five edges.

Simulators implement exactly this: at each active edge, compute every flip-flop's next state from the *old* values of all states, then update them all at once. That simultaneous update is what Verilog's nonblocking assignment `<=` models.

## Level-sensitive latches versus edge-triggered flip-flops

A D latch with enable $G$ is transparent while $G = 1$ ($Q$ follows $D$) and holds while $G = 0$. A positive-edge D flip-flop is two latches in series, a master transparent while the clock is low and a slave transparent while it is high, so the output changes only at the rising edge, to the value $D$ had just before it.

The difference is visible when $D$ changes while the clock is high: the latch passes every change through, the flip-flop ignores it until the next rising edge. Latches are smaller and allow time borrowing (a slow stage can use part of the next phase), which some high-performance designs exploit. They also make timing analysis harder, because a path through a transparent latch does not stop there. Unintended latches, inferred from an incompletely specified `always` block, are a standard lint warning.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 3; Weste and Harris, "CMOS VLSI Design" (4th ed., 2010), chapter 10 on sequential circuits.
