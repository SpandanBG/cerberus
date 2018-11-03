package main

import (
	"fmt"

	"./modules/database"
)

func main() {
	database.StartService("storage")
	fmt.Println(database.Create("config", map[string]interface{}{
		"name": "Spandan",
		"box": map[string]interface{} {
			"age": 21,
			"gender": "male",
		},
	}))
	fmt.Println(
		database.Read("config"),
	)
	fmt.Println(
		database.Update("config", map[string]interface{}{
			"name": "SudoKid",
			"box": map[string]interface{} {
				"age": 22,
			},
		}),
	)
	fmt.Println(
		database.Read("config"),
	)
}
