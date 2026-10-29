# The GCD design as a first run

GCD, a small greatest-common-divisor circuit, is the design every new ORFS user runs first. It is small enough to finish in minutes yet exercises every stage, so it validates the installation and gives a baseline. A first run teaches the configuration file, the floorplan arithmetic that turns utilization into a core, and where the results land, which is the foundation for every later experiment on larger designs.

## The design configuration

Each design in ORFS has a config.mk per platform, under designs/<platform>/<design>/. For GCD it sets DESIGN_NAME (the top module), PLATFORM, VERILOG_FILES (the RTL), and SDC_FILE (the constraints, which define the clock). Two knobs shape the physical result. CORE_UTILIZATION sets how much of the core area the standard cells should fill, as a percentage; the floorplan derives the core from it. PLACE_DENSITY sets the target density that global placement spreads cells toward. Everything about the process, such as layers, rows, PDN and libraries, comes from the platform's own files, which is why the same GCD runs on several platforms by changing one path. Running the flow is a make invocation that names the design's config; it then executes the six stages in order.

## Sizing the core from utilization

Utilization is cell area divided by core area, so the core area is $A_{core} = A_{cells}/u$. Example: GCD on sky130hd synthesizes to 2400 square microns of cells and CORE_UTILIZATION is 45, so $A_{core} = 2400/0.45 = 5333.3$ square microns. A square core has side $\sqrt{5333.3} = 73.03$ microns. The sky130hd row height is 2.72 microns, so $\lfloor 73.03/2.72 \rfloor = 26$ whole rows fit. Rows on nangate45 are 1.4 microns tall and on asap7 0.27 microns, so the same arithmetic gives very different row counts. Common slips are multiplying by the utilization instead of dividing, which shrinks the core below the cells' own area, and forgetting the square root. The real floorplan also snaps the core to the site grid and leaves a margin between core and die for the power ring and IO.

## Reading the first results

ORFS writes everything per platform and design: logs/ holds each step's log, reports/ holds timing, area and power reports, results/ holds stage databases, netlists and the final GDS, and objects/ holds intermediate files such as the prepared libraries. Read the final report first: worst and total negative slack, design area and utilization, power, and DRC count after routing. If setup slack is negative, the target being missed is the clock period in the SDC; relax it or let the repair steps work harder. If global routing reports overflow, demand exceeds track supply, so lower CORE_UTILIZATION to give the router space. Finally, run it twice: identical inputs should give identical results, and that reproducible baseline is what every later change is judged against.

## Going deeper

The OpenROAD-flow-scripts tutorial and documentation for the GCD design and config.mk variables.
