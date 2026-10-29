package infra

import (
	"math"
	"sort"
	"strings"

	"github.com/chuma-beep/mathua/internal/generator"
	"github.com/chuma-beep/mathua/internal/generator/eda/edakit"
)

func init() { registrars = append(registrars, registerPar) }

func registerPar(reg *generator.Registry) {
	reg.Register("infra.par.amdahl", edakit.Func(amdahl))
	reg.Register("infra.par.thread_pool_scheduling", edakit.Func(threadPoolScheduling))
	reg.Register("infra.par.deterministic_parallel", &edakit.Bank{Items: deterministicParallelItems})
	reg.Register("infra.par.conflict_free_batches", edakit.Func(conflictFreeBatches))
	reg.Register("infra.par.containers_cloud", &edakit.Mix{Parts: []generator.Generator{&edakit.Bank{Items: containersCloudItems}, edakit.Func(containersCloud)}})
	reg.Register("infra.par.license_management", edakit.Func(licenseManagement))
}

var parStepNames = []string{"global placement", "detailed routing", "parasitic extraction", "timing analysis", "DRC", "global routing"}

// amdahl asks for speedup, the speedup ceiling and the serial fraction a target implies.
func amdahl(ctx generator.GeneratorContext) generator.Problem {
	step := edakit.Pick(parStepNames)
	switch variant(ctx) {
	case 1:
		s := float64(edakit.Int(2, 30)) / 100
		p := edakit.Pick([]int{4, 8, 16, 32, 64})
		sp := 1 / (s + (1-s)/float64(p))
		ans := edakit.Num(sp, 2)
		q := edakit.Sprintf("Profiling shows that %s of a %s run is serial (reading the database, building data structures, writing results) and the rest parallelizes perfectly. What speedup does the run reach on %d threads? Round to 2 decimal places.", edakit.Num(s*100, 0)+"%", step, p)
		return edakit.Numeric(q, ans, edakit.Sprintf("Amdahl's law: S = 1 / (s + (1 - s)/p) = 1 / (%s + %s/%d) = %s.", edakit.Num(s, 2), edakit.Num(1-s, 2), p, ans),
			map[string]string{"speedup": ans})
	case 2:
		serial := float64(edakit.Int(2, 30))
		par := float64(10 * edakit.Int(5, 40))
		s := serial / (serial + par)
		ans := edakit.Num(1/s, 2)
		q := edakit.Sprintf("A single-threaded %s run takes %s s, of which %s s is inherently serial and the rest parallelizes perfectly. What is the largest speedup any number of threads can deliver? Round to 2 decimal places.", step, edakit.Num(serial+par, 0), edakit.Num(serial, 0))
		return edakit.Numeric(q, ans, edakit.Sprintf("As p grows, the parallel term vanishes and S approaches 1/s. Here s = %s/%s = %s, so the ceiling is %s.", edakit.Num(serial, 0), edakit.Num(serial+par, 0), edakit.Num(s, 4), ans),
			map[string]string{"serial_fraction": edakit.Num(s, 4), "ceiling": ans})
	}
	p := edakit.Pick([]int{8, 16, 32, 64})
	target := float64(p) * float64(edakit.Int(30, 80)) / 100
	target = math.Round(target*2) / 2
	if target < 2 {
		target = 2
	}
	s := (float64(p)/target - 1) / float64(p-1)
	ans := edakit.Num(s*100, 2)
	q := edakit.Sprintf("A team wants a %s engine to run %s times faster on %d threads than on one. Under Amdahl's law, at most what percentage of the single-threaded run may be serial? Round to 2 decimal places.", step, edakit.Num(target, 1), p)
	return edakit.Numeric(q, ans, edakit.Sprintf("Solve S = 1/(s + (1-s)/p) for s: s = (p/S - 1)/(p - 1) = (%d/%s - 1)/%d = %s, which is %s%%.", p, edakit.Num(target, 1), p-1, edakit.Num(s, 5), ans),
		map[string]string{"serial_percent": ans})
}

// parListSchedule assigns tasks in the given order to the worker that becomes free first
// (lowest index on ties) and returns the makespan and per-worker loads.
func parListSchedule(tasks []int, p int) (int, []int) {
	load := make([]int, p)
	for _, t := range tasks {
		best := 0
		for w := 1; w < p; w++ {
			if load[w] < load[best] {
				best = w
			}
		}
		load[best] += t
	}
	m := 0
	for _, l := range load {
		if l > m {
			m = l
		}
	}
	return m, load
}

