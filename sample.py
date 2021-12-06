import boto3
import time


client = boto3.Session(region_name="us-east-1").client("cloudformation", endpoint_url='http://localhost:5000')
paramters = [{
    "ParameterKey": "TopicName",
    "ParameterValue": "ParameterValue"
}]
template_body = open("fixtures/templates/valid_template.yaml").read()
stack_name = "TestStack"
# client.create_stack(StackName=stack_name, TemplateBody=template_body, Parameters=paramters)
# time.sleep(1)

# Should fail since I am updating the stack with the same paramters and template body
result = client.update_stack(StackName=stack_name, TemplateBody=template_body, Parameters=paramters)
print(result)
