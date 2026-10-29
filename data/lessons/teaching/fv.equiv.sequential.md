# Sequential equivalence checking

When synthesis retimes registers, changes state encodings or removes redundant state, the two designs no longer have matching registers, and combinational checking cannot compare them. Sequential equivalence checking (SEC) asks the harder question: starting from their reset states, do the two machines produce the same outputs for every input sequence? It is a model-checking problem on the product of the two machines, and its basic quantities are reachable product states and the depth of the first difference.

## The product machine

Run both machines side by side on the same input stream. A product state is a pair (state of $A$, state of $B$), the initial product state is the pair of reset states, and on each input both machines step together. The machines are equivalent from reset if every reachable product state has equal outputs (for Moore machines; for Mealy machines, equal outputs on every input). To count reachable product states, search breadth-first from the reset pair. Example: $A$ has states A0 (output 0) and A1 (output 1), moving to A1 on input 1 and to A0 on input 0; $B$ has B0 (output 0, on 0 stay, on 1 to B1), B1 (output 1, on 0 to B2, on 1 stay) and B2 (output 0, on 0 to B0, on 1 to B1). From (A0, B0) the search reaches (A1, B1) and (A0, B2): 3 reachable pairs, all with equal outputs, so the machines are equivalent.

## Shortest distinguishing sequence

If some reachable product state has different outputs, the machines are not equivalent, and the breadth-first distance to the nearest such pair is the length of the shortest distinguishing input sequence. That sequence is the counterexample SEC reports, and BMC on the product machine finds it at exactly that depth. In the example above, change B2's output to 1. The pair (A0, B2) is reached by the inputs 1 then 0, and there $A$ outputs 0 while $B$ outputs 1, so the shortest distinguishing sequence has length 2. No pair at distance 1 differs: from the reset pair, input 0 stays at (A0, B0) and input 1 reaches (A1, B1), both with equal outputs. Breadth-first order matters: the first differing pair found is the nearest one.

## Sequential equivalence on larger machines

In practice SEC is a safety property on the miter of the two sequential designs: assert that the outputs are always equal, and run BMC to find differences and induction or IC3 to prove equality. On larger machines the computation is the same, only bigger: breadth-first search over pairs until a differing pair appears or the reachable set closes. Plain induction usually fails because unreachable pairs of states can differ, so SEC engines strengthen it with candidate equivalences between registers of the two designs, found by simulation and proved inductively together (the van Eijk method), much as SAT sweeping does combinationally. Retiming moves registers across logic, and re-encoding (for example one-hot versus binary) changes the meaning of every register, so the useful equivalences then relate signals rather than registers. Reset matters: designs may be equivalent only from reset, not from arbitrary states.

## Going deeper

van Eijk, "Sequential Equivalence Checking Based on Structural Similarities" (IEEE TCAD, 2000); Mishchenko, Case, Brayton and Jang, "Scalable and Scalably-Verifiable Sequential Synthesis" (ICCAD 2008).
