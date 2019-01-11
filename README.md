# genderizer
*golang package for genderize.io*

---

[https://genderize.io/](https://genderize.io/)

---
### example

```go
package main

import (
	"fmt"
	"github.com/chaseisabelle/genderizer"
	"os"
)

func main() {
	genderizations, err := genderizer.Genderize("", os.Args[1:]...)

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
```
*running the example...*
```
$ go run -race main.go chase isabelle
genders:

	name: chase
	gender: male
	probability: 0.960000
	count: 306

	name: isabelle
	gender: female
	probability: 1.000000
	count: 867
```
