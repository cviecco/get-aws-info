package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/yaml.v2"
)

/*
credentials:
    - account_id: "6865596DDDD"
      access_key_id: "AKIAIQCPR"
      secret_access_key: "AKIAIQC"
*/

type awsCredentialsConfig struct {
	Account_id        string
	Access_Key_ID     string
	Secret_Access_Key string
}

type AWSConfigFile struct {
	Credentials []awsCredentialsConfig
}

var (
	Version       = "No version provided"
	debug         = flag.Bool("debug", false, "Enable debug mode.")
	timeoutString = flag.String("httptimeout", "4s", "Duration of the http client timeouts")
)

func getRegionString() (string, error) {
	svc := ec2metadata.New(session.New())
	region, err := svc.Region()
	return region, err
}

func getAccountIdInstanceId() (string, string, error) {
	svc := ec2metadata.New(session.New())
	identityDocument, err := svc.GetInstanceIdentityDocument()
	if err != nil {
		return "", "", err
	}
	return identityDocument.AccountID, identityDocument.InstanceID, nil
}

func getTag(requestedTagName string, regionString string, credentials *credentials.Credentials) (string, error) {
	_, instanceID, err := getAccountIdInstanceId()
	if err != nil {
		return "", err
	}

	//timeout := time.Duration(5 * time.Second)
	timeout, err := time.ParseDuration(*timeoutString)
	if err != nil {
		return "", err
	}

	svc := ec2.New(session.New(&aws.Config{Region: aws.String(regionString),
		Credentials: credentials,
		HTTPClient:  &http.Client{Timeout: timeout}}))

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

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s (version %s):\n", os.Args[0], Version)
	flag.PrintDefaults()
}

func main() {
	var configFilename = flag.String("config", "credentials.yml", "credentials file")
	var tagName = flag.String("tagname", "Name", "Name of the tag of interest")
	flag.Usage = Usage
	flag.Parse()

	var config AWSConfigFile
	if *debug {
		log.Printf("using config=%s\n", *configFilename)
	}
	source, err := ioutil.ReadFile(*configFilename)
	if err != nil {
		log.Printf("Cannot read config file: %s. Err=%s\n", *configFilename, err.Error())
		os.Exit(1)
	}
	err = yaml.Unmarshal(source, &config)
	if err != nil {
		log.Printf("Cannot parse config file: %s. Err=%s\n", *configFilename, err.Error())
		os.Exit(1)
	}
	if *debug {
		fmt.Printf("%+v\n", config)
	}

	var creds *credentials.Credentials

	accountID, _, err := getAccountIdInstanceId()
	if err != nil {
		log.Fatalf("Cannot get accountID, %s", err)
	}

	for _, account := range config.Credentials {
		if account.Account_id == accountID {
			if *debug {
				log.Printf("%+v\n", account)
			}
			creds = credentials.NewStaticCredentials(account.Access_Key_ID, account.Secret_Access_Key, "")
		}
	}

	regionString, err := getRegionString()
	if err != nil {
		log.Fatalf("Cannot get Region string %s", err)
	}

	value, err := getTag(*tagName, regionString, creds)
	if err != nil {
		log.Fatalf("Cannot get Tag Name %s", err)
	}
	fmt.Printf("%s\n", value)

}
