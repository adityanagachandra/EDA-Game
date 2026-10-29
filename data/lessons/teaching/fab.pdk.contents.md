# What a process design kit contains

A process design kit (PDK) is the foundry's description of a process, packaged for design tools. It contains device models, rules, technology files, layer definitions and parameterized device layouts, versioned together. Every EDA engine is configured from it. For a tool designer, the PDK is the boundary between algorithm and process: code that hard-codes a process fact breaks on the next node.

## The pieces of a PDK

A typical PDK ships: SPICE-level device models (BSIM-CMG for FinFETs, with corners and statistical variants); a design-rule manual and signoff decks for DRC, LVS, antenna and density; extraction rule files describing layer thicknesses, permittivities and resistivities; a layer map assigning names and purposes to GDS layer and datatype numbers; technology files for layout editors and routers, such as technology LEF; and parameterized cells (PCells) that generate device layouts from width, length and finger count. Standard-cell libraries, memories and I/O are built on top of the PDK, usually as separate IP. Customer design data is never part of it.

## Which tool reads which file

Each engine reads its own slice. Circuit simulators read device models. Routers read technology LEF: layer directions, pitches, widths, via definitions and a subset of spacing rules. Extractors read the interconnect technology description. DRC and LVS engines read rule decks: the LVS deck also contains device recognition rules that extract transistors from layer combinations. Layout editors read the layer map and display settings. Because GDS stores only integers, a wrong layer map silently mislabels data, so every tool must use the same one. Open PDKs let researchers publish full flows on real processes, which closed PDKs under non-disclosure agreements do not.

## Versions and maturity

PDKs evolve. Early versions use projected rules and pre-silicon models; later releases re-fit models to silicon and adjust rules as yield learning progresses. A change of one spacing value can make a previously clean layout fail, so design teams record the PDK version with every result and re-run signoff after updates. For tool developers this means treating PDK content as configuration: the same engine must run many processes and versions, read process facts from files, and validate that the files it reads are consistent, for example that the router's abstract rules agree with the signoff deck.

## Going deeper

The SkyWater SKY130 and GlobalFoundries GF180MCU open PDK documentation; the Si2 OpenPDK specifications.
