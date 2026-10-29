# Reading a timing path report

A path report is the timer's explanation of one slack. It lists the launch clock, every cell and net increment along the data path, and then the capture clock with its adjustments. A tool designer reads reports to debug delay calculation and constraints; an engineer reads them to decide on fixes. The arithmetic is simple, but each line has a precise meaning, and misreading one, such as the sign of a pessimism credit, gives a wrong conclusion.

## Data arrival time

The launch section starts at the launching clock edge, usually 0, adds the clock network delay to the launch flop, then the flop's clock-to-Q, then each net and cell increment in order, ending at the capture flop's data pin. The running total is the data arrival time. With an edge at 0, 115 ps of clock network, an 85 ps clock-to-Q, cell and net increments summing to 371 ps and a final 10 ps net, the arrival is 581 ps. Each increment is one arc of the timing graph, marked with its transition, rise or fall, which lets a reader check unateness and see which delay table was used.

## Data required time

The capture section starts at the capturing edge, one period later for a single-cycle path, adds the capture clock network delay, adds the clock reconvergence pessimism credit, and subtracts the clock uncertainty and the library setup time. With an edge at 800 ps, 130 ps of network, a 9 ps credit, 40 ps of uncertainty and a 58 ps setup, the required time is $800 + 130 + 9 - 40 - 58 = 841$ ps. The pessimism credit is added because it removes double-counted variation on the shared clock path. Under OCV, the capture clock network is the early value and the launch network the late value, which is why the two networks can differ even for neighboring flops.

## Slack and what to read next

Slack is required minus arrival: with 782 ps required and 541 ps of arrival, the path has 241 ps of slack. When it is negative, the report shows where the time went. A single large cell increment points at a weak driver or a huge load; a large net increment at a long or resistive wire; a large difference between launch and capture network delay at skew; a big uncertainty line at an over-tight constraint. The report also names the start and end points, the path group and the check type, which tell whether the constraint itself is right. OpenSTA prints this format with report_checks, and -fields adds slew, capacitance and fanout for each stage.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8 on timing reports; the OpenSTA documentation for report_checks.
