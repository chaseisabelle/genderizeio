package main

import (
	"fmt"
	"github.com/chaseisabelle/genderizer"
	"os"
)

func main() {
	genderizations, err := genderizer.Genderize(os.Args[1:]...)

	if err != nil {
		panic(err)
	}

	fmt.Println("genderizations:")

	for _, gender := range genderizations {
		fmt.Println()
		fmt.Println(fmt.Sprintf("\tname: %s", gender.Name))
		fmt.Println(fmt.Sprintf("\tgender: %s", gender.Gender))
		fmt.Println(fmt.Sprintf("\tprobability: %f", gender.Probability))
		fmt.Println(fmt.Sprintf("\tcount: %d", gender.Count))
	}
}
