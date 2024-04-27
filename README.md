# FrontEnd Masters course on Go & Aws by Melkey Dev

## Tech Stack

- Go
- AWS
- CDK

## CDK Usage

This repo was created starting from a CDK template.

The `cdk.json` file tells the CDK toolkit how to execute your app.

### Useful Commands:

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `go test` run unit tests

## Project Structure

In the application root you can find our main app that tells CDK what to deploy on AWS

In the `/lambda` folder you can find the lambda function that will act as a REST API Server with three endpoints:

- `/register`: where you can register a user on `DynamoDB`, an AWS fully managed NoSQL Database
- `/login`: where you can login and make the lambda assign you a JWT token that will identifiy the user in input
- `/protected`: where you can send the retrieved JWT Token and can check if you have the rights to have a successful response from this endpoint

Big thanks to [FrontEnd Masters](https://frontendmasters.com/) and [MelkeyDev](https://github.com/Melkeydev/) for this amazing course.
