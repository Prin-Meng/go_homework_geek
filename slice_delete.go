package main

import "fmt"

func deleteSliceElement1[T any](slice []T, index int) ([]T, T, error) {
	// 判断下标是否在切片的范围内
	length := len(slice)
	if index >= length || index < 0 {
		var zero T
		return nil, zero, indexOutRangeErr(length, index)
	}
	delElement := slice[index]
	// 使用子切片完成
	return append(slice[:index], slice[index+1:]...), delElement, nil

}

func deleteSliceElement2[T any](slice []T, index int) ([]T, T, error) {
	// 判断下标是否在切片的范围内
	length := len(slice)
	if index >= length || index < 0 {
		var zero T
		return nil, zero, indexOutRangeErr(length, index)
	}
	// 使用移位置完成
	delElement := slice[index]
	for i := index; i+1 < length; i++ {
		slice[i] = slice[i+1]
	}
	// 去掉最后一个元素
	slice = slice[:length-1]
	// 缩容操作
	slice = sliceShrink[T](slice)
	return slice, delElement, nil
}

func indexOutRangeErr(length int, index int) error {
	return fmt.Errorf("delete index out of range, the length is: %d and the index is %d", length, index)
}

func sliceShrink[T any](slice []T) []T {
	// 获取容量以及长度
	capacity, length := cap(slice), len(slice)

	// 默认不变化
	factor := float64(1)
	// 容量小于32不变
	if capacity <= 32 {
		factor = float64(1)
	} else if capacity <= 256 && capacity/length >= 4 {
		// 容量在33-256, 且使用率不足1/4
		factor = float64(0.5)
	} else if capacity > 256 && capacity/length >= 2 {
		// 容量大于256并且超过所使用长度的一半
		factor = float64(0.625)
	}
	capacity = int(float64(capacity) * factor)

	// 没有变化就返回
	if factor == float64(1) {
		return slice
	}

	sliceNew := make([]T, capacity)
	sliceNew = append(sliceNew, slice...)
	return sliceNew
}

func main() {
	sliceForTest := make([]int, 25, 1025)
	for i := 0; i < 25; i++ {
		sliceForTest[i] = i
	}
	//fmt.Println(sliceForTest)
	//fmt.Println(deleteSliceElement1[int](sliceForTest, 2))
	//fmt.Println(deleteSliceElement2[int](sliceForTest, 2))
	slice, _, _ := deleteSliceElement2[int](sliceForTest, 2)
	fmt.Printf("cureent slicne length is: %d", len(slice))
}
