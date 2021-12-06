package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultKeypairLocation = os.Getenv("DEFAULT_KEYPAIR")
)

//importKeypair is used to import existing keypair
func importKeypair() bool {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := ec2.New(sess)
	if err != nil {
		exitErrorf("Something went wrong when creating aws session, %v", err)
	}

	describeKeyPairsOutput, err := svc.DescribeKeyPairs(nil)
	if err != nil {
		exitErrorf("Unable to get key pairs, %v", err)
	}

	if len(describeKeyPairsOutput.KeyPairs) == 0 {
		fmt.Println("No keypair found, importing default keypair")
		path, err := filepath.Abs(DefaultKeypairLocation)
		if err != nil {
			exitErrorf("Keypair public key file not found, %v", err)
		}
		fmt.Println(path)
		fileContent, _ := ioutil.ReadFile(path)
		keypair := ec2.ImportKeyPairInput{
			KeyName: aws.String("DefaultKeypair"), PublicKeyMaterial: fileContent,
		}

		importKeypairOutput, err := svc.ImportKeyPair(&keypair)

		if err != nil {
			exitErrorf("Something went wrong when creating keypair, %v", err)
		}
		fmt.Println(importKeypairOutput)
	}

	describeKeyPairsOutput, err = svc.DescribeKeyPairs(nil)
	if err != nil {
		exitErrorf("Unable to get key pairs, %v", err)
	}
	fmt.Println("Key Pairs:")
	for _, pair := range describeKeyPairsOutput.KeyPairs {
		fmt.Printf("%s: %s\n", *pair.KeyName, *pair.KeyFingerprint)
	}
	return true
}

type Params struct {
	Type    string `yaml:"Type"`
	Default string `yaml:"Default"`
}

type Resource struct {
}
type Template struct {
	// Version    string            `yaml:"AWSTemplateFormatVersion"`
	Parameters map[string]Params `yaml:"Parameters"`
	// Resources  []Resource
	// Conditions []Condition
	// Outputs: []Output
}

func main() {
	var SESS, _ = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	stackParameters := make(map[string]string)
	stackParameters["TopicName"] = "MyTopic"
	stack := CloudformationStack{
		TemplatePath:          "./fixtures/templates/valid_template.yaml",
		StackName:             "Existing",
		TemplateParameterPath: "./fixtures/parameters/parameters.yaml",
		AwsSession:            SESS,
		Parameters:            stackParameters,
		Timeout:               1,
		Capabilities:          &[]string{},
		DisableRollback:       false,
	}
	if !stack.Exists() {
		stack.CreateStack()
	}

	res, err := stack.UpdateStack()
	fmt.Printf("res: %v\n err: %v", res, err)
}
