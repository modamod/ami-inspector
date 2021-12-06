package main

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"gopkg.in/yaml.v3"
)

//CloudformationStack is a class describing the required attributes and methods to manage a cloudformation template.
type CloudformationStack struct {
	TemplatePath          string
	TemplateParameterPath string
	StackName             string
	AwsSession            *session.Session
	Capabilities          *[]string
	DisableRollback       bool
	Parameters            map[string]string
	Timeout               int64
}

func MakeCloudformationStack() CloudformationStack {
	cf := CloudformationStack{
		Capabilities: new([]string),
		Parameters:   make(map[string]string),
	}

	return cf
}

//GetTemplateBody returns a reference to a string containing the template body.
func (cf *CloudformationStack) GetTemplateBody() *string {
	templateBody, err := ioutil.ReadFile(cf.TemplatePath)
	if err != nil {
		printErrorf("Something went wrong when trying to open template body. %v", err)
		return nil
	}
	return aws.String(string(templateBody))

}

//ValidateTemplate validates a CF template returns false when the validations fails.
func (cf *CloudformationStack) ValidateTemplate() (isValid bool, err error) {
	templateBody, err := ioutil.ReadFile(cf.TemplatePath)
	if err != nil {
		printErrorf("Something went wrong when trying to open template body. %v", err)
		return false, err
	}

	_, err = cf.Client().ValidateTemplate(&cloudformation.ValidateTemplateInput{TemplateBody: aws.String(string(templateBody))})
	if err != nil {
		return false, err
	}
	return true, err
}

//Client is a function that returns a cloudformation client object
func (cf *CloudformationStack) Client() *cloudformation.CloudFormation {
	return cloudformation.New(cf.AwsSession)
}

//Exists a function that tests if the cloudformation exists or not.
func (cf *CloudformationStack) Exists() bool {
	stack, _ := cf.Client().DescribeStacks(&cloudformation.DescribeStacksInput{StackName: aws.String(cf.StackName)})
	return stack.Stacks != nil
}

func (cf *CloudformationStack) GetParameters() (map[string]string, error) {
	params := make(map[string]string)
	parametersFile, err := ioutil.ReadFile(cf.TemplateParameterPath)
	if err != nil {
		printErrorf("Something went wrong when trying to open parameters file. %v", err)
		return nil, err
	}
	err = yaml.Unmarshal(parametersFile, &params)
	if err != nil {
		printErrorf("Something went wrong when loading yaml content from paramters file. %v", err)
		return nil, err
	}
	return params, nil
}
func (cf *CloudformationStack) GenerateParameters() ([]*cloudformation.Parameter, error) {
	var Parameters []*cloudformation.Parameter
	parametersMap, err := cf.GetParameters()
	for key, element := range parametersMap {
		Parameters = append(Parameters, &cloudformation.Parameter{ParameterKey: aws.String(key), ParameterValue: aws.String(element)})
	}
	return Parameters, err
}

func (cf *CloudformationStack) getCreateStackInput() (*cloudformation.CreateStackInput, error) {
	parameters, err := cf.GenerateParameters()
	stackInput := cloudformation.CreateStackInput{
		Capabilities:     aws.StringSlice(*cf.Capabilities),
		DisableRollback:  aws.Bool(cf.DisableRollback),
		Parameters:       parameters,
		StackName:        aws.String(cf.StackName),
		TemplateBody:     cf.GetTemplateBody(),
		TimeoutInMinutes: aws.Int64(cf.Timeout),
	}

	return &stackInput, err
}

func (cf *CloudformationStack) CreateStack() (*cloudformation.CreateStackOutput, error) {
	stackInput, err := cf.getCreateStackInput()
	if err != nil {
		return nil, err
	}

	res, err := cf.Client().CreateStack(stackInput)
	return res, err
}

func (cf *CloudformationStack) DeleteStack() (*cloudformation.DeleteStackOutput, error) {
	deleteStackInput := cloudformation.DeleteStackInput{StackName: aws.String(cf.StackName)}
	res, err := cf.Client().DeleteStack(&deleteStackInput)
	return res, err
}

func (cf *CloudformationStack) getUpdateStackInput() (*cloudformation.UpdateStackInput, error) {
	parameters, err := cf.GenerateParameters()
	updateStackInput := cloudformation.UpdateStackInput{
		Capabilities: aws.StringSlice(*cf.Capabilities),
		Parameters:   parameters,
		StackName:    aws.String(cf.StackName),
		TemplateBody: cf.GetTemplateBody(),
	}
	return &updateStackInput, err
}

func (cf *CloudformationStack) UpdateStack() (*cloudformation.UpdateStackOutput, error) {
	updateStackInput, err := cf.getUpdateStackInput()
	if err != nil {
		return nil, err
	}

	updateStackOutput, err := cf.Client().UpdateStack(updateStackInput)

	return updateStackOutput, err
}

func (cf *CloudformationStack) DescribeStack() (*cloudformation.DescribeStacksOutput, error) {
	describeStacksInput := cloudformation.DescribeStacksInput{StackName: aws.String(cf.StackName)}
	describeStacksOutpt, err := cf.Client().DescribeStacks(&describeStacksInput)

	return describeStacksOutpt, err
}

func (cf *CloudformationStack) GetStatus() (string, error) {
	describeStacksoutput, err := cf.DescribeStack()

	return *describeStacksoutput.Stacks[0].StackStatus, err
}
