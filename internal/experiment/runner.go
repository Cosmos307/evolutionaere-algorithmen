// Package experiment contains the run function to compare the 3 evolutionary algorithms
package experiment

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Cosmos307/scp-ea/internal/algorithm"
	"github.com/Cosmos307/scp-ea/internal/scp"
)

func resolvePMut(cfg Config, numCols int) float64 {
	if cfg.PMutFactor > 0 {
		return cfg.PMutFactor / float64(numCols)
	}
	return 0
}

// Run executes all algorithms × seeds × instances, prints results, then prints summary.
func Run(cfg Config) {
	fmt.Println("algorithm,instance,seed,repair,budget,best_cost,uncovered,redundancy,runtime_ms,gap_percent")

	totals := map[string]int{"hillclimber": 0, "stochastic_hillclimber": 0, "ga": 0}
	counts := map[string]int{"hillclimber": 0, "stochastic_hillclimber": 0, "ga": 0}

	for _, path := range cfg.InstancePaths {
		inst := scp.ParseScpDataFile(path)
		inst.Name = filepath.Base(path)

		for _, seed := range cfg.Seeds {
			r := runHillclimber(inst, cfg, seed)
			totals["hillclimber"] += r.BestCost
			counts["hillclimber"]++

			r = runStochasticHillclimber(inst, cfg, seed)
			totals["stochastic_hillclimber"] += r.BestCost
			counts["stochastic_hillclimber"]++

			r = runGA(inst, cfg, seed)
			totals["ga"] += r.BestCost
			counts["ga"]++
		}
	}

	fmt.Println()
	fmt.Println("=== SUMMARY (avg best_cost) ===")
	algs := []string{"hillclimber", "stochastic_hillclimber", "ga"}
	bestAlg := ""
	bestAvg := -1.0
	for _, alg := range algs {
		avg := float64(totals[alg]) / float64(counts[alg])
		fmt.Printf("%-15s avg_cost=%.2f\n", alg, avg)
		if bestAvg < 0 || avg < bestAvg {
			bestAvg = avg
			bestAlg = alg
		}
	}
	fmt.Printf("\nBest algorithm: %s (avg_cost=%.2f)\n", bestAlg, bestAvg)
}

func printResult(r Result) {
	fmt.Printf("%s,%s,%d,%v,%d,%d,%d,%d,%d,%.2f\n",
		r.Algorithm, r.Instance, r.Seed, r.UseRepair,
		r.Budget, r.BestCost, r.Uncovered, r.Redundancy,
		r.RuntimeMs, r.GapPercent,
	)
}

func calcGap(bestCost, bestKnown int) float64 {
	if bestKnown <= 0 {
		return -1
	}
	return 100.0 * float64(bestCost-bestKnown) / float64(bestKnown)
}

func writeConvergenceCSV(cfg Config, algName, instName string, seed int64, conv []algorithm.ConvergencePoint) {
	if !cfg.LogConvergence || len(conv) == 0 {
		return
	}
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output dir: %v\n", err)
		return
	}
	fname := fmt.Sprintf("convergence_%s_%s_seed%d.csv", instName, algName, seed)
	path := filepath.Join(cfg.OutputDir, fname)
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create convergence file: %v\n", err)
		return
	}
	defer f.Close()

	fmt.Fprintln(f, "eval,best_cost,uncovered,redundancy,score")
	interval := cfg.LogInterval
	if interval <= 0 {
		interval = 1
	}
	for _, p := range conv {
		if p.Eval%interval == 0 || p.Eval == 1 || p.Eval == len(conv) {
			fmt.Fprintf(f, "%d,%d,%d,%d,%.4f\n", p.Eval, p.BestCost, p.Uncovered, p.Redundancy, p.Score)
		}
	}
}

func runHillclimber(inst *scp.Instance, cfg Config, seed int64) Result {
	hcfg := algorithm.HillclimberConfig{
		PInit:       cfg.PInit,
		AcceptEqual: cfg.AcceptEqual,
		UseRepair:   cfg.UseRepair,
		Budget:      cfg.Budget,
	}
	rng := rand.New(rand.NewSource(seed))
	start := time.Now()
	_, fit, conv := algorithm.RunHillclimber(inst, hcfg, rng)
	ms := time.Since(start).Milliseconds()

	writeConvergenceCSV(cfg, "hillclimber", inst.Name, seed, conv)

	gap := -1.0
	if fit.Uncovered == 0 {
		gap = calcGap(fit.Cost, inst.BestKnown)
	}
	r := Result{
		Algorithm:  "hillclimber",
		Instance:   inst.Name,
		Seed:       seed,
		UseRepair:  cfg.UseRepair,
		Budget:     cfg.Budget,
		BestCost:   fit.Cost,
		Uncovered:  fit.Uncovered,
		Redundancy: fit.Redundancy,
		RuntimeMs:  ms,
		GapPercent: gap,
	}
	printResult(r)
	return r
}

func runStochasticHillclimber(inst *scp.Instance, cfg Config, seed int64) Result {
	ecfg := algorithm.StochasticHillclimberConfig{
		PInit:     cfg.PInit,
		PMut:      resolvePMut(cfg, inst.NumCols),
		UseRepair: cfg.UseRepair,
		Budget:    cfg.Budget,
	}
	rng := rand.New(rand.NewSource(seed))
	start := time.Now()
	_, fit, conv := algorithm.RunStochasticHillclimber(inst, ecfg, rng)
	ms := time.Since(start).Milliseconds()

	writeConvergenceCSV(cfg, "stochastic_hillclimber", inst.Name, seed, conv)

	gap := -1.0
	if fit.Uncovered == 0 {
		gap = calcGap(fit.Cost, inst.BestKnown)
	}
	r := Result{
		Algorithm:  "stochastic_hillclimber",
		Instance:   inst.Name,
		Seed:       seed,
		UseRepair:  cfg.UseRepair,
		Budget:     cfg.Budget,
		BestCost:   fit.Cost,
		Uncovered:  fit.Uncovered,
		Redundancy: fit.Redundancy,
		RuntimeMs:  ms,
		GapPercent: gap,
	}
	printResult(r)
	return r
}

func runGA(inst *scp.Instance, cfg Config, seed int64) Result {
	gcfg := algorithm.GAConfig{
		PInit:          cfg.PInit,
		PopSize:        cfg.PopSize,
		TournamentSize: cfg.TournamentSize,
		CrossoverProb:  cfg.CrossoverProb,
		PMut:           resolvePMut(cfg, inst.NumCols),
		UseRepair:      cfg.UseRepair,
		Budget:         cfg.Budget,
	}
	rng := rand.New(rand.NewSource(seed))
	start := time.Now()
	_, fit, conv := algorithm.RunGA(inst, gcfg, rng)
	ms := time.Since(start).Milliseconds()

	writeConvergenceCSV(cfg, "ga", inst.Name, seed, conv)

	gap := -1.0
	if fit.Uncovered == 0 {
		gap = calcGap(fit.Cost, inst.BestKnown)
	}
	r := Result{
		Algorithm:  "ga",
		Instance:   inst.Name,
		Seed:       seed,
		UseRepair:  cfg.UseRepair,
		Budget:     cfg.Budget,
		BestCost:   fit.Cost,
		Uncovered:  fit.Uncovered,
		Redundancy: fit.Redundancy,
		RuntimeMs:  ms,
		GapPercent: gap,
	}
	printResult(r)
	return r
}
