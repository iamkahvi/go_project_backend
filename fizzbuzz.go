package main

import (
	"fmt"
)

func main2() {
	s := make([]string, 3)
	s[0] = "kahvi"
	s[1] = "patel"
	s[2] = "joe"
	s = append(s, "roger")
	fmt.Println(s)
	fmt.Println(s[:1])

	n := make(map[string]int)
	n["satish"] = 3
	n["joe"] = 4
	n["steve"] = 30
	res, err := n["kahvi"]
	fmt.Println(n)
	if err {
		fmt.Println(res)
	}

	for k, v := range n {
		fmt.Printf("%s -> %d\n", k, v)
	}

}

func add(a int, b int) int {
	return a + b
}

func fizzbuzz() {
	for i := 1; i <= 100; i++ {
		x, y := i%3 == 0, i%5 == 0
		if x {
			fmt.Print("Fizz")
		}
		if y {
			fmt.Print("Buzz")
		}
		if !x && !y {
			fmt.Print(i)
		}
		fmt.Print("\n")
	}
}
