# Specification, architecture and RTL

The front of the flow decides what the chip is before any tool decides how to build it. A specification states function, interfaces and targets; an architecture chooses the blocks, the memory system and the pipelines that meet those targets; RTL writes each block as registers and the logic between them. Decisions made here bound everything downstream, and no back-end tool can recover performance that the architecture did not leave room for.

## The specification

A specification fixes the function (instruction set, protocol, algorithm), the interfaces (bus standards, pins, clocks), and the targets: clock frequency, throughput, latency, power budget, area or cost, and reliability requirements. It is the reference that verification checks against, so ambiguity here becomes bugs later. Good specifications are partly executable: a reference model in C++ or Python that produces the expected outputs, from which tests and later formal properties are derived. Targets need margins and priorities: when frequency and power conflict later, the specification should say which one wins. Reliability, safety and security requirements, such as error-correcting memories or a secure boot path, also start here.

## Architecture: blocks, pipelines and budgets

The architect picks the micro-architecture: how many pipeline stages, which memories and how wide, which functions get dedicated hardware, which clock domains exist. Each choice trades power, performance and area. Adding a pipeline stage shortens the logic between registers so the clock can be faster, at the cost of latency, registers and power. The architect also writes budgets: each block gets an area, a power share and timing at its ports. These budgets are what floorplanning and timing closure later try to honour. Architects evaluate these choices with performance models long before RTL exists.

## RTL as the hand-off

RTL in Verilog, SystemVerilog or VHDL describes behaviour per clock cycle: what each register holds after the edge, as a function of inputs and current state. It is precise about cycles and silent about gates. Writing RTL is design, not transcription: the coding style decides whether synthesis infers a multiplexer or a priority chain, whether a memory becomes an SRAM macro or flip-flops, and whether the logic between two registers is short enough to meet the clock. Lint and clock-domain-crossing checks run on RTL before synthesis, because a structural problem found here costs minutes; found after layout it costs weeks.

## Going deeper

Hennessy and Patterson, "Computer Architecture: A Quantitative Approach" (6th ed., 2017); Harris and Harris, "Digital Design and Computer Architecture".
