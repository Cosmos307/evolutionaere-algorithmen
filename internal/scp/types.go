// Package scp contains all of the scp domain
package scp

type Instance struct {
	Name      string
	NumRows   int
	NumCols   int
	Costs     []int
	RowToCols [][]int
	ColToRows [][]int
	BestKnown int
}

type Solution struct {
	SelectedCols []bool
}

type BitSolution struct {
	Bits []bool
}

type Fitness struct {
	Uncovered  int
	Cost       int
	Redundancy int
	Score      float64
}
