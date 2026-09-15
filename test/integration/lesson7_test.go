package integration

import (
	"os/exec"
	"strings"
	"testing"
)

func TestLesson7Commands(t *testing.T) {
	commands := []struct {
		name string
		want string
	}{
		{name: "01_complexity", want: "O(1): value=40"},
		{name: "02_search", want: "binary index=5 comparisons=2"},
		{name: "03_maps", want: "duplicate=true"},
		{name: "04_insert_delete", want: "slice=[30 15]"},
		{name: "05_simple_sorts", want: "sorted=[1 2 4 5]"},
		{name: "06_quicksort", want: "Quick Sort: [2 3 5 7 9]"},
		{name: "07_benchmark_debug", want: "Измерения: make bench-search / make bench-sort"},
		{name: "08_leetcode", want: "LeetCode 704: target=9 -> 4"},
	}

	for _, test := range commands {
		t.Run(test.name, func(t *testing.T) {
			command := exec.Command("go", "run", "./cmd/"+test.name)
			command.Dir = "../.."
			output, err := command.CombinedOutput()
			if err != nil {
				t.Fatalf("go run ./cmd/%s: %v\n%s", test.name, err, output)
			}
			if !strings.Contains(string(output), test.want) {
				t.Fatalf("output of %s does not contain %q:\n%s", test.name, test.want, output)
			}
		})
	}
}
