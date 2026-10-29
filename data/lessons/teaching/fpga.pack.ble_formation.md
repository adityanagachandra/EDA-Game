# Pairing LUTs and flip-flops into BLEs

Packing begins by pairing each flip-flop with the LUT that drives it, so the pair shares one basic logic element. How many pairs form depends on the BLE's outputs: whether the LUT output and the flip-flop output can both leave the element. The pairing sets the BLE count, which is the input to clustering and the first area number the flow can trust.

## Pairing in a single-output BLE

In the classic BLE the output multiplexer selects either the LUT output or the flip-flop output, so only one of them leaves the element. A LUT and a flip-flop can therefore share a BLE only when the LUT's sole sink is that flip-flop's D input; otherwise the LUT's other sinks would have no way to read it. Every LUT and every flip-flop that cannot pair needs a BLE of its own. The count is $B = L + F - P$, for $L$ LUTs, $F$ flip-flops and $P$ legal pairs. With 5 LUTs and 3 flip-flops, where one LUT drives only a flip-flop, one drives a flip-flop and two other gates, and the third flip-flop comes from an input pin, $P = 1$ and $B = 7$.

## From BLEs to clusters

Clusters hold $N$ BLEs, so a lower bound on the cluster count is $\lceil B/N \rceil$. With $B = 7$ and $N = 2$, at least 4 clusters are needed. The bound is reached only when cluster input pins never block a BLE from joining, which is optimistic for large clusters with few inputs. Pairing first matters because a missed pair costs a whole BLE: a flip-flop alone in a BLE wastes its LUT, and a LUT alone wastes its flip-flop. Packers such as VPack and T-VPack form BLEs in this first step and only then grow clusters, so the BLE count is fixed before any clustering heuristic runs.

## Elements with two outputs

Modern logic elements expose both the LUT output and the flip-flop output, so a LUT can share an element with a flip-flop it drives even when it has other sinks. The pairing rule relaxes to: each LUT that drives at least one flip-flop absorbs one of them. With 6 LUTs and 2 flip-flops, each flip-flop driven by a LUT that also has other sinks, the single-output rule gives 0 pairs and 8 elements, while the two-output rule gives 2 pairs and 6 elements. Commercial elements go further, with a flip-flop input that bypasses the LUT and several flip-flops per LUT. The architecture description must state these options precisely, because the packer's legality check is only as good as its model of the element.

## Going deeper

Betz, Rose and Marquardt (1999), chapter 3 on VPack; Marquardt, Betz and Rose, "Using Cluster-Based Logic Blocks and Timing-Driven Packing to Improve FPGA Speed and Density", FPGA (1999).
