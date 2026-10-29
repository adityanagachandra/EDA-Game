// Package analog generates analog, mixed-signal and RF problems for an EDA tool designer:
// small-signal design with gm/Id, amplifier building blocks and stability, matching and analog
// layout, sizing as optimization, data converters, PLLs, RF cascades, mixed-signal modeling,
// specialized branches and the agent-first frontier. Files are split by subdomain.
package analog

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the analog/mixed-signal/RF generators to reg.
func Register(reg *generator.Registry) {
	// Small-signal foundations.
	reg.Register("analog.ss.automation_gap", &edakit.Bank{Items: automationGapItems})
	reg.Register("analog.ss.small_signal_model", edakit.Func(smallSignal))
	reg.Register("analog.ss.db_conversions", edakit.Func(dbConversions))
	reg.Register("analog.ss.intrinsic_gain", edakit.Func(intrinsicGain))
	reg.Register("analog.ss.gm_id", edakit.Func(gmID))
	reg.Register("analog.ss.lookup_tables", edakit.Func(lookupTables))

	// Amplifiers.
	reg.Register("analog.amp.common_source", edakit.Func(commonSource))
	reg.Register("analog.amp.current_mirror", edakit.Func(currentMirror))
	reg.Register("analog.amp.diff_pair", edakit.Func(diffPair))
	reg.Register("analog.amp.opamp_two_stage", edakit.Func(twoStage))
	reg.Register("analog.amp.miller_compensation", edakit.Func(miller))
	reg.Register("analog.amp.stability_phase_margin", edakit.Func(phaseMargin))
	reg.Register("analog.amp.comparator", edakit.Func(comparator))

	// Matching and layout.
	reg.Register("analog.match.pelgrom", edakit.Func(pelgrom))
	reg.Register("analog.match.offset_budget", edakit.Func(offsetBudget))
	reg.Register("analog.layout.common_centroid", edakit.Func(commonCentroid))
	reg.Register("analog.layout.interdigitation", edakit.Func(interdig))
	reg.Register("analog.layout.parasitics", edakit.Func(parasitics))
	reg.Register("analog.layout.automation", &edakit.Bank{Items: layoutAutoItems})

	// Sizing.
	reg.Register("analog.size.optimization", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: sizeOptItems}, edakit.Func(feasibleChoice), &edakit.Bank{Items: sizeOptItems},
	}})
	reg.Register("analog.size.equation_based", edakit.Func(equationBased))
	reg.Register("analog.size.simulation_loop", edakit.Func(simLoop))
	reg.Register("analog.size.bayesian_opt", edakit.Func(bayesOpt))
	reg.Register("analog.size.pvt_robustness", edakit.Func(pvtRobust))

	// Data converters.
	reg.Register("analog.conv.adc", edakit.Func(adc))
	reg.Register("analog.conv.quantization_snr", edakit.Func(quantSNR))
	reg.Register("analog.conv.sar", edakit.Func(sar))
	reg.Register("analog.conv.sampling", edakit.Func(sampling))
	reg.Register("analog.conv.dac", edakit.Func(dac))
	reg.Register("analog.conv.enob_fom", edakit.Func(enobFOM))

	// PLLs.
	reg.Register("analog.pll.vco", edakit.Func(vco))
	reg.Register("analog.pll.loop", edakit.Func(pllLoop))
	reg.Register("analog.pll.loop_dynamics", edakit.Func(loopDynamics))
	reg.Register("analog.pll.phase_noise", edakit.Func(phaseNoise))

	// RF.
	reg.Register("analog.rf.s_parameters", edakit.Func(sParams))
	reg.Register("analog.rf.noise_figure", edakit.Func(noiseFigure))
	reg.Register("analog.rf.friis", edakit.Func(friis))
	reg.Register("analog.rf.iip3", edakit.Func(iip3))

	// Mixed-signal.
	reg.Register("analog.ms.simulation", &edakit.Bank{Items: msSimItems})
	reg.Register("analog.ms.verilog_a", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: vaItems}, edakit.Func(vaCompute), &edakit.Bank{Items: vaItems},
	}})
	reg.Register("analog.ms.real_number", edakit.Func(realNumber))
	reg.Register("analog.ms.open_flow", &edakit.Bank{Items: openFlowItems})

	// Specialized branches.
	reg.Register("analog.special.photonics", edakit.Func(photonics))
	reg.Register("analog.special.mems", edakit.Func(mems))
	reg.Register("analog.special.power", edakit.Func(power))
	reg.Register("analog.special.neuromorphic", edakit.Func(neuromorphic))

	// Frontier.
	reg.Register("analog.frontier.sizing_agent", &edakit.Mix{Parts: []generator.Generator{
		&edakit.Bank{Items: sizingAgentItems}, edakit.Func(agentLogChoice), &edakit.Bank{Items: sizingAgentItems},
	}})
	reg.Register("analog.frontier.surrogate_spice", edakit.Func(surrogateSpice))
	reg.Register("analog.frontier.layout_agent", &edakit.Bank{Items: layoutAgentItems})
	reg.Register("analog.frontier.migration_agent", edakit.Func(migration))
	reg.Register("analog.frontier.reward", edakit.Func(reward))
}
