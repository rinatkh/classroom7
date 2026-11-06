// Package testutil хранит общие входы тестов, а не реализации алгоритмов.
package testutil

import "math"

type Case struct {
	Name   string
	Values []int
}

func Cases() []Case {
	return []Case{
		{"nil", nil}, {"empty", []int{}}, {"one", []int{7}},
		{"two-sorted", []int{1, 2}}, {"two-reversed", []int{2, 1}},
		{"equal", []int{4, 4, 4, 4}}, {"duplicates", []int{3, 1, 3, 2, 1, 3}},
		{"negative", []int{-9, -1, -5, -3}}, {"mixed", []int{-3, 0, 8, -1, 8}},
		{"sorted", []int{-4, -1, 0, 3, 9}}, {"reverse", []int{9, 3, 0, -1, -4}},
		{"extremes", []int{math.MaxInt, 0, math.MinInt, -1, math.MaxInt}},
		{"nearly", []int{1, 2, 4, 3, 5, 6}}, {"zero", []int{0, 0, 0}},
		{"odd-length", []int{8, 2, 7, 1, 6, 0, 5}},
	}
}

type SearchCase struct {
	Name   string
	Values []int
	Target int
}

func SearchCases() []SearchCase {
	return []SearchCase{
		{"nil", nil, 3}, {"empty", []int{}, 0}, {"one-found", []int{7}, 7},
		{"one-missing", []int{7}, 5}, {"two-first", []int{1, 3}, 1},
		{"two-last", []int{1, 3}, 3}, {"two-gap", []int{1, 3}, 2},
		{"first", []int{-8, -2, 0, 3, 10}, -8}, {"middle", []int{-8, -2, 0, 3, 10}, 0},
		{"last", []int{-8, -2, 0, 3, 10}, 10}, {"below", []int{1, 3, 5}, -1},
		{"above", []int{1, 3, 5}, 8}, {"duplicates", []int{1, 2, 2, 2, 5}, 2},
		{"all-equal", []int{4, 4, 4, 4}, 4}, {"equal-missing", []int{4, 4, 4}, 3},
		{"extreme-low", []int{math.MinInt, 0, math.MaxInt}, math.MinInt},
		{"extreme-high", []int{math.MinInt, 0, math.MaxInt}, math.MaxInt},
		{"negative", []int{-9, -5, -2}, -5},
	}
}
