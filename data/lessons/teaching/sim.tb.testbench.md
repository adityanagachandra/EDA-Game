# Testbench architecture

A testbench is the program that exercises a design under test (DUT) and decides whether it behaved correctly. Modern testbenches are layered software systems: stimulus is generated as transactions, converted to pin activity, observed, and checked automatically against a model. The architecture determines how much of the testbench survives from one test to the next and from block level to chip level, and its signal-level discipline determines whether results are deterministic.

## Components of a layered testbench

A **stimulus generator** produces transactions such as "write 0x40 to address 0x1000". A **driver** converts each transaction into signal activity on the DUT interface, following the protocol cycle by cycle. A **monitor** watches the same interface passively and reconstructs transactions from pin activity. A **reference model** predicts what the DUT should produce for each input transaction. A **scoreboard** (or checker) compares observed outputs against the predictions, and a **coverage collector** records which scenarios occurred. Only the driver and monitor know the pin protocol; everything above them works on transactions. This split means a new test usually changes only the stimulus, while checking and coverage remain the same and apply to every test.

## Directed, self-checking and random testbenches

A **directed** test applies a hand-written sequence for one feature and often checks its own expected values. Directed tests are precise but only find bugs someone anticipated, and their number grows with the feature list. A **self-checking** testbench checks results automatically, so a regression of thousands of runs reports pass or fail without anyone reading waveforms; checking lives in the environment, not in each test. A **constrained-random** testbench generates legal but random transactions, reaching combinations no one wrote down, and depends on the self-checking environment to judge them and on coverage to show what was reached. In practice teams combine all three: random tests for breadth, directed tests for specific corners and bring-up, with one shared checking and coverage environment.

## Driving and sampling without races

The testbench and the DUT run in the same simulator, so their interaction at clock edges is subject to the scheduling rules. If a driver changes a DUT input with a blocking assignment exactly at the edge where the DUT samples it, the DUT may see the old or the new value depending on process order. The discipline is to drive inputs after the active edge (as a real flop's output would change) and to sample outputs just before it. SystemVerilog **clocking blocks** encode this: an input skew samples DUT outputs a little before the edge and an output skew drives DUT inputs a little after it. Monitors must also follow the protocol cycle by cycle, or they miss back-to-back transfers.

## Going deeper

Bergeron, "Writing Testbenches Using SystemVerilog" (2006); Spear and Tumbush, "SystemVerilog for Verification" (3rd ed., 2012), chapters on testbench architecture and clocking blocks.
