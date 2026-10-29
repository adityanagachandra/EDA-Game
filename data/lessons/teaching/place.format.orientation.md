# DEF units, orientations and pin coordinates

Every pin coordinate a placer uses comes from combining a cell's LEF pin offset with the instance's DEF location and orientation. Getting this transform wrong silently corrupts wirelength, timing and routing. The arithmetic is simple, but a tool designer must know exactly what the DEF placement point means and how each of the eight orientations maps local coordinates.

## Database units and the placement point

DEF coordinates are integers in database units (DBU). UNITS DISTANCE MICRONS 2000 means one micron is 2000 DBU, so PLACED ( 62640 124800 ) puts the instance origin at (31.32, 62.4) um. The placement point is the lower-left corner of the instance's bounding box after the orientation is applied, not the cell's own LEF origin before transformation. A cell of SIZE 2.17 BY 1.6 placed at that point occupies x from 31.32 to 33.49 um and y from 62.4 to 64.0 um in every orientation that keeps its width horizontal. Integers make legality exact: a site grid of 0.19 um at 2000 DBU per um is a step of 380 units, and any location off that grid is illegal.

## Transforming a pin offset

A pin at offset $(p_x, p_y)$ inside a $W \times H$ macro maps as follows. N leaves it at $(p_x, p_y)$. FN mirrors about the vertical axis: $(W - p_x, p_y)$. FS mirrors about the horizontal axis: $(p_x, H - p_y)$. S rotates by 180 degrees: $(W - p_x, H - p_y)$. The rotations by 90 degrees (E, W and their flipped forms) also swap width and height, and standard cells in rows never use them. The absolute pin is the placement point plus the transformed offset. A cell 2.22 by 1.4 um with a pin at (1.46, 1.35), placed FN at (43.32, 30.8), has its pin at $(43.32 + 0.76, 30.8 + 1.35) = (44.08, 32.15)$ um. Note that the mirror subtracts from the macro width, not from the site width or the row height, and that the placement point itself does not move when the orientation changes.

## Why orientation matters to the placer

Rows alternate N and FS so that power rails line up, so legalization must flip a cell when it moves it between an even and an odd row, and every pin position changes when it does. Detailed placers also flip cells deliberately: mirroring a cell to FN can bring its output pin closer to its sink and shorten a net by a fraction of the cell width, a cheap move with no displacement. Wirelength code therefore stores pin offsets per orientation, or transforms on the fly, and incremental updates must recompute a net's bounding box after a flip. Off-by-one errors here, such as forgetting to subtract from $W$ rather than from the site width, produce placements that look legal but have wrong pin locations.

## Going deeper

The LEF/DEF Language Reference (version 5.8), sections on COMPONENTS and orientation; the OpenDB documentation in OpenROAD.
