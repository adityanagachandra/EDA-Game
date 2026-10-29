# The Verilog event scheduler

Verilog simulation semantics are defined by an event scheduler: processes wake on events, run, and schedule further events into ordered regions of each time step. Nonblocking assignment, races, testbench timing and the behavior of every compiled simulator follow from this model. A tool designer who builds a simulator, a linter or a compiled model must reproduce it or deliberately pick one of the orders it allows.

## Regions of a time step

Each simulation time step is processed in ordered regions. The active region holds ordinary work: blocking assignments, updates of continuous assignments, evaluation of the right-hand sides of nonblocking assignments and resumption of processes. The inactive region holds events delayed by #0. The NBA region holds the left-hand-side updates of nonblocking assignments. Finally, the monitor region (postponed in SystemVerilog, which adds more regions such as observed and reactive) runs monitor and strobe output.

The simulator empties the active region, moves inactive events in, and only when both are empty does it apply the NBA updates, which may wake more processes and refill the active region. This ordering is why a flip-flop modeled with <= samples old values: every right-hand side is read in the active region before any NBA update lands. It also explains why strobe shows settled values while display, which executes immediately, may show intermediate ones.

## Processes and sensitivity

Every always block, continuous assignment and gate is a process that waits for events. A continuous assignment assign y = a & b; re-evaluates when a or b changes value. An always block resumes when its event control fires. When no events remain at the current time, the simulator advances time to the earliest scheduled future event rather than stepping through empty time units, which makes event-driven simulation efficient when activity is sparse: only the fan-out of a change is evaluated.

An incomplete sensitivity list breaks the model. With always @(a) y = a & b;, a change on b alone does not wake the block, so y keeps its old value. Synthesis ignores sensitivity lists and builds an AND gate, so simulation and hardware disagree. Writing always @(*) or always_comb makes the simulator derive the list from the code and removes this class of mismatch.

## Races

Processes triggered by the same event execute in an order the standard leaves unspecified. Two clocked blocks, one doing a = b; and the other b = a; with blocking assignments, may end with both holding a's old value or both holding b's: the result is a race, and different simulators may legitimately differ. The standard fix is a coding rule: use nonblocking assignments for every variable written in a clocked block and read by another, so every process reads pre-edge values and order no longer matters.

Testbenches race too. If a testbench drives an input with a blocking assignment at the same posedge where the design samples it, one simulator may capture the new value and another the old. Driving on the opposite edge, using nonblocking assignments or using SystemVerilog clocking blocks fixes it. Compiled simulators such as Verilator choose one static order, which is just one of the legal ones, so racy code is a portability bug, not a simulator bug.

## Going deeper

IEEE Std 1800-2017, clause 4 (scheduling semantics); Cummings, "Nonblocking Assignments in Verilog Synthesis, Coding Styles That Kill!" (SNUG 2000).
