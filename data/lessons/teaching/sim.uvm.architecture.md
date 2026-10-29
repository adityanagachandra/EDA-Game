# UVM architecture

The Universal Verification Methodology (UVM) is a SystemVerilog class library and a set of conventions for building reusable testbenches. Its value is structural: a standard component hierarchy, standard communication through transaction-level ports, and standard ways to configure and vary it. A tool designer meets UVM as the dominant workload of commercial simulators and as the vocabulary in which verification teams describe their environments.

## The component hierarchy

A UVM testbench is a tree of components. At the top is a **test**, which builds an **environment** and decides configuration and stimulus. The environment contains one **agent** per DUT interface, plus a **scoreboard** and coverage collectors. An active agent contains a **sequencer**, which arbitrates and supplies transactions; a **driver**, which turns each transaction into pin activity through a virtual interface; and a **monitor**, which observes the pins and publishes reconstructed transactions through an analysis port. Monitors feed the scoreboard and coverage, so checking sees exactly what happened on the pins. Every component knows only its own job and talks to others through ports, which is what lets components be swapped and reused.

## Active and passive agents and reuse

An agent can be **active**, with sequencer and driver creating traffic, or **passive**, with only a monitor. At block level, the agent for a bus interface is active: the testbench drives the DUT. At SoC level, the same interface may be driven by another block of real RTL; the agent is reconfigured as passive, and its monitor, protocol checks and coverage keep working unchanged. This vertical reuse is a central goal of the methodology: protocol knowledge is written once in an agent and carried from block to subsystem to chip. Coordinated traffic across several agents is produced by virtual sequences that start sequences on each agent's sequencer, rather than by wiring drivers together.

## Components versus objects

UVM separates **components** (derived from `uvm_component`) from **objects** (derived from `uvm_object`). Components form the static hierarchy: drivers, monitors, agents, environments, scoreboards. They are created during the build phase at time 0, have a place in the tree and a hierarchical name, and take part in phasing. Objects are transient data: transactions (`uvm_sequence_item`), sequences, configuration objects. They are created, randomized, passed around and discarded many times during a run. A component created late, for example with `new()` in the middle of the run phase, misses the build and connect phases and is never configured or connected. Keeping this distinction clear is what makes the hierarchy predictable.

## Going deeper

Accellera, "Universal Verification Methodology (UVM) 1.2 User's Guide" (2015); IEEE Std 1800.2-2020 (UVM class reference); Salemi, "The UVM Primer" (2013).
