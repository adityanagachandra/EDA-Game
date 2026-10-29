# Congestion prediction

Routing congestion is known only after global routing, but the placer needs it while cells can still move. Congestion prediction maps placement-stage features to the router's overflow map, usually as an image-to-image problem on the gcell grid. The fast estimate RUDY is both a baseline and a standard input feature, and the evaluation metric decides whether the model helps where it matters.

## RUDY as a feature

RUDY (rectangular uniform wire density) spreads each net's estimated wirelength, its half-perimeter $w + h$, uniformly over its bounding box. The density it adds at every point of the box is $(w + h)/(wh)$, in um of wire per square um. A 10 by 20 um net adds $30/200 = 0.15$ per um; a 2 by 2 um net adds $4/4 = 1$, so tight pin clusters show up as hot. RUDY is cheap, needs no router, and correlates with routing demand, which is why learned congestion models take it as an input channel together with pin density, cell density and macro blockage maps.

## From nets to a tile map

A gcell's RUDY value sums each net's density times the area of overlap between the net's box and the gcell, divided by the gcell area. For a 10 by 10 um tile fully inside the 10 by 20 um box above, the net contributes $0.15 \times 100 = 15$ um of wire, a density of 0.15. A second net of density 0.4 overlapping 25 square um adds 10 um, raising the tile to $25/100 = 0.25$. Rasterizing every feature this way turns the placement into a stack of 2-D maps, and the router's overflow per gcell into a target map. A convolutional encoder-decoder then learns the image-to-image mapping, which captures spatial context that RUDY alone ignores, such as detours around macros.

## Judging a congestion map

Map-wide metrics can hide the failures that matter. NRMSE divides the root-mean-square error by the range of the true values. True overflows 0, 2, 4 and 8 against predictions 1, 2, 3 and 6 give squared errors summing to 6, RMSE $\sqrt{6/4} = 1.22$ and NRMSE $1.22/8 = 0.153$. Most tiles are uncongested, so a model that predicts low everywhere scores well while missing the few hot tiles a placer must spread. Useful evaluations add per-tile classification metrics on hotspots, ranking metrics such as correlation of hot regions, and above all whether the placer's routed result improves when it uses the prediction.

## Going deeper

Spindler and Johannes, "Fast and Accurate Routing Demand Estimation for Efficient Routability-driven Placement" (DATE 2007), for RUDY; Xie et al., "RouteNet: Routability Prediction for Mixed-Size Designs Using Convolutional Neural Network" (ICCAD 2018).
