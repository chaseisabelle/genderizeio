package main

import (
	"fmt"
	"github.com/chaseisabelle/genderizer"
	"net/http"
	"os"
)

func main() {
	genderizations, err := (&genderizer.Genderizer{
		Client: &http.Client{},
		Endpoint: genderizer.ENDPOINT,
	}).Genderize(os.Args[1:]...)

	if err != nil {
		panic(err)
	}

	fmt.Println("genderizations:")

	for _, genderization := range genderizations {
		fmt.Println()
		fmt.Println(fmt.Sprintf("\tname: %s", genderization.Name))
		fmt.Println(fmt.Sprintf("\tgender: %s", genderization.Gender))
		fmt.Println(fmt.Sprintf("\tprobability: %f", genderization.Probability))
		fmt.Println(fmt.Sprintf("\tcount: %d", genderization.Count))
	}
}
