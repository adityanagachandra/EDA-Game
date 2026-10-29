# The foundry and fabless model

Most chips today are designed by companies that do not own a factory. A fabless company designs chips and sells them; a foundry manufactures chips that others design. Between them sit EDA vendors, IP vendors and the PDK, the contract that lets a design team use a foundry's process without seeing inside it. This division of labour shapes what EDA tools must do and who pays for them.

## Integrated, fabless and foundry companies

An integrated device manufacturer (IDM) designs and fabricates its own chips. A fabless company designs chips and outsources fabrication, packaging and often test. A pure-play foundry, of which TSMC is the largest, manufactures for many customers and designs no competing products. Some IDMs also offer foundry services. The model became possible when design and manufacturing could be separated by a clean interface: design rules and device models delivered as a PDK. Outsourced assembly and test (OSAT) companies handle packaging and testing for fabless firms. For a fabless company the foundry choice is a major decision, because libraries, IP and signoff methodology are all tied to that foundry's process.

## The PDK as the interface

The foundry gives qualified customers a PDK for each process: device models, rule decks, extraction data and layer maps, plus foundry-validated standard-cell and memory libraries or references to third-party ones. EDA tools must be certified against those decks: a foundry qualifies specific tool versions for signoff on its process. The design team never sees the full process recipe, and the foundry never needs to see the RTL. The cost of entry is reduced to licences and mask sets, which made thousands of fabless companies possible. The PDK therefore works like an API between two companies: as long as both sides honour it, design and manufacturing can evolve independently.

## Economics and consequences for EDA

Building a leading-edge fab costs tens of billions of dollars, which only a few companies can fund. Fabless companies trade that capital cost for per-wafer prices and less control over the process. Shuttle runs (multi-project wafers) let many small designs share one mask set, lowering the cost of prototypes and of university tape-outs. For EDA, the model means tools must work with many foundries' rule decks and libraries, signoff is defined by the foundry's certified flow, and a new process node requires the tool, library and PDK to be developed together, years before volume production.

## Going deeper

Miller, "Chip War" (2022), for the industry history; foundry design-enablement documentation on PDKs and certified reference flows.
