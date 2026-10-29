# Photomasks, reduction and MEEF

A photomask, or reticle, is a quartz plate carrying the pattern of one layer, written by an electron-beam or laser writer. Scanners image it at a reduction, so mask errors shrink on the wafer, but near the resolution limit they can be amplified again. Masks are a large share of the cost of a new design, and their geometry fixes how many dies one exposure prints.

## Reduction and mask scale

Modern scanners demagnify the mask by 4. A 20 nm wafer line is drawn as 80 nm on the mask before OPC bias, and a full scanner field of 26 mm by 33 mm on the wafer is 104 mm by 132 mm on a 6-inch mask. Reduction relaxes mask making: the mask writer works at four times the wafer dimension, and mask defects and size errors are divided by four in the image. DUV masks are chrome on quartz, often with an attenuating phase-shift film; EUV masks are reflective multilayer mirrors of molybdenum and silicon with an absorber pattern, used at a small angle of incidence.

## The mask error enhancement factor

Near the resolution limit the image does not respond linearly to mask changes. The mask error enhancement factor (MEEF) measures the amplification: $\Delta CD_{wafer} = MEEF \cdot \Delta CD_{mask}/M$, where $M = 4$ is the reduction. A MEEF of 1 means errors simply shrink by four. Dense lines at low $k_1$, contacts and line ends commonly have MEEF of 2 to 4. A 2 nm mask-scale error with MEEF 3 prints a $3 \times 2/4 = 1.5$ nm wafer error, three times what reduction alone predicts. High MEEF forces tighter mask specifications and makes OPC solutions more sensitive.

## Fields, dies and mask cost

One exposure prints one field; the scanner steps it across the wafer. A full field is 26 mm by 33 mm. A die of 12 mm by 15 mm fits $\lfloor 26/12 \rfloor \times \lfloor 33/15 \rfloor = 2 \times 2 = 4$ times per field. High-NA EUV uses anamorphic optics with 4x reduction in one direction and 8x in the other, so the field halves to 26 mm by 16.5 mm; a die taller than 16.5 mm then needs stitching across two exposures, which reaches the floorplanner. A leading-edge mask set has dozens of masks, so mask cost is a design-scale concern and a reason DTCO counts masks.

## Going deeper

Rizvi (ed.), "Handbook of Photomask Manufacturing Technology"; Wong, "Resolution Enhancement Techniques in Optical Lithography", on MEEF.
