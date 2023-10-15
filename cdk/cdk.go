package main

import (
	"cdk/helpers"
	"fmt"
	"log"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func createDynamoDBTable(stack awscdk.Stack, id string, tableName string, attributeName string, attributeType awsdynamodb.AttributeType) awsdynamodb.Table {
	return awsdynamodb.NewTable(stack, jsii.String(id), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String(attributeName),
			Type: attributeType,
		},
		TableName:   jsii.String(tableName),
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
		// If you have other common configurations like billing mode,
		// read/write capacity, etc., you can add them here.
	})
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	// ------ Define DynamoDB tables -------- //
	// table1 := awsdynamodb.NewTable(stack, jsii.String("MetricsTable"), &awsdynamodb.TableProps{
	// 	PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("ID"), Type: awsdynamodb.AttributeType_STRING},
	// 	TableName:    jsii.String("Table1"),
	// })

	table1 := createDynamoDBTable(stack, "MetricsTable", "UserMatrix", "id", awsdynamodb.AttributeType_STRING)
	table2 := createDynamoDBTable(stack, "ContributorsTable", "Contributors", "id", awsdynamodb.AttributeType_NUMBER)
	table3 := createDynamoDBTable(stack, "CommitsTable", "Github_Commits", "id", awsdynamodb.AttributeType_STRING)
	table4 := createDynamoDBTable(stack, "CommentsTable", "Github_IssueComments", "id", awsdynamodb.AttributeType_NUMBER)
	table5 := createDynamoDBTable(stack, "PullsTable", "PullRequests", "pr_id", awsdynamodb.AttributeType_NUMBER)

	mastersRole := awsiam.NewRole(stack, jsii.String("MastersRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewArnPrincipal(jsii.String("arn:aws:iam::727433422324:user/dynamo_access")),
	})

	fargateRole := awsiam.NewRole(stack, jsii.String("CustomFargateRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks-fargate-pods.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
		RoleName:  jsii.String("CustomFargateRole"),
	})

	// Granting permissions to the role to access DynamoDB
	fargateRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: jsii.Strings(
			"dynamodb:PutItem",
			"dynamodb:GetItem",
			"dynamodb:UpdateItem",
			"dynamodb:Query",
			"dynamodb:Scan",
		),
		Resources: jsii.Strings(
			*table1.TableArn(),
			*table2.TableArn(),
			*table3.TableArn(),
			*table4.TableArn(),
			*table5.TableArn(),

			// Add table2.TableArn(), table3.TableArn() etc. as you uncomment and create those tables
		),
	}))

	// Granting permissions to the role to access ECR
	fargateRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions: jsii.Strings(
			"ecr:GetAuthorizationToken",
			"ecr:GetDownloadUrlForLayer",
			"ecr:BatchGetImage",
			"ecr:BatchCheckLayerAvailability",
		),
		Resources: jsii.Strings("*"),
	}))

	// ------ EKS Cluster -------- //

	// Note: This is a basic setup for Fargate cluster. Adjust it according to your needs.
	cluster := awseks.NewFargateCluster(stack, jsii.String("UserEKSCluster"), &awseks.FargateClusterProps{
		Version:     awseks.KubernetesVersion_V1_27(), // Version:     awseks.KubernetesVersion_V1_26, // Adjust the version if needed
		MastersRole: mastersRole,
		// More properties can be added here as needed
	})

	// Add Fargate profile
	// Note: This function call might differ in the actual Go CDK API, as the API might evolve.
	cluster.AddFargateProfile(jsii.String("UserDevFargateProfile"), &awseks.FargateProfileOptions{
		Selectors: &[]*awseks.Selector{{
			Namespace: jsii.String("default"),
		}},
		PodExecutionRole: fargateRole,
	})

	// Create the ALB
	// alb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("MyALB"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
	// 	Vpc:            cluster.Vpc(),
	// 	InternetFacing: jsii.Bool(true), // Set to false if the ALB should be internal
	// })

	// Define a target group for your ALB
	// targetGroup := alb.AddListener(jsii.String("MyHTTPListener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
	// 	Port: jsii.Int(80), // The port on which the ALB listens
	// }).AddTargets(jsii.String("FargateTargets"), &awselasticloadbalancingv2.AddApplicationTargetsProps{
	// 	Port:    jsii.Int(80),                                                         // The port on which the target group forwards traffic
	// 	Targets: &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{cluster}, // Use your Fargate service as a target
	// })

	// awseks.KubernetesManifest()
	// Define the folder path where your Kubernetes manifest files are located.
	manifestsFolder := "./cofigs"

	manifests, err := helpers.LoadKubernetesManifests(manifestsFolder)
	if err != nil {
		log.Printf("Failed retrieving manifest file: %v\n", err)
		return stack // Or you can handle the error accordingly
	}

	var manifestList []*map[string]interface{}

	for i, manifest := range manifests {
		identifier := fmt.Sprintf("Manifest%d", i)
		manifestList = append(manifestList, &manifest)
		log.Printf("Creating KubernetesManifest resource: %s\n", identifier)
	}

	awseks.NewKubernetesManifest(stack, jsii.String("AllManifests"), &awseks.KubernetesManifestProps{
		Cluster:  cluster,
		Manifest: &manifestList,
	})

	// Load Kubernetes manifests from the folder.
	// for i, manifest := range manifests {
	// 	identifier := fmt.Sprintf("Manifest%d", i)

	// 	manifestList := []*map[string]interface{}{&manifest}

	// 	awseks.NewKubernetesManifest(stack, jsii.String(identifier), &awseks.KubernetesManifestProps{
	// 		Cluster:  cluster,
	// 		Manifest: &manifestList,
	// 	})
	// }

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "CdkStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
