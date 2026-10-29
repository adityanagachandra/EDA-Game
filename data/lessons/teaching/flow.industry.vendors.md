# EDA vendors and industry structure

The EDA industry is small in revenue next to the semiconductor industry it serves, but it is indispensable: no advanced chip can be built without its tools. A few large vendors supply complete flows, a ring of specialists supplies point tools, and an open-source ecosystem has grown around academic research and open PDKs. Knowing who builds what helps a tool designer place their own work.

## The large digital vendors

Synopsys and Cadence Design Systems are the two largest EDA companies, and each sells a complete flow: simulation and verification, synthesis, place and route, signoff timing and extraction, physical verification and analog tools, plus a large IP business. Siemens EDA, formerly Mentor Graphics until Siemens acquired it in 2017, is the third, strongest in physical verification (Calibre), test (Tessent), verification (Questa) and PCB design. The industry consolidated through decades of acquisitions, and a startup with a strong point tool is often bought by one of these vendors. Each large vendor's tools share a database and scripting environment, which is a major part of what customers buy.

## Specialists by domain

Several companies lead in adjacent domains. Ansys is known for multiphysics simulation, including chip power integrity, electromigration and thermal analysis; Synopsys agreed in 2024 to acquire it. Keysight's EEsof tools focus on RF and microwave design and electromagnetic simulation. Silvaco sells TCAD, analog and custom design tools. Zuken and Altium focus on PCB design. FPGA vendors, AMD (Xilinx) and Altera, write the implementation tools for their own devices, because those tools depend on undisclosed details of the fabric. Specialists often survive by owning a physics domain, such as electromagnetics, thermal or process simulation, where general digital vendors have less depth.

## How the business works

EDA tools are licensed rather than sold, usually as time-based licences counted by simultaneous users or jobs, so a design team's compute farm and licence pool limit how many runs it can do in parallel. Vendors work closely with foundries to certify tools on each new node, and with large customers to tune tools on their designs. Switching costs are high: scripts, methodology and signoff correlation are built around one vendor's flow. That lock-in, the need for foundry certification, and the size of the code bases are the barriers that new entrants, including open-source projects, face.

## Going deeper

The ESD Alliance (a SEMI technology community) market statistics reports for the industry's size and segments.