func parInts(xs []int) string {
	var s []string
	for _, x := range xs {
		s = append(s, edakit.IntStr(x))
	}
	return strings.Join(s, ", ")
}

// threadPoolScheduling asks for greedy and longest-first makespans and the lower bound.
func threadPoolScheduling(ctx generator.GeneratorContext) generator.Problem {
	p := edakit.Int(2, 4)
	n := edakit.Scaled(ctx, p+3, p+7)
	tasks := make([]int, n)
	sum, mx := 0, 0
	for i := range tasks {
		tasks[i] = edakit.Int(1, 12)
		sum += tasks[i]
		if tasks[i] > mx {
			mx = tasks[i]
		}
	}
	q := edakit.Sprintf("A router splits a design into %d independent regions whose routing times are %s ms, and runs them on a pool of %d worker threads. ", n, parInts(tasks), p)
	switch variant(ctx) {
	case 1:
		m, load := parListSchedule(tasks, p)
		return edakit.Numeric(q+"Tasks are taken in the listed order, and each goes to the worker that becomes free first (the lowest-numbered one on a tie). What is the makespan, in ms?", edakit.IntStr(m),
			edakit.Sprintf("Greedy list scheduling gives worker loads %s ms; the makespan is the largest, %d ms.", parInts(load), m), map[string]string{"makespan": edakit.IntStr(m)})
	case 2:
		sorted := append([]int(nil), tasks...)
		sort.Sort(sort.Reverse(sort.IntSlice(sorted)))
		m, load := parListSchedule(sorted, p)
		g, _ := parListSchedule(tasks, p)
		return edakit.Numeric(q+"The pool first sorts tasks from longest to shortest, then gives each to the worker that becomes free first (the lowest-numbered one on a tie). What is the makespan, in ms?", edakit.IntStr(m),
			edakit.Sprintf("Sorted order is %s. Longest-processing-time-first gives worker loads %s ms, so the makespan is %d ms; the unsorted greedy order gives %d ms.", parInts(sorted), parInts(load), m, g),
			map[string]string{"makespan": edakit.IntStr(m), "greedy": edakit.IntStr(g)})
	}
	lb := math.Max(float64(sum)/float64(p), float64(mx))
	ans := edakit.Num(lb, 2)
	return edakit.Numeric(q+"What lower bound on the makespan holds for every possible schedule, in ms? Use the larger of the total work divided by the worker count and the longest single task. Round to 2 decimal places.", ans,
		edakit.Sprintf("Total work is %d ms, so %d workers need at least %s ms; the longest task alone needs %d ms. The bound is the larger, %s ms.", sum, p, edakit.Num(float64(sum)/float64(p), 2), mx, ans),
		map[string]string{"sum": edakit.IntStr(sum), "longest": edakit.IntStr(mx), "bound": ans})
}

