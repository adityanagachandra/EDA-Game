# Depth-optimal LUT mapping

An FPGA implements logic in $K$-input lookup tables, each of which can compute any function of $K$ inputs. Mapping an AIG to LUTs is then a choice of cuts: every LUT is a $K$-feasible cut of some node. Depth-optimal mapping, first solved in polynomial time by FlowMap, labels each node with its minimum possible LUT depth.

## LUTs are cuts

A $K$-input LUT can implement node $n$ over any $K$-feasible cut $C$ of $n$, since the cone between $C$ and $n$ is just a function of at most $K$ inputs. A mapping chooses one cut per mapped node, starting from the outputs; the leaves of each chosen cut must themselves be mapped (or be primary inputs). The LUT count is area and the longest chain of LUTs is depth. Unlike cell mapping there is no library to match: any function fits, so the only question is which cuts to choose. With $K = 2$ each AND node is its own LUT and the LUT depth equals the AIG depth.

## Labeling for minimum depth

Process nodes in topological order. A primary input gets label 0. An AND node $n$ gets $\mathrm{label}(n) = \min_C (1 + \max_{v \in C} \mathrm{label}(v))$ over its nontrivial $K$-feasible cuts $C$. The label is the minimum depth of any LUT mapping of $n$'s cone, and the cut achieving it is recorded. The output labels give the depth of the mapping. FlowMap computes the same labels without enumerating cuts, by testing with a max-flow computation whether a cut of height $\mathrm{label} - 1$ exists. Cut-enumeration mappers, such as ABC's if, enumerate cuts explicitly, which also lets them optimize area among the depth-optimal choices.

## Depth gains and area recovery

With 3-input LUTs, an AIG of depth 3 often maps to depth 2, because a LUT absorbs several AND levels at once. Larger $K$ absorbs more levels per LUT, which is why 6-input LUT architectures map logic to far fewer levels than its AIG depth. Depth-optimal mapping alone tends to duplicate logic across overlapping cones, so mappers follow it with **area recovery**: among cuts that keep each node's required depth, prefer those with low area flow (area shared among fanouts) or low exact local area. Required times propagate from the outputs, so non-critical nodes get slack to choose cheaper cuts. This is the standard depth-then-area flow of modern FPGA mappers.

## Going deeper

Cong and Ding, "FlowMap: an optimal technology mapping algorithm for delay optimization in lookup-table based FPGA designs", IEEE TCAD (1994); Mishchenko, Chatterjee and Brayton, "Improvements to technology mapping for LUT-based FPGAs", IEEE TCAD (2007).
