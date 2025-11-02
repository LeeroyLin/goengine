package pool

import "github.com/LeeroyLin/goengine/core/elog"

// 获得等于或大于该值的 2的整数次幂 的值
func getNearGreaterPower(val int) int {
	result := 1
	for result < val {
		result <<= 1
	}
	return result
}

// 获得等于或小于该值的 2的整数次幂 的值
func getNearLessPower(val int) int {
	if val < 1 {
		elog.Error("[Pool] wrong val to getNearLessPower.", val)
		return 0
	}

	lastResult := 1
	result := 1
	for {
		if result > val {
			return lastResult
		}

		lastResult = result

		result <<= 1
	}
}
