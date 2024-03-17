package functools

// Reverse 反转切片
//
// 参数：
//  - slice []T：需要反转的切片
//
// 返回值：
//  - []T：反转后的切片
func Reverse[T int16 | int32 | int64 | uint | uint32 | uint64](slice []T) []T {
	reversed := make([]T, len(slice))
	for index, value := range slice {
		reversed[len(slice)-1-index] = value
	}
	return reversed
}
