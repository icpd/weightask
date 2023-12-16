package weightask

import "sort"

type WeightSlice []int

func (p *WeightSlice) Sort() {
	sort.Sort(sort.Reverse(sort.IntSlice(*p)))
}

func (p *WeightSlice) Add(val int) {
	*p = append(*p, val)
}

func (p *WeightSlice) Remove(val int) {
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

func (p *WeightSlice) GetTopWeight() int {
	return (*p)[0]
}
