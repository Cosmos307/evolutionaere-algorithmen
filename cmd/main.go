package main

import (
	"fmt"

	"github.com/Cosmos307/scp-ea/internal/experiment"
)

func main() {
	baseCfg := experiment.Config{
		InstancePaths:  []string{"data/scp41.txt"},
		Seeds:          []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30},
		Budget:         50000,
		UseRepair:      true,
		PInit:          0.1,
		AcceptEqual:    false,
		PMutFactor:     1.0, // use the standard rate 1/l
		PopSize:        30,
		TournamentSize: 3,
		CrossoverProb:  0.7,
		LogConvergence: true,
		LogInterval:    100,
		OutputDir:      "results/raw/baseline",
	}

	studies := []struct {
		name string
		cfg  experiment.Config
	}{
		{name: "baseline", cfg: baseCfg},
		{name: "pmut_0.5_over_l", cfg: withOutputDir(baseCfg, "results/raw/pmut_0.5_over_l", func(cfg *experiment.Config) {
			cfg.PMutFactor = 0.5
		})},
		{name: "pmut_2.0_over_l", cfg: withOutputDir(baseCfg, "results/raw/pmut_2.0_over_l", func(cfg *experiment.Config) {
			cfg.PMutFactor = 2.0
		})},
		{name: "ga_pop_50", cfg: withOutputDir(baseCfg, "results/raw/ga_pop_50", func(cfg *experiment.Config) {
			cfg.PopSize = 50
		})},
		{name: "ga_tournament_2", cfg: withOutputDir(baseCfg, "results/raw/ga_tournament_2", func(cfg *experiment.Config) {
			cfg.TournamentSize = 2
		})},
		{name: "ga_crossover_0.9", cfg: withOutputDir(baseCfg, "results/raw/ga_crossover_0.9", func(cfg *experiment.Config) {
			cfg.CrossoverProb = 0.9
		})},
	}

	for _, study := range studies {
		fmt.Printf("\n=== STUDY: %s ===\n", study.name)
		experiment.Run(study.cfg)
	}
}

func withOutputDir(base experiment.Config, outputDir string, update func(*experiment.Config)) experiment.Config {
	cfg := base
	cfg.OutputDir = outputDir
	update(&cfg)
	return cfg
}
