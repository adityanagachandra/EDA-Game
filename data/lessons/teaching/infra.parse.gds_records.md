# GDSII stream records

GDSII is the binary stream format in which layouts reach the mask shop, and KLayout, Magic and every DRC engine read and write it. Unlike LEF and DEF it is a flat sequence of typed binary records, so a reader is a loop over record headers. Knowing the record layout lets a tool designer size files, write a minimal reader, and convert stored integers into physical coordinates.

## Record headers and the XY record

Every GDSII record starts with a 4-byte header: a 2-byte big-endian length that counts the whole record including the header, a 1-byte record type (HEADER, BGNLIB, UNITS, BGNSTR, BOUNDARY, XY, ENDEL and so on) and a 1-byte data type (no data, 2-byte integer, 4-byte integer, 8-byte real, ASCII string). A reader reads the header, then length minus 4 bytes of payload, and dispatches on the type. Coordinates live in XY records as pairs of 4-byte signed integers, 8 bytes per point. A BOUNDARY is closed by repeating its first point, so a rectangle with 4 vertices stores 5 points: $4 + 8 \times 5 = 44$ bytes. An L-shaped polygon with 6 vertices stores 7 points, $4 + 56 = 60$ bytes.

## Bytes in a layout element

An element is a sequence of records ending in ENDEL. A BOUNDARY element is BOUNDARY (4 bytes, no data), LAYER (4 + 2), DATATYPE (4 + 2), XY (4 + 8(n + 1)) and ENDEL (4), so an n-vertex boundary costs $24 + 8(n + 1)$ bytes: 64 bytes for a rectangle. A PATH element is a wire with a centerline: PATH, LAYER, DATATYPE, WIDTH (4 + 4, a 4-byte integer), the XY of its k points (not closed, so no repeated point) and ENDEL, which is $28 + 8k$ bytes. A structure holding two rectangles and a 3-point path is $64 + 64 + 52 = 180$ bytes. Sizes like these explain why flat GDS of a large chip runs to many gigabytes, and why hierarchy, with SREF and AREF references to shared structures, and the more compact OASIS format matter.

## Database units to microns

Coordinates are integers in database units. The UNITS record holds two 8-byte reals: the size of a database unit in user units, and its size in meters. A typical library stores 0.001 and 1e-9: one database unit is 1 nm, and the user unit is the micron. An XY integer of 123456 is then $123456 \times 10^{-9}$ m $= 123.456$ um. With a database unit of 5e-10 m the same integer is 61.728 um. A reader must use the meter value, not assume nanometers, and a tool that merges GDS files from different sources, as the final step of a flow merges cell GDS into the design, must rescale any file whose database unit differs, or reject it. The 8-byte reals use the GDSII excess-64 base-16 floating-point format, not IEEE 754, which a hand-written reader must decode explicitly.

## Going deeper

The GDSII Stream Format manual (Calma, release 6.0); the KLayout documentation on stream formats.
