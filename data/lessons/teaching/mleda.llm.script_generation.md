# LLMs for flow-script generation

EDA tools are driven by Tcl and Python scripts against documented command APIs. Translating a request such as "report the ten worst setup paths after placement" into commands is a natural LLM task, and the tool itself can check the result. The risks are invented commands, scripts that do more than asked, and scripts that run cleanly and do the wrong thing.

## Scripts as tool programs

A flow script is a program against a tool's command set and database API, run in a particular flow state: design loaded, constraints read, placement done. A correct answer to a reporting request calls the right timing-report command with the right options in that state. Editing scripts, such as Python that resizes cells through the database API, must also respect database invariants and keep dependent analyses, such as incremental timing, up to date. Commands and options change between tool versions, so retrieving the installed version's command reference and example scripts into the prompt is the most effective grounding.

## Hallucinated commands and silent errors

Invented options and commands are a common failure, and the tool is the authority on what exists: a harness validates each command against the tool's own registry or help output, or runs a fast dry run, before launching a long flow. A script that runs without errors still has to be checked for its effect: the reports or design state should show the requested result, such as the expected utilization or path count. Scripts generated for one design often hardcode a die size or macro names; reused elsewhere they fail or silently do the wrong thing, so values should come from the database or parameters.

## Sandboxing and checking a generated flow

A generated script can run arbitrary shell commands, so it runs in a sandbox with no network access, limited file-system scope and CPU and time limits: design data must not leave and the workspace must not be damaged. The harness also flags scripts whose blast radius exceeds the request, such as rebuilding a whole power grid to change one strap pitch. Comparisons between script variants need fixed seeds or several runs, since run-to-run noise can look like an improvement. A script that only sets parameters is a Tier 1 configuration change, reversible and without recompilation, but it still needs a run and a QoR comparison before adoption.

## Going deeper

The OpenROAD documentation for its Tcl and Python APIs; the OpenROAD-flow-scripts documentation for flow structure.
