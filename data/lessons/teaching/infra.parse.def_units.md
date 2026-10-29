# DEF database units

DEF stores every coordinate as an integer in database units, declared once by `UNITS DISTANCE MICRONS n`. Integers make geometry exact and comparisons cheap, but every reader and writer must convert correctly between microns and units, and LEF dimensions arrive in microns. Unit conversion is the first computation any physical-design tool performs, and the source of some of its most confusing bugs.

## From DEF integers to microns

With `UNITS DISTANCE MICRONS 2000`, one micron is 2000 database units (dbu), so one dbu is 0.5 nm. The component line `- u12 NAND2_X1 + PLACED ( 38000 57600 ) N ;` puts the instance origin at $x = 38000 / 2000 = 19$ um and $y = 57600 / 2000 = 28.8$ um. The value of n is chosen by the writer from a fixed set (100, 200, 1000, 2000 and so on); it must be no finer than the LEF `DATABASE MICRONS`. Inside a database the integers are usually kept as they are, and conversion to microns happens only for reports and user input, so the engines work in exact integer arithmetic. A coordinate that is not a whole number of database units cannot be written at all, which is why a writer that rescales a design to a coarser resolution must snap every shape and report what moved.

## Die area from DIEAREA

`DIEAREA ( 0 0 ) ( 600000 400000 ) ;` with 2000 units per micron is a die of $600000 / 2000 = 300$ um by $400000 / 2000 = 200$ um, so $300 \times 200 = 60000$ square microns. Area scales with the square of the unit: a product of two dbu lengths is in dbu squared, so the conversion divides by $n^2 = 4{,}000{,}000$, not by n. Dividing once is the classic mistake, and gives an area 2000 times too large. When the lower-left corner is not at the origin, as in `DIEAREA ( 2000 2000 ) ( 602000 402000 )`, subtract corners first. DIEAREA can also be a rectilinear polygon with more points, in which case the area comes from the polygon, not the bounding box.

## LEF microns into DEF rows

LEF gives the SITE width in microns, for example 0.19 um. A DEF row `ROW ROW_0 core 3800 5600 N DO 400 BY 1 STEP 380 0 ;` with 2000 units per micron uses a step of $0.19 \times 2000 = 380$ dbu, and its right edge is at $3800 + 400 \times 380 = 155800$ dbu. Every LEF dimension must land on an integer at the DEF resolution, which is one reason resolutions are chosen to match the manufacturing grid: a 0.005 um grid needs at least 200 units per micron. If a conversion produces a fraction, the writer must round, and a rounded site width accumulates error across hundreds of sites in a row, which breaks legality checks in the placer.

## Going deeper

The LEF/DEF Language Reference, sections on UNITS, DIEAREA and ROW; the OpenDB documentation on database units.
