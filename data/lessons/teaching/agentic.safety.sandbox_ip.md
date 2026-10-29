# Sandboxing and the IP boundary

Agent-written code runs inside engines that hold some of the most valuable data a company owns: unreleased designs and foundry process information. Two architectural rules make that acceptable. Every skill and every candidate change runs in a sandbox with hard limits and no network. And design and process data never leave the site where they live; only skills and know-how move. Both are properties of the system, enforced by code, not policies in a document.

## What the sandbox enforces

A sandbox is an execution environment whose capabilities are restricted by the operating system, not by the code inside it. For an EDA skill it enforces: no network access; file access only to the skill's workspace and the read-only inputs it was given; CPU, memory and wall-time limits; no ability to spawn arbitrary processes or load unapproved libraries. Containers with resource controls, seccomp filters or a microVM provide these. The point is that a buggy plug-in that loops forever is killed at its timeout, one that allocates without bound hits its memory limit, and one that tries to open a socket fails, whatever the code says.

## The IP boundary as architecture

The IP boundary separates the customer's data (RTL, netlists, layouts, timing constraints, the foundry's PDK and models) from everything outside the site. In an agent-first deployment it shapes the whole system: the model runs on-site on hardware inside the boundary; reward computation and evaluation run where the data is; logs and traces that leave the site are scrubbed of design content; and the engines that are themselves the vendor's sensitive IP can ship as compiled libraries, exposing only extension points and manifests. The question for each data flow is which side of the boundary it crosses, and every crossing is deliberate and reviewed.

## What flows back and what never does

When a skill is created or improved at a site, the vendor wants the skill and the know-how: the plug-in's code, its manifest, its golden tests on synthetic or public designs, and the pattern of when it helped. The customer's design never flows back, and neither do the customer's process data or their proprietary constraints. Golden tests for a promoted skill must therefore be rebuilt on designs the vendor may hold, such as open designs on an open PDK, rather than copied from the customer's suite. A skill whose usefulness can be demonstrated only on customer data stays at that customer.

## Going deeper

The Linux kernel documentation on cgroups and seccomp; the agent-first EDA notes in this curriculum, section on verification and the IP boundary.
