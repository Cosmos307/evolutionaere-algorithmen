package algorithm

// ConvergencePoint records best fitness at a given evaluation count.
type ConvergencePoint struct {
	Eval       int
	BestCost   int
	Uncovered  int
	Redundancy int
	Score      float64
}
