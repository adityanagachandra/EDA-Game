# Test pattern languages: STIL

ATPG tools and testers come from different vendors and speak different languages, so patterns travel in a standard interchange format. STIL, the Standard Test Interface Language (IEEE 1450), describes signals, timing, scan structures and pattern data in a tester-independent text form. Knowing what it represents, and what a converter must translate, matters to anyone building ATPG back ends or debugging a pattern that passes simulation and fails on the tester.

## What a STIL file describes

A STIL file declares **Signals** (the device pins and their directions) and **SignalGroups** (named buses such as all scan-in pins). A **Timing** block holds **WaveformTables**, which map waveform characters to timed events within a tester cycle: drive 0 might mean force low at 0 ns, H might mean compare high at a 40 ns strobe. **ScanStructures** describe the chains, their scan-in and scan-out pins and lengths. **Procedures** such as load_unload and capture describe how to shift and how to pulse clocks; **Pattern** blocks then call those procedures with chain data and parallel vectors. **PatternBurst** and **PatternExec** say which patterns run, in what order and with which timing.

## Stimulus, response and masks

Vector characters distinguish what the tester drives from what it checks: 0 and 1 drive, L and H compare low and high, and X (or N in some tables) means no compare. When ATPG simulation predicts an unknown output, the pattern must carry X there, or good parts fail at random. Describing scan with procedures rather than expanding every shift cycle keeps files compact and lets the converter map chain data into the tester's deep **scan memory** while surrounding vectors go into main vector memory.

## Translation to the tester

A tester converter reads STIL, maps signals to tester channels, turns waveform tables into the tester's timing sets, and loads vectors and scan data. Errors here are systematic: if a strobe is placed before the device's output settles, every die fails in the same cycle, which is the signature of a timing translation problem rather than a defect. Older formats such as WGL served the same role in vendor-specific ways; STIL's open standard and explicit scan constructs reduced that translation risk. CTL (IEEE 1450.6) extends STIL with core test information for IEEE 1500 test reuse.

## Going deeper

IEEE Std 1450, "Standard Test Interface Language"; IEEE Std 1450.6, "Core Test Language"; Wang, Wu and Wen (2006), chapter on test data formats.
