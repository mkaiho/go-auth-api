import * as cdk from 'aws-cdk-lib';
import { Tags } from 'aws-cdk-lib';
import { Subnet, Vpc } from 'aws-cdk-lib/aws-ec2';
import { Construct } from 'constructs';
// import * as sqs from 'aws-cdk-lib/aws-sqs';

interface StageContext {
  name: string;
}

export class GoAuthApiStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const env: string = this.node.tryGetContext("env");
    const context: StageContext = this.node.tryGetContext(env);
    const revision = require("child_process")
      .execSync("git rev-parse HEAD")
      .toString()
      .trim();

    /**
     * VPC
     */
    const vpc = new Vpc(this, `${context.name}-vpc`, {
      vpcName: `${context.name}-vpc`,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      natGateways: 0,
      subnetConfiguration: [],
      cidr: "10.0.0.0/16",
    });

    /**
     * Subnet
     */
    const appPublicSubnet = new Subnet(this, `${context.name}-app-public-subnet`, {
      vpcId: vpc.vpcId,
      cidrBlock: "10.0.0.0/24",
      availabilityZone: "ap-northeast-1a",
      mapPublicIpOnLaunch: false,
    });
    Tags.of(appPublicSubnet).add('Name', `${context.name}-app-public-subnet`)
    const apPrivateSubnet = new Subnet(this, `${context.name}-app-private-subnet`, {
      vpcId: vpc.vpcId,
      cidrBlock: "10.0.1.0/24",
      availabilityZone: "ap-northeast-1a",
      mapPublicIpOnLaunch: false,
    });
    Tags.of(apPrivateSubnet).add('Name', `${context.name}-app-private-subnet`)

    /**
     * ECR
     */
    const image = this.synthesizer.addDockerImageAsset({
      sourceHash: revision,
      directoryName: `${__dirname}/../../../`,
    });
  }
}
