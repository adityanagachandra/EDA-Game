# VHDL std_logic and resolution functions

When several drivers drive one VHDL signal, a resolution function computes the value the signal takes. IEEE 1164's std_logic uses nine values and one resolution table to model strong and weak drivers, high impedance and unknowns. Simulators implement this table on every multiply driven net, so a tool designer must be able to apply it and to know which types allow multiple drivers at all.

## Nine values and their strengths

std_ulogic has nine values: 'U' (uninitialized), 'X' (forcing unknown), '0' and '1' (forcing values), 'Z' (high impedance), 'W' (weak unknown), 'L' and 'H' (weak 0 and weak 1) and '-' (don't care). Strength is the key idea. A forcing value comes from an active driver; a weak value comes from something like a pull-up resistor; 'Z' means not driving at all.

The resolution rule follows from strength. A forcing value beats a weak one, a weak value beats 'Z', and 'Z' with anything gives the other value. So '1' and 'Z' resolve to '1', and 'H' and 'Z' resolve to 'H'. Most std_logic signals in RTL have a single driver, and the resolution function returns that driver's value unchanged.

## Conflicts between drivers

Two different values of the same strength conflict, and the result is the unknown of that strength: '0' and '1' give 'X'; 'L' and 'H' give 'W'. A weak conflict loses to any forcing value, so 'W' and '0' give '0'. 'X' dominates every value except 'U', which dominates everything, and '-' resolves like 'X'.

With 'L' and 'H' driving the same bus the result is 'W', telling the designer that two pull resistors fight. With a pull-up 'H' and an active '0', the result is '0', which is how open-drain buses such as I2C are modeled: any device pulling low wins over the pull-up.

## Resolving many drivers

The resolution table is commutative and associative, so the result for any number of drivers is obtained by folding the table over them in any order. For drivers '0', 'H' and 'Z': '0' with 'H' gives '0', and '0' with 'Z' gives '0'. For '1', 'L' and 'X', the 'X' dominates and the result is 'X'.

Which types may have several drivers is part of the language. std_logic is resolved and allows them; std_ulogic is unresolved, so a second driver is a compile-time error, which catches accidental multiple drivers early. Synthesis maps a resolved signal with several tristate drivers to tristate buffers or, inside a chip where internal tristates are not allowed, rejects it or converts it to multiplexers.

## Going deeper

IEEE Std 1164-1993, the resolution table of std_logic; Ashenden, "The Designer's Guide to VHDL", the chapter on resolved signals.
