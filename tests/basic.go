package tests

import "fmt"

func main()  {
	to := []string{"joni","dinda","fitri"}
	cc := []string{"siti","lala"}
	result := append(to, cc...)
	fmt.Println(result)
}