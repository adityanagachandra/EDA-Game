# How poor decisions propagate downstream

A flow is a directed sequence of stages, and an error made in one stage is inherited by every stage after it. Fixing it means changing the stage that made it and re-running all of its descendants. The cost of a decision is therefore not the cost of its own stage but of everything downstream that depends on it. That arithmetic explains ECOs, margins and the value of catching problems early.

## Re-running the descendants

Treat the flow as a chain: synthesis, scan insertion, floorplan, placement, CTS, routing, extraction, signoff. A change at stage $k$ invalidates the outputs of stage $k$ and all stages after it. In a chain of 8 stages, a change at floorplan (stage 3) forces 6 stages to re-run: floorplan, placement, CTS, routing, extraction and signoff. A change at routing forces 3. In a real flow the dependency structure is a DAG rather than a chain, and the set to re-run is everything reachable from the changed node, which a topological order of the DAG gives directly.

## Adding up the cost

If each stage has a run time, the cost of a change is the sum over the stages that must re-run. Suppose floorplan 2 h, placement 6 h, CTS 3 h, routing 10 h, extraction 2 h and signoff 5 h. A floorplan change costs $2 + 6 + 3 + 10 + 2 + 5 = 28$ h of compute before the result is known, while a routing-level ECO costs $10 + 2 + 5 = 17$ h, or less with incremental routing. That gap is why engineers push fixes as late and as locally as they can.

## Iterations that may fail

A re-run does not always succeed. If each attempt closes with probability $q$, independently, the number of attempts is geometric with mean $1/q$. With $q = 0.4$ and 28 h per attempt, the expected time to closure is $28/0.4 = 70$ h. Poor early decisions lower $q$: a floorplan with too little white space makes every placement attempt likelier to fail. The expected cost of a decision is its downstream run time divided by the success probability it leaves, which is why good floorplans and realistic constraints pay off many times.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 1 on the design flow and its iterations.
