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

    const bucketRole = new cdk.aws_iam.Role(this, "BucketRole", {
      assumedBy: new cdk.aws_iam.ServicePrincipal("lambda.amazonaws.com"),
      managedPolicies: [
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "service-role/AWSLambdaBasicExecutionRole"
        ),
        cdk.aws_iam.ManagedPolicy.fromAwsManagedPolicyName(
          "CloudWatchFullAccess"
        ),
      ],
    });

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
        },
        role: bucketRole,
      }
    );

    documentUploadBucket.grantRead(processDocumentsFunction);
    processingCompleteBucket.grantWrite(processDocumentsFunction);
  }
}
