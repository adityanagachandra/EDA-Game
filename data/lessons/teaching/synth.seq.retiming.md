# Retiming

Retiming moves registers across combinational logic without changing the circuit's input-output behavior. Leiserson and Saxe modeled a circuit as a graph with register counts on edges and showed that every legal retiming is described by one integer per node. Retiming can shorten the clock period or reduce the register count, and it is one of the few sequential optimizations used in production flows.

## The retiming graph and lags

Each combinational element is a node $v$ with a delay $d(v)$. Each connection is an edge $e: u \to v$ with a weight $w(e)$, the number of registers on it. A **retiming** assigns an integer lag $r(v)$ to each node, meaning $r(v)$ registers are moved from the outputs of $v$ to its inputs. The retimed weight of $e: u \to v$ is $w_r(e) = w(e) + r(v) - r(u)$. With $w(v_5 \to v_1) = 0$, $r(v_1) = 1$ and $r(v_5) = 0$, the edge gains one register: $w_r = 0 + 1 - 0 = 1$. A retiming is **legal** when every $w_r(e) \ge 0$. The register count of every cycle is unchanged, which is why the function is preserved.

## Register count after retiming

Summing $w_r$ over edges gives the new register count. Grouping by node, the change is $\sum_v r(v)(|FI(v)| - |FO(v)|)$, where $FI$ and $FO$ are the fanin and fanout edges. A node with equal fanin and fanout counts can move registers without changing the total; a node with one input and four fanouts trades four output registers for one input register when $r(v) = 1$. In a ring with one chord, nodes on the chord have unequal counts, so their lags change the total. Registers on edges leaving the same node can be shared physically, which the plain edge sum overcounts; minimum-area retiming models that sharing explicitly.

## Clock period of a retimed circuit

The clock period is the largest total delay along a path that crosses no register: the longest path in the subgraph of zero-weight edges, adding node delays. It is computed in linear time by a topological pass, since a legal circuit has no zero-weight cycle. In a four-node ring with delays 2, 1, 3, 3 where the two register-free edges are $v_4 \to v_1$ and $v_1 \to v_2$, the period is $3 + 2 + 1 = 6$. A retiming that puts the registers elsewhere can raise or lower it: the lags $r = (-1, -1, 0, -1)$ leave $v_3 \to v_4 \to v_1 \to v_2$ register-free, giving $3 + 3 + 2 + 1 = 9$. Choosing lags to minimize the period is the next problem.

## Going deeper

Leiserson and Saxe, "Retiming synchronous circuitry", Algorithmica (1991); Shenoy and Rudell, "Efficient implementation of retiming", ICCAD (1994).
