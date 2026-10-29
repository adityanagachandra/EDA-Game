# Chiplets and UCIe

Universal Chiplet Interconnect Express (UCIe) is an open die-to-die interconnect standard. It specifies a physical layer, a die-to-die adapter and protocol mappings, so chiplets from different sources can be combined. A designer uses it as a set of fixed building blocks: modules of a defined lane count at a defined data rate, for standard or advanced packages. The arithmetic of lanes, rates and modules sets the bandwidth a chiplet edge can deliver.

## Lanes, rate and module bandwidth

UCIe organizes the physical layer into modules. The specification defines a standard-package module of 16 data lanes per direction for organic substrates and an advanced-package module of 64 lanes per direction for interposers and bridges, with data rates per lane from 4 up to 32 GT/s in its published versions. Raw bandwidth per direction is lanes times rate: an advanced module of 64 lanes at 16 GT/s carries $64 \times 16 = 1024$ Gb/s each way, before protocol overhead. A standard module at the same rate carries 256 Gb/s. The advanced module packs more lanes into a narrow edge because the microbump pitch is finer, which is the whole reason advanced packaging exists.

## Shoreline bandwidth

Modules are placed side by side along a die edge. If a module occupies $w$ um of edge, an edge of $E$ mm holds $\lfloor 1000E/w \rfloor$ modules, and the edge bandwidth is that count times the module bandwidth. With 500 um modules on a 10 mm edge, 20 modules at 1024 Gb/s give 20.48 Tb/s per direction. Dividing by the edge length gives the shoreline density, here about 2 Tb/s per millimetre. Modules can also be stacked in depth to raise bandwidth further at the cost of die area. These numbers drive architecture: an accelerator needing 40 Tb/s to a neighbouring chiplet needs two edges or a faster rate, and the floorplan must provide them.

## Sizing for a target

Bandwidth targets are often given in bytes per second, while link rates are in bits. To deliver $T$ TB/s per direction with modules of $B$ Gb/s, the module count is $\lceil 8000 T / B \rceil$. For 1.5 TB/s with 64-lane modules at 32 GT/s, $B = 2048$ Gb/s and the count is $\lceil 12000/2048 \rceil = 6$. Forgetting the factor of eight is the classic error. A real plan then adds protocol and retry overhead, reserves redundant lanes for repair, and checks power: every module draws power in proportion to its traffic, so the energy per bit becomes the next constraint, alongside the bump map the modules impose.

## Going deeper

The UCIe Consortium specification and white papers; Das Sharma and others, "Universal Chiplet Interconnect Express (UCIe): an open industry standard for innovations with chiplets at package level", IEEE Transactions on Components, Packaging and Manufacturing Technology (2022).
