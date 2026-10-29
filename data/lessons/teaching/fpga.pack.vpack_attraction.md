# VPack and T-VPack attraction

VPack and T-VPack grow one cluster at a time: choose a seed BLE, then repeatedly add the unclustered BLE most attracted to the cluster until it is full or no legal candidate remains. VPack's attraction counts shared nets; T-VPack adds timing criticality. The attraction function is the policy of the packer and the easiest place to change its behavior.

## VPack attraction

VPack seeds a cluster with the unclustered BLE that has the most used inputs, then scores each candidate $B$ by the number of nets it shares with the cluster $C$: $\mathrm{Attraction}(B) = |\mathrm{Nets}(B) \cap \mathrm{Nets}(C)|$. A candidate touching $\{n_1, n_3, n_7, n_9\}$ next to a cluster touching $\{n_1, n_2, n_4, n_8, n_9, n_{12}\}$ shares $n_1$ and $n_9$, so its attraction is 2. Shared nets become absorbed or at least need fewer cluster pins, so maximizing them reduces both inter-cluster wiring and pin pressure. The greedy choice is fast, linear in the number of candidate pins per step, which is why it scales to large netlists.

## T-VPack attraction

T-VPack scores $\mathrm{Attraction}(B) = \lambda \cdot \mathrm{Crit}(B) + (1 - \lambda) \cdot |\mathrm{Nets}(B) \cap \mathrm{Nets}(C)| / G$, where $\mathrm{Crit}(B)$ is the criticality of the most critical connection between $B$ and the cluster, from a timing analysis of the unpacked netlist, and $G$ normalizes the shared-net term to $[0, 1]$, for example the most nets a BLE can touch. With $\lambda = 0.9$, criticality 0.05 and 2 shared nets with $G = 5$, the score is $0.9 \times 0.05 + 0.1 \times 2/5 = 0.085$. Seeds are also chosen by criticality, so the most critical paths are packed first, with their connections inside clusters where the local interconnect is fast.

## What the weighting changes

A large $\lambda$ lets criticality dominate, so a highly critical candidate with few shared nets can beat a non-critical one that shares many. With $\lambda = 0.75$, candidate X2 with criticality 0.95 and 1 shared net scores 0.762, ahead of X3 with criticality 0.6 and 2 shared nets at 0.5. Marquardt, Betz and Rose reported faster circuits from timing-driven packing than from share-only packing. The criticalities come from a pre-placement delay model, so they are estimates; packers that iterate with placement feedback, or that unpack and repack, address that. Because the attraction function is self-contained, it is the natural seam for new packing policies.

## Going deeper

Marquardt, Betz and Rose, "Using Cluster-Based Logic Blocks and Timing-Driven Packing to Improve FPGA Speed and Density", FPGA (1999); Betz, Rose and Marquardt (1999), chapter 3.
