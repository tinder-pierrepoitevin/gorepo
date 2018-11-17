package main

import "fmt"
import "strconv"

func convertToBase7(num int) string {
	if num == 0 {
		return "0"
	}
	return convertToBase7rec(num)
}

func convertToBase7rec(num int) string {
	if num == 0 {
		return ""
	}
	if num < 0 {
		return "-" + convertToBase7rec(-num)
	}
	return convertToBase7rec(num/7) + strconv.Itoa(num%7)
}

func main() {
	fmt.Printf("-1: %s \n", convertToBase7(-1))
	fmt.Printf("1: %s \n", convertToBase7(1))
	fmt.Printf("10: %s \n", convertToBase7(10))
	fmt.Printf("-1000000: %s \n", convertToBase7(-1000000))
	fmt.Printf("1000000: %s \n", convertToBase7(1000000))
	fmt.Printf("777777: %s \n", convertToBase7(777777))
}
