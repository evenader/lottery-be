// Copyright (c) 2026 evenader. All rights reserved.
package fastpicker

import (
	"math/rand"
	"sort"
	"sync"
)

// 定义一个全局池，用于复用 []int64 切片
// 这里的关键是：我们要复用的是 cdf 切片
var cdfPool = sync.Pool{
	New: func() interface{} {
		// 初始容量可以根据业务场景设定，比如 10000
		return make([]int64, 0, 10000)
	},
}

type FastPicker struct {
	ratios      []int64
	cdf         []int64
	totalWeight int64
}

// NewFastPicker 借用池中的内存进行初始化
func NewFastPicker(ratios []int64) *FastPicker {
	n := len(ratios)
	// 从池中获取切片
	cdf := cdfPool.Get().([]int64)

	// 调整切片长度以匹配当前需求
	if cap(cdf) < n {
		// 如果池里的切片太小，扩容
		cdf = make([]int64, n)
	} else {
		cdf = cdf[:n]
	}

	var sum int64
	for i, v := range ratios {
		sum += v
		cdf[i] = sum
	}

	return &FastPicker{
		ratios:      ratios,
		cdf:         cdf,
		totalWeight: sum,
	}
}

// Recycle 必须手动调用，将内存还回池中
func (p *FastPicker) Recycle() {
	// 重置切片，但保留底层数组
	p.cdf = p.cdf[:0]
	cdfPool.Put(p.cdf)
}

func (p *FastPicker) PickAndRemove(rng *rand.Rand) int {
	if p.totalWeight <= 0 || len(p.cdf) == 0 {
		return -1
	}

	target := rng.Int63n(p.totalWeight)
	idx := sort.Search(len(p.cdf), func(i int) bool {
		return p.cdf[i] > target
	})

	pickedWeight := p.ratios[idx]

	// 增量更新 CDF
	for i := idx; i < len(p.cdf); i++ {
		p.cdf[i] -= pickedWeight
	}

	p.totalWeight -= pickedWeight
	p.cdf = append(p.cdf[:idx], p.cdf[idx+1:]...)
	p.ratios = append(p.ratios[:idx], p.ratios[idx+1:]...)

	return idx
}
