# cocotb and open-source simulators

cocotb is an open-source framework for writing testbenches in Python. It does not simulate anything itself: it attaches to an HDL simulator through that simulator's procedural interface and runs Python coroutines alongside the simulation. With Icarus Verilog, Verilator or GHDL underneath, it gives a complete open-source verification flow, and its design exposes clearly how a testbench synchronizes with a simulator's scheduler.

## Python coroutines over the simulator interface

cocotb loads into the simulator as a shared library. For Verilog simulators it uses VPI, the Verilog procedural interface; for VHDL it uses VHPI or FLI, depending on the simulator. Through that interface it reads and writes signals and registers callbacks for time and value changes. Tests are Python `async` functions. A test drives `dut.a.value = 5`, then waits on a **trigger**: `await RisingEdge(dut.clk)` suspends the coroutine until the simulator reports that edge, and `await Timer(10, 'ns')` until 10 ns of simulated time pass. While Python waits, the simulator runs; when the callback fires, cocotb resumes the coroutine. Many coroutines can run concurrently (a clock generator, a driver, a monitor), interleaved cooperatively with simulated time.

## The open-source simulators underneath

**Icarus Verilog** compiles Verilog into an intermediate form run by its `vvp` runtime, an event-driven, four-state simulator with delays, good for debugging RTL behavior. **Verilator** compiles synthesizable SystemVerilog into a C++ model; it is usually much faster but two-state, stricter about constructs, and statically scheduled, so code that runs in Icarus may need changes. **GHDL** simulates VHDL. The same cocotb test can often run on more than one of them, which is useful in itself: a difference between simulators points to a race, an X dependence or a tool bug. Choosing among them is the usual trade-off between debug fidelity (four-state, delays) and throughput (compiled, two-state).

## Reading values at the right moment

A coroutine resumed by `RisingEdge(dut.clk)` runs at the edge, before the simulator has processed the nonblocking updates of that time step. Reading a register output there returns the value from before the edge, which surprises new users and is exactly the scheduling-region behavior of the simulator. Awaiting `ReadOnly()` resumes at the end of the time step, when values have settled, but writes are not allowed there; `ReadWrite()` gives a point where writes are still permitted. The other cost is speed: each crossing between simulator and Python, and every Python statement executed per cycle, adds overhead that can dominate run time. Efficient cocotb testbenches work at transaction level and keep per-cycle Python work small.

## Going deeper

The cocotb documentation (triggers, coroutines and the simulator support pages); the Icarus Verilog, Verilator and GHDL user manuals.
