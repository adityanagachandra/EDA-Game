# Cholesky for conductance matrices

Resistive networks, RC interconnect and power grids give symmetric positive definite conductance matrices. For them, Cholesky factorization $G = LL^T$ halves the work and storage of LU and needs no pivoting for stability. Power-grid IR-drop solvers and thermal solvers use sparse Cholesky directly or as the model for incomplete Cholesky preconditioners.

## Why grounded conductance matrices are SPD

For a resistor network, the quadratic form of the nodal matrix is the dissipated power:
$$v^T G v = \sum_{\text{edges } (i,j)} g_{ij}(v_i - v_j)^2 + \sum_i g_{i0}\, v_i^2 \ \ge 0.$$
Example: nodes 1 and 2 joined by 4 mS, 2 mS from node 1 to ground, 3 mS from node 2 to ground, $v = (0.5, -0.2)$ V. Then $v^TGv = 4(0.7)^2 + 2(0.25) + 3(0.04) = 1.96 + 0.5 + 0.12 = 2.58$ mW.

The form is zero only if every edge term and ground term vanishes, which forces $v = 0$ when every connected piece reaches ground through some conductance. So $G$ is positive definite. Remove all paths to ground from one piece and that piece can float at any voltage: $G$ becomes singular, which is the floating-node failure that gmin prevents.

## Computing the Cholesky factor

Cholesky computes $L$ column by column: $l_{jj} = \sqrt{a_{jj} - \sum_{k<j} l_{jk}^2}$ and $l_{ij} = (a_{ij} - \sum_{k<j} l_{ik}l_{jk})/l_{jj}$ for $i > j$. For
$$G = \begin{pmatrix} 4 & -2 \\ -2 & 5 \end{pmatrix}: \quad l_{11} = 2, \quad l_{21} = -2/2 = -1, \quad l_{22} = \sqrt{5 - 1} = 2.$$
Check: $LL^T = \begin{pmatrix} 4 & -2 \\ -2 & 5 \end{pmatrix}$.

Positive definiteness guarantees that every square root has a positive argument; a non-positive one is a proof that the matrix is not SPD, which makes Cholesky a cheap test. The fill pattern of $L$ is that of symmetric elimination, so the same orderings (AMD, nested dissection) apply.

## Half the work and no pivoting

Dense Cholesky costs about $n^3/3$ flops against $2n^3/3$ for LU, and stores only $L$. For a dense block of $n = 3000$ unknowns, that is $9 \times 10^9$ flops, 90 ms at 100 GFLOP/s, versus 180 ms for LU.

More important for sparse codes, the pivot order can be chosen purely to reduce fill, because diagonal pivots are always stable for SPD matrices: elimination never grows entries beyond the original diagonal. The symbolic analysis is then exact and reusable. CHOLMOD and similar libraries exploit this with supernodal kernels.

The limit is structure. Voltage sources in MNA, controlled sources and inductors with branch currents make the matrix unsymmetric or indefinite, so general circuit simulation needs LU. Power grids modelled as resistors with current loads, and thermal networks, stay SPD and suit Cholesky or conjugate gradients.

## Going deeper

Davis, "Direct Methods for Sparse Linear Systems" (SIAM, 2006), chapter 4; Chen, Davis, Hager and Rajamanickam, "Algorithm 887: CHOLMOD, supernodal sparse Cholesky factorization and update/downdate", ACM TOMS (2008).
