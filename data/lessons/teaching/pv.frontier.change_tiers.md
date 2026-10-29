# Change tiers for verification engines

Verification engines follow the same four-layer model as implementation engines: an immutable geometry and matching core, policies such as check plug-ins and device extractors, parameters such as rule values and switches, and observation hooks over results. Requests arrive constantly from foundries, customers and agents. Landing each at the smallest safe tier keeps the core stable and the engine trustworthy.

## Tier 1 changes in verification

Most verification changes are data. A new deck version raising metal 2 spacing from 40 to 42 nm changes a rule value. A waiver for a known-good IP violation is an entry in a waiver file. Running only antenna checks is a rule selection, and a new density window size and step are check parameters. None of these touch engine code, and none need a rebuild. In a well-designed verification engine the deck is the configuration layer, so a foundry update should be Tier 1 by construction.

## Tier 2 plug-ins

Some requests need new logic but not new core. A via spacing rule that depends on how many neighbouring vias lie within 100 nm, which the deck language cannot express, can be a custom check written against the engine's spatial query API. A custom capacitor needs a device-extractor plug-in that names its recognition layers and terminals. A team-specific property comparison is a hook in LVS. Exporting results as typed records for an agent is an observation hook. Each plug-in sees plain data, runs in a sandbox, and must pass golden tests on reference layouts.

## Tier 3 core edits

Core edits change what every check relies on: the Boolean and sizing kernel, exact predicates, hierarchy and context handling, the coordinate width, the matching algorithm. A wrong result on 45-degree edges is a core bug, fixed with a human-approved diff and full regression, never patched by a plug-in that post-processes outputs. A seam is in the right place when Tier 3 edits essentially never happen for a class of request. If customers keep needing core edits for variants of end-of-line rules, the engine is missing a parameterized rule family or a check extension point.
