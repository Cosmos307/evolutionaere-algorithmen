package algorithm

import (
	"math/rand"

	"github.com/Cosmos307/scp-ea/internal/scp"
)

type StochasticHillclimberConfig struct {
	PInit     float64 // probability each col is selected in initial solution
	PMut      float64 // mutation rate per bit, standard = 1/NumCols
	UseRepair bool    // apply repair after each mutation
	Budget    int     // max number of fitness evaluations
}

// mutateBits flips each bit independently with probability pMut.
func mutateBits(sol *scp.Solution, pMut float64, rng *rand.Rand) {
	for i := range sol.SelectedCols {
		if rng.Float64() < pMut {
			sol.SelectedCols[i] = !sol.SelectedCols[i]
		}
	}
}

// RunStochasticHillclimber runs hillclimbing with independent per-bit mutation.
// Returns best solution, its fitness, and convergence log (one point per eval).
func RunStochasticHillclimber(inst *scp.Instance, cfg StochasticHillclimberConfig, rng *rand.Rand) (*scp.Solution, *scp.Fitness, []ConvergencePoint) {
	pMut := cfg.PMut
	if pMut <= 0 {
		pMut = 1.0 / float64(inst.NumCols) // standard rate
	}

	current := NewRandomSolution(inst, cfg.PInit, rng)
	if cfg.UseRepair {
		scp.Repair(current, inst)
	}
	currentFit := scp.Evaluate(current, inst)
	best := copySolution(current)
	bestFit := currentFit
	evals := 1

	conv := []ConvergencePoint{{Eval: evals, BestCost: bestFit.Cost, Uncovered: bestFit.Uncovered, Redundancy: bestFit.Redundancy, Score: bestFit.Score}}

	for evals < cfg.Budget {
		candidate := copySolution(current)
		mutateBits(candidate, pMut, rng)
		if cfg.UseRepair {
			scp.Repair(candidate, inst)
		}
		candidateFit := scp.Evaluate(candidate, inst)
		evals++

		// accept if not worse (equal or better)
		if !scp.BetterThan(currentFit, candidateFit) {
			current = candidate
			currentFit = candidateFit
			if scp.BetterThan(candidateFit, bestFit) {
				best = candidate
				bestFit = candidateFit
			}
		}
		conv = append(conv, ConvergencePoint{Eval: evals, BestCost: bestFit.Cost, Uncovered: bestFit.Uncovered, Redundancy: bestFit.Redundancy, Score: bestFit.Score})
	}

	return best, bestFit, conv
}
