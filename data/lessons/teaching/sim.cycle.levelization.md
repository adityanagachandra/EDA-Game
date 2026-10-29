# Levelization of combinational logic

Levelization assigns every combinational gate a level so that evaluating gates in level order guarantees each gate sees settled inputs. It is the step that turns a netlist into a schedule: cycle-based simulators, compiled simulators such as Verilator, and parallel simulators all start from it. It is a longest-path computation on a directed acyclic graph, and its failure cases (combinational loops) are exactly the designs that need special handling in a fast simulator.

## Levels as longest paths from sources

Cut the design at its registers. The sources of the combinational graph are primary inputs and flip-flop outputs, which hold their values for the whole cycle; give them level 0. A gate's level is one more than the highest level among its inputs: $L(g) = 1 + \max_{u \in \text{fanin}(g)} L(u)$. Equivalently, $L(g)$ is the number of gates on the longest path from any source to $g$. Flip-flop inputs are sinks: they are read at the next clock edge, so they do not feed back into this cycle's graph. That cut is what makes a synchronous design's combinational logic acyclic, and it is why levelization is defined per cycle rather than over all time.

## Computing levels in topological order

Process gates in a topological order, so every gate's inputs are leveled before the gate itself. Example: sources $a$, $b$ and flop output $q$ at level 0; $n_1 = \text{AND}(a, b)$ is level 1; $n_2 = \text{OR}(n_1, q)$ is level 2; $n_3 = \text{XOR}(a, n_2)$ is level 3; $n_4 = \text{NAND}(n_1, b)$ is level 2. The design has 3 levels, with one gate at level 1, two at level 2 and one at level 3. The computation is linear in the number of edges, one pass of Kahn's algorithm or a depth-first post-order. The maximum level is the logic depth, which bounds how many sequential steps one cycle's evaluation needs even with unlimited parallelism.

## Why level order works and when it fails

If gates are evaluated in nondecreasing level order, every input of a gate was computed earlier in the same pass, so each gate is evaluated exactly once per cycle and its result is final. That removes the event queue entirely: no scheduling, no delta cycles, no repeated evaluation from glitches. Gates within one level are independent, so a level is also a unit of parallel work. The scheme fails when the graph is not acyclic. A combinational loop (through a mux path, a latch, or a false loop that is never sensitized) has no topological order. Simulators break such loops by iterating that strongly connected component until it settles, or by splitting signals so that the remaining graph is acyclic; Verilator warns about such cases with its UNOPTFLAT message.

## Going deeper

Abramovici, Breuer and Friedman, "Digital Systems Testing and Testable Design", chapter 3; Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms", the chapter on topological sort and DAG shortest paths.