// conflictFreeBatches colors a tile grid so that concurrently processed tiles never touch.
func conflictFreeBatches(ctx generator.GeneratorContext) generator.Problem {
	r := edakit.Scaled(ctx, 3, 9)
	c := edakit.Scaled(ctx, 3, 10)
	p := edakit.Int(2, 6)
	ceil := func(a, b int) int { return (a + b - 1) / b }
	classes := []int{ceil(r, 2) * ceil(c, 2), ceil(r, 2) * (c / 2), (r / 2) * ceil(c, 2), (r / 2) * (c / 2)}
	rounds := 0
	var per []string
	for _, k := range classes {
		rounds += ceil(k, p)
		per = append(per, edakit.Sprintf("%d tiles → %d rounds", k, ceil(k, p)))
	}
	q := edakit.Sprintf("A detailed router divides the die into a %d × %d grid of tiles (%d rows, %d columns). Two tiles that share an edge or a corner may touch the same wires, so they must never be routed at the same time. The tiles are colored by (row mod 2, column mod 2), giving four batches; batches run one after another, and inside a batch any tiles may run together. ", r, c, r, c)
	facts := map[string]string{"rounds": edakit.IntStr(rounds), "largest_batch": edakit.IntStr(classes[0])}
	switch variant(ctx) {
	case 1:
		return edakit.Numeric(q+"Rows and columns are numbered from 0. How many tiles are in the batch with even row and even column?", edakit.IntStr(classes[0]),
			edakit.Sprintf("Even rows: ceil(%d/2) = %d. Even columns: ceil(%d/2) = %d. The batch holds %d × %d = %d tiles, and no two of them touch, since any two differ by at least two in row or in column.", r, ceil(r, 2), c, ceil(c, 2), ceil(r, 2), ceil(c, 2), classes[0]), facts)
	case 2:
		return edakit.Numeric(q+edakit.Sprintf("With %d worker threads, each routing one tile per round, how many rounds does the whole grid take?", p), edakit.IntStr(rounds),
			edakit.Sprintf("The four batches hold %s tiles. Each batch needs ceil(size/%d) rounds: %s, for %d rounds in all.", parInts(classes), p, strings.Join(per, ", "), rounds), facts)
	}
	t := edakit.Int(2, 9) * 10
	serial := float64(r*c) * float64(t)
	par := float64(rounds) * float64(t)
	ans := edakit.Num(serial/par, 2)
	facts["speedup"] = ans
	return edakit.Numeric(q+edakit.Sprintf("Every tile takes %d ms. With %d worker threads, each routing one tile per round, what is the speedup over routing all tiles on one thread? Round to 2 decimal places.", t, p), ans,
		edakit.Sprintf("Serially: %d tiles × %d ms = %s ms. Batched: %d rounds (%s) × %d ms = %s ms. Speedup = %s, below %d because rounds at the end of each batch are only partly full.", r*c, t, edakit.Num(serial, 0), rounds, strings.Join(per, ", "), t, edakit.Num(par, 0), ans, p), facts)
}

// containersCloud is the computed part of the containers concept: sizing a parameter sweep.
func containersCloud(ctx generator.GeneratorContext) generator.Problem {
	j := edakit.Scaled(ctx, 20, 200)
	m := edakit.Int(4, 32)
	t := float64(edakit.Int(1, 8)) / 2
	rate := float64(edakit.Int(2, 12)) / 4
	waves := (j + m - 1) / m
	wall := float64(waves) * t
	switch edakit.Int(1, 3) {
	case 1:
		ans, d := numOptions(wall, 1, "h", float64(j)*t/float64(m), float64(j)*t, float64(waves+1)*t, float64(waves-1)*t)
		q := edakit.Sprintf("A design-space exploration launches %d independent flow runs, one per parameter set, on %d identical cloud machines that each run one job at a time. Every run takes %s h. What is the wall-clock time of the sweep?", j, m, edakit.Num(t, 1))
		return edakit.Choice(q, ans, edakit.Sprintf("The jobs run in ceil(%d/%d) = %d waves of %s h, so the sweep takes %s.", j, m, waves, edakit.Num(t, 1), ans), d...)
	case 2:
		held := float64(waves*m) * t * rate
		used := float64(j) * t * rate
		ans, d := numOptions(held, 2, "credits", used, wall*rate, held/float64(m)*float64(m+1))
		q := edakit.Sprintf("A sweep of %d flow runs, %s h each, runs on %d cloud machines reserved for the whole sweep and billed at %s credits per machine-hour, idle or not. Each machine runs one job at a time. What does the sweep cost?", j, edakit.Num(t, 1), m, edakit.Num(rate, 2))
		return edakit.Choice(q, ans, edakit.Sprintf("The sweep takes ceil(%d/%d) = %d waves, %s h, and all %d machines are billed for it: %d × %s × %s = %s. Releasing idle machines in the last wave would cost %s.", j, m, waves, edakit.Num(wall, 1), m, m, edakit.Num(wall, 1), edakit.Num(rate, 2), ans, edakit.Num(used, 2)+" credits"), d...)
	}
	used := float64(j) * t * rate
	held := float64(waves*m) * t * rate
	ans, d := numOptions(used, 2, "credits", held, wall*rate, used*2)
	q := edakit.Sprintf("A sweep of %d flow runs, %s h each, runs on an autoscaling pool of up to %d machines; a machine is billed only while it runs a job, at %s credits per machine-hour. What does the sweep cost?", j, edakit.Num(t, 1), m, edakit.Num(rate, 2))
	return edakit.Choice(q, ans, edakit.Sprintf("Billing follows the work, not the wall time: %d jobs × %s h × %s = %s. The pool size changes the %s h wall time, not the cost.", j, edakit.Num(t, 1), edakit.Num(rate, 2), ans, edakit.Num(wall, 1)), d...)
}

