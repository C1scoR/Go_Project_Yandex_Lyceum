package main

import "fmt"

func AddSalt(s string) func(string) string {
	return func(input string) string {
		return s + input
	}
}
func main() {
	salter := AddSalt("pepper_")
	fmt.Println(salter("steak"))
}
