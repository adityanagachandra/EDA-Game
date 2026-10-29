# Thermomechanical stress from CTE mismatch

Every package joins materials that expand at different rates with temperature: silicon at about 2.6 ppm/K, copper near 17, organic substrates around 14 to 17, mold compounds from under 10 to over 20, and solder near 20. When the assembly is cooled from the temperature where it was bonded, or cycled in use, the mismatch becomes strain and stress. Those stresses drive warpage, joint fatigue, cracking and delamination, so they are computed from the first package concept onward.

## Mismatch strain

A free body of CTE $\alpha$ strains by $\alpha \Delta T$ when its temperature changes by $\Delta T$. Two bonded materials want different strains, and the difference $\Delta\alpha \Delta T$ is the mismatch strain that the bond must accommodate. Silicon on an organic substrate at 15 ppm/K, cooled 200 K from a solder reflow, has a mismatch of $12.4 \times 200 = 2480$ ppm, or 0.248%. That sounds small, but it is imposed over the whole bonded length and is comparable to the elastic strain limit of many materials. The reference temperature is the stress-free temperature, roughly where the bond solidified or the adhesive cured, not room temperature, which is why process history matters.

## Biaxial stress in a constrained film

If a thin layer is bonded to a much thicker, stiffer one, the layer is forced to take the substrate's strain, and the mismatch becomes stress in the layer. For an isotropic film in equal biaxial strain, $\sigma = \frac{E}{1-\nu}\Delta\alpha\Delta T$, where $E/(1-\nu)$ is the biaxial modulus. A copper film ($E = 120$ GPa, $\nu = 0.34$) on silicon with $\Delta\alpha = 14.4$ ppm/K cooled 100 K carries $120000/0.66 \times 0.00144 = 261.8$ MPa, beyond the yield stress of bulk annealed copper, so real films partly yield. The $1-\nu$ factor is the difference between a film constrained in two directions and a bar constrained in one, where $\sigma = E\Delta\alpha\Delta T$.

## Temperature swing to a stress limit

Inverting the film formula gives the temperature swing that reaches an allowable stress: $\Delta T = \sigma_{allow}(1-\nu)/(E\Delta\alpha)$. For a dielectric layer with $E = 70$ GPa, $\nu = 0.2$, $\Delta\alpha = 10$ ppm/K and a 300 MPa allowable stress, $\Delta T = 300 \times 0.8/(70000 \times 10^{-5}) = 342.9$ K, comfortably above a 200 K process swing. These simple formulas govern thin films far from edges. Near free edges and corners, stress concentrates and changes character, with peeling and shear components that drive delamination, so edge regions are always analysed with finite elements rather than the film formula.

## Going deeper

Freund and Suresh, "Thin Film Materials: Stress, Defect Formation and Surface Evolution" (2003); Suhir, "Structural Analysis in Microelectronic and Fiber-Optic Systems" (1991).
