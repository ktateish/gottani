package main

import (
	"fmt"
	"sort"
)

//line example.com/lib/lib.go:3
type Entries struct {
	ids   []int
	names []string
}

func (e *Entries) Append(id int, name string) {
	e.ids = append(e.ids, id)
	e.names = append(e.names, name)
}

func (e *Entries) ForEachEntry(f func(id int, name string)) {
	for i, id := range e.ids {
		f(id, e.names[i])
	}
}

func (e *Entries) Len() int {
	return len(e.ids)
}

func (e *Entries) Less(i, j int) bool {
	return e.ids[i] < e.ids[j]
}

func (e *Entries) Swap(i, j int) {
	e.ids[i], e.ids[j] = e.ids[j], e.ids[i]
	e.names[i], e.names[j] = e.names[j], e.names[i]
}

type SortByNames Entries

func (e *SortByNames) Len() int {
	ee := (*Entries)(e)
	return ee.Len()
}

func (e *SortByNames) Less(i, j int) bool {
	return e.names[i] < e.names[j]
}

func (e *SortByNames) Swap(i, j int) {
	ee := (*Entries)(e)
	ee.Swap(i, j)
}

//line main.go:10
func main() {
	data := []struct {
		id   int
		name string
	}{
		{5, "vvv"},
		{2, "yyy"},
		{4, "www"},
		{1, "zzz"},
		{3, "xxx"},
	}

	e := &Entries{}
	for _, d := range data {
		e.Append(d.id, d.name)
	}

	sort.Sort(e)
	fmt.Println("Sorted by ID")
	e.ForEachEntry(func(id int, name string) {
		fmt.Printf("Entry: id:%d name:%s\n", id, name)
	})
	sort.Sort((*SortByNames)(e))
	fmt.Println("Sorted by Name")
	e.ForEachEntry(func(id int, name string) {
		fmt.Printf("Entry: id:%d name:%s\n", id, name)
	})
}
