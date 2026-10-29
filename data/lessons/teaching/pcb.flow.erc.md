# Electrical rule checking of schematics

Electrical rule checking (ERC) runs on the schematic before layout. It reads each net's pins and their declared electrical types and flags combinations that cannot work: two outputs fighting, an input that nothing drives, a power input with no source, a net with a single pin. ERC is a set of cheap, exact predicates over the netlist hypergraph. It catches connectivity mistakes early, when they cost nothing to fix, but it knows only pin types, not the parts' datasheets.

## Pin types and the conflict matrix

Every symbol pin carries a type: input, output, bidirectional, tri-state, passive, open collector, power input, power output, or unconnected. ERC tools encode which types may share a net in a **pin conflict matrix**: output with output is an error (two drivers fight), output with input is fine, passive with anything is usually fine. The matrix is configuration, so a team can tighten or relax it. A net that touches only one pin is flagged too, since it connects nothing; a pin deliberately left open gets a no-connect marker so the check knows the intent. Power pins are special: KiCad, for example, reports a power input that no power output drives, which is why designs fed from a connector add a power flag symbol to declare the source.

## Evaluating the rules net by net

ERC is a pass over nets. For each net, count pins by type and evaluate each rule as a predicate on the counts. With the rules (a) exactly one pin, (b) two or more outputs or power outputs, (c) an input but no output, bidirectional, passive or power-output pin, and (d) a power input but no power output, consider three nets. N1 has an output and two inputs: no rule fires. N2 has two outputs and an input: rule (b) fires. N3 has two inputs and a bidirectional pin: rule (c) does not fire, since a bidirectional pin can drive. Count N2 once even if it broke two rules. The work is linear in the number of pins, so ERC runs interactively after every edit.

## What ERC cannot see

ERC reasons only about declared pin types, so it inherits every library mistake. A symbol that marks an output pin as passive hides a real conflict. It cannot see that a regulator's output capacitor is below the datasheet minimum, that an I2C bus lacks pull-ups (a passive pull-up and an open-drain pin are both legal), or that a 3.3 V output drives a 1.8 V input. Those are semantic checks that need part parameters and datasheet rules. Treating ERC warnings as noise is the common failure: a design with hundreds of waived warnings hides the one real error. Good flows keep the conflict matrix strict, annotate intentional exceptions explicitly, and add parameter-aware checks on top, which is where review agents can help.

## Going deeper

The KiCad Schematic Editor documentation, section on the electrical rules checker and its pin conflict map.
