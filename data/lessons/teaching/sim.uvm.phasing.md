# UVM phasing and objections

A UVM testbench is a program with many cooperating components, and they must be built, connected, run and shut down in a coordinated order. UVM phasing is that order. It also defines when the test ends, through objections. Most UVM failures that look like "the test did nothing" or "the test never ends" are phasing errors, and the traversal rules are a small application of tree traversal over the component hierarchy.

## The phase sequence and traversal order

The common phases are, in order: **build**, **connect**, **end_of_elaboration**, **start_of_simulation**, **run**, **extract**, **check**, **report** and **final**. All except run are functions that take no simulation time. Build executes **top-down**: a parent's build_phase runs before its children's, because the parent creates and configures its children there, a pre-order traversal of the tree. Connect and the later function phases execute **bottom-up**, like a post-order traversal; by then the whole tree exists, so every port a component connects is already built. The run phase is a task, started in parallel in every component, and it is the only common phase that consumes simulation time.

## Run-time phases and objections

The run phase ends when no component holds an **objection**. A test raises an objection before starting its main sequence and drops it when the sequence completes; drivers and scoreboards may also object while work is pending. When the count of raised objections reaches zero, after any configured drain time, run ends and the cleanup phases follow. If no objection is ever raised, the run phase ends at time 0 and the test finishes before its sequence does anything. If an objection is never dropped, the test hangs until a timeout. In parallel with run_phase, UVM also schedules run-time sub-phases (reset, configure, main, shutdown and their pre and post variants), each ending by the same objection mechanism.

## Using the phases correctly

Each kind of work has a phase. Read configuration with `uvm_config_db::get` in build_phase, before creating or configuring children. Create components only in build_phase, never in connect or run. Make port connections in connect_phase. Do final checks, such as reporting expected transactions that never arrived, in check_phase, after all traffic has ended. Report summaries in report_phase. Because function phases cannot block, setup is guaranteed to complete before simulation time advances. The phase schedule is also where a simulator's performance shows: build of a large hierarchy can take noticeable time, and the run phase is where the DUT and the testbench actually compete for CPU.

## Going deeper

Accellera, "UVM 1.2 User's Guide" (2015), chapter on phasing; IEEE Std 1800.2-2020, the uvm_phase and uvm_objection classes.
