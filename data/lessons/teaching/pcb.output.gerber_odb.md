# Gerber and ODB++

Fabrication data leaves the board tool in one of a few formats. Gerber (RS-274X, and its attribute-carrying extension X2) describes each layer as a 2D vector image, with drill data in separate Excellon files. ODB++ and IPC-2581 package the whole job, layers, drills, netlist, components and stackup, in one structured database. A tool designer must know what each format preserves, because whatever a format drops must be carried by side documents that can disagree with it.

## The Gerber image model

A Gerber file is a sequence of commands that draw one layer's image. **Apertures** are defined with statements such as %ADD10C,0.25*%, aperture D10 a circle 0.25 units in diameter. Operations use the current aperture: **D01** draws (interpolates) from the current point to a new coordinate, linearly or circularly depending on G01, G02 or G03; **D02** moves without drawing; **D03** flashes the aperture, which is how pads are usually drawn. **Regions** (G36 and G37) fill arbitrary polygons, as for copper pours. Dark and clear polarity lets later objects erase earlier ones. One file per layer is the convention. The image carries no stackup, no layer order and, in classic RS-274X, no connectivity.

## Coordinates and units

Coordinates are integers interpreted by a format statement. %FSLAX46Y46*% declares absolute coordinates with leading zeros omitted, 4 integer and 6 decimal digits; %MOMM*% sets millimetres. So X12500000 means 12.5 mm and Y-3250000 means -3.25 mm. Example: X5000000Y2000000D02* moves to (5, 2), and X17500000Y2000000D01* draws a horizontal segment to (17.5, 2), 12.5 mm long. Misreading the decimal digits scales the whole board by a power of ten, which is why the format statement must precede coordinates and why viewers check board extents. Excellon drill files have their own unit and zero-suppression settings, and a mismatch between Gerber and drill units is a classic CAM error.

## Attributes and unified databases

Classic Gerber leaves the fabricator to reconstruct meaning: which file is which layer (from file names), the layer order and stackup (from a drawing), and the netlist for bare-board test (from an IPC-D-356 file). **Gerber X2** adds attributes: .FileFunction identifies a file's layer, .AperFunction marks objects as pads, vias or conductors, and newer extensions add net and component information. **ODB++** and **IPC-2581** go further, storing layers, drills, netlist, components, BOM and stackup in one job, so CAM import is automatic and internally consistent. Whatever the format, fabricators run their own DFM checks against their real process, because the designer's rule deck is only a model of it. A missing outline layer or drill file stops the job.

## Going deeper

Ucamco, "The Gerber Layer Format Specification" (current revision). The IPC-2581 specification. The ODB++ specification published by Siemens.
