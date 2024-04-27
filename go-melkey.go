package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GoMelkeyStackProps struct {
	awscdk.StackProps
}

func NewGoMelkeyStack(scope constructs.Construct, id string, props *GoMelkeyStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	// New Dynamo DB Table
	table := awsdynamodb.NewTable(stack, jsii.String("Go Melkey FEM User Table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"), // Based the logic on `username`
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("user_table"),
	})

	// New Lambda
	backendFunction := awslambda.NewFunction(stack, jsii.String("Go Melkey FEM Lambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
	})

	// Grant the lambda function access to the created table
	table.GrantReadWriteData(backendFunction)

	apiGateway := awsapigateway.NewRestApi(stack, jsii.String("Go Melkey FEM Api Gateway"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "OPTIONS"),
			AllowOrigins: jsii.Strings("localhost"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})

	integration := awsapigateway.NewLambdaIntegration(backendFunction, nil)

	// Define routes
	// Register
	apiGateway.Root().AddResource(jsii.String("register"), nil).AddMethod(jsii.String("POST"), integration, nil)

	// Login
	apiGateway.Root().AddResource(jsii.String("login"), nil).AddMethod(jsii.String("POST"), integration, nil)

	// Protected resource
	apiGateway.Root().AddResource(jsii.String("protected"), nil).AddMethod(jsii.String("GET"), integration, nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoMelkeyStack(app, "GoMelkeyStack", &GoMelkeyStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
