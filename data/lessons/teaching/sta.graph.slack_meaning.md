# What slack measures

Static timing analysis reduces a chip's timing to one signed number per endpoint: slack. Every optimization step downstream, gate sizing, buffering, placement moves, useful skew, reads slack and tries to raise its worst value. Knowing exactly what the number means, and what it does not, is the difference between fixing a timing problem and moving it somewhere else.

## Required minus arrival

The arrival time at an endpoint is the latest time its data input settles, found by propagating delays forward from the launching clock edge through the combinational logic. The required time is the latest moment the capturing flip-flop can accept that data: the next capturing clock edge, minus the flop's setup time, adjusted for clock skew and uncertainty. Setup slack is required time minus arrival time. A path with a 1000 ps clock, 940 ps of arrival and 40 ps of setup has slack $1000 - 40 - 940 = 20$ ps.

## The sign of slack

Positive slack is margin: the path could get that much slower and still meet timing. Negative slack is a violation, and its size says how much faster the path must become. Two summary numbers describe a whole design. Worst negative slack (WNS) is the most negative endpoint slack, the single worst path. Total negative slack (TNS) adds every negative endpoint slack, so it measures how widespread the problem is. A design with a WNS of -5 ps and a TNS of -3000 ps has a broad problem that local fixes cannot close.

## Hold slack runs the other way

Hold checks the opposite failure: new data arriving so early that it corrupts the value the capture flop is still holding from the previous cycle. Hold slack is the earliest arrival time minus the hold requirement, the same-edge clock arrival plus the hold time. Here a larger arrival time is better. That asymmetry is why timing fixes interact: upsizing a gate to speed a setup-critical path also speeds every short path through that gate, and can turn a passing hold check into a failing one.
