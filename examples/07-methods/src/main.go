package main

import (
	"fmt"
	"sort"

	"example.com/lib"
)

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

	e := &lib.Entries{}
	for _, d := range data {
		e.Append(d.id, d.name)
	}

	sort.Sort(e)
	fmt.Println("Sorted by ID")
	e.ForEachEntry(func(id int, name string) {
		fmt.Printf("Entry: id:%d name:%s\n", id, name)
	})
	sort.Sort((*lib.SortByNames)(e))
	fmt.Println("Sorted by Name")
	e.ForEachEntry(func(id int, name string) {
		fmt.Printf("Entry: id:%d name:%s\n", id, name)
	})
}
