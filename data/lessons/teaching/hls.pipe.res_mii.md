# The resource-constrained minimum II

ResMII is the lowest II that the available resources allow, whatever the dependences. In every window of II cycles, each resource must serve all of one iteration's uses of it. It is one of the two lower bounds every modulo scheduler computes first, and it explains most II violations an HLS report shows, especially those caused by memory ports.

## One resource type

In steady state a pipelined loop starts one iteration every II cycles, so every II cycles each resource must perform one iteration's worth of work. If an iteration uses resource $r$ for $u_r$ cycles in total and there are $k_r$ copies, then $k_r \cdot \mathrm{II} \ge u_r$, so $\mathrm{II} \ge \lceil u_r / k_r \rceil$. Ten multiplications per iteration on four fully pipelined multipliers give $\lceil 10 / 4 \rceil = 3$: at II = 2 the multipliers would offer only 8 slots for 10 uses. For a non-pipelined unit each use occupies it for its whole latency, so $u_r$ is the number of uses times the occupancy. The bound depends only on counts, so it is computed before scheduling and costs nothing.

## ResMII over all resources

Every resource gives its own bound, and ResMII is the maximum: $\mathrm{ResMII} = \max_r \lceil u_r / k_r \rceil$. With 6 multiplications on 4 multipliers, 2 additions on 4 adders and 1 division on 1 divider, the bounds are 2, 1 and 1, so ResMII = 2. The resource with the largest bound is the binding resource; adding copies of any other resource does nothing for the II. That is why HLS reports name the limiting resource. The bound can be loose when uses of different resources must line up in particular cycles, and modulo scheduling may then need a larger II, but in practice ResMII is usually achieved when it dominates. Reducing it means more copies of the binding resource, or fewer uses through algebraic rewriting or operator sharing across iterations.

## Memory ports as resources

Memories are the most common binding resource. A block RAM has one or two ports, each performing one read or one write per cycle. An iteration that reads array A six times and writes it twice makes 8 accesses; in a single-port RAM that forces $\mathrm{II} \ge 8$, in a dual-port RAM $\mathrm{II} \ge 4$. If the same iteration needs only $\lceil 8 / 2 \rceil = 4$ cycles of multiplier time, memory is still the binding resource. The remedies are array partitioning, which splits A into banks with their own ports, reshaping, which packs several elements into one wide word, and restructuring the code so that values are read once and reused from registers. A report saying the II is limited by a memory port is therefore a request to change the memory architecture, not the arithmetic.

## Going deeper

Rau, "Iterative modulo scheduling: an algorithm for software pipelining loops" (MICRO-27, 1994); Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018).
