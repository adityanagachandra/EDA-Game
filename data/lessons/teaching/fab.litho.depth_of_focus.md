# Depth of focus and the focus budget

A pattern prints correctly only within a narrow range of focus. That range, the depth of focus, shrinks as numerical aperture rises, and everything that moves the wafer surface out of the focal plane must fit inside it: wafer flatness, leftover topography from polishing, and the scanner's own focus control. Density and fill rules in a PDK are one way this budget reaches layout.

## The depth-of-focus relation

The usable focus range for a given feature is usually written $DOF = k_2 \lambda / NA^2$, where $k_2$ depends on the feature, the illumination and how much dose and size error is tolerated. For 193 nm immersion at NA 1.35 with $k_2 = 0.7$: $DOF = 0.7 \times 193 / 1.8225 = 74.1$ nm. For EUV at NA 0.33 with the same $k_2$: $0.7 \times 13.5 / 0.1089 = 86.8$ nm. A process window is the region in dose and focus where the printed size stays within specification; its focus extent at acceptable dose latitude is the working depth of focus.

## The price of numerical aperture

Resolution improves as $1/NA$ but depth of focus shrinks as $1/NA^2$. Raising NA from 0.33 to 0.55 at 13.5 nm improves resolution by a factor of 1.67 but cuts depth of focus by $(0.55/0.33)^2 = 2.78$. The same trade appeared when dry ArF moved from NA 0.75 to 0.93: resolution improved 24% while depth of focus fell by a factor of 1.54. High-NA EUV therefore needs thinner resists, flatter wafers and tighter focus control, and it is one reason design teams cannot assume a new tool solves a pitch without side effects.

## Budgeting focus

The depth of focus is a budget shared by every source of height error across the exposure field: wafer and chuck flatness, scanner focus and levelling errors, lens field curvature, and the topography left on the wafer by earlier steps. If a layer has 90 nm of usable depth of focus, flatness takes 20 nm and focus control takes 25 nm, the underlying topography may use at most 45 nm. That remainder is a requirement on CMP, and CMP's result depends on pattern density, so it becomes density rules and dummy fill in the PDK. A tool that places fill or routes metal is, indirectly, managing the lithographer's focus budget.

## Going deeper

Mack, "Fundamental Principles of Optical Lithography", chapters on focus and process windows; Burn J. Lin's papers on depth of focus scaling.
