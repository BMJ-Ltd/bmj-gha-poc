package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var major, minor, patch int

// function increment patch number
func incrementPatch() {
	patch++
}
func incrementMinor() {
	minor++
	patch = 0
}
func incrementMajor() {
	major++
	minor = 0
	patch = 0
}

// function to append a string to a slice of strings
func appendString(slice []string, data string) []string {
	return append(slice, data)
}

// function to parse version number
func parseVn(n string) (string, error) {
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	match := re.FindStringSubmatch(n)
	if match == nil {
		return "", fmt.Errorf("invalid version number: %s", n)
	}
	// remove all but numers from string and convert to int
	major, _ = strconv.Atoi(regexp.MustCompile(`\D`).ReplaceAllString(match[1], ""))
	minor, _ = strconv.Atoi(regexp.MustCompile(`\D`).ReplaceAllString(match[2], ""))
	patch, _ = strconv.Atoi(regexp.MustCompile(`\D`).ReplaceAllString(match[3], ""))
	return fmt.Sprintf("Incrementing: %v.%v.%v", major, minor, patch), nil

}

func main() {
	//variables
	repositoryName := os.Getenv("INPUT_ECR_NAME")
	versionType := os.Getenv("INPUT_VERSION_TYPE")

	//repositoryName := "activity-api"
	//versionType := "patch"

	slice := []string{}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create ECR Service Client
	svc := ecr.New(sess)

	result, err := svc.ListImages(&ecr.ListImagesInput{
		RepositoryName: aws.String(repositoryName),
	})
	fmt.Print(result)

	// if we have an error print it and exit the program
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecr.ErrCodeServerException:
				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
			case ecr.ErrCodeInvalidParameterException:
				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
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
	//loop through the images and append the version number to a slice
	for _, image := range result.ImageIds {

		if *image.ImageTag == "<untagged>" {
			continue
		}

		slice = appendString(slice, *image.ImageTag)
	}
	// if lenght slice == 0 then exit the program
	if len(slice) == 0 {
		fmt.Println("No images found, creating a 0.0.1 version")
		major = 0
		minor = 0
		patch = 1
	} else {
		//sort the slice
		sort.Strings(slice)
		//parse the version number
		fmt.Println(parseVn(slice[len(slice)-1]))
		//increment the version number
		switch versionType {
		case "patch":
			incrementPatch()
		case "minor":
			incrementMinor()
		case "major":
			incrementMajor()
		}

	}

	// create the new version number
	//newVersion := fmt.Sprintf("%v.%v.%v", major, minor, patch)

	// if lenght imageIds is 0 then exit
	// if len(result.ImageIds) == 0 {
	// fmt.Println("No images found")
	// return
	// }

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, fmt.Sprintf("%v.%v.%v", major, minor, patch)))

}
