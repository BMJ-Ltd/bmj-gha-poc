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

	repositoryName := os.Getenv("INPUT_ECR_NAME")
	versionType := os.Getenv("INPUT_VERSION_TYPE")
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

				slice := []string{}

				// iterate over the images and print them out
				for _, image := range result.ImageIds {
					slice = appendString(slice, *image.ImageTag)
					fmt.Println(*image.ImageTag)
				}
				sort.Strings(slice)
				fmt.Println(parseVn(slice[len(slice)-1]))

				// do required increment
				if versionType == "major" {
					incrementMajor()
				}
				if versionType == "minor" {
					incrementMinor()

				}
				if versionType == "patch" {
					incrementPatch()
				}

			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(fmt.Sprintf(`::set-output name=myOutput::%s`, fmt.Sprintf("%v.%v.%v", major, minor, patch)))

}
