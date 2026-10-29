# Basic logic elements and clusters

A basic logic element (BLE) pairs a LUT with a flip-flop and a multiplexer that chooses the combinational or registered output. FPGAs group N BLEs into a cluster (a CLB, slice or logic array block) with a fast local interconnect. Cluster size and input count trade logic density against routing demand, and every packer, placer and area model in VPR is built on these parameters.

## The BLE and its configuration

A classic BLE has a $k$-LUT, a D flip-flop fed by the LUT, and a 2:1 output multiplexer whose select bit chooses the LUT output (combinational) or the flip-flop output (registered). The configuration of one BLE is therefore $2^k$ LUT bits plus one select bit. A cluster of $N$ BLEs holds $N(2^k + 1)$ such bits; for $N = 10$ and $k = 6$ that is $10 \times 65 = 650$ bits, before the local interconnect and the global routing are counted. Real logic elements add flip-flop clock-enable and set or reset controls, fracturable LUTs and carry logic, but the BLE abstraction is what VPR's architecture files describe first and what most FPGA CAD research assumes.

## Sizing the cluster inputs

A cluster with $N$ BLEs of $k$-LUTs has $kN$ LUT inputs, but it does not need $kN$ cluster input pins. BLEs share signals, and BLE outputs feed back locally. Betz and Rose found that about $I = \frac{k}{2}(N + 1)$ inputs give close to full logic utilization. For $k = 4$, $N = 4$ this is $I = 10$, not 16. Inside, a fully connected local crossbar lets each of the $kN$ LUT inputs select any of the $I$ cluster inputs or $N$ feedback outputs, which costs $kN(I + N)$ switches: for $k = 4$, $N = 4$, $I = 10$, that is $16 \times 14 = 224$. Fewer cluster inputs cut the global connection-box area; a sparser crossbar cuts local area but makes packing harder.

## Cluster counts and device sizing

Once LUTs and flip-flops have been paired into BLEs, a lower bound on the clusters a design needs is $\lceil B / N \rceil$ for $B$ BLEs. Architecture experiments often size the smallest square array that fits: the side $n$ is the least integer with $n^2 \ge \lceil B/N \rceil$. For 2616 BLEs in clusters of 10, the design needs at least 262 clusters, and since $16^2 = 256 < 262 \le 289 = 17^2$, the array is $17 \times 17$. The bound is reached only if input-pin limits never stop a cluster from filling. Larger $N$ absorbs more nets inside clusters and so reduces inter-cluster routing, but each cluster's local crossbar grows roughly quadratically, which is why commercial cluster sizes settled around 8 to 10.

## Going deeper

Betz and Rose, "How Much Logic Should Go in an FPGA Logic Block?", IEEE Design and Test (1998); Betz, Rose and Marquardt (1999), chapter 3; the VTR architecture-description documentation.
