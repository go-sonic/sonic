package util

func RainbowPage(page, total, display int) []int {
	isEven := display%2 == 0
	left := display / 2
	right := display / 2
	length := display
	if isEven {
		right++
	}
	if total < display {
		length = total
	}
	result := make([]int, length, length)
	if total >= display {
		if page <= left {
			for i := 0; i < length; i++ {
				result[i] = i + 1
			}
		} else if page > total-right {
			for i := 0; i < length; i++ {
				result[i] = i + total - display + 1
			}
		} else {
			for i := 0; i < length; i++ {
				if isEven {
					result[i] = i + page - length + 1
				} else {
					result[i] = i + page - length
				}
			}
		}
	} else {
		for i := 0; i < length; i++ {
			result[i] = i + 1
		}
	}
	return result
}
