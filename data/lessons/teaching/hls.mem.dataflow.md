# Task-level dataflow

Loop pipelining overlaps iterations of one loop. Dataflow overlaps whole tasks: a function body made of several sub-functions or loops runs as a pipeline of processes connected by channels, so the second task works on frame $k$ while the first starts frame $k + 1$. It is how HLS builds streaming architectures, and its throughput is set by the slowest task, not the sum of all of them.

## Tasks, channels and the slowest stage

In a dataflow region each task is a process with its own controller, and tasks communicate through channels: FIFOs for streamed data, or ping-pong buffers, two memories that the producer and consumer swap between frames. A task's interval is how often it can accept a new frame; for a non-pipelined task it equals its latency, and for a pipelined task it can be smaller. In steady state the region accepts a new frame only as fast as its slowest task, so the interval between frames at the output is the maximum task interval. With tasks of intervals 240, 310 and 130 cycles, frames emerge every 310 cycles. Speeding up any task other than the slowest one does nothing for throughput; the designer, or a tool, must find and fix the bottleneck task.

## Total time for many frames

A simple model of a dataflow region charges the first frame the sum of the task latencies, since it must pass through every task, and every further frame the maximum interval. For $N$ frames the total is $\sum_t L_t + (N - 1)\max_t I_t$. With latencies 270, 580 and 410 and intervals 270, 170 and 410, 101 frames take $1260 + 100 \times 410 = 42260$ cycles. Sequential execution without dataflow would take $N \sum_t L_t = 101 \times 1260 = 127260$. The formula assumes channels never fill and tasks never wait on each other except for data; real regions can stall when a FIFO is too shallow for a burst, which the tool's co-simulation reveals. It is the same structure as loop pipelining, one level up: frames are iterations and tasks are stages.

## Speedup and the rules that make it possible

The speedup of dataflow over sequential execution is $\frac{N \sum L_t}{\sum L_t + (N - 1)\max I_t}$, which tends to $\sum L_t / \max I_t$ for many frames. With latencies summing to 700 and a maximum interval of 90, 49 frames give $34300 / 5020 = 6.83$, approaching 7.8. Dataflow requires discipline in the code: each channel must have a single producer and a single consumer, data must flow forward without feedback between tasks, and tasks must not share arrays except through channels. Violations make the tool fall back to sequential execution or insert conservative synchronization. Channel depth is a design parameter: ping-pong buffers double memory, and FIFOs must be deep enough to absorb rate mismatches, or the region deadlocks or stalls.

## Going deeper

The dataflow optimization chapter of the AMD Vitis HLS user guide (UG1399); Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), chapter on streaming.
