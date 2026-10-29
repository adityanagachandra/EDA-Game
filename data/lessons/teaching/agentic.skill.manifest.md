# Skills and manifests

A skill is an engine, or a piece of one, packaged so that an agent can find it, run it, change it and prove the change correct. The package is ordinary: source code, a build, golden tests and a short plain-language manifest. The manifest is the part that makes the skill usable by an agent, because it is what the agent reads instead of the source.

## What a skill folder contains

A skill folder holds four things. The **source** is the engine code or the plug-in, for example a detailed-placement move set or a derate provider. The **build** compiles it as a separate unit, typically a shared library at a defined extension point, so that changing it rebuilds only that unit. The **golden tests** are known-good inputs and outputs: small reference designs with the expected result, plus a QoR baseline for differential comparison. The **manifest** is a short text file that says what the skill does, what it consumes and produces, and which parts are meant to be tuned. Nothing here is exotic. It is the code an engine developer already writes, with a description and tests next to it.

## What a manifest must say

A good manifest answers the questions an agent asks when choosing and customizing. Purpose: "detailed placement for designs with regular datapath cells; abuts cells in the same bit slice". Inputs and outputs: a legal placement in, a legal placement out, with HPWL and displacement reported. Applicability and limits: "assumes a single row height; not tested with multi-height cells". Tunable surface: named parameters with types and ranges, such as `max_displacement` in sites from 1 to 20, and the extension points it offers. Cost: typical runtime per million cells. A manifest that lists parameters without ranges, or omits its limits, invites the agent to call the skill where it does not apply.

## Manifests that mislead

Because selection is made from manifests alone, a wrong manifest is as harmful as a wrong algorithm. Three failures recur. A **stale** manifest describes behavior that a later code change removed. An **overclaiming** manifest says "works for all designs" when the golden tests cover only one style. A manifest with **missing limits** omits the condition under which the skill degrades, such as high utilization. The remedy is to tie the manifest to the tests: every claimed capability is exercised by a golden test, the manifest is versioned with the code, and a change that alters behavior must update both.

## Going deeper

The agent-first EDA notes in this curriculum, section on skills, manifests and selection.
