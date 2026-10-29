# Boundary scan and JTAG

Boundary scan, standardized as IEEE 1149.1 and known as JTAG, puts a scan cell on every digital pin of a chip and links them into a shift register reached through a four- or five-wire test access port. Its first purpose was testing the solder joints and wires of a board without physical probes; it has since become the universal low-pin-count access path to on-chip test, debug and configuration logic.

## The architecture

A compliant device has a **test access port** (TAP) with TCK (clock), TMS (mode select), TDI (data in), TDO (data out) and optionally TRST. A 16-state TAP controller, steered by TMS, selects between scanning the **instruction register** and scanning one of several **data registers**: the **boundary register** of pin cells, the 1-bit **bypass register**, and usually a 32-bit **device identification register**. On a board, devices are chained TDI to TDO. The data-register path length is the sum over devices of whichever register each device has selected: a 510-cell boundary register in EXTEST plus a device in BYPASS contributes $510 + 1$.

## Interconnect testing

In EXTEST, the boundary cells of driving devices put values on board nets and the cells of receiving devices capture them. Each net receives one bit per parallel test vector, so after $v$ vectors each net has a $v$-bit serial code. To detect opens and shorts, every net needs a unique code, and none may be all 0s or all 1s, since a stuck net or a wired short could mimic those. The minimum is
$$v = \lceil \log_2(N + 2) \rceil$$
vectors for $N$ nets, the modified counting sequence: 1165 nets need 11 vectors. Wired-AND and wired-OR shorts then produce the AND or OR of the two codes, which differs from both and identifies the short.

## Test time

Every vector is a full data-register scan, and the scan overhead comes from the TAP state machine: an $L$-bit DR scan from Run-Test/Idle back to Run-Test/Idle takes $L + 5$ TCK cycles and an IR scan of $I$ bits takes $I + 6$. An EXTEST session that loads a 32-bit instruction chain once and applies 16 vectors through a 146-bit boundary chain takes $32 + 6 + 16 \times 151 = 2454$ cycles, 490.8 microseconds at 5 MHz. TCK is slow, often 10 to 50 MHz on boards, so boundary scan is for connectivity and access, not for high-volume internal scan test.

## Going deeper

IEEE Std 1149.1, "Test Access Port and Boundary-Scan Architecture"; Parker, "The Boundary-Scan Handbook" (4th ed., 2016); Bushnell and Agrawal (2000), chapter 16.
