# Distributed runs, containers and the cloud

A single flow run fits on one machine, but EDA work rarely comes as one run: regressions run dozens of designs, sweeps run hundreds of parameter sets, characterization runs thousands of independent simulations. Teams spread that work across clusters and clouds, which raises three problems a tool builder must solve: making every machine run the same software, sizing and paying for the fleet, and moving large design data to where the computation happens without breaking the IP boundary.

## Containers pin the tool environment

A container image packages a flow's user space: tool binaries such as Yosys, OpenROAD and KLayout, the libraries they link, the Python interpreter and its packages. Every host that runs the image sees the same versions, so a flow that passed in CI behaves the same on a laptop and on a cluster node. Containers share the host's Linux kernel; isolation comes from kernel namespaces, and resource limits from cgroups, which can cap a job's memory and CPU so that one runaway router does not take down the other jobs on the node.

Reproducibility depends on naming the image exactly. A mutable tag such as latest can be moved to a new build; pinning an image by its content digest names one immutable image. The PDK and design data are usually mounted into the container rather than baked in, which keeps images small and keeps proprietary data out of shared registries.

## Sizing a parameter sweep

For $J$ independent jobs of $t$ hours each on $M$ machines that run one job at a time, the jobs run in $\lceil J/M \rceil$ waves, so the wall-clock time is $\lceil J/M \rceil\, t$. Cost depends on billing. If machines are reserved for the whole sweep, you pay for $M \lceil J/M \rceil t$ machine-hours, idle ones in the last wave included. If an autoscaling pool bills only while a job runs, you pay for $J t$ machine-hours whatever $M$ is: the pool size changes the wall time, not the cost.

Example: 50 runs of 2 h on 8 machines take $\lceil 50/8 \rceil = 7$ waves, 14 h. Reserved at 1.5 credits per machine-hour, the cost is $8 \times 14 \times 1.5 = 168$ credits; billed per job, it is $50 \times 2 \times 1.5 = 150$ credits. Adding machines to 50 would cut the wall time to 2 h at the same per-job cost.

## Scheduling, storage and data locality

Independent jobs scale out almost perfectly; one tightly coupled run that shares a large in-memory database does not, because splitting it across machines turns memory accesses into network messages. So EDA distribution is mostly about many independent jobs, orchestrated as a graph of steps by a scheduler or workflow engine. Two engineering points dominate. Storage: hundreds of jobs reading a 30 GB database each at the same moment saturate a shared file system, so inputs are staged or cached near the compute. Restartability: if each step writes its outputs to durable storage keyed by its inputs, a failed sweep resumes from where it stopped, and cheap preemptible machines become usable for short jobs.

Cloud use also raises an architectural question that is not about price: foundry agreements and IP policy constrain where design and PDK data may live and who may access it, so tenancy, encryption and data residency are part of the design.

## Going deeper

The Linux kernel documentation on cgroups and namespaces; the OpenROAD-flow-scripts documentation on running the flow in Docker; Dean and Barroso, "The Tail at Scale" (CACM 2013), on large-fleet job behaviour.
