# On-site, offline deployment

Many design teams will not let their data leave their network, and some sites have no outside connection at all. An agent-first EDA system must therefore build, run, test and learn entirely inside the customer's environment. That requirement reaches into the toolchain, the model, the packaging of engines and the update process, and it must be designed in from the start.

## Pinned toolchains and offline builds

The harness rebuilds plug-ins on-site, so the compiler, standard library, build system, compiler cache and every third-party dependency must be present at the site in pinned versions. A build that downloads a package at configure time fails at an air-gapped site and is not reproducible anywhere. Pinning also matters for correctness: a different compiler version can change floating-point code generation and with it the results that goldens check. The deliverable is a self-contained environment, such as a container image or a versioned toolchain bundle, whose identity is recorded in every run's provenance.

## Hosting the model inside the boundary

If the model runs as a hosted API, every prompt (including source excerpts, design summaries and reports) leaves the site. For an on-site deployment the model must run on hardware inside the boundary, which requires a model whose weights can be obtained and run locally under a license that permits it, served by an inference runtime such as vLLM or SGLang. The harness is unchanged: it calls the same tool-call interface, now pointed at a local server. Model updates arrive as reviewed artifacts like any other software, not as silent changes behind an API.

## Packaging engines for a site

Engines whose source is the vendor's sensitive IP can ship as compiled libraries that expose only extension points, manifests and headers for the seams. The customer, and the customer's model, customize through configuration and plug-ins without seeing the core source. Engines that are not sensitive may ship as source. The vendor decides this per engine, and the decision fixes what the on-site model can edit: Tier 2 plug-ins against published headers are always possible, while Tier 3 edits to a compiled-only core happen only at the vendor. Updates are delivered as signed, versioned bundles that the site installs, tests against its own goldens and can roll back.

## Going deeper

The vLLM and SGLang documentation on self-hosted serving; the agent-first EDA notes in this curriculum, section on verification and the IP boundary.
