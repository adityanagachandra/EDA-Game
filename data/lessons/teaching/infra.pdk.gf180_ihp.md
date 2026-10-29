# GF180MCU and IHP SG13G2

SKY130 is not the only open process. GF180MCU from GlobalFoundries and SG13G2 from IHP widen the choice: one aimed at mixed-voltage microcontroller-class designs, the other at SiGe BiCMOS for RF and high-speed analog. A tool builder who supports more than one PDK learns quickly which parts of a flow are process-neutral and which must be replaced, which is exactly the separation a good flow architecture enforces.

## GF180MCU for mixed-voltage designs

GF180MCU is a 180 nm process for microcontroller-class designs, released as an open PDK by GlobalFoundries. Its distinguishing trait is device options for several supply voltages, including higher-voltage devices useful for interfacing with 5 V signals, which the newer, thinner-oxide nodes cannot tolerate. It offers standard-cell libraries at more than one track height. Taller cells hold wider transistors and more internal routing room, so they drive larger loads and route more easily, but each cell costs more area. At 180 nm, digital density is low compared with modern nodes, so it suits control logic, sensor interfaces and mixed-signal glue rather than large datapaths. Moving a flow from SKY130 to GF180MCU breaks anything tied to the old site: row height, site name, tapcell spacing, PDN pitch and routing-layer names all change.

## IHP SG13G2 for SiGe BiCMOS and RF

IHP SG13G2 is a 130 nm SiGe BiCMOS process with an open PDK. BiCMOS means it offers both CMOS transistors and bipolar transistors on one die. The bipolar devices are silicon-germanium heterojunction bipolar transistors (HBTs), which switch far faster than the process's MOSFETs and are the reason the process exists: low-noise amplifiers, mixers, oscillators and high-speed drivers at frequencies well into the tens of gigahertz. Its metal stack includes thick top metals with low sheet resistance, which raise inductor quality factor and lower loss on RF lines and power routes. Because it is BiCMOS, it also has a CMOS standard-cell library, so a digital flow runs on it and RF and control logic can share a chip. The digital density is that of a 130 nm process, so a large digital block costs real area.

## Choosing and porting between open PDKs

Choose a PDK by its hardest requirement. Dense digital logic at modest speed: SKY130. Higher-voltage interfaces: GF180MCU. RF or very high-speed analog: SG13G2, accepting lower digital density. Then port the flow. Every physical and timing view comes from the PDK: technology LEF, cell LEF, Liberty per corner, cell GDS, SPICE models, and the DRC and LVS decks. On top of these, a flow keeps per-platform settings such as layer directions, PDN straps, tapcell rules, fill rules and the default routing layers. Well-structured flows keep these in a platform directory, separate from each design's configuration, which only names the platform. Then rerunning a design on another process is a one-line change in the design config, plus whatever the design itself must adapt to (a new clock target, a different IO ring).

## Going deeper

The GF180MCU and IHP SG13G2 open PDK repositories and their documentation; the OpenROAD-flow-scripts platform directories.
