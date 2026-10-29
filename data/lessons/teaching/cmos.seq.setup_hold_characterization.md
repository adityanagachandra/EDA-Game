# Characterizing setup, hold and clock-to-q

A flip-flop has no single correct setup time. Its clock-to-q grows smoothly as data arrives closer to the clock edge, until it fails outright. Characterization chooses a criterion that turns that curve into the setup and hold numbers written into Liberty. A tool designer must know the criterion, because it decides how much margin every timing check in the design inherits.

## The push-out curve

Sweep the data-to-clock skew, the time by which D changes before the clock edge, and measure clock-to-q at each point. For large skews clock-to-q is flat at its nominal value. As the skew shrinks, the master has less time to settle and the slave receives a marginal value, so clock-to-q rises: first gently, then steeply, until at some skew the flop captures the wrong value or never resolves. The curve is continuous; there is no physical cliff at which "setup" happens. Characterization therefore needs a rule for where on the curve to draw the line, and the rule must be applied identically to every flop, slew and corner.

## The push-out criterion

A common criterion defines setup time as the smallest skew at which clock-to-q has degraded by no more than a set percentage, often 10%, from its nominal value. If nominal clock-to-q is 50 ps, the limit is 55 ps, and the setup time is the smallest skew whose measured clock-to-q is at most 55 ps. Characterizers search for this point by bisection on the skew, running one transient per probe, which is why constraint characterization costs many more simulations than delay. Hold time is found the same way from the other side, with data changing after the edge. Some flows instead use a pass-fail criterion, the smallest skew at which the correct value is still captured, which gives a smaller, less conservative number with no push-out guard band.

## Why the criterion matters downstream

Setup time and clock-to-q are coupled: data arriving later than the library setup time would still be captured, just with a longer clock-to-q. The total data-to-output delay, skew plus clock-to-q, is smallest at a skew a little inside the push-out region, which is why some high-performance flows characterize at that point. A static timer, though, treats setup and clock-to-q as independent constants, so the library must pick one consistent operating point, and the push-out percentage decides how much hidden margin each check carries. A change of criterion silently shifts every constraint table in the library. That is why a characterization engine treats the criterion as part of its core definition, not as a tunable parameter.

## Going deeper

Stojanovic and Oklobdzija, "Comparative analysis of master-slave latches and flip-flops" (IEEE JSSC, 1999); Srivastava and Roychowdhury, work on interdependent setup and hold characterization; the Liberty Reference Manual, timing constraint groups.
