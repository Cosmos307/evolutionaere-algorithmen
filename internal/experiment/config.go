package experiment

type Config struct {
	// shared experiment setup
	InstancePaths []string
	Seeds         []int64
	Budget        int     // max fitness evaluations per run
	UseRepair     bool    // repair infeasible solutions after mutation/crossover
	PInit         float64 // initial probability that a column is selected

	// shared mutation parameters for stochastic hillclimber and GA
	PMutFactor float64 // per-bit mutation rate = PMutFactor / n, e.g. 1.0 => 1/n

	// hillclimber only
	AcceptEqual bool // accept neutral moves

	// GA only
	PopSize        int
	TournamentSize int
	CrossoverProb  float64

	// convergence logging
	LogConvergence bool   // write per-run convergence CSV files
	LogInterval    int    // log every N evals, 0 = every eval
	OutputDir      string // directory for CSV output
}
