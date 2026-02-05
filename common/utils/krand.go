// Copyright (c) 2026 evenader. All rights reserved.

package utils

import "crypto/rand"

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

var (
	kinds = [][]byte{
		[]byte("0123456789"),
		[]byte("abcdefghijklmnopqrstuvwxyz"),
		[]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	allKinds = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// 训练营自带的，有性能问题
// ### 1. 致命缺陷分析
// A. 每次调用都 Seed（性能与随机性的双重自杀）
// Go
//
// rand.Seed(time.Now().UnixNano())
// 性能问题： rand.Seed 是一个重量级操作，它会重置随机生成器的状态。在高并发下，CPU 会浪费大量时间在初始化种子，而不是生成随机数。
//
// 随机性失效： 如果你在一个极短的时间内（比如 1 微秒内）并发调用这个函数两次，time.Now().UnixNano() 返回的值可能是完全一样的。
//
// 结果： 这会导致这两个并发请求生成的“随机”字符串完全相同。在注册场景下，这会直接导致唯一索引冲突。
//
// B. 全局锁竞争（并发瓶颈）
// rand.Intn 内部使用了一个全局的随机对象，它是带锁的。
//
// 在高并发（如 5G 场景下的信令处理）中，所有 Goroutine 都会去争抢这把锁。你的服务吞吐量会因为这个 Krand 函数而遇到天花板。
//
// C. 伪随机的可预测性
// math/rand 生成的是伪随机数。如果攻击者连续观察到几个随机生成的昵称，他们有可能反推出你的种子，从而预测出下一个注册用户的随机昵称。
// 随机字符串
//
//	func Krand(size int, kind int) string {
//		ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
//		is_all := kind > 2 || kind < 0
//		rand.Seed(time.Now().UnixNano())
//		for i := 0; i < size; i++ {
//			if is_all { // random ikind
//				ikind = rand.Intn(3)
//			}
//			scope, base := kinds[ikind][0], kinds[ikind][1]
//			result[i] = uint8(base + rand.Intn(scope))
//		}
//		return string(result)
//	}

func Krand(size int, kind int) string {
	var source []byte
	switch kind {
	case KC_RAND_KIND_NUM:
		source = kinds[KC_RAND_KIND_NUM]
	case KC_RAND_KIND_LOWER:
		source = kinds[KC_RAND_KIND_LOWER]
	case KC_RAND_KIND_UPPER:
		source = kinds[KC_RAND_KIND_UPPER]
	case KC_RAND_KIND_ALL:
		source = allKinds
	default:
		source = allKinds
	}
	// 2. 一次性获取足够的随机字节（核心优化！）
	// 这里我们直接向内核要 size 个随机字节
	randomBytes := make([]byte, size)
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果内核熵池炸了，必须处理
		panic("crypto/rand failed: " + err.Error())
	}

	// 3. 在内存中完成映射（不涉及系统调用）
	result := make([]byte, size)
	sourceLen := byte(len(source))
	for i := 0; i < size; i++ {
		// 使用取模运算将 0-255 映射到字符集索引
		// 注意：取模在数学上会有微小的分布不均，但在昵称场景下完全可以忽略
		idx := randomBytes[i] % sourceLen
		result[i] = source[idx]
	}

	return string(result)
}
