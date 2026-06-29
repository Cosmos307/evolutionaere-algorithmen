package scp

func Evaluate(sol *Solution, inst *Instance) *Fitness {
	fit := Fitness{}
	coveredRows := make([]int, inst.NumRows)

	// Pass 1: mark all coverage, count uncovered, sum costs.
	for i, colSelected := range sol.SelectedCols {
		if colSelected {
			fit.Cost += inst.Costs[i]
			for _, row := range inst.ColToRows[i] {
				coveredRows[row]++
			}
		}
	}
	for _, count := range coveredRows {
		if count == 0 {
			fit.Uncovered++
		}
	}

	// Pass 2: count redundant cols.
	for i, colSelected := range sol.SelectedCols {
		if colSelected {
			redundant := true
			for _, row := range inst.ColToRows[i] {
				if coveredRows[row] == 1 {
					redundant = false
					break
				}
			}
			if redundant {
				fit.Redundancy++
			}
		}
	}

	penaltyUncovered := 1
	for _, c := range inst.Costs {
		penaltyUncovered += c
	}
	fit.Score = float64(penaltyUncovered*fit.Uncovered) + float64(fit.Cost) + 0.01*float64(fit.Redundancy)

	return &fit
}

func BetterThan(a, b *Fitness) bool {
	if a.Uncovered < b.Uncovered {
		return true
	} else if a.Uncovered == b.Uncovered {
		if a.Cost < b.Cost {
			return true
		} else if a.Cost == b.Cost {
			if a.Redundancy < b.Redundancy {
				return true
			}
		}
	}
	return false
}
