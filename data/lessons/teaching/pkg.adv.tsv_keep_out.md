# TSV keep-out zones

Copper expands about six times more than silicon per kelvin, so after the anneal and cool-down a TSV leaves residual stress in the silicon around it. That stress changes carrier mobility and so the speed of nearby transistors. Placement tools therefore apply a keep-out zone (KOZ) around each TSV where no cells, or no sensitive cells, may be placed. The KOZ size is a direct area cost and depends on how fast the stress decays.

## Blocked area per TSV

A placer usually models the KOZ as a square or circle around the TSV. With TSV diameter $d$ and keep-out distance $k$ measured from the TSV edge, a square blockage has side $d + 2k$ and area $(d+2k)^2$. A 5 um TSV with a 5 um KOZ blocks $15^2 = 225$ um^2, nine times the TSV's own square. This is why the KOZ, not the via, dominates the area cost of TSVs on an active die. A KOZ can be relaxed for cells that tolerate a few percent of speed change, or applied per device orientation, because the piezoresistive response of silicon differs between electrons and holes and with the transistor's direction relative to the TSV.

## Fraction of a block lost

Multiply the per-TSV blockage by the count and divide by the block area: $f = N (d + 2k)^2 / A$. Two thousand TSVs at 225 um^2 each in a 4 mm^2 block lose $450000/4000000 = 11.25\%$ of the area. That loss lowers placement density and lengthens wires, which can erase the wirelength benefit of going 3D. TSV arrays reduce the loss because neighbouring keep-out zones overlap when TSVs are clustered at small pitch, so ten TSVs in a row block much less than ten isolated TSVs. Placement-aware 3D flows therefore weigh TSV clustering against the longer lateral wires needed to reach a cluster.

## Sizing the keep-out from stress decay

For a long cylinder with a thermal misfit in an infinite plate, the classic Lamé solution gives radial and hoop stresses that fall as $(R/r)^2$ with distance $r$ from the axis, where $R$ is the TSV radius. If the stress at the TSV edge is $\sigma_0$ and transistors tolerate $\sigma_{th}$, the safe radius is $r = R\sqrt{\sigma_0/\sigma_{th}}$, and the KOZ measured from the edge is $r - R$. For $R = 2.5$ um, $\sigma_0 = 200$ MPa and $\sigma_{th} = 20$ MPa, $r = 2.5\sqrt{10} = 7.91$ um, a KOZ of 5.41 um. The square-law decay means that a tenfold tighter tolerance costs only about three times the radius, and the real field near the surface and liner differs, so foundries calibrate KOZ rules by measurement and finite-element models.

## Going deeper

Yang, Athikulwongse, Lim and Pan, "TSV stress aware timing analysis with applications to 3D-IC layout optimization", DAC (2010); Timoshenko and Goodier, "Theory of Elasticity", on the thick-walled cylinder.
