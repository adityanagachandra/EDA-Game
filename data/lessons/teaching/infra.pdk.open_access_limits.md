# What open access does not give you

Open PDKs and open tools let anyone design a chip and verify it. They do not, on their own, turn that design into a product. Knowing exactly where the open release stops matters for anyone who builds tools on it, because a tool that reports clean DRC or positive slack on an open PDK is reporting evidence under assumptions, not a guarantee of working, manufacturable, qualified silicon.

## What the open PDK files do not include

An open PDK releases rules, device models, layer data and basic libraries. It usually does not release the foundry's production yield data, statistical correlation between models and silicon lots, or a portfolio of complex, silicon-qualified IP: high-speed SerDes, PLLs, qualified SRAM compilers, and certified IO and ESD libraries. Some of these exist as open projects, for example open SRAM generators, but their models may not be validated across corners. Open DRC and LVS decks may also lag or simplify the foundry's internal signoff decks, so a layout clean under the open deck may still violate a rule the foundry checks. Treat open results as strong evidence within the encoded rules and models. When a design depends on something the release lacks, the risk is the designer's, and a tool should say which inputs were assumed.

## From open files to fabricated silicon

A verified GDS is not a chip. Mask making and wafer processing happen at a mask shop and the foundry, and they need a fabrication run. Open designs usually reach silicon through multi-project shuttles: many small designs share one mask set and one wafer run, which divides the dominant fixed cost, the masks, among all participants. Each project gets a handful of packaged or bare dies back. Prototypes that work at room temperature prove the design concept. A product still needs reliability qualification across temperature and voltage, a production test program, packaging qualification, and a supply agreement for volume. None of these comes from open files. A tool builder working on open flows should be clear that the flow ends at a verified handoff; the rest of the path is a business and engineering process with the foundry.

## Engineering around the gaps

Good practice on an open PDK is conservative. Reuse the IO cells the PDK provides rather than drawing new pad circuits, since ESD design depends on silicon data you do not have. Add timing margin beyond the reported slack: if the libraries have not been correlated to measured silicon, model error is unknown, and a guard band of a few percent of the clock period is cheap insurance. Choose modest clock targets and layouts that stay well inside the characterized range. Cross-check physical verification with two independent engines, such as Magic and KLayout. For research that should inform a commercial product, carry over methods and relative comparisons (algorithm A beats algorithm B by 4% wirelength on the same process) rather than absolute numbers, which do not transfer between processes and decks.

## Going deeper

The documentation of open shuttle programs and of the SKY130, GF180MCU and SG13G2 PDKs; Weste and Harris, "CMOS VLSI Design" (4th ed.), chapters on manufacturing and test.
