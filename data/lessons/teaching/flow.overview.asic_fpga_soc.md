# ASICs, FPGAs and systems on chip

Three product kinds dominate digital hardware. An application-specific integrated circuit (ASIC) is fabricated for one purpose. A field-programmable gate array (FPGA) is a fixed chip whose logic and wiring are configured after manufacture. A system on chip (SoC) integrates processors, memories, accelerators and interfaces on one die, and may be built as an ASIC or prototyped on an FPGA. Their economics decide which flow a team runs.

## Three ways to build digital logic

An ASIC is designed down to masks: the team owns every transistor's placement and pays for a mask set and fabrication. It gives the best power, performance and area, but mistakes are fixed only by a new mask set. An FPGA is an array of lookup tables, flip-flops, memories and DSP blocks with programmable routing. Its flow ends in a bitstream, not masks, so a bug fix is a reconfiguration. The programmability costs area, speed and power relative to an ASIC of the same function. An SoC is an architecture, not a fabric: it names what is integrated (CPU cores, caches, a network on chip, peripherals), and it is usually an ASIC built largely from IP blocks.

## Non-recurring and unit cost

The choice between ASIC and FPGA is a cost model. An ASIC has a large non-recurring engineering cost $F_A$ (design, verification, masks) and a small unit cost $u_A$. An FPGA has a small non-recurring cost $F_F$ and a larger unit cost $u_F$. Total costs are equal at the break-even volume
$$N^* = \frac{F_A - F_F}{u_F - u_A}.$$
With $F_A = 5{,}000{,}000$, $F_F = 200{,}000$, $u_A = 4$ and $u_F = 64$ (all in dollars), $N^* = 4{,}800{,}000 / 60 = 80{,}000$ units. Above that volume the ASIC is cheaper. Below that volume the FPGA wins.

## Choosing and combining

Volume is not the only input. Time to market favours FPGAs, since there is no fabrication wait. Power or performance limits favour ASICs. Standards that are still changing favour reprogrammability. Many teams do both: they prototype and emulate an SoC on FPGAs to run software early, then tape out the ASIC. Some products embed programmable logic inside an ASIC, or pair a processor SoC with FPGA fabric on one die. For an EDA engineer the consequence is two families of back-end tools: placement and routing onto a fixed FPGA fabric, and full-custom placement and routing onto silicon.

## Going deeper

Kuon and Rose, "Measuring the Gap Between FPGAs and ASICs", IEEE Transactions on CAD (2007).
