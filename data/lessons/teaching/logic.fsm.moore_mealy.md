# Moore and Mealy machines

A finite-state machine (FSM) is the standard model of control logic: a state register, next-state logic and output logic. In a Moore machine the outputs depend only on the state; in a Mealy machine they depend on the state and the present inputs. The distinction decides when outputs change, how many states a behavior needs, and whether an output path is combinational from input to output, which matters for timing.

## The Moore machine

A Moore machine is a tuple $(S, I, O, \delta, \lambda, s_0)$ with next-state function $\delta: S \times I \to S$ and output function $\lambda: S \to O$. Each clock edge moves the state; the output is a function of the state alone, so it changes only after a clock edge and is glitch-free if $\lambda$ is registered.

Trace a machine with states A (output 0), B (output 0), C (output 1), with transitions A: $0 \to$ A, $1 \to$ B; B: $0 \to$ A, $1 \to$ C; C: $0 \to$ A, $1 \to$ C. From A, the inputs 1, 1, 0, 1, 1, 1 visit B, C, A, B, C, C, so the outputs after each edge are 0, 1, 0, 0, 1, 1: three 1s. This machine raises its output after two consecutive 1s.

## The Mealy machine

In a Mealy machine the output function is $\lambda: S \times I \to O$: outputs belong to transitions, not states. The output is a combinational function of the present input, so it responds in the same cycle the input arrives, one cycle earlier than an equivalent Moore machine.

A Mealy version of the two-1s detector needs only states A and B': from A, input 1 goes to B' with output 0; from B', input 1 stays in B' with output 1; input 0 always returns to A with output 0. On inputs 1, 1, 0, 1, 1, 1 it emits 0, 1, 0, 0, 1, 1, the same sequence as the Moore machine's outputs above, but each output is available during the cycle of its input rather than after the next edge. The price is a combinational path from input to output, which a timer must time through the FSM and which can glitch.

## Pattern detectors with overlap

Sequence detectors are the classic FSM exercise and a common verification target. A detector for pattern $p$ asserts its output whenever the last $|p|$ inputs equal $p$. With overlap allowed, occurrences can share bits: in the stream 1011011, the pattern 1011 occurs twice, at positions 1 to 4 and 4 to 7, sharing the middle 1.

To count occurrences, slide a window of length $|p|$ along the stream and count matches. The FSM does the same incrementally: its state records the longest prefix of $p$ that is a suffix of the inputs seen so far. A detector that resets to the start after each match counts only non-overlapping occurrences, and confusing the two is a classic specification ambiguity that a formal property check will expose.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 3.4; Hopcroft, Motwani and Ullman, "Introduction to Automata Theory, Languages, and Computation" (3rd ed., 2006), chapter 2.
