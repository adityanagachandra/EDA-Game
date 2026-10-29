# Node slack and the critical path

Slack is defined at every pin, not only at endpoints. Once arrivals and required times are known everywhere, the slack of a pin is the slack of the worst path through it. That one number per pin is what lets an optimizer find critical regions without enumerating paths, and it gives a clean characterization of the critical path: the chain of pins that all share the worst slack.

## Slack at an endpoint

At an endpoint, setup slack is its required time minus its latest arrival, $S = RAT - AT$. Negative slack is a violation. With $RAT = 89$ ps and $AT = 99$ ps the endpoint fails by 10 ps. Each endpoint gets its own slack because it has its own clock, its own setup time and its own arrival. Endpoints of different clock domains or path groups are compared only through these slacks, never through raw arrival times, since the arrival time alone says nothing without the deadline it must meet. In a report, the endpoint list sorted by slack is the first thing a timing engineer reads.

## Slack at an internal node

At any pin $v$, $S(v) = RAT(v) - AT(v)$, where $AT$ came from the forward max pass and $RAT$ from the backward min pass. This equals the worst slack of any startpoint-to-endpoint path through $v$: the arrival is the latest path into $v$, and the required time is the tightest path out. Suppose a node has $AT = 35$ ps and $RAT = 50$ ps; every path through it has at least 15 ps of slack, and at least one path has exactly 15. Optimizers sort pins by this slack to decide where to work, because speeding up a pin with large slack cannot help the worst endpoint.

## The critical path shares one slack

Follow the worst path from its startpoint to its endpoint. Every pin on it has a slack equal to the endpoint's slack, because along that path the latest arrival and the tightest deadline come from the same path. So the critical path is exactly a chain of nodes whose slack equals the design's worst slack, and a timer can trace it by walking backward from the worst endpoint, at each step choosing the fanin that produced the arrival. When several paths tie, more nodes share the worst value. Counting nodes at the worst slack is a quick check: if a graph's worst slack is -10 ps and four nodes carry -10 ps, those four form the critical path or paths.

## Going deeper

Sapatnekar, "Timing" (2004), chapters on static timing analysis and critical paths; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8.
