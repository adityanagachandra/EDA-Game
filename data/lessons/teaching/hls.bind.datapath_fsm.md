# The datapath and its FSM controller

The RTL an HLS tool emits has two halves. The datapath holds the functional units, registers and multiplexers chosen by allocation and binding. The controller is a finite-state machine that walks the schedule and, in every state, drives the multiplexer selects, register load enables and unit function selects. Knowing how one maps to the other lets a tool designer read generated RTL and predict where its critical paths and area go.

## Datapath and controller

The datapath computes and stores values: units, registers, multiplexers and the wires between them. The controller decides what the datapath does each cycle. In a statically scheduled design it is a Moore machine with one state per control step; in state $s$ it asserts the control word for step $s$: which source each multiplexer selects, which registers load, and which function each multi-function unit performs, such as add or subtract. A register that must capture a new value only in steps 2 and 5 has its load enable asserted in exactly those states, and holds its value otherwise. A loop whose body takes 4 steps becomes a cycle of 4 states with a conditional transition from the last state back to the first or out to the exit. The state register, next-state logic and output decoding are all controller logic.

## State count and encoding

A schedule of $S$ steps plus an idle state gives $S + 1$ states. Binary encoding needs $\lceil \log_2 (S + 1) \rceil$ flip-flops; one-hot encoding needs $S + 1$, one per state. A 16-step schedule with an idle state has 17 states, which binary encodes in 5 flip-flops and one-hot in 17. Binary saves flip-flops but every control signal must decode several state bits; one-hot makes each signal an OR of the few states that assert it, which is usually faster and is the common choice on FPGAs, where flip-flops are plentiful. Nested loops multiply control: each loop level adds its own states and counter, and HLS tools often generate a hierarchy of FSMs, one per loop or function, communicating through start and done signals.

## Where the controller shapes timing and area

Sharing moves cost from units into multiplexers and control: one adder shared by six additions needs input multiplexers and select signals decoded from the state. Heavy sharing adds multiplexer levels in front of every shared unit and lengthens its paths, which can force a slower clock. A large FSM with many states and high-fanout enable signals can itself become the critical path. Pipelined loops change the controller's shape: instead of one state per step of one iteration, the kernel repeats every II cycles, and stage-valid bits track which iterations are in flight during the prologue and epilogue. Variable-latency operations, such as a memory read that may stall, add wait states that loop until a valid or done signal arrives, so the schedule is no longer purely static.

## Going deeper

Gajski, Dutt, Wu and Lin, "High-Level Synthesis" (1992), chapter on control synthesis; De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 4.
