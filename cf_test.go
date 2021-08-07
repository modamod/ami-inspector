package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

var SESS, _ = session.NewSession(&aws.Config{
	Region: aws.String("us-east-1"), Endpoint: aws.String("http://localhost:5000")},
)

func TestValidateNonExistingTemplate(t *testing.T) {
	stack := CloudformationTemplate{TemplatePath: "nonExistingPath", StackName: "NonExisting", TemplateParameterPath: "NonExisting", AwsSession: SESS}
	isValid, _ := stack.ValidateTemplate()
	assert.False(t, isValid)
}

func TestValidateInvalidTemplate(t *testing.T) {
	stack := CloudformationTemplate{TemplatePath: "./fixtures/templates/invalid_template.yaml", StackName: "", TemplateParameterPath: "", AwsSession: SESS}
	isValid, _ := stack.ValidateTemplate()
	assert.False(t, isValid)
}

func TestValidateTemplateOnValidTemplate(t *testing.T) {
	stack := CloudformationTemplate{TemplatePath: "./fixtures/templates/valid_template.yaml", StackName: "", TemplateParameterPath: "", AwsSession: SESS}
	isValid, _ := stack.ValidateTemplate()
	assert.True(t, isValid)
}

func TestExistsWithNonExistingStack(t *testing.T) {
	stack := CloudformationTemplate{TemplatePath: "./fixtures/templates/valid_template.yaml", StackName: "NonExisting", TemplateParameterPath: "./fixtures/parameters/parameters.yaml", AwsSession: SESS}
	assert.False(t, stack.Exists())
}

func TestExistsWithExistingStack(t *testing.T) {
	stackParameters := make(map[string]string)
	stackParameters["TopicName"] = "MyTopic"
	stack := CloudformationTemplate{
		TemplatePath:          "./fixtures/templates/valid_template.yaml",
		StackName:             "Existing",
		TemplateParameterPath: "./fixtures/parameters/parameters.yaml",
		AwsSession:            SESS,
		Parameters:            stackParameters,
		Timeout:               1,
		Capabilities:          &[]string{},
		DisableRollback:       false,
	}

	stack.CreateStack()
	assert.True(t, stack.Exists())
}

func TestDeleteStackWithNonExistingStack(t *testing.T) {
	stack := CloudformationTemplate{
		TemplatePath:          "./fixtures/templates/valid_template.yaml",
		StackName:             "NonExisting",
		TemplateParameterPath: "./fixtures/parameters/parameters.yaml",
		AwsSession:            SESS,
		Parameters:            map[string]string{},
		Timeout:               1,
		Capabilities:          &[]string{},
		DisableRollback:       false,
	}
	stack.DeleteStack()
	assert.False(t, stack.Exists())
}

func TestDeleteStackWithExistingStack(t *testing.T) {
	stackParameters := make(map[string]string)
	stackParameters["TopicName"] = "MyTopic"
	stack := CloudformationTemplate{
		TemplatePath:          "./fixtures/templates/valid_template.yaml",
		StackName:             "Existing",
		TemplateParameterPath: "./fixtures/parameters/parameters.yaml",
		AwsSession:            SESS,
		Parameters:            stackParameters,
		Timeout:               1,
		Capabilities:          &[]string{},
		DisableRollback:       false,
	}
	// To ensure that each test is decoupled from any other test creating a stack is required
	// Used the same stack name that I created earlier just to make sure that no new resources
	// are created if they already exist.
	stack.CreateStack()
	assert.True(t, stack.Exists())
	stack.DeleteStack()
	assert.False(t, stack.Exists())
}
