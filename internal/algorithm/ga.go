package algorithm

import (
	"math/rand"

	"github.com/Cosmos307/scp-ea/internal/scp"
)

type GAConfig struct {
	PInit          float64 // probability each col selected in initial solution
	PopSize        int     // population size, e.g. 30 or 50
	TournamentSize int     // tournament selection size, e.g. 2 or 3
	CrossoverProb  float64 // probability of crossover vs cloning, e.g. 0.7
	PMut           float64 // mutation rate per bit, 0 = use 1/n
	UseRepair      bool    // apply repair after crossover+mutation
	Budget         int     // max number of fitness evaluations
}

// tournamentSelect picks the best individual from a random tournament of size k.
func tournamentSelect(pop []*scp.Solution, fits []*scp.Fitness, k int, rng *rand.Rand) int {
	best := rng.Intn(len(pop))
	for i := 1; i < k; i++ {
		challenger := rng.Intn(len(pop))
		if scp.BetterThan(fits[challenger], fits[best]) {
			best = challenger
		}
	}
	return best
}

// uniformCrossover creates a child by randomly picking each bit from parent a or b.
func uniformCrossover(a, b *scp.Solution, rng *rand.Rand) *scp.Solution {
	child := &scp.Solution{SelectedCols: make([]bool, len(a.SelectedCols))}
	for i := range child.SelectedCols {
		if rng.Float64() < 0.5 {
			child.SelectedCols[i] = a.SelectedCols[i]
		} else {
			child.SelectedCols[i] = b.SelectedCols[i]
		}
	}
	return child
}

// worstIndex returns index of worst individual in population.
func worstIndex(fits []*scp.Fitness) int {
	worst := 0
	for i := 1; i < len(fits); i++ {
		if scp.BetterThan(fits[worst], fits[i]) {
			worst = i
		}
	}
	return worst
}

// RunGA runs the genetic algorithm on inst with the given config and rng.
// Returns best solution, its fitness, and convergence log (one point per eval).
func RunGA(inst *scp.Instance, cfg GAConfig, rng *rand.Rand) (*scp.Solution, *scp.Fitness, []ConvergencePoint) {
	pMut := cfg.PMut
	if pMut <= 0 {
		pMut = 1.0 / float64(inst.NumCols)
	}

	// initialize population
	pop := make([]*scp.Solution, cfg.PopSize)
	fits := make([]*scp.Fitness, cfg.PopSize)
	for i := range pop {
		pop[i] = NewRandomSolution(inst, cfg.PInit, rng)
		if cfg.UseRepair {
			scp.Repair(pop[i], inst)
		}
		fits[i] = scp.Evaluate(pop[i], inst)
	}
	evals := cfg.PopSize

	// track best
	bestIdx := 0
	for i := 1; i < cfg.PopSize; i++ {
		if scp.BetterThan(fits[i], fits[bestIdx]) {
			bestIdx = i
		}
	}
	best := copySolution(pop[bestIdx])
	bestFit := fits[bestIdx]

	conv := []ConvergencePoint{}
	for i := 0; i < evals; i++ {
		conv = append(conv, ConvergencePoint{Eval: i + 1, BestCost: bestFit.Cost, Uncovered: bestFit.Uncovered, Redundancy: bestFit.Redundancy, Score: bestFit.Score})
	}

	for evals < cfg.Budget {
		// select parents
		aIdx := tournamentSelect(pop, fits, cfg.TournamentSize, rng)
		bIdx := tournamentSelect(pop, fits, cfg.TournamentSize, rng)

		// crossover or clone
		var child *scp.Solution
		if rng.Float64() < cfg.CrossoverProb {
			child = uniformCrossover(pop[aIdx], pop[bIdx], rng)
		} else {
			child = copySolution(pop[aIdx])
		}

		// mutate
		mutateBits(child, pMut, rng)

		if cfg.UseRepair {
			scp.Repair(child, inst)
		}
		childFit := scp.Evaluate(child, inst)
		evals++

		// replace worst individual
		w := worstIndex(fits)
		if scp.BetterThan(childFit, fits[w]) {
			pop[w] = child
			fits[w] = childFit

			if scp.BetterThan(childFit, bestFit) {
				best = child
				bestFit = childFit
			}
		}
		conv = append(conv, ConvergencePoint{Eval: evals, BestCost: bestFit.Cost, Uncovered: bestFit.Uncovered, Redundancy: bestFit.Redundancy, Score: bestFit.Score})
	}

	return best, bestFit, conv
}
