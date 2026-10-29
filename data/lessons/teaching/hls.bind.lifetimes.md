# Value lifetimes and register count

Once operations are scheduled, every intermediate value that crosses a clock edge needs a register for some interval of steps: its lifetime. The number of registers a datapath needs is the largest number of lifetimes that overlap in any step. Computing lifetimes is the first step of register binding, and it is how a scheduling decision turns into register area.

## From a schedule to lifetimes

A value written at the end of step $d$ is first readable in step $d + 1$, and it must be held through its last read in step $u$. Its register lifetime is the interval of steps $[d + 1, u]$. A value written in step 1 and read in step 5 occupies a register during steps 2 to 5; a value written in step 3 and read in step 4 occupies one only during step 4. Values chained within a cycle never reach a register and have no lifetime at all. Primary inputs and outputs are usually held in interface registers and handled separately. Lifetimes depend on the schedule: moving a consumer earlier or a producer later shortens a lifetime, which is one reason schedulers sometimes prefer ALAP placement for operations that consume long-lived inputs.

## Counting live values

A value is live in step $t$ when $t$ lies in its lifetime. The number of values live in step $t$ is the register demand of that step, and the datapath needs at least the maximum demand over all steps. With lifetimes $[2, 5]$, $[8, 8]$, $[3, 5]$, $[4, 7]$ and $[7, 7]$, step 4 has three live values and so does step 5, while step 7 has two; the maximum is 3. For values that each live in one contiguous interval, this maximum is also achievable: the left-edge algorithm binds the values into exactly that many registers. That is a special property of intervals; values that wrap around a loop's back edge form circular arcs, and then the maximum overlap is only a lower bound.

## Values with several readers

A value read by several operations lives until its last reader, not its first. If $v_1$ is written in step 1 and read in steps 3 and 5, its lifetime is $[2, 5]$, even though one reader finished early. Fan-out therefore stretches lifetimes, and a scheduler that places one late consumer of a heavily used value can raise the register count of the whole block. Tools have two remedies: reschedule the late consumer earlier inside its frame, or recompute the value near the late consumer when the operation is cheap, which trades a unit cycle for a register. Loop-carried values live from their definition in one iteration to their use in the next, so in a pipelined loop their lifetimes are measured modulo the initiation interval, which is where rotating register files come from in software pipelining.

## Going deeper

Kurdahi and Parker, "REAL: a program for register allocation" (DAC, 1987); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 6.
