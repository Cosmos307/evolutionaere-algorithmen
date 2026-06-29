package experiment

type Result struct {
	Algorithm  string
	Instance   string
	Seed       int64
	UseRepair  bool
	Budget     int
	BestCost   int
	Uncovered  int
	Redundancy int
	RuntimeMs  int64
	GapPercent float64 // -1 if BestKnown unknown or solution invalid
}
