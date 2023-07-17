package main

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func main() {

	models, err := getModels()
	if err != nil {
		panic(err)
	}

	prompt := promptui.Select{
		Label: "Select a Model",
		Items: models,
	}

	i, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	downloadModel(models[i])

	fmt.Printf("You choose %q\n", result)
}
