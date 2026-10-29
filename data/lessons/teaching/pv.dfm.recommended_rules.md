# Recommended rules and DFM scoring

Beyond mandatory rules, foundries publish recommended rules: larger spacings, longer extensions, redundant vias. A layout that ignores them is legal but yields worse. DFM tools score how well a layout follows them and apply them where space allows. The engineering problem is choosing where the yield benefit justifies the area or routing cost.

## Scoring recommended rules

Recommended rules are scored, not passed or failed. A redundant-via score is the fraction of via locations with two or more cuts: 1534 of 2284 locations gives 67.2 percent. A spacing score is the fraction of neighbouring-wire gaps at or above the recommended value. Better scores weight each location by its yield impact, for example its critical area, so a long tightly spaced bus counts more than a short stub. The score is a metric to optimize, which makes it a natural objective for automated DFM passes and for agents.

## Wire spreading and critical area

Under the model where a gap's average short critical area is $L x_0^2 / s$, spreading two wires from 80 to 90 nm reduces it to $80 / 90$ of its value, an 11.1 percent reduction. The gain per nanometre is largest for the tightest gaps, so spreading targets minimum-spaced long runs where free tracks exist. Spreading can create new problems: it moves wires toward other neighbours, changes coupling, and may conflict with double-patterning colours. A DFM pass therefore re-runs DRC and extraction on what it changes.

## How tools use recommended rules

Routers apply recommended spacings and extensions opportunistically, where tracks are free, and fall back to minimum rules where routing is tight. Redundant-via insertion runs after routing, using leftover space. Recommended rules guard against specific failures: a longer metal extension past a via protects the enclosure from line-end shortening. Applying everything everywhere would cost area; the right policy applies rules where they buy the most yield, which is a cost-benefit decision the engine exposes as parameters.
