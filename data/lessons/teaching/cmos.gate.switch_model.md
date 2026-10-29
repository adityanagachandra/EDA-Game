# The MOSFET as a switch

Before a tool can time, size or characterize a gate, it needs the simplest faithful picture of a transistor: a switch controlled by its gate voltage. The switch model is what switch-level simulators, logic extractors and netlist checkers run on, and it is the first thing a designer reasons with when drawing a new cell. Knowing exactly where the model is right, and where it breaks, tells you which questions it can answer and which need resistance, capacitance and thresholds.

## Gate voltage opens and closes the channel

An nMOS conducts when its gate is high and is open when its gate is low. A pMOS is the complement: it conducts when its gate is low. In circuit terms, the device turns on when its gate-to-source voltage passes the threshold $V_t$, and in a gate the nMOS sources sit near ground and the pMOS sources near $V_{DD}$, so "high" and "low" are measured against those rails. A gate input therefore drives a pair of opposite switches: when input A is 1, every nMOS it controls is closed and every pMOS it controls is open. That single fact is what makes complementary logic work, and it is why a CMOS gate input looks to the driver like a capacitor rather than a resistor: no DC current flows into a gate.

## Series means AND, parallel means OR

Switches compose. Two switches in series conduct only if both are closed; two in parallel conduct if either is. For nMOS devices with gates A and B, a series pair conducts on $AB$ and a parallel pair on $A + B$. For pMOS devices the same topology conducts on the complemented inputs: a series pair on $\bar A \bar B$ and a parallel pair on $\bar A + \bar B$. A static CMOS gate puts an nMOS network between the output and ground and a pMOS network between $V_{DD}$ and the output, built so that for every input vector exactly one network conducts. Example: two series nMOS from Y to ground with two parallel pMOS from $V_{DD}$ to Y pull Y low exactly when $AB = 1$, and high otherwise, which is a NAND2.

## Where the switch model stops

The ideal switch has no resistance, no capacitance and no threshold loss, so it predicts logic values but never delay or power. Two refinements recover most of what tools need. First, strengths: an nMOS passes a strong 0 but only a weak 1, rising to $V_{DD} - V_{tn}$, and a pMOS passes a strong 1 but a weak 0. Switch-level simulators track strong, weak and undriven levels so they can flag degraded nodes. Second, charge storage: a node driven by neither network keeps its value on its capacitance for a while, a state such simulators model explicitly. For timing, the next step up replaces each closed switch by an effective resistance and each node by a capacitance, which is the RC model.

## Going deeper

Weste and Harris, "CMOS VLSI Design: A Circuits and Systems Perspective" (4th ed.), chapter 1; Bryant's work on switch-level simulation (the MOSSIM simulators).
