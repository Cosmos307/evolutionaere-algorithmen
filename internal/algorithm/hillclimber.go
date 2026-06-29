// Package algorithm implements the evolutionary algorithms to tackle the set coverage problem
package algorithm

import (
	"math/rand"

	"github.com/Cosmos307/scp-ea/internal/scp"
)

type HillclimberConfig struct {
	PInit       float64 // probability each col is selected in initial solution
	AcceptEqual bool    // accept equally good neighbors (neutral moves)
	UseRepair   bool    // apply repair after each mutation
	Budget      int     // max number of fitness evaluations
}

// NewRandomSolution creates a solution where each col is selected with probability pInit
func NewRandomSolution(inst *scp.Instance, pInit float64, rng *rand.Rand) *scp.Solution {
	sol := &scp.Solution{SelectedCols: make([]bool, inst.NumCols)}
	for i := range sol.SelectedCols {
		if rng.Float64() < pInit {
			sol.SelectedCols[i] = true
		}
	}
	return sol
}

// copySolution returns a deep copy of sol
func copySolution(sol *scp.Solution) *scp.Solution {
	cp := &scp.Solution{SelectedCols: make([]bool, len(sol.SelectedCols))}
	copy(cp.SelectedCols, sol.SelectedCols)
	return cp
}

// flipRandomBit flips exactly one random bit in sol
func flipRandomBit(sol *scp.Solution, rng *rand.Rand) {
	i := rng.Intn(len(sol.SelectedCols))
	sol.SelectedCols[i] = !sol.SelectedCols[i]
}

// RunHillclimber runs a hillclimber on inst with the given config and rng.
// Returns best solution, its fitness, and convergence log (one point per eval).
func RunHillclimber(inst *scp.Instance, cfg HillclimberConfig, rng *rand.Rand) (*scp.Solution, *scp.Fitness, []ConvergencePoint) {
	current := NewRandomSolution(inst, cfg.PInit, rng)
	if cfg.UseRepair {
		scp.Repair(current, inst)
	}
	bestFit := scp.Evaluate(current, inst)
	best := copySolution(current)
	evals := 1

	conv := []ConvergencePoint{{Eval: evals, BestCost: bestFit.Cost, Uncovered: bestFit.Uncovered, Redundancy: bestFit.Redundancy, Score: bestFit.Score}}

	for evals < cfg.Budget {
		candidate := copySolution(current)
		flipRandomBit(candidate, rng)
		if cfg.UseRepair {
			scp.Repair(candidate, inst)
		}
		candidateFit := scp.Evaluate(candidate, inst)
		evals++

		// accept if strictly better, or equal when AcceptEqual is on
		if scp.BetterThan(candidateFit, bestFit) || (cfg.AcceptEqual && !scp.BetterThan(bestFit, candidateFit)) {
			current = candidate
			if scp.BetterThan(candidateFit, bestFit) {
				best = candidate
				bestFit = candidateFit
			}
		}
		conv = append(conv, ConvergencePoint{Eval: evals, BestCost: bestFit.Cost, Uncovered: bestFit.Uncovered, Redundancy: bestFit.Redundancy, Score: bestFit.Score})
	}

	return best, bestFit, conv
}