// licenseManagement asks how licenses cap concurrency, how well they are used, and what early
// check-in saves.
func licenseManagement(ctx generator.GeneratorContext) generator.Problem {
	j := edakit.Scaled(ctx, 6, 40)
	l := edakit.Int(2, 8)
	if l >= j {
		l = j - 1
	}
	waves := (j + l - 1) / l
	switch variant(ctx) {
	case 1:
		t := edakit.Int(1, 6)
		wall := waves * t
		q := edakit.Sprintf("A team must run %d signoff timing jobs, each holding one tool license for its whole %d h run. Machines are plentiful, but only %d licenses exist. What is the wall-clock time to finish all jobs, in hours?", j, t, l)
		return edakit.Numeric(q, edakit.IntStr(wall), edakit.Sprintf("At most %d jobs run at once, so the jobs run in ceil(%d/%d) = %d waves of %d h: %d h.", l, j, l, waves, t, wall),
			map[string]string{"waves": edakit.IntStr(waves), "wall": edakit.IntStr(wall)})
	case 2:
		t := edakit.Int(1, 6)
		wall := waves * t
		u := float64(j*t) / float64(l*wall) * 100
		ans := edakit.Num(u, 2)
		q := edakit.Sprintf("%d jobs of %d h each share %d licenses, each job holding one license for its whole run, and run in back-to-back waves. Over the wall-clock time of the batch, what percentage of the available license-hours is in use? Round to 2 decimal places.", j, t, l)
		return edakit.Numeric(q, ans, edakit.Sprintf("The batch takes %d waves × %d h = %d h, offering %d × %d = %d license-hours. Jobs use %d × %d = %d, so utilization is %s%%. The last wave is only partly full.", waves, t, wall, l, wall, l*wall, j, t, j*t, ans),
			map[string]string{"utilization": ans})
	}
	b := edakit.Int(1, 5)
	c := edakit.Int(1, 4)
	hold := waves * (b + c)
	early := waves*b + c
	q := edakit.Sprintf("Each of %d jobs has a %d h phase that needs the tool license, followed by a %d h post-processing phase (report parsing and plotting) that does not. There are %d licenses and plenty of machines. How many hours of wall-clock time are saved if each job checks its license back in when the licensed phase ends, instead of holding it until the job finishes?", j, b, c, l)
	return edakit.Numeric(q, edakit.IntStr(hold-early),
		edakit.Sprintf("Holding: %d waves × (%d + %d) h = %d h. Checking in: the licensed phases run back to back in %d waves of %d h, and the last wave's post-processing adds %d h, so %d h. Saved: %d h.", waves, b, c, hold, waves, b, c, early, hold-early),
		map[string]string{"hold": edakit.IntStr(hold), "early": edakit.IntStr(early)})
}

