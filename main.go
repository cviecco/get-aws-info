package main

import (
	"flag"
	"fmt"
	"log"
	//"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"gopkg.in/yaml.v2"
)

var (
	debug = flag.Bool("debug", false, "Enable debug mode.")
)

func getRegionString() (string, error) {
	svc := ec2metadata.New(session.New())
	region, err := svc.Region()
	return region, err
}

func getInstanceId() (string, error) {
	svc := ec2metadata.New(session.New())
	identityDocument, err := svc.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return identityDocument.InstanceID, nil
}

func getTag(requestedTagName string, regionString string, credentials *credentials.Credentials) (string, error) {
	//svc := ec2.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
	instanceID, err := getInstanceId()
	if err != nil {
		return "", err
	}

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(regionString)}))

	describeTagsparams := &ec2.DescribeTagsInput{
		//DryRun: aws.Bool(true),
		Filters: []*ec2.Filter{
			{ // Required
				Name: aws.String("resource-type"),
				Values: []*string{
					aws.String("instance"), // Required
					// More values...
				},
			},
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID), // Required
					// More values...
				},
			},
			// More values...
		},
		MaxResults: aws.Int64(10),
		NextToken:  aws.String("String"),
	}

	resp2, err := svc.DescribeTags(describeTagsparams)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return "", err
	}

	// Pretty-print the response data.
	if *debug {
		fmt.Println(resp2)
	}
	for _, tag := range resp2.Tags {
		if *tag.Key == requestedTagName {
			if *debug {
				log.Printf("found %s: %s\n", *tag.Key, *tag.Value)
			}
			return *tag.Value, nil
		}
	}

	return "", nil
}

func main() {
	//var credentialsFilePath = flag.String("credentials", "credentials.yml", "credentials file")
	var tagName = flag.String("tagname", "Name", "Name of the tag of interest")
	flag.Parse()

	regionString, err := getRegionString()
	if err != nil {
		log.Fatalf("Cannot get Region string %s", err)
	}

	//creds := credentials.NewStatic

	value, err := getTag(*tagName, regionString, nil)
	if err != nil {
		log.Fatalf("Cannot get Tag Name %s", err)
	}
	fmt.Printf("%s\n", value)

}
