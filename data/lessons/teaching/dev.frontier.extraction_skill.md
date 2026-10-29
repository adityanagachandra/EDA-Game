# Parameter extraction as a Tier 1 skill

Compact-model parameter extraction is a near-ideal agent task: the inputs are data, the outputs are parameter values, the model equations are fixed, and a simulator scores the result. In the four-layer engine view it is configuration, Tier 1. Packaging it as a skill with a manifest, golden tests and a clear boundary for when a request stops being a fit and becomes a code change is what makes it safe to automate.

## Which tier a modeling request lands in

Change tiers apply directly to device modeling.

Tier 1, configuration: new values in a model card, with model equations unchanged. Refitting BSIM-CMG parameters to new wafer data, retargeting a corner card, or adjusting temperature coefficients are all Tier 1. No recompile, and the most common request.

Tier 2, plug-in: a new behavior at a defined extension point, such as a threshold-shift term that depends on a layout-proximity quantity the model lacks, or a custom aging law. It is written as a model extension, for example a Verilog-A module or a compiled plug-in, and the simulator core is untouched.

Tier 3, core: a change to the simulator itself, such as Newton step limiting, matrix assembly or time-step control. It affects every circuit, needs human approval and full regression, and should be rare.

An agent fitting parameters that discovers the fit only works with a new equation must stop and escalate to the plug-in path, with its own tests, rather than quietly editing the model during a Tier 1 task.

## The extraction skill and its manifest

A skill packages the extraction engine with a plain-language manifest that an agent can read without opening the source. For extraction it should declare: the supported model and version (BSIM4 4.8, BSIM-CMG 110, and so on); the required measurements (I-V families, C-V, temperatures, geometries); which parameters it fits, in what stages, with what bounds; the tunables (stage weights, error metrics, bias ranges); and the fit report it returns.

Selection is a cheap text decision: match the device type and model version in a request against manifests. Running a planar extraction on FinFET data would fit the wrong equations, however low the error.

Good extraction skills fit in stages, threshold first, then mobility and short-channel terms, then parasitics, using the bias region where each parameter dominates, and expose bounds so that fits stay physical. These choices are exactly the tunables a manifest should advertise, since different processes need different ranges.

## Verifying and promoting extraction skills

Every extraction skill ships golden tests: a reference dataset whose fit must reproduce a known-good card and error report within tolerance. The harness reruns them on every edit to the skill.

A Tier 1 refit that lowers fit error still needs a differential check before promotion: benchmark circuits (ring oscillators, inverter chains, a few cells) simulated with the old and new cards must move only as intended. A better device fit can shift circuit behavior in a region nobody looked at.

Measured wafer data is customer IP. The skill runs on-site, inside a sandbox with no network, and only the reusable improvement (a better staging sequence, a better bound set) can leave as know-how. The customer's data and cards do not.

An improved sequence found for one tricky process goes through the flywheel: iterate as a plug-in, validate with golden and differential tests, then promote it into the shared skill library, where its manifest makes it selectable for future requests.

## Going deeper

The repository's docs/eda/agent-first-concepts.md, sections 2, 4 and 7; Schroder, "Semiconductor Material and Device Characterization" (3rd ed.), chapter 4, for the underlying extraction methods.
