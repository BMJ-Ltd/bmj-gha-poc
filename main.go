package main

import (
	"fmt"
	"os"
)

func main() {
	myInput := os.Getenv("INPUT_ECR_NAME")

	output := fmt.Sprintf("Hello Mike  %s", myInput)

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, output))
}
