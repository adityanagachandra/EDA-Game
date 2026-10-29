package numer

import (
	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

// Register adds the numerical-methods-for-EDA generators to reg.
func Register(reg *generator.Registry) {
	fn := func(id string, g func(generator.GeneratorContext) generator.Problem) {
		reg.Register(id, edakit.Func(g))
	}
	// sparse
	fn("numer.sparse.mna_pattern", mnaPattern)
	fn("numer.sparse.storage", storage)
	fn("numer.sparse.csr_matvec", csrMatvec)
	fn("numer.sparse.fill_in", fillIn)
	fn("numer.sparse.markowitz", markowitz)
	fn("numer.sparse.minimum_degree", minimumDegree)
	reg.Register("numer.sparse.nested_dissection", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(ndGrid), edakit.Func(ndDense), &edakit.Bank{Items: nestedDissectionItems}}})
	// direct
	fn("numer.direct.lu", lu)
	fn("numer.direct.small_pivot", smallPivot)
	fn("numer.direct.threshold_pivoting", thresholdPivoting)
	fn("numer.direct.condition", condition)
	fn("numer.direct.cholesky", cholesky)
	// nonlinear
	fn("numer.nonlinear.newton_systems", newtonSystems)
	fn("numer.nonlinear.exponential", exponential)
	fn("numer.nonlinear.damping", damping)
	fn("numer.nonlinear.convergence_test", convergenceTest)
	reg.Register("numer.nonlinear.continuation", &edakit.Bank{Items: continuationItems})
	// ode
	fn("numer.ode.backward_euler", backwardEuler)
	fn("numer.ode.trapezoidal", trapezoidal)
	fn("numer.ode.companion", companion)
	fn("numer.ode.bdf2", bdf2)
	fn("numer.ode.stability_regions", stabilityRegions)
	fn("numer.ode.stiffness", stiffness)
	fn("numer.ode.lte", lte)
	fn("numer.ode.step_control", stepControl)
	reg.Register("numer.ode.dae", &edakit.Bank{Items: daeItems})
	// iter
	fn("numer.iter.jacobi_gs", jacobiGS)
	fn("numer.iter.conjugate_gradient", conjugateGradient)
	fn("numer.iter.cg_convergence", cgConvergence)
	fn("numer.iter.gmres", gmres)
	fn("numer.iter.preconditioning", preconditioning)
	fn("numer.iter.ilu", ilu)
	fn("numer.iter.multigrid", multigrid)
	// pde
	fn("numer.pde.finite_difference", finiteDifference)
	fn("numer.pde.boundary_conditions", boundaryConditions)
	fn("numer.pde.explicit_heat", explicitHeat)
	fn("numer.pde.fem", fem)
	fn("numer.pde.fvm", fvm)
	reg.Register("numer.pde.bem", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(bemPanels), edakit.Func(bemCost), &edakit.Bank{Items: bemItems}}})
	// mor
	fn("numer.mor.moments", moments)
	reg.Register("numer.mor.krylov_projection", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(krylovVector), edakit.Func(krylovMoments), &edakit.Bank{Items: krylovItems}}})
	fn("numer.mor.svd_truncation", svdTruncation)
	// hpc
	fn("numer.hpc.autodiff_forward", autodiffForward)
	fn("numer.hpc.adjoint", adjoint)
	fn("numer.hpc.gpu", gpu)
	fn("numer.hpc.parallel_scaling", parallelScaling)
	fn("numer.hpc.reduction_order", reductionOrder)
	// frontier
	reg.Register("numer.frontier.kernel_core", &edakit.Bank{Items: kernelCoreItems})
	reg.Register("numer.frontier.solver_seam", &edakit.Bank{Items: solverSeamItems})
	reg.Register("numer.frontier.differential_tests", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(backwardErrorChoice), edakit.Func(toleranceChoice), &edakit.Bank{Items: diffTestItems}}})
	reg.Register("numer.frontier.gpu_agent_loop", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(batchChoice), edakit.Func(budgetChoice), &edakit.Bank{Items: gpuLoopItems}}})
	reg.Register("numer.frontier.fast_accurate", &edakit.Mix{Parts: []generator.Generator{
		edakit.Func(toleranceIterChoice), edakit.Func(rankingChoice), &edakit.Bank{Items: fastAccurateItems}}})
}
