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

	//nolint:gosimple
	result := make([]int, length, length)
	if total >= display {
		switch {
		case page <= left:
			for i := 0; i < length; i++ {
				result[i] = i + 1
			}
		case page > total-right:
			for i := 0; i < length; i++ {
				result[i] = i + total - display + 1
			}
		default:
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
