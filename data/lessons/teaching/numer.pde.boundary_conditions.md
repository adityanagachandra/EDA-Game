# Boundary conditions in thermal grids

A thermal model is only as good as its boundaries: a heat sink, an insulated package edge, a convective surface. Each kind of boundary condition becomes a specific row in the finite-difference system. Getting them right decides whether the absolute temperatures are meaningful and whether the matrix is even nonsingular.

## Dirichlet: fixed temperatures

A Dirichlet condition fixes $T$ at a boundary node, for example an idealized heat-sink contact at ambient. The node is not an unknown; its known value moves to the right-hand side of its neighbours' equations.

For a rod of length $L$ with $-T'' = Q$ ($Q = q/k$ uniform) and ends held at $T_a$ and $T_b$, the exact solution is
$$T(x) = T_a + (T_b - T_a)\frac{x}{L} + \frac{Q}{2}x(L - x).$$
It is quadratic, so the three-point stencil reproduces it exactly at every node. With $L = 4$ mm, $T_a = 30$, $T_b = 40$ °C, $Q = 2$ K/mm$^2$ and 4 intervals, the node at $x = 1$ mm has $T = 30 + 2.5 + 3 = 35.5$ °C. The heating adds a parabola peaking at the middle, $QL^2/8 = 4$ K above the linear profile.

## Neumann: insulated edges and ghost nodes

An insulated boundary has zero heat flux, $T' = 0$. At the end node $N$, introduce a ghost node $T_{N+1}$ and impose symmetry $T_{N+1} = T_{N-1}$, which makes the centred derivative zero. The stencil at $N$ then reads
$$\frac{2T_{N-1} - 2T_N}{h^2} + Q = 0 \quad\Rightarrow\quad T_N = T_{N-1} + \frac{Qh^2}{2}.$$
With $T_{N-1} = 70$ °C, $Q = 2$ K/mm$^2$ and $h = 0.5$ mm: $T_N = 70.25$ °C. The insulated end is the hottest point because heat cannot leave there.

A pure Neumann problem, insulated on every side, is singular: temperature is determined only up to a constant, and with net heat input there is no steady state at all. Every thermal model needs at least one Dirichlet or Robin boundary to anchor it, just as a circuit needs ground.

## Robin: convection to a heat sink

A convective surface loses heat at rate $h_c(T - T_{amb})$ per unit area, with heat-transfer coefficient $h_c$ in W/(m$^2$ K). Balancing it against conduction from the first interior node, one-sided, gives
$$k\frac{T_1 - T_0}{\Delta x} = h_c(T_0 - T_{amb}) \quad\Rightarrow\quad T_0 = \frac{(k/\Delta x)T_1 + h_cT_{amb}}{k/\Delta x + h_c}.$$
This is a series combination of a conduction conductance $k/\Delta x$ and a convection conductance $h_c$. With $k = 1$ W/(m K), $\Delta x = 1$ mm, $h_c = 1000$ W/(m$^2$ K), $T_1 = 80$ °C and $T_{amb} = 25$ °C: both conductances equal 1000, so $T_0 = 52.5$ °C, halfway.

When $h_c \gg k/\Delta x$ the condition approaches Dirichlet; when $h_c \to 0$ it becomes insulating. Package thermal models often use Robin conditions with an effective $h_c$ to represent the heat sink and board without meshing them.

## Going deeper

Incropera, DeWitt, Bergman and Lavine, "Fundamentals of Heat and Mass Transfer" (Wiley), chapters on finite-difference conduction; LeVeque, "Finite Difference Methods for Ordinary and Partial Differential Equations" (SIAM, 2007), chapter 2.
