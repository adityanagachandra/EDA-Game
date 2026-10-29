// Package csim generates circuit-simulation problems: modified nodal analysis, the DC
// operating point by Newton-Raphson, transient integration and step control, AC and noise
// analysis, periodic steady state, sensitivity and variation, model order reduction, fast
// SPICE and the agent-first frontier of the simulator.
package csim

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the circuit-simulation generators to reg.
func Register(reg *generator.Registry) {
	// Netlists and decks.
	reg.Register("csim.tool.spice_netlist", edakit.Func(spiceNetlist))
	reg.Register("csim.tool.scale_suffixes", edakit.Func(scaleSuffixes))
	reg.Register("csim.tool.model_cards", &edakit.Bank{Items: modelCardItems})
	reg.Register("csim.tool.ngspice_xyce", &edakit.Bank{Items: ngspiceXyceItems})

	// Modified nodal analysis.
	reg.Register("csim.mna.stamps", edakit.Func(mnaStamps))
	reg.Register("csim.mna.current_sources", edakit.Func(currentSources))
	reg.Register("csim.mna.voltage_sources", edakit.Func(voltageSources))
	reg.Register("csim.mna.vccs", edakit.Func(vccsStamps))
	reg.Register("csim.mna.vcvs", edakit.Func(vcvsStamps))
	reg.Register("csim.mna.solve_small", edakit.Func(solveSmall))
	reg.Register("csim.mna.structure", edakit.Func(mnaStructure))

	// DC operating point.
	reg.Register("csim.dc.op_reduction", edakit.Func(opReduction))
	reg.Register("csim.dc.diode_companion", edakit.Func(diodeCompanion))
	reg.Register("csim.dc.newton_companion", edakit.Func(newtonCompanion))
	reg.Register("csim.dc.newton_convergence", edakit.Func(newtonConvergence))
	reg.Register("csim.dc.mosfet_linearization", edakit.Func(mosfetLinearization))
	reg.Register("csim.dc.junction_limiting", edakit.Func(junctionLimiting))
	reg.Register("csim.dc.damping", edakit.Func(damping))
	reg.Register("csim.dc.gmin_stepping", edakit.Func(gminStepping))
	reg.Register("csim.dc.convergence_aids", &edakit.Bank{Items: convergenceAidItems})

	// Transient analysis.
	reg.Register("csim.tran.dae", &edakit.Bank{Items: daeItems})
	reg.Register("csim.tran.companion_c", edakit.Func(companionC))
	reg.Register("csim.tran.companion_l", edakit.Func(companionL))
	reg.Register("csim.tran.be_step", edakit.Func(beStep))
	reg.Register("csim.tran.trap_step", edakit.Func(trapStep))
	reg.Register("csim.tran.trap_ringing", edakit.Func(trapRinging))
	reg.Register("csim.tran.gear_bdf", edakit.Func(gearBDF))
	reg.Register("csim.tran.lte_estimate", edakit.Func(lteEstimate))
	reg.Register("csim.tran.timestep_control", edakit.Func(timestepControl))
	reg.Register("csim.tran.breakpoints", &edakit.Bank{Items: breakpointItems})

	// AC and noise.
	reg.Register("csim.ac.linearize_op", edakit.Func(linearizeOP))
	reg.Register("csim.ac.complex_mna", edakit.Func(complexMNA))
	reg.Register("csim.ac.small_signal", edakit.Func(smallSignal))
	reg.Register("csim.ac.noise", edakit.Func(noise))

	// Periodic steady state.
	reg.Register("csim.pss.shooting", edakit.Func(shooting))
	reg.Register("csim.pss.harmonic_balance", edakit.Func(harmonicBalance))
	reg.Register("csim.pss.rf_methods", &edakit.Bank{Items: rfMethodItems})

	// Sensitivity and variation.
	reg.Register("csim.var.sensitivity", edakit.Func(sensitivity))
	reg.Register("csim.var.adjoint", edakit.Func(adjoint))
	reg.Register("csim.var.mismatch", edakit.Func(mismatch))
	reg.Register("csim.var.monte_carlo_corners", edakit.Func(monteCarloCorners))
	reg.Register("csim.var.process_corners", &edakit.Bank{Items: processCornerItems})

	// Model order reduction.
	reg.Register("csim.mor.moments", edakit.Func(moments))
	reg.Register("csim.mor.awe", edakit.Func(awe))
	reg.Register("csim.mor.prima", edakit.Func(prima))
	reg.Register("csim.mor.macromodel", &edakit.Bank{Items: macromodelItems})

	// Fast SPICE and mixed-signal.
	reg.Register("csim.fast.multirate", edakit.Func(multirate))
	reg.Register("csim.fast.table_models", &edakit.Bank{Items: tableModelItems})
	reg.Register("csim.fast.mixed_signal", &edakit.Bank{Items: mixedSignalItems})

	// Frontier.
	reg.Register("csim.frontier.convergence_agent", &edakit.Bank{Items: convergenceAgentItems})
	reg.Register("csim.frontier.option_tiers", &edakit.Bank{Items: optionTierItems})
	reg.Register("csim.frontier.fast_accurate_sizing", edakit.Func(fastAccurateSizing))
	reg.Register("csim.frontier.surrogate_check", edakit.Func(surrogateCheck))
	reg.Register("csim.frontier.device_seam", &edakit.Bank{Items: deviceSeamItems})
}
