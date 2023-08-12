package priortask

import "sort"

type PrioritySlice []int

func (p *PrioritySlice) Sort() {
	sort.Sort(sort.Reverse(sort.IntSlice(*p)))
}

func (p *PrioritySlice) Add(val int) {
	*p = append(*p, val)
}

func (p *PrioritySlice) Remove(val int) {
	index := -1
	for i, v := range *p {
		if v == val {
			index = i
			break
		}
	}

	if index >= 0 {
		*p = append((*p)[:index], (*p)[index+1:]...)
	}
}

func (p *PrioritySlice) GetTopPriority() int {
	return (*p)[0]
}
