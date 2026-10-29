# Amortized runtime of skill updates

A common objection to an agent that edits and recompiles engines is runtime overhead. The answer is an amortization argument: a flow is set up once, a skill update is small and occasional, and the fast modeling loop that does the design work dominates total compute. The argument is only as good as its numbers, so a tool designer should be able to run them.

## One-time flow setup

Standing up an agentic flow for a new kind of run takes several iterations while the model works out the sequence of skills and which ones to customize. If setup takes $I$ iterations of cost $c$ each and the flow then serves $R$ production runs, the setup cost per run is $I c / R$. With 12 setup iterations of 2 hours each and 400 runs over a deployment, setup adds $24 / 400 = 0.06$ hours per run, about 3.6 minutes against runs that take hours. Setup is a fixed cost; it matters for a flow used three times and vanishes for a flow used hundreds of times.

## Fraction of compute on skill updates

A skill update (a plug-in edit, a few harness attempts, a validation run) costs $U$ compute-hours and happens $u$ times a month; the modeling loop and signoff runs consume $L$ compute-hours a month. The fraction spent on updates is $uU / (uU + L)$. With 6 updates a month at 5 compute-hours each and 2,000 compute-hours of loop work, the fraction is $30 / 2030 = 1.48\%$. Selection, a text decision over manifests, is smaller still. The overhead would become significant only if updates were frequent and expensive, such as a core edit with full regression every day.

## Break-even for a customization

A customization is worth making if it saves more than it costs. If a skill update costing $U$ hours makes each subsequent run $s$ hours faster, or saves the equivalent effort through better QoR, it breaks even after $\lceil U / s \rceil$ runs. An update costing 5 hours that saves 0.4 hours per run breaks even after 13 runs. If the flow will run only 5 more times, the update is not worth it on runtime grounds. An agent that customizes should estimate this before starting, from the expected number of remaining runs and the measured gain on the fast loop.

## Going deeper

The agent-first EDA notes in this curriculum, section on fast and accurate engines and the runtime economics of skill updates.
