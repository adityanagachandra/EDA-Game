# Reflow soldering profiles

Reflow soldering melts printed solder paste in an oven that moves the board through heated zones. The temperature the board experiences over time, the profile, must preheat evenly, activate the flux, exceed the solder's liquidus long enough to wet, peak below the parts' limits and cool at a controlled rate. Profiles are measured with thermocouples on a real board and summarized by ramp rates, soak time, time above liquidus and peak temperature, all computed from the piecewise curve.

## Profile stages and ramp rates

A typical lead-free profile has four stages: a **preheat** ramp from room temperature to about 150 C, a **soak** that rises slowly toward the high 100s C while the flux activates and the board's temperature evens out, a **reflow** ramp to a peak typically in the 235 to 250 C range for SAC alloys, and **cooling**. Ramp rate is the slope of each segment: $(150 - 25)/90 \approx 1.39$ C per second for a 90 s preheat. Ramps much faster than a few C per second risk thermal shock that cracks ceramic capacitors, and the cooling rate affects joint grain structure. Recipes follow the paste vendor's window and the parts' limits, the most sensitive part setting the maximum.

## Time above liquidus

The joint forms only while the solder is molten. For SAC305 the liquidus is about 217 C. **Time above liquidus** (TAL) is measured from the up-crossing of 217 C to the down-crossing, each found by linear interpolation within its segment. Example: a ramp from (160 s, 185 C) to (210 s, 245 C) crosses 217 C at $160 + 32/60 \cdot 50 \approx 186.7$ s; a fall from (210 s, 245 C) to (250 s, 205 C) crosses it at $210 + 28/40 \cdot 40 = 238$ s; TAL is about 51.3 s. Too short a TAL gives poor wetting and cold joints; too long grows brittle intermetallic layers and stresses parts. Paste vendors specify a TAL window, often tens of seconds to about a minute and a half.

## Soak, thermal mass and profiling

A board is not one temperature. Large parts and planes heat slowly, small parts quickly, so the **soak** gives the board time to equalize before the reflow ramp; the time from 150 C to about 200 C is a common measure of it. Too short a soak leaves cold spots under large BGAs; too long exhausts the flux. Thermocouples are placed on the largest thermal masses and the smallest parts, and the profile must satisfy the window at all of them, which is why the board design itself (copper balance, part placement) affects solderability. Double-sided boards see two reflows, and parts on the first side must survive the second. Profiling is a measurement loop: adjust zone temperatures and conveyor speed until every thermocouple trace lands inside the window.

## Going deeper

IPC J-STD-020, on peak reflow temperatures for package classification. Paste vendors' application notes on SAC305 profiles.
