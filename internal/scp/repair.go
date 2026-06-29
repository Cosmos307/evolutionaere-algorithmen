package scp

import "math"

// Repair ensures all rows are covered, then removes redundant cols.
func Repair(sol *Solution, inst *Instance) {
	greedyAdd(sol, inst)
	prune(sol, inst)
}

// greedyAdd adds cols until all rows covered.
// Picks col with best cost/newly-covered-rows ratio each step.
func greedyAdd(sol *Solution, inst *Instance) {
	coveredRows := make([]bool, inst.NumRows)
	for i, selected := range sol.SelectedCols {
		if selected {
			for _, row := range inst.ColToRows[i] {
				coveredRows[row] = true
			}
		}
	}

	// for each uncovered row, select best col to cover it
	for row, covered := range coveredRows {
		if covered {
			continue
		}

		// find best col: lowest cost per newly-covered-row ratio
		bestCol := -1
		bestRatio := math.MaxFloat64
		for _, col := range inst.RowToCols[row] {
			if sol.SelectedCols[col] {
				// selected col covers this row but row marked uncovered — coverage tracking bug
				panic("greedyAdd: selected col covers row marked uncovered — coverage tracking bug")
			}
			// count how many currently uncovered rows this col would fix
			newlyCovered := 0
			for _, covRow := range inst.ColToRows[col] {
				if !coveredRows[covRow] {
					newlyCovered++
				}
			}
			if newlyCovered == 0 {
				continue // col covers nothing new
			}
			// cheapest col per newly covered row wins
			ratio := float64(inst.Costs[col]) / float64(newlyCovered)
			if ratio < bestRatio {
				bestRatio = ratio
				bestCol = col
			}
		}

		if bestCol == -1 {
			// no col can cover remaining rows (should not happen)
			panic("greedyAdd: uncovered row has no candidate col — instance data corrupt")
		}

		// select col and update coverage
		sol.SelectedCols[bestCol] = true
		for _, covRow := range inst.ColToRows[bestCol] {
			coveredRows[covRow] = true
		}
	}
}

// prune removes selected cols that are redundant.
func prune(sol *Solution, inst *Instance) {
	coveredRows := make([]int, inst.NumRows)
	for i, selected := range sol.SelectedCols {
		if selected {
			for _, row := range inst.ColToRows[i] {
				coveredRows[row]++
			}
		}
	}

	for i, selected := range sol.SelectedCols {
		if !selected {
			continue
		}
		redundant := true
		for _, row := range inst.ColToRows[i] {
			if coveredRows[row] == 1 {
				redundant = false
				break
			}
		}
		if redundant {
			sol.SelectedCols[i] = false
			for _, row := range inst.ColToRows[i] {
				coveredRows[row]--
			}
		}
	}
}
