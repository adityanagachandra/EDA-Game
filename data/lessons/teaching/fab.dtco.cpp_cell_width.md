# Contacted poly pitch and cell width

The width of a standard cell is counted in contacted poly pitches (CPP): the pitch of the gate lines, including room for a source-drain contact between neighbouring gates. CPP times the number of gate positions, plus whatever the cell boundary needs, sets cell width. Together with metal pitch and track height it gives the basic area scaling metric of a node.

## Cell width in gate pitches

Gates in a modern library sit on a fixed grid. A cell with $g$ gate positions spans about $g$ CPPs, plus space at its boundary to isolate it from its neighbour. With a single diffusion break, which cuts the active region under one dummy gate shared between abutting cells, each cell adds one CPP: width $= (g + 1)\,\mathrm{CPP}$. With a double diffusion break the isolation takes two gate positions: $(g + 2)\,\mathrm{CPP}$. A NAND2 at 48 nm CPP is $3 \times 48 = 144$ nm wide with single diffusion break and 192 nm with double, a 33% difference.

## Row length from a cell list

Because every cell is an integer number of CPPs, row utilization is arithmetic on CPP counts. A row holding an INV (1 gate), a NAND2 (2), an AOI22 (4) and a DFF (about 12) with single diffusion breaks uses $2 + 3 + 5 + 13 = 23$ CPPs, which at 48 nm is 1.104 um. Placement legalizers snap cells to these sites, and filler cells fill the gaps. The diffusion-break style matters for density: single diffusion break saves one CPP per cell boundary, which across a block of small cells is several percent of area, at the cost of a harder isolation process.

## Area scaling between nodes

A common first-order logic density metric multiplies CPP by the cell height, itself track count times metal pitch. Node A at CPP 54 nm, 6 tracks and 36 nm pitch gives $54 \times 6 \times 36$; node B at CPP 48 nm, 6 tracks and 30 nm pitch gives a ratio of $(48 \times 30)/(54 \times 36) = 74.1\%$. Cutting a half track adds another factor of $5.5/6$. Real density gains also depend on diffusion-break style, gate contact placement, and how well blocks route, so foundries quote density on representative mixes of cells. For a tool designer, CPP and metal pitch are the grid every placer and router works on.

## Going deeper

Clark et al., "ASAP7: A 7-nm finFET predictive process design kit" (2016); the IRDS "More Moore" tables of CPP and metal pitch by node.
