# Lookup-table design with gm/Id

Lookup-table design replaces device equations by tables generated once from SPICE sweeps of the foundry model. Indexed by $g_m/I_D$ and channel length, they give current density, intrinsic gain, transit frequency and capacitances. Sizing becomes table lookups and interpolation, accurate to the model, fast enough for scripts and optimizers, and portable to a new process by regenerating the tables.

## Current density sets the width

For a fixed channel length and bias conditions, the drain current scales linearly with width, so the useful quantity is the current density $I_D/W$. A table sweeps $V_{GS}$ on a unit-width device and records $g_m/I_D$, $I_D/W$, $g_m/g_{ds}$, $f_T$ and $C_{gg}/W$ at each point. Reindexing by $g_m/I_D$ gives rows such as: at $g_m/I_D = 14$, $I_D/W = 9.5\ \mu$A/$\mu$m and $g_m/g_{ds} = 32$. Once the required current is known, the width is a division: a device carrying 76 microamperes at that point needs $W = 76/9.5 = 8\ \mu$m. Because density depends on $L$ and $V_{DS}$, real tables are multi-dimensional, but the procedure is the same.

## Interpolating between rows

A design point rarely lands on a table row. Linear interpolation between neighbouring rows is the simplest estimate: for $g_m/I_D = 15$ between rows at 14 and 18 with densities 9.5 and 4.1, the fraction is $(15 - 14)/(18 - 14) = 0.25$ and the estimate is $9.5 + 0.25(4.1 - 9.5) = 8.15\ \mu$A/$\mu$m. Current density falls roughly exponentially with $g_m/I_D$ toward weak inversion, so interpolating $\log(I_D/W)$ is more accurate, and production scripts use splines over fine sweeps. The table resolution is a design decision: a coarse table is fast to generate but adds error that an optimizer will treat as real. This is ordinary numerical interpolation, with its error term, applied to device data.

## A complete sizing step

A lookup-table sizing step chains the pieces. A spec gives $g_m$; the designer chooses $g_m/I_D$ and $L$ from the trade-offs (gain from $g_m/g_{ds}$, speed from $f_T$); the current is $g_m/(g_m/I_D)$; the width is the current over the density. For $g_m = 600\ \mu$S at $g_m/I_D = 12$: $I_D = 50\ \mu$A, and with density 20 the width is 2.5 micrometres. The same row gives the intrinsic gain and capacitance without another simulation, so the stage's gain and poles can be predicted at once. A sizing script that does this for every device lands close to the simulated result on the first pass, which is why gm/Id tables are both a designer's tool and the inner model of several automated sizers and migration flows.

## Going deeper

Jespers and Murmann, "Systematic Design of Analog CMOS Circuits Using Pre-Computed Lookup Tables" (2017), chapters 2 to 4.
