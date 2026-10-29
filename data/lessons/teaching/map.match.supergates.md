# Supergates

Boolean matching sees one cell per cut. Supergates widen the view: small compositions of library gates, generated once per library, each treated as a single matchable function with its own area and pin delays. They reduce the structural bias of the subject graph and let the mapper find multi-gate implementations in one matching step. ABC's standard-cell mapper is built on them.

## What a supergate is

A supergate is a fixed tree of library gates: a root gate whose input pins are driven either by inputs of the supergate or by other gates. Its function is the composition, computed once as a truth table, and it is stored in the matching library like a cell. Its area is the sum of the component areas. A NAND2 whose pin A is driven by an INV, with area 3 + 2 = 5, is a supergate computing NOT(NOT a AND b). Generation enumerates compositions up to a bound on inputs (for example 5) and levels, discarding those dominated by an existing cell or supergate with the same function.

## Pin-to-output delays

Each supergate input has a delay to the output equal to the sum of the pin delays along its path. In NAND2 with pin A at 14 ps driven by INV at 6 ps, input a has delay 6 + 14 = 20 ps and input b, on pin B at 14 ps, has 14 ps. With a arriving at 3 ps and b at 20 ps, the output settles at max(3 + 20, 20 + 14) = 34 ps. The mapper uses these per-pin delays with the matching transformation, so it can place a late signal on a supergate input with a short path, just as it would with a single cell.

## Inputs, bounds and library size

A supergate's input count is the root's pin count, minus the pins driven internally, plus the inputs of the driving gates: NAND2 with one pin driven by a NAND2 has 2 - 1 + 2 = 3 inputs. Bounding the input count keeps the cut functions small enough to canonicalize, and bounding levels keeps the number of supergates manageable, since it grows combinatorially with library size. Dominance pruning removes a supergate if another with the same function is no worse in area and in every pin delay. The result is a library of thousands of functions with known costs, matched by canonical form.

## Going deeper

Mishchenko, Chatterjee and Brayton, "Technology mapping with Boolean matching, supergates and choices", UC Berkeley technical report, 2005; Chatterjee, Mishchenko, Brayton, Wang and Kam, "Reducing structural bias in technology mapping", ICCAD 2005.
