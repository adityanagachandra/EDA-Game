# Packing

Packing groups BLEs into clusters subject to the cluster's size $N$, its input-pin count $I$ and its clock limits. Good packing absorbs nets inside clusters, where fast local interconnect replaces global routing, and keeps related logic together for the placer. It is the step that turns a LUT netlist into the blocks that placement and routing see.

## Counting cluster inputs

A cluster needs one input pin for every distinct net that its BLEs read but that no BLE inside the cluster produces. Nets produced inside reach other BLEs through local feedback in the crossbar. The count is $|\bigcup_b \mathrm{in}(b) \setminus \{\mathrm{out}(b)\}|$ over the cluster's BLEs. If B1 reads $\{e_1, e_2, e_5, e_8\}$, B2 reads $\{e_6, e_7, e_9, o_1\}$ and B3 reads $\{e_1, e_3, e_4, e_6\}$, where $o_1$ is B1's output, the external nets are $e_1$ to $e_9$, so the cluster needs 9 inputs. The input limit $I$ is the usual binding constraint: a cluster often has free BLE slots but no free pins for a candidate's inputs.

## Legality of adding a BLE

Adding a candidate BLE is legal when a slot is free and the recomputed input count is at most $I$. The recount can go down as well as up. The candidate's new external inputs add pins, but if it produces a net the cluster currently imports, that pin is freed, because the net is now generated locally. A candidate that reads $\{e_7, e_9\}$ and drives $e_4$, added to a cluster importing $\{e_1, e_2, e_3, e_4, e_5, e_7, e_8, e_9\}$, needs no new pin and frees one: the count drops from 8 to 7. Packers check this incrementally for every candidate, which is why their data structures track, per cluster, which nets are imported, produced and absorbed.

## Absorbed nets

A net is absorbed when its driver and all its sinks lie in one cluster. It then uses only local interconnect and disappears from the inter-cluster routing problem. A BLE output used only inside the cluster and read there by another BLE is absorbed; one with any sink outside is not, even though it is also read locally. Absorption is the main benefit of clustering: the more nets a cluster swallows, the fewer the router must carry, and T-VPack's attraction function rewards shared nets for exactly this reason. The cluster's output-pin count rarely binds in this model, because every BLE output has its own cluster output; inputs are the scarce resource.

## Going deeper

Betz, Rose and Marquardt (1999), chapter 3; Marquardt, Betz and Rose, FPGA (1999); the VTR documentation for the VPR packer.
