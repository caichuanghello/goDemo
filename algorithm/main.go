package main

import (
)

func main(){



}

//冒泡排序 : 每次都只选出一个最大或者小的值
func bubbleSort(arr []int)[]int{
	count:=len(arr)

	if count <=1 {
		return arr
	}
	for i:=0;i<count-1;i++ {
		for j:=1+i;j<count;j++{
			if arr[i] > arr[j] {
				arr[i],arr[j] = arr[j],arr[i]
			}
		}
	}
	return arr

}



