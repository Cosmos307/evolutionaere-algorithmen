package scp

import (
	"fmt"
	"os"
)

func ParseScpDataFile(path string) *Instance {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var inst Instance
	parseNums(&inst, file)
	parseColumCost(&inst, file)
	parseRowToCols(&inst, file)
	buildColToRows(&inst)

	return &inst
}

// Parse row and colum Number
func parseNums(inst *Instance, file *os.File) {
	_, err := fmt.Fscan(file, &inst.NumRows, &inst.NumCols)
	if err != nil {
		panic(err)
	}
}

// Parse cost of columns
func parseColumCost(inst *Instance, file *os.File) {
	inst.Costs = make([]int, inst.NumCols)
	for i := range inst.NumCols {
		tokensRead, err := fmt.Fscan(file, &inst.Costs[i])
		if err != nil || tokensRead == 0 {
			panic(err)
		}
	}
}

// Parse rows with NumCols which cover the row
func parseRowToCols(inst *Instance, file *os.File) {
	inst.RowToCols = make([][]int, inst.NumRows)
	for i := 0; i < inst.NumRows; i++ {
		var countColCoverRow int
		_, err := fmt.Fscan(file, &countColCoverRow)
		if err != nil {
			panic(err)
		}

		inst.RowToCols[i] = make([]int, countColCoverRow)
		for j := range countColCoverRow {
			_, err := fmt.Fscan(file, &inst.RowToCols[i][j])
			if err != nil {
				panic(err)
			}
			inst.RowToCols[i][j]--
		}
	}
}

func buildColToRows(inst *Instance) {
	inst.ColToRows = make([][]int, inst.NumCols)
	for row, cols := range inst.RowToCols {
		for _, col := range cols {
			inst.ColToRows[col] = append(inst.ColToRows[col], row)
		}
	}
}
