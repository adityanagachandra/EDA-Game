# Through-hole versus surface-mount technology

Components attach to a board in two ways. Through-hole technology (THT) passes leads through plated holes and solders them, usually by wave or selective soldering. Surface-mount technology (SMT) places parts on lands printed with solder paste and melts the paste in a reflow oven. SMT dominates because it is denser and faster to place, while THT survives where mechanical strength matters. The choice changes routing capacity, assembly steps and failure modes, all of which board tools must model.

## Two ways to attach a part

A **through-hole** part's leads go through plated barrels; solder fills the hole and forms fillets on both sides, anchoring the part mechanically. Wave soldering passes the board's underside over a standing wave of molten solder; selective soldering does the same with a small nozzle per joint. A **surface-mount** part sits on lands: a stencil prints paste, a machine places parts on the wet paste, and a reflow oven melts it, with surface tension pulling parts onto their lands. SMT won because parts are smaller, both sides of the board can be populated, no hole is drilled per lead, and machines place tens of thousands of parts per hour. THT remains for large connectors, power parts and anything that takes mechanical stress.

## Consequences for layout and process

Each plated hole passes every layer, with a pad or an antipad on each, so a through-hole part costs routing space on all layers while a surface land occupies only the outer one. Mixed-technology boards need extra process steps. **Pin-in-paste** (intrusive reflow) prints paste over the hole and reflows the THT part with the SMT parts, avoiding wave soldering if the part survives reflow temperature. Wave soldering a board with bottom-side SMT parts requires gluing them and following orientation and spacing rules to avoid solder shadows and bridges. Package choice also changes inspection: a ball grid array hides its joints under the body, so it needs X-ray inspection, while leaded parts can be inspected optically.

## Assembly defects and their causes

Each technology has its characteristic failures, many traceable to layout. **Tombstoning** lifts a small chip part onto one end when one land wets first, from asymmetric land size, thermal mass (one land on a plane, one on a thin trace) or paste volume; balanced land design and thermal reliefs prevent it. **Bridging** at fine pitch happens when paste volume is large relative to the gap, so fine-pitch parts get thinner stencils. Moisture absorbed by plastic packages turns to steam in reflow and cracks them, which is why moisture-sensitivity levels are classified and parts are dry-stored or baked. THT joints carry insertion force through the barrel; SMT connector joints transfer it to the land-to-laminate bond, which can lift pads, so SMT connectors often add hold-down tabs.

## Going deeper

Prasad, "Surface Mount Technology: Principles and Practice" (2nd ed., 1997). IPC J-STD-020, on moisture sensitivity classification.
