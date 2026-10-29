# Resolution enhancement techniques

At low $k_1$ an ordinary mask under ordinary illumination prints a blurred, low-contrast image. Resolution enhancement techniques (RET) change the illumination, the mask's phase and the mask's surroundings so that the image of the intended patterns improves. They shape layout in turn: RET tuned for a few pitches and orientations is why modern design rules are gridded, unidirectional and full of forbidden pitches.

## Off-axis illumination and phase shifting

With on-axis illumination, a grating images only if both first diffraction orders enter the lens. Tilting the illumination (annular, quadrupole or dipole sources) lets the zeroth order and one first order both pass at about half that pitch, which is how $k_1$ approaches its two-beam floor of 0.25. A dipole aligned to a grating images that orientation very well and the perpendicular one poorly, pushing critical metal layers to be unidirectional. Phase-shift masks add a 180-degree phase to some light: attenuated PSM leaks weak, phase-reversed light through dark areas to steepen edges, and alternating PSM gives neighbouring openings opposite phase.

## Assist features and source-mask optimization

An isolated line diffracts differently from a dense one, so it has a different best focus and size, the iso-dense bias. Sub-resolution assist features (SRAFs) are mask shapes too small to print that, placed beside an isolated line, make its diffraction pattern resemble a dense one. They need room: some spacings leave no space for an assist and print poorly. Source-mask optimization (SMO) treats the illumination pixels and the mask shapes as joint variables, optimized on a set of critical clips; the chosen source is then frozen and full-chip OPC runs against it.

## How RET constrains layout

RET is tuned for the patterns that matter most, so other patterns suffer. Strong off-axis illumination can leave forbidden pitches between the dense and isolated cases where contrast collapses. Alternating phase masks require a two-colouring of phase regions, and odd cycles in the conflict graph are unprintable without a layout change. The foundry's response is restricted design rules: a small set of allowed pitches, one direction per layer, fixed line widths and end-of-line rules. Those restrictions simplify the router's world but take away jogs and arbitrary spacing, and they explain many entries in the design-rule manual.

## Going deeper

Wong, "Resolution Enhancement Techniques in Optical Lithography" (2001); Rosenbluth et al., "Optimum mask and source patterns to print a given shape", Journal of Microlithography, Microfabrication and Microsystems (2002).
