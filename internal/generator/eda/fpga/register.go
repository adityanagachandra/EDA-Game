package fpga

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the FPGA CAD generators to reg.
func Register(reg *generator.Registry) {
	// Architecture.
	reg.Register("fpga.arch.lut", edakit.Func(lutBits))
	reg.Register("fpga.arch.ble_clb", edakit.Func(bleClb))
	reg.Register("fpga.arch.carry_chain", edakit.Func(carryChain))
	reg.Register("fpga.arch.hard_blocks", edakit.Func(hardBlocks))
	reg.Register("fpga.arch.routing_architecture", edakit.Func(routingArch))
	reg.Register("fpga.arch.segments", edakit.Func(segments))

	// Synthesis and LUT mapping.
	reg.Register("fpga.map.flow", &edakit.Bank{Items: mapFlowItems})
	reg.Register("fpga.map.wide_functions", edakit.Func(wideFunctions))
	reg.Register("fpga.map.lut_cuts", edakit.Func(lutCuts))
	reg.Register("fpga.map.area_flow", edakit.Func(areaFlow))
	reg.Register("fpga.map.inference", &edakit.Bank{Items: inferenceItems})

	// Packing.
	reg.Register("fpga.pack.ble_formation", edakit.Func(bleFormation))
	reg.Register("fpga.pack.clustering", edakit.Func(clustering))
	reg.Register("fpga.pack.vpack_attraction", edakit.Func(vpackAttraction))

	// Placement.
	reg.Register("fpga.place.bbox_cost", edakit.Func(bboxCost))
	reg.Register("fpga.place.vpr_annealing", edakit.Func(vprAnnealing))
	reg.Register("fpga.place.temperature_schedule", edakit.Func(temperatureSchedule))
	reg.Register("fpga.place.range_limit", edakit.Func(rangeLimit))
	reg.Register("fpga.place.timing_cost", edakit.Func(timingCost))
	reg.Register("fpga.place.analytical", edakit.Func(analytical))

	// Routing.
	reg.Register("fpga.route.rr_graph", edakit.Func(rrGraph))
	reg.Register("fpga.route.maze_expansion", edakit.Func(mazeExpansion))
	reg.Register("fpga.route.pathfinder", edakit.Func(pathfinder))
	reg.Register("fpga.route.history", edakit.Func(history))
	reg.Register("fpga.route.timing_driven", edakit.Func(timingDriven))
	reg.Register("fpga.route.min_channel_width", edakit.Func(minChannelWidth))

	// Timing, bitstreams and reconfiguration.
	reg.Register("fpga.impl.delay_model", edakit.Func(delayModel))
	reg.Register("fpga.impl.bitstream", edakit.Func(bitstream))
	reg.Register("fpga.impl.partial_reconfig", edakit.Func(partialReconfig))

	// Toolchains.
	reg.Register("fpga.flow.open_toolchain", &edakit.Bank{Items: openToolchainItems})
	reg.Register("fpga.flow.bitstream_docs", &edakit.Bank{Items: bitstreamDocsItems})
	reg.Register("fpga.flow.vendor_flows", &edakit.Bank{Items: vendorFlowItems})

	// Frontier.
	reg.Register("fpga.frontier.seams", &edakit.Bank{Items: seamItems})
	reg.Register("fpga.frontier.option_tuning", edakit.Func(optionTuning))
	reg.Register("fpga.frontier.ai_placement", &edakit.Bank{Items: aiPlacementItems})
	reg.Register("fpga.frontier.bitstream_golden", &edakit.Bank{Items: bitstreamGoldenItems})
	reg.Register("fpga.frontier.fast_timing", edakit.Func(fastTiming))
}