var deterministicParallelItems = []edakit.Item{
	{KP: 1, Q: "A placer sums the wirelength gradient over all nets with a parallel reduction. Run twice with the same input and 16 threads, it reports total wirelengths that differ in the 12th digit. What is the cause?",
		A: "Floating-point addition is not associative, and the reduction order differed", D: []string{"A data race corrupted memory in one run", "The random seed is reset differently on each run", "Double precision is too short to hold a wirelength"},
		Why: "With dynamic scheduling, partial sums are combined in whatever order threads finish. Rounding depends on order, so the last digits move even with no race."},
	{KP: 1, Q: "In decimal arithmetic that rounds every result to 4 significant digits, a tool computes (1000 + 0.4) + 0.4 and then 1000 + (0.4 + 0.4). What are the two results?",
		A: "1000 for the first and 1001 for the second", D: []string{"1000.8 for both", "1000 for both", "1001 for the first and 1000 for the second"},
		Why: "1000 + 0.4 rounds back to 1000, twice. But 0.4 + 0.4 = 0.8 survives, and 1000.8 rounds to 1001. Only the order of the additions changed."},
	{KP: 1, Q: "Why do tiny reduction-order differences matter in a placer more than in a single report?",
		A: "The engine iterates, so small differences change later decisions and the final placement diverges", D: []string{"They do not matter; the final placement is always identical", "Placers use integer arithmetic only, so differences cannot occur", "Only the final HPWL report is affected, never the cell positions"},
		Why: "An optimizer's decisions (a move accepted, a cell ordered first) depend on comparisons, and a last-digit change can flip one; the divergence then compounds across iterations."},
	{KP: 1, Q: "Which change makes a parallel float sum bitwise reproducible across runs with the same thread count?",
		A: "Statically assign index ranges to threads and combine partial sums in a fixed order", D: []string{"Use dynamic scheduling so every thread stays busy", "Use an atomic float add into one shared accumulator", "Increase the number of threads until the result stabilizes"},
		Why: "A fixed partition and a fixed combine order fix every rounding step. Atomic adds and dynamic scheduling make the order depend on timing."},
	{KP: 2, Q: "A router processes nets in parallel and each thread appends its finished nets to a shared output list. The list order varies from run to run. What is the standard fix?",
		A: "Give each work item an index and merge results in index order", D: []string{"Protect the list with a mutex", "Sort threads by their operating-system ID", "Run the router with a single thread in production"},
		Why: "A mutex removes the race but not the order dependence. Writing results to slots keyed by item index, then merging in that order, gives the same output whatever the timing."},
	{KP: 2, Q: "A team wants results to be identical whether a job runs on 4 threads or 32. Which design achieves that?",
		A: "Partition the work into a fixed number of chunks independent of the thread count and reduce chunks in order", D: []string{"Give each thread an equal share of the work, sized by the thread count", "Seed each thread's random generator with its thread ID", "Let threads steal work and merge whichever partial result arrives first"},
		Why: "If the chunking depends on the thread count, the rounding and the random streams change with it. Fixed chunks with fixed seeds and ordered reduction make the thread count a pure performance knob."},
	{KP: 2, Q: "A parallel simulated-annealing step needs random numbers in every worker. Which seeding scheme keeps runs reproducible?",
		A: "Derive each work item's stream from the run seed and the item index", D: []string{"Share one global generator protected by a lock", "Seed each worker from the wall-clock time", "Seed each worker with its thread ID"},
		Why: "A shared locked generator hands out numbers in timing-dependent order. A stream keyed by run seed and item index gives every item the same numbers no matter which thread runs it."},
	{KP: 2, Q: "A partitioned parallel legalizer gives each thread a region and then must fix cells near region borders. How can the border fix stay deterministic?",
		A: "Process border cells in a second, ordered pass after all regions finish", D: []string{"Let whichever thread reaches a border cell first move it", "Skip border cells and leave them overlapping", "Lock each border cell and process in lock-acquisition order"},
		Why: "First-come handling ties the result to thread timing. A barrier followed by an ordered pass over the shared region removes the dependence."},
	{KP: 3, Q: "Two candidate moves in a detailed placer have exactly equal gain. The code picks whichever appears first in a hash map iteration. What is wrong?",
		A: "Hash-map iteration order can vary, so the tie-break is not deterministic", D: []string{"Nothing, because equal gains give identical placements", "Hash maps are slower than arrays for two elements", "Equal gains mean both moves must be applied"},
		Why: "Iteration order may depend on insertion history, pointer values or a randomized hash seed. Break ties by a stable key such as instance index."},
	{KP: 3, Q: "A tool sorts instances by pointer address to get a stable order for a heuristic. Results differ between runs. Why?",
		A: "Allocation addresses differ from run to run, so the order changes", D: []string{"Pointer comparison is undefined for any two objects", "Sorting is not stable for equal keys", "Pointers are always sorted in allocation order"},
		Why: "Address layout depends on the allocator, address-space randomization and earlier allocations. Use a persistent ID or name as the key."},
	{KP: 3, Q: "A regression flags a 0.1% wirelength change on a commit that only touched comments. The engine is supposedly deterministic. What should the team suspect first?",
		A: "A hidden source of nondeterminism such as unordered iteration, timing-dependent merging or uninitialized memory", D: []string{"The compiler optimized the comments", "Wirelength always varies by 0.1% between runs", "The test design was changed by the commit"},
		Why: "Deterministic engines should produce bit-identical results from identical inputs. A change without a cause is a determinism bug, and it would hide real regressions."},
	{KP: 3, Q: "Why is 'same result regardless of thread count' worth paying a few percent of speed for in a production EDA engine?",
		A: "Bugs found on a big machine can be reproduced on a laptop, and QoR comparisons are not confounded by run noise", D: []string{"It makes the engine run faster on many threads", "It removes the need for regression tests", "It guarantees the result is optimal"},
		Why: "Reproducibility makes debugging and differential testing possible. Without it, every QoR delta carries noise from scheduling."},
}

