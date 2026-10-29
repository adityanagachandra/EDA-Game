# Pipeline hazards and stalls

In a processor pipeline, instructions depend on each other, and a dependent instruction may need a value its predecessor has not yet written. These data hazards, together with structural and control hazards, force the pipeline to stall or to forward values. Counting stall cycles for a code sequence is the core exercise of pipeline design, and the interlock and forwarding logic is a classic target for formal verification.

## The five-stage pipeline and data hazards

The classic RISC pipeline has stages IF (fetch), ID (decode and register read), EX (ALU), MEM (data memory) and WB (register write). Instruction $j$ reads its sources in ID; instruction $i$ writes its result in WB, three cycles after its own ID. With the register file written in the first half of a cycle and read in the second, a reader may sit in ID in the same cycle as the writer is in WB.

Without forwarding, a dependent instruction must therefore reach ID no earlier than 3 cycles after its producer's ID. Back-to-back dependent instructions (distance 1) stall 2 cycles; at distance 2 they stall 1; at distance 3 or more, none. A stall delays the instruction and everything behind it, which changes later distances.

## Forwarding and the load-use hazard

Forwarding routes the ALU result from the EX/MEM or MEM/WB pipeline register straight to the ALU inputs, so a dependent ALU instruction can use a result in the cycle right after it is computed, with no stall. A load is different: its data appears only at the end of MEM, so an instruction that uses a load result immediately must stall 1 cycle (the load-use hazard) and then receive the value by forwarding.

Count total cycles for `lw r1, 0(r2)`, `add r3, r1, r4`, `sub r5, r3, r1`. With forwarding: the add stalls 1 cycle on r1; the sub then forwards r3 with no stall. Three instructions take $5 + 2 = 7$ cycles without stalls, plus 1: 8 cycles. Without forwarding: the add stalls 2 cycles, and the sub, now at distance 1 from the add, stalls 2 more, giving 11 cycles.

## Structural and control hazards

A structural hazard arises when two instructions need the same hardware in the same cycle: one memory port shared by fetch and a load, or a register file with too few write ports. The fix is duplication (separate instruction and data caches) or a stall.

A control hazard arises at branches: the pipeline does not know the next fetch address until the branch resolves. If it resolves in EX, two wrongly fetched instructions must be flushed when the branch is taken, a 2-cycle penalty. Hazard detection is a small amount of logic with a large state space (every combination of instructions in flight), which is why processor verification uses formal methods on the interlock and forwarding control, and why compilers schedule independent instructions into the delay slots that hazards create.

## Going deeper

Patterson and Hennessy, "Computer Organization and Design" (RISC-V ed., 2017), sections 4.6 to 4.8; Hennessy and Patterson, "Computer Architecture: A Quantitative Approach" (6th ed., 2017), appendix C.
