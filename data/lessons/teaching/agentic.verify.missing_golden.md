# Gaps in golden coverage

A change that passes its tests is only as safe as its tests are complete. Agent-written plug-ins make this sharper: a model will produce code that passes exactly the tests it is shown, so any gap in the golden and differential suite is a place where a wrong change can be promoted. Spotting the missing test is a core skill of anyone who reviews agent changes.

## Coverage of the change surface

Every condition in a change needs a test that exercises it both ways. A derate rule "add 3% on HVT data arcs when logic depth is below 5" has three conditions (Vt class, data versus clock, depth threshold), so the reference set needs paths that satisfy each condition and paths that fail each one, including the boundary depth of exactly 5. A router plug-in that penalizes vias within 20 um of an analog block needs nets at 19, 20 and 21 um. A placement rule for multi-row cells needs a design that has them. If the only reference design has depths of 8 to 12, the depth condition has never been tested.

## What the intended effect cannot catch

A test that checks only the intended effect misses collateral damage. Checking that targeted nets lost vias does not show that untargeted nets kept theirs; checking that WNS improved does not show that hold slack did not degrade or that runtime did not double. A differential must list unintended movers across all objects, not just the target set, and cover the metrics the change might plausibly affect: setup and hold, early and late, every corner the design uses, DRC counts after routing, runtime and memory. Legality checks must run on every output, because a plug-in that improves a metric by producing an illegal result is the classic failure.

## Goldens that rot

Goldens also decay. A golden regenerated from the engine's current output after a change no longer tests anything; it only records the change. A baseline recorded on an older engine version mixes unrelated changes into every differential. A golden whose expected value is too loose, such as "WNS within 50 ps", passes almost anything. Good practice regenerates goldens only through a reviewed change with a stated reason, pins the baseline to an engine version, and keeps tolerances near the measured noise. When a bug escapes to a user, the first fix is a new golden that would have caught it.

## Going deeper

Myers, Sandler and Badgett, "The Art of Software Testing" (3rd ed., 2011), on boundary-value analysis.
