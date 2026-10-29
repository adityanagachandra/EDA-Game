# Dogleg channel routing

A dogleg breaks a net's trunk into pieces at its intermediate pin columns, letting each piece sit on a different track. Deutsch's dogleg router (1976) applies this systematically: split every multipin net into two-pin subnets and route the subnets with the constrained left-edge algorithm. Doglegs break many vertical-constraint cycles and shorten constraint chains, usually saving tracks.

## Splitting into subnets

A net with pins in $m$ distinct columns $c_1 < \dots < c_m$ becomes $m - 1$ subnets with spans $[c_1, c_2], [c_2, c_3], \dots$. Adjacent subnets of one net meet at a pin column, where a vertical branch can join them even if they are on different tracks.

Example: nets with 2, 3, 3, 2 and 3 distinct pin columns become $1 + 2 + 2 + 1 + 2 = 8$ subnets. Two-pin nets are unchanged. The split happens only at pin columns, which is the restricted dogleg of Deutsch's algorithm; unrestricted doglegs may jog at any column but enlarge the search.

## The subnet VCG

Vertical constraints move to subnets. At a column with net $a$ on top and net $b$ at the bottom, only the subnets of $a$ and of $b$ that cover that column are constrained, with an edge from each such subnet of $a$ to each such subnet of $b$. Subnets of a net that do not reach the column are free of the constraint.

Example: if net 3 must be above net 2 at column 4 and net 2 above net 1 at column 9, but net 2's subnet at column 4 is not the one at column 9, the chain $3 \to 2 \to 1$ falls apart into two shorter chains. A path in the subnet VCG maps to a path in the net VCG, so splitting can never lengthen the longest chain.

## What doglegs buy

The longest-chain lower bound drops when splitting separates constraints that hit different parts of a net. Example: a net VCG whose longest chain has 3 nets may have a subnet VCG whose longest chain has 2, saving a track. Splitting also breaks many cycles, so the constrained left-edge algorithm can finish where it would otherwise fail.

Doglegs cost vias (each jog between tracks is a pair of vias in a two-layer channel) and leave cycles through two-pin nets untouched. Density is unaffected: doglegs do not change how many nets cross a column. So the dogleg router narrows the gap between the chain bound and density, but cannot go below density.

## Going deeper

D. N. Deutsch, "A Dogleg Channel Router" (1976); Sherwani, "Algorithms for VLSI Physical Design Automation".
