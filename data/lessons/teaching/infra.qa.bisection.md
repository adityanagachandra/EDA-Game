# Version control and bisecting regressions

Version control gives every state of the tool a name, a commit, and lets anyone rebuild any of them. That turns a regression into a search problem: if last week's build was good and today's is bad, the culprit is one of the commits in between. Bisection finds it by binary search over history. For EDA tools, where each test means building the tool and running a design for many minutes, the logarithmic cost of bisection is the difference between an hour and a week.

## Halving the suspect range

Suppose $n$ commits separate a known good build from a known bad one, so the first bad commit is one of $n$ suspects, the bad endpoint included. Bisection tests the middle suspect: if it is bad, the culprit is at or before it; if good, after it. Each test halves the range, so identifying the culprit takes at most
$$\lceil \log_2 n \rceil$$
tests. Example: 200 commits between a good and a bad nightly need $\lceil \log_2 200 \rceil = 8$ tests, since $2^7 = 128 < 200 \le 256 = 2^8$.

Git automates the bookkeeping: git bisect checks out the midpoint, the engineer or a script marks it good or bad, and it proposes the next commit. A script that builds the tool, runs the design and compares a QoR threshold can drive the whole search unattended. The test must be deterministic, or a noisy QoR measure will steer bisection the wrong way.

## Time to find the culprit

Each bisection step costs a build and a test run, so the worst-case time is
$$T = \lceil \log_2 n \rceil\,(t_{build} + t_{test}).$$
Example: 200 suspects, a 15 min build and a 40 min design run: $8 \times 55 = 440$ min, just over 7 hours, unattended. Reducing either term pays off 8 times: a compiler cache that cuts the build to 5 min saves 80 min, and a smaller design that still shows the regression saves far more.

Two practical complications. Some commits do not build, or fail for an unrelated reason; bisection can skip them, at the cost of possibly ending with a small range rather than one commit. And a regression may need two commits to interact, which a single bisection cannot isolate. Keeping commits small and individually buildable is what makes history searchable.

## Bisection against a linear scan

A linear scan tests suspects one by one, oldest first. If the first $n - 1$ all pass, the last one is the culprit without testing it, so the worst case is $n - 1$ tests, against $\lceil \log_2 n \rceil$ for bisection.

For 200 suspects, that is 199 tests against 8: 191 extra build-and-run cycles, or several days at 55 min each. The gap grows quickly, since doubling the range adds only one bisection step. The comparison also explains why per-commit CI matters even though bisection is efficient: a regression caught on the commit that introduced it needs no search at all, while a regression found weeks later, across thousands of commits by many authors, still costs a dozen expensive steps and a confusing history. Nightly QoR tracking sits in between: it bounds the search to one day's commits, typically a few dozen.

## Going deeper

Chacon and Straub, "Pro Git" (2nd ed.), section on git bisect; Zeller, "Why Programs Fail" (2nd ed., 2009), chapter on delta debugging.
