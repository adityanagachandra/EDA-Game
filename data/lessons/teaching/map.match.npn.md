# NPN classes and Boolean matching

Boolean matching recognizes a library cell by its function rather than its shape. A cell's pins can be connected in any order, and inverters are cheap, so functions that differ by input permutation, input negation or output negation should match the same cell. These transformations form a group, and its orbits, the NPN classes, organize every modern matcher.

## The NPN transformation group

An NPN transformation of an n-input function permutes the inputs, optionally negates each input, and optionally negates the output: $g(x) = f(\pi(x) \oplus \nu) \oplus o$. There are $n!$ permutations, $2^n$ input-negation patterns and 2 output polarities, so the group has $2^{n+1} n!$ elements: 16 for n = 2, 96 for n = 3, 768 for n = 4. The number of functions is $2^{2^n}$: 256 for n = 3, 65536 for n = 4. Two functions are NPN-equivalent if some transformation maps one to the other; a cell matches a cut function if the cut function is NPN-equivalent to the cell's function, the transformation telling the mapper which pin gets which signal and where inverters go.

## Orbits and orbit sizes

The orbit of f is the set of functions reachable from it. By the orbit-stabilizer theorem, orbit size times the number of transformations that fix f equals the group order. Symmetric functions have small orbits: XOR3 has orbit {XOR3, XNOR3} of size 2 under NPN, because permutations do nothing and each negation only flips the output. AND3 has orbit size 16: negating inputs gives the 8 functions with a single one, and output negation gives their 8 complements. A function with no symmetry has orbit size 96. Orbit sizes tell a matcher how many truth tables point to the same cell.

## How many classes, and why it matters

The NPN classes of n-input functions number 2, 4, 14 and 222 for n = 1, 2, 3, 4, against 4, 16, 256 and 65536 functions; for n = 5 there are 616126. For 3 inputs, 4 of the 14 classes are degenerate (constants and functions of fewer inputs) and 10 depend on all three inputs. Because there are only 222 classes of 4-input functions, a mapper can precompute, for each class, the best library implementation, and then match any 4-input cut by computing its class. For 6 inputs the classes are too many to tabulate fully, so matchers canonicalize on the fly and store only the classes that the library and its supergates actually realize.

## Going deeper

Hurst, Miller and Muzio, "Spectral Techniques in Digital Logic" (1985), for NPN classification; Mishchenko, Chatterjee and Brayton, "Technology mapping with Boolean matching, supergates and choices", ERL technical report, UC Berkeley, 2005.
