# DREAMPlace and GPU-accelerated placement

DREAMPlace recast electrostatics-based global placement as training a neural network: cell positions are the trainable parameters, wirelength plus density is the loss, and a deep-learning toolkit (PyTorch) provides the optimisation loop and GPU execution. With custom GPU operators it ran global placement more than an order of magnitude faster than CPU implementations with similar quality, and it became a platform for placement research.

## Placement as training

In DREAMPlace, the x and y coordinates of movable cells form a parameter tensor. The forward pass computes weighted-average wirelength and electrostatic density energy; the backward pass computes their gradients; an optimiser such as Nesterov's method updates the parameters. The algorithm is ePlace's, with the same smooth wirelength, spectral density solve, penalty schedule and fillers, so quality is similar to RePlAce; what changes is the execution. Expressing placement in a deep-learning framework brings automatic differentiation for new cost terms, GPU memory management, and easy experimentation, at the price of depending on the framework's runtime.

## Custom operators and GPU kernels

A naive implementation using generic tensor operations would be slow and memory-hungry, because wirelength is a sum over irregular nets. DREAMPlace writes key operators by hand, with forward and backward CUDA kernels: wirelength over nets and pins, density accumulation onto bins, and the Poisson solve through fast cosine transforms. Net-level and bin-level parallelism map well to thousands of GPU threads. Scatter operations, adding cell contributions into bins or pin gradients into cells, use atomic additions, and since floating-point addition is not associative, their nondeterministic order can make two runs differ slightly unless deterministic reductions are used. Later versions also moved legalization and detailed placement onto the GPU, since those became the bottleneck once global placement was fast.

## What it changed for tool designers

DREAMPlace showed that a placer can be built from differentiable operators, so a new objective (a routability term, a timing term, a new density model) is added by writing one operator with a gradient, without rewriting the optimiser. Speed changes how placement is used: when global placement takes a minute instead of half an hour, an agent or a design-space exploration loop can try many parameter settings and floorplans. End-to-end speedup is limited by the stages not accelerated, Amdahl's law in practice. Its results also illustrate that implementation details, such as precision, kernel approximations and stopping rules, shift QoR slightly, so comparisons between placers must use the same metrics, legalizer and detailed placer.

## Going deeper

Lin, Dhar, Li, Ren, Khailany and Pan, "DREAMPlace: deep learning toolkit-enabled GPU acceleration for modern VLSI placement" (DAC 2019; IEEE TCAD 2021); the DREAMPlace open-source repository documentation.
