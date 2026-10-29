# SPICE netlists

SPICE, written at Berkeley in the early 1970s, set the input format that nearly every circuit simulator still reads: the netlist. A netlist is a list of element cards, one per line, each naming an element, the nodes it connects and its value. The simulator's first job is to turn those cards into unknowns and equations. A tool designer must know exactly which cards create unknowns, because the count fixes the size of every matrix the engine will factor.

## Elements, nodes and the ground node

Each card begins with a letter that sets the element type: R resistor, C capacitor, L inductor, V independent voltage source, I independent current source, D diode, M MOSFET, E voltage-controlled voltage source, G voltage-controlled current source, X subcircuit instance. The fields after the name are node names, then the value or model. Node names are arbitrary strings; node 0 is always ground, the reference whose voltage is zero by definition. The card R3 in mid 2.2k connects nodes in and mid. The first line of a deck is a title, lines starting with * are comments, and dot cards such as .tran request analyses. To count nodes, collect every name that appears in a node field and drop 0. A deck with cards V1 in 0 1.8, R1 in out 1k and C1 out 0 1p has two non-ground nodes, in and out.

## Which elements add branch currents

Modified nodal analysis makes every non-ground node voltage an unknown. Elements whose current is a function of their terminal voltages, such as resistors ($i = g v$), current sources and, at DC, capacitors (open), stamp directly into node equations and add nothing. Elements whose current cannot be written that way need their current as an extra unknown: independent voltage sources, and inductors, which SPICE keeps as branch currents so that the same matrix structure serves DC (a 0 V short) and transient. So the DC system size is nodes plus voltage sources plus inductors. A deck with 4 nodes, 2 voltage sources and 1 inductor has $4 + 2 + 1 = 7$ unknowns, however many resistors and capacitors it contains.

## From netlist lines to unknowns

Controlled sources follow the same rule. A VCVS card E1 n+ n- nc+ nc- gain forces a voltage, so its current is unknown and it adds a branch. A VCCS card G1 n+ n- nc+ nc- gm defines a current from voltages, so it only stamps conductance-like entries. Current-controlled sources (F and H cards) need the current of a controlling voltage source, which is why SPICE requires one, often a 0 V source used as an ammeter. Count carefully: a new node created only by a controlled source is still a node. The deck V1 in 0 1, R1 in a 1k, E1 buf 0 a 0 10, R2 buf 0 1k has nodes in, a and buf (3), plus branches for V1 and E1, so 5 unknowns.

## Going deeper

Nagel, "SPICE2: A Computer Program to Simulate Semiconductor Circuits" (Berkeley ERL Memo M520, 1975); the ngspice manual, chapters on circuit description; Najm, "Circuit Simulation" (Wiley, 2010), chapter 2.
