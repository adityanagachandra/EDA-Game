# The TAP controller state machine

Every IEEE 1149.1 device contains the same 16-state finite-state machine, clocked by TCK and steered only by the value of TMS on each rising edge. Because all devices on a board share TMS and TCK, they move through the states in lockstep, and the host drives the whole chain by choosing a TMS sequence. Knowing the state graph is how one writes, checks or debugs any JTAG driver.

## The state graph

From **Test-Logic-Reset**, TMS = 0 enters **Run-Test/Idle**, where the machine waits. TMS = 1 moves to **Select-DR-Scan**; another 1 moves to **Select-IR-Scan**, and a third 1 returns to Test-Logic-Reset. From either Select state, TMS = 0 enters the corresponding column: **Capture**, then **Shift** (TMS = 0 stays there, shifting one bit per edge), **Exit1**, an optional **Pause** loop through **Exit2**, and **Update**, which latches the shifted value. From Update, TMS = 0 returns to Run-Test/Idle and TMS = 1 goes to Select-DR-Scan. Holding TMS at 1 for five edges reaches Test-Logic-Reset from any state, a guaranteed synchronization.

## Following a TMS sequence

To find where a sequence leads, apply one transition per TCK edge. From Shift-DR, TMS = 0, 0 stays in Shift-DR; 1 goes to Exit1-DR; 1 to Update-DR; 1 to Select-DR-Scan; and 0 to Capture-DR. The fewest edges between two states is a shortest path in this 16-node graph: Select-DR-Scan to Shift-DR takes 2 (via Capture-DR), Run-Test/Idle to Shift-IR takes 4. Drivers precompute these paths, and a breadth-first search over the graph finds them.

## Scan cycle counts

A full scan of $k$ bits returning to Run-Test/Idle costs a fixed overhead plus $k$. For the DR: 3 edges reach Shift-DR (Select-DR-Scan, Capture-DR, Shift-DR), $k$ edges shift, the last with TMS = 1 so that it exits to Exit1-DR while shifting the final bit, then 2 edges go through Update-DR back to Run-Test/Idle: $k + 5$. The IR path needs one more edge through Select-IR-Scan: $k + 6$. A 128-bit IR scan takes 134 cycles. Forgetting that the exiting edge also shifts is the classic off-by-one in JTAG drivers.

## Going deeper

IEEE Std 1149.1; Parker, "The Boundary-Scan Handbook" (2016), chapter 1.
