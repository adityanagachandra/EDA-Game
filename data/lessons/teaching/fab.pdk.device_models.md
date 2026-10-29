# Device models and corners in a PDK

A PDK's device models describe how each transistor behaves across voltage, temperature and process variation. Corners bound the global spread, statistical models describe mismatch, and layout-dependent parameters carry context. Every timing library, every circuit simulation and every power estimate is computed from these models, so their structure decides what margin a design must carry.

## Corners for global variation

Transistors vary from lot to lot and wafer to wafer. Foundries summarize this with process corners named by NMOS then PMOS speed: TT (typical), SS (slow-slow), FF (fast-fast), and the skewed SF and FS. Skewed corners matter for circuits that depend on the ratio of NMOS to PMOS strength, such as SRAM cells, level shifters and ratioed logic. Each corner is combined with supply voltage and temperature to define a PVT corner, and Liberty libraries are characterized at each. Slow corners at low voltage bound setup timing; fast corners at high voltage and temperature bound leakage and hold.

## Statistical models and mismatch

Corners stack worst cases and can be pessimistic, or miss combinations. Statistical models describe parameter distributions so Monte Carlo simulation can sample real behaviour. They separate global variation, shared by all devices on a die, from local mismatch between nominally identical neighbours, which comes from random dopant fluctuation, line-edge roughness and work-function granularity. Pelgrom's law says mismatch standard deviation scales as $1/\sqrt{WL}$: quadrupling gate area halves mismatch sigma. Library characterization uses these models to produce variation-aware timing data, such as Liberty variation format tables. Those tables let a timer carry sigma per arc instead of a flat margin.

## Layout dependence, aging and model maturity

A transistor's threshold and mobility depend on its surroundings: distance to the well edge (well proximity effect), length of diffusion and distance to diffusion breaks (stress effects), and neighbouring gate density. Extraction measures these distances and passes them as instance parameters so simulation sees the layout context. Reliability models add end-of-life shifts from bias temperature instability and hot-carrier injection, which feed aging-aware timing margins. Device corners and interconnect RC corners are separate because FEOL and BEOL vary from different steps. Early models are projections and are re-fitted as silicon data arrives, so timing results carry the model version.

## Going deeper

Pelgrom, Duinmaijer and Welbers, "Matching properties of MOS transistors", IEEE Journal of Solid-State Circuits (1989); the BSIM-CMG technical manual from the BSIM group at UC Berkeley.
