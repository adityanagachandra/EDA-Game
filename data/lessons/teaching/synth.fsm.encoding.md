# State assignment for minimal logic

After minimization, each symbolic state gets a binary code, and that choice changes the next-state and output logic dramatically. Minimum-length binary codes save flip-flops; one-hot codes simplify logic and suit FPGAs; heuristic encoders place codes so that states that behave alike get adjacent codes and their logic merges into fewer, larger cubes.

## Flip-flop counts of common encodings

With $n$ states, a minimum-length binary encoding needs $\lceil \log_2 n \rceil$ flip-flops and one-hot needs $n$. For 21 states that is 5 against 21, a difference of 16. Gray-like encodings use the minimum length but give consecutive states adjacent codes, which helps counters and low-power designs. One-hot looks wasteful, but each next-state bit becomes a small OR of the transitions into that state, and decoding the current state is free, so the logic is often smaller and faster; FPGA flows favor it because flip-flops are plentiful. Unused codes of a compact encoding become don't cares only if they are unreachable.

## Adjacency cost of an assignment

Heuristic encoders such as MUSTANG compute an affinity weight $w_{ij}$ for each pair of states, for example from shared next states or shared outputs, and seek codes that make high-weight pairs adjacent. A simple objective is $\sum_{i<j} w_{ij} \cdot \mathrm{dist}(code_i, code_j)$ with Hamming distance. With codes S0 = 001, S1 = 000, S2 = 100, S3 = 111, S4 = 010 and weights $w(S0, S3) = 4$, $w(S1, S2) = 3$ among others, the pair S0, S3 at distance 2 costs 8 and the pair S1, S2 at distance 1 costs 3. Summing all weighted distances gives the cost to minimize. Adjacent codes let cubes for the two states merge in two-level minimization.

## Exact choice for small machines

For four states in two bits there are 24 assignments, but the cost depends only on which two pairs become diagonal (distance 2), since every other pair is at distance 1. There are just three ways to split four states into two diagonal pairs, so the minimum cost is the total weight plus the smallest combined weight of a diagonal split. Larger machines make exact assignment intractable, since it contains graph embedding problems, so tools use heuristics: KISS and NOVA derive face-embedding constraints from symbolic minimization, MUSTANG uses weighted adjacency, and Yosys's fsm pass extracts the machine and re-encodes it. The encoding decision interacts with retiming and with power, since codes determine flip-flop toggling.

## Going deeper

Devadas, Ma, Newton and Sangiovanni-Vincentelli, "MUSTANG: state assignment of finite state machines targeting multilevel logic implementations", IEEE TCAD (1988); De Micheli, Brayton and Sangiovanni-Vincentelli, "Optimal state assignment for finite state machines", IEEE TCAD (1985).
