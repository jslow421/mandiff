import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const documentUploadBucket = new cdk.aws_s3.Bucket(
      this,
      "MedManualUploadBucket",
      {
        removalPolicy: cdk.RemovalPolicy.DESTROY,
        autoDeleteObjects: true,
        lifecycleRules: [
          {
            expiration: cdk.Duration.days(2),
          },
        ],
        bucketName: "med-manual-upload-bucket",
      }
    );

    const processingCompleteBucket = new cdk.aws_s3.Bucket(
      this,
      "MedManualCompleteBucket",
      {
        removalPolicy: cdk.RemovalPolicy.DESTROY,
        autoDeleteObjects: true,
        lifecycleRules: [
          {
            expiration: cdk.Duration.days(30),
          },
        ],
        bucketName: "med-manual-complete-bucket",
      }
    );
    //   this,
    //   "ExtractedEnglishBucket",
    //   {
    //     removalPolicy: cdk.RemovalPolicy.DESTROY,
    //     autoDeleteObjects: true,
    //     lifecycleRules: [
    //       {
    //         expiration: cdk.Duration.days(30),
    //       },
    //     ],
    //     bucketName: "med-manual-extracted-english-bucket",
    //   }
    // );

    const functionRole = new cdk.aws_iam.Role(this, "BucketRole", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaBasicExecutionRole"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "CloudWatchFullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonS3FullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonTextractFullAccess"
        ),
      ],
    });

    functionRole.addToPolicy(
      new cdk.aws_iam.PolicyStatement({
        actions: ["textract:StartDocumentAnalysis"],
        resources: ["*"],
      })
    );

    const completeNotificationQueue = new cdk.aws_sns.Topic(
      this,
      "CompleteNotificationQueue",
      {
        topicName: "manual-text-extraction-complete-notification-queue",
      }
    );

    functionRole.addToPolicy(
      new cdk.aws_iam.PolicyStatement({
        actions: ["sns:Publish"],
        resources: [completeNotificationQueue.topicArn],
      })
    );
    //   new cdk.aws_iam.PolicyStatement({
    //     actions: [
    //       "comprehend:DetectDominantLanguage",
    //       "comprehend:DetectEntities",
    //     ],
    //     resources: ["*"],
    //   })
    // );

    const snsRole = new cdk.aws_iam.Role(this, "SNSRole", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "CloudWatchFullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonS3FullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonTextractFullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonSQSFullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonSNSFullAccess"
        ),
      ],
    });

    const bedrockLambdaRole = new cdk.aws_iam.Role(this, "BedrockLambdaRole", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaBasicExecutionRole"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "CloudWatchFullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonS3FullAccess"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "AmazonBedrockFullAccess"
        ),
      ],
    });

    completeNotificationQueue.grantPublish(functionRole);

    const processDocumentsFunction = new cdk.aws_lambda.Function(
      this,
      "ProcessDocuments",
      {
        runtime: cdk.aws_lambda.Runtime.PROVIDED_AL2023,
        architecture: cdk.aws_lambda.Architecture.ARM_64,
        handler: "bootstrap",
        code: cdk.aws_lambda.Code.fromAsset("../out/lambda/process-documents"),
        environment: {
          UPLOAD_BUCKET: documentUploadBucket.bucketName,
          COMPLETE_BUCKET: processingCompleteBucket.bucketName,
          SNS_TOPIC_ARN: completeNotificationQueue.topicArn,
          ROLE_ARN: snsRole.roleArn,
        },
        role: functionRole,
        timeout: cdk.Duration.minutes(10),
      }
    );

    const createTextFileFunction = new cdk.aws_lambda.Function(
      this,
      "CreateTextFile",
      {
        runtime: cdk.aws_lambda.Runtime.PROVIDED_AL2023,
        architecture: cdk.aws_lambda.Architecture.ARM_64,
        memorySize: 1024,
        handler: "bootstrap",
        code: cdk.aws_lambda.Code.fromAsset("../out/lambda/convert-files"),
        environment: {
          COMPLETE_BUCKET: processingCompleteBucket.bucketName,
        },
        role: functionRole,
        timeout: cdk.Duration.minutes(10),
      }
    );

    const llmFunction = new cdk.aws_lambda.Function(this, "LlmFunction", {
      runtime: cdk.aws_lambda.Runtime.PROVIDED_AL2023,
      architecture: cdk.aws_lambda.Architecture.ARM_64,
      memorySize: 1024,
      handler: "bootstrap",
      code: cdk.aws_lambda.Code.fromAsset("../out/lambda/llm"),
      environment: {
        COMPLETE_BUCKET: processingCompleteBucket.bucketName,
      },
      role: bedrockLambdaRole,
      timeout: cdk.Duration.minutes(3),
    });
    //   this,
    //   "ExtractEnglishLanguage",
    //   {
    //     runtime: cdk.aws_lambda.Runtime.PROVIDED_AL2023,
    //     architecture: cdk.aws_lambda.Architecture.ARM_64,
    //     handler: "bootstrap",
    //     code: cdk.aws_lambda.Code.fromAsset("../out/lambda/extract-english"),
    //     environment: {
    //       SOURCE_BUCKET: processingCompleteBucket.bucketName,
    //       TARGET_BUCKET: extractedEnglishBucket.bucketName,
    //     },
    //     role: functionRole,
    //   }
    // );

    documentUploadBucket.grantRead(processDocumentsFunction);
    processingCompleteBucket.grantWrite(processDocumentsFunction);
    processingCompleteBucket.grantReadWrite(createTextFileFunction);
    processingCompleteBucket.grantReadWrite(llmFunction);
  }
}
