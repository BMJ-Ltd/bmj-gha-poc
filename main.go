package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

// function to append a string to a slice of strings
func appendString(slice []string, data string) []string {
	return append(slice, data)
}

// create a functio to list all ecr images in a repository
func listImages(repositoryName string) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create ECR Service Client
	svc := ecr.New(sess)

	result, err := svc.ListImages(&ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeRepositoryNotFoundException:
				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	slice := []string{}

	// iterate over the images and print them out
	for _, image := range result.ImageIds {
		slice = appendString(slice, *image.ImageTag)
		fmt.Println(*image.ImageTag)
	}
	sort.Strings(slice)
	fmt.Println("----------------------------------------------------")

	//iterate over the sorted images and print them out
	for _, image := range slice {
		fmt.Println(image)
	}

	//fmt.Println(slice)

}

func main() {
	myInput := os.Getenv("INPUT_ECR_NAME")

	output := fmt.Sprintf("Hello Mike  %s", myInput)

	listImages(myInput)

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, output))
}
