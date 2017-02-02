package main

import (
	"flag"
	"fmt"
	"log"
	//"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var debug = false

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

func getTag(name string, regionString string) (string, error) {
	//svc := ec2.New(mySession, aws.NewConfig().WithRegion("us-west-2"))
	instanceID, err := getInstanceId()
	if err != nil {
		return "", err
	}

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(regionString)}))

	//svc := ec2.New(sess)

	params := &ec2.DescribeInstancesInput{
		//DryRun: aws.Bool(true),
		Filters: []*ec2.Filter{
			{ /*
				// Required
				Name: aws.String("String"),
				Values: []*string{
					aws.String("String"), // Required
					// More values...
				},
			*/
			},
			// More values...
		},
		InstanceIds: []*string{
			//aws.String("i-0c8b43c46be4852d5"), // Required
			aws.String(instanceID),
			// More values...
		},
		//MaxResults: aws.Int64(1),
		//NextToken: aws.String("String"),
	}
	resp, err := svc.DescribeInstances(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return "", err
	}

	// Pretty-print the response data.
	fmt.Println(resp)

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
	getTag(*tagName, regionString)

}