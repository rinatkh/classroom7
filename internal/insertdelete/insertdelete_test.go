package insertdelete

import (
	"reflect"
	"testing"
)

func TestInsertAt(t *testing.T) {
	for _, test := range []struct {
		name   string
		values []int
		index  int
		value  int
		want   []int
		ok     bool
	}{
		{name: "begin", values: []int{2, 3}, index: 0, value: 1, want: []int{1, 2, 3}, ok: true},
		{name: "middle", values: []int{1, 3}, index: 1, value: 2, want: []int{1, 2, 3}, ok: true},
		{name: "end", values: []int{1, 2}, index: 2, value: 3, want: []int{1, 2, 3}, ok: true},
		{name: "invalid", values: []int{1, 2}, index: 3, value: 3, want: []int{1, 2}, ok: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			got, ok := InsertAt(test.values, test.index, test.value)
			if ok != test.ok || !reflect.DeepEqual(got, test.want) {
				t.Fatalf("InsertAt() = %v, %t; want %v, %t", got, ok, test.want, test.ok)
			}
		})
	}
}

func TestDeleteAt(t *testing.T) {
	values := make([]*int, 3, 3)
	for index := range values {
		value := index + 1
		values[index] = &value
	}
	got, ok := DeleteAt(values, 1)
	if !ok || len(got) != 2 || *got[0] != 1 || *got[1] != 3 {
		t.Fatalf("DeleteAt() = %v, %t", got, ok)
	}
	if values[2] != nil {
		t.Fatal("DeleteAt() не обнулила последнюю ячейку базового массива")
	}
	if same, valid := DeleteAt([]int{1}, -1); valid || !reflect.DeepEqual(same, []int{1}) {
		t.Fatalf("DeleteAt(invalid) = %v, %t", same, valid)
	}
}

func TestDeleteUnordered(t *testing.T) {
	got, ok := DeleteUnordered([]int{10, 20, 30, 40}, 1)
	if !ok || !reflect.DeepEqual(got, []int{10, 40, 30}) {
		t.Fatalf("DeleteUnordered() = %v, %t; want [10 40 30], true", got, ok)
	}
	if same, valid := DeleteUnordered([]int{1}, 1); valid || !reflect.DeepEqual(same, []int{1}) {
		t.Fatalf("DeleteUnordered(invalid) = %v, %t", same, valid)
	}
}

func TestMapOperations(t *testing.T) {
	values := make(map[string]int)
	Put(values, "ключ", 1)
	Put(values, "ключ", 2)
	if values["ключ"] != 2 {
		t.Fatalf("Put() left %v; want map[ключ:2]", values)
	}
	Remove(values, "ключ")
	Remove(values, "отсутствует")
	Remove(nil, "безопасно")
	if len(values) != 0 {
		t.Fatalf("Remove() left %v; want empty map", values)
	}
}

func TestExample(t *testing.T) {
	const want = "slice=[30 15]; map[Аня]=21"
	if got := Example(); got != want {
		t.Fatalf("Example() = %q; want %q", got, want)
	}
}
