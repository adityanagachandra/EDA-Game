// Package agentic generates problems for agent-first EDA: engine architecture, skills and
// selection, the harness loop, change tiers and seams, verification of engine changes, safety,
// the skill flywheel and fast-accurate loops, models and their evaluation, tool-grounded
// learning, and orchestration. Files are split by subdomain; this one holds Register.
package agentic

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func bank(items []edakit.Item) *edakit.Bank { return &edakit.Bank{Items: items} }

func mix(parts ...generator.Generator) *edakit.Mix { return &edakit.Mix{Parts: parts} }

// Register adds the agent-first EDA generators to reg.
func Register(reg *generator.Registry) {
	reg.Register("agentic.engine.callable_engines", bank(callableItems))
	reg.Register("agentic.engine.four_layer", bank(fourLayerItems))
	reg.Register("agentic.engine.core_vs_surface", bank(coreSurfaceItems))
	reg.Register("agentic.engine.observation_hooks", mix(bank(hookItems), edakit.Func(hookOverhead), bank(hookItems)))
	reg.Register("agentic.skill.manifest", bank(manifestItems))
	reg.Register("agentic.skill.selection", mix(bank(selectionItems), edakit.Func(selectionScore), bank(selectionItems)))
	reg.Register("agentic.skill.selection_cost", edakit.Func(selectionCost))

	reg.Register("agentic.harness.loop", bank(harnessItems))
	reg.Register("agentic.harness.build_budget", edakit.Func(buildBudget))
	reg.Register("agentic.harness.retry_expectation", edakit.Func(retryExpectation))
	reg.Register("agentic.harness.feedback", bank(feedbackItems))
	reg.Register("agentic.harness.check_ordering", edakit.Func(checkOrdering))

	reg.Register("agentic.tiers.change", bank(tierItems))
	reg.Register("agentic.tiers.turnaround", edakit.Func(tierTurnaround))
	reg.Register("agentic.seam.extension_point", bank(extensionItems))
	reg.Register("agentic.seam.granularity", edakit.Func(seamGranularity))
	reg.Register("agentic.seam.context_design", bank(contextItems))
	reg.Register("agentic.seam.abi_versioning", bank(abiItems))

	reg.Register("agentic.verify.golden_differential", mix(bank(goldenItems), edakit.Func(goldenDiffCalc), bank(goldenItems)))
	reg.Register("agentic.verify.qor_baseline", edakit.Func(qorBaseline))
	reg.Register("agentic.verify.seed_noise", edakit.Func(seedNoise))
	reg.Register("agentic.verify.missing_golden", bank(missingGoldenItems))
	reg.Register("agentic.verify.promotion_gate", edakit.Func(promotionGate))

	reg.Register("agentic.safety.sandbox_ip", bank(sandboxItems))
	reg.Register("agentic.safety.resource_budget", edakit.Func(resourceBudget))
	reg.Register("agentic.safety.determinism", bank(determinismItems))
	reg.Register("agentic.safety.onsite_offline", bank(onsiteItems))
	reg.Register("agentic.safety.untrusted_inputs", bank(untrustedItems))

	reg.Register("agentic.flywheel.skill_promotion", bank(promotionItems))
	reg.Register("agentic.flywheel.amortization", edakit.Func(amortization))
	reg.Register("agentic.brain.domain_model", bank(brainItems))
	reg.Register("agentic.loop.fast_accurate", mix(bank(fastAccItems), edakit.Func(loopBudget), bank(fastAccItems)))
	reg.Register("agentic.loop.fast_mode_calibration", edakit.Func(calibration))
	reg.Register("agentic.loop.recursion_decision", edakit.Func(recursion))

	reg.Register("agentic.model.three_roles", bank(rolesItems))
	reg.Register("agentic.model.abstraction_layer", bank(abstractionItems))
	reg.Register("agentic.eval.pass_at_k", edakit.Func(passAtK))
	reg.Register("agentic.eval.edit_correctness", edakit.Func(editCorrectness))
	reg.Register("agentic.eval.skill_benchmarks", bank(benchItems))

	reg.Register("agentic.rl.tool_grounded", bank(rlItems))
	reg.Register("agentic.rl.reward_design", edakit.Func(rewardDesign))
	reg.Register("agentic.rl.reward_hacking", bank(hackItems))
	reg.Register("agentic.rl.proxy_signoff_gap", edakit.Func(proxyGap))

	reg.Register("agentic.mcp.tool_interface", bank(mcpItems))
	reg.Register("agentic.orch.multi_agent", bank(multiItems))
	reg.Register("agentic.signoff.eco_loop", bank(ecoItems))
	reg.Register("agentic.signoff.eco_convergence", edakit.Func(ecoConvergence))
	reg.Register("agentic.dse.exploration_agent", edakit.Func(exploration))
	reg.Register("agentic.landscape.industry", bank(landscapeItems))
}