var containersCloudItems = []edakit.Item{
	{KP: 1, Q: "A flow that ran cleanly last month now fails at synthesis on a new machine. The team ships the flow as a container image built from a fixed recipe. What does the image guarantee?",
		A: "The same tool binaries, libraries and versions on every host that runs it", D: []string{"The same CPU microarchitecture on every host", "That the design's QoR improves", "That no license server is needed"},
		Why: "A container packages the user-space environment, so tool and library versions are pinned. The kernel and hardware still come from the host."},
	{KP: 1, Q: "Which part of the execution environment does a Linux container NOT bring with it?",
		A: "The kernel", D: []string{"The C++ runtime library", "The Python interpreter version", "The tool binaries"},
		Why: "Containers share the host kernel and isolate processes with namespaces and cgroups. Everything in user space, including runtimes and tools, comes from the image."},
	{KP: 1, Q: "A container image tag such as latest is used in a regression flow. Why is that risky?",
		A: "The tag can point to a different image later, so the flow is not reproducible", D: []string{"Tags slow down image pulls", "Containers ignore tags at run time", "A tag forces the container to run as root"},
		Why: "A mutable tag can be moved to a new build. Pinning an image by its content digest names exactly one image."},
	{KP: 1, Q: "A memory-hungry routing job is run inside a container with a cgroup memory limit of 64 GB and exceeds it. What happens?",
		A: "The kernel kills the process in that cgroup when it cannot reclaim memory", D: []string{"The job swaps to the network file system", "The container automatically grows its limit", "Other containers on the host are killed first"},
		Why: "cgroups enforce resource limits per group. Exceeding the memory limit triggers the out-of-memory killer for that group, which protects the other jobs on the host."},
	{KP: 1, Q: "Why do EDA teams often build containers for the open-source flow rather than installing tools on each machine?",
		A: "Many tools and exact library versions must match, and a container builds them once for every host", D: []string{"Open-source tools only run inside containers", "Containers make the tools run faster", "Containers remove the need for a PDK"},
		Why: "A flow chains Yosys, OpenROAD, KLayout and more, each with dependencies. One image keeps all of them consistent across laptops, CI and clusters."},
	{KP: 3, Q: "A cluster scheduler places 500 flow jobs. Each reads a 30 GB design database from shared storage. What is the likely bottleneck?",
		A: "Shared storage bandwidth while all jobs load their inputs", D: []string{"CPU cycles spent parsing the scheduler's job file", "The number of Linux users on the cluster", "The container image size on each host"},
		Why: "Hundreds of jobs reading large inputs at once saturate the file system. Staging data locally or caching inputs near the compute relieves it."},
	{KP: 3, Q: "Which job type benefits least from being split across many machines?",
		A: "One tightly coupled detailed-routing run that shares one in-memory database", D: []string{"A sweep of 200 independent flow runs with different parameters", "Running each regression design as its own job", "Characterizing each standard cell of a library separately"},
		Why: "Independent jobs scale out trivially. A single run that needs fast access to one shared database pays heavy communication costs when distributed."},
	{KP: 3, Q: "A company wants to run proprietary designs on a public cloud. Which concern is architectural rather than a cost question?",
		A: "Where design and PDK data may reside and who can access it", D: []string{"The hourly price of a large instance", "The number of CPU cores per instance", "How long instances take to boot"},
		Why: "Foundry agreements and IP policy constrain where data can live; the architecture must enforce that boundary, for example with isolated tenancy and encryption."},
	{KP: 3, Q: "A workflow orchestrator runs a flow as a graph of steps, each in its own container. A step fails halfway through a 100-step sweep. What property lets the sweep resume without redoing finished steps?",
		A: "Each step writes its outputs to durable storage keyed by its inputs", D: []string{"Every step runs in the same container", "The orchestrator restarts the sweep from the beginning", "Steps share global variables in memory"},
		Why: "If finished outputs are recorded against their inputs, the orchestrator can skip completed steps and rerun only what failed or depends on it."},
	{KP: 3, Q: "Spot or preemptible cloud machines cost much less but can be reclaimed at short notice. Which EDA workload suits them best?",
		A: "Many short independent jobs that can simply be rerun if interrupted", D: []string{"A 3-day single-machine signoff run without checkpoints", "The license server for the whole team", "An interactive debugging session"},
		Why: "Preemption kills jobs. Short, independent, restartable jobs lose little when it happens; long jobs need checkpointing to use such capacity."},
}
