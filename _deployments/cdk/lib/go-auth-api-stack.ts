import * as cdk from 'aws-cdk-lib';
import { Stack, Tags } from 'aws-cdk-lib';
import {
  CfnInternetGateway,
  CfnVPCGatewayAttachment,
  IpAddresses,
  Peer,
  Port,
  SecurityGroup,
  Subnet,
  Vpc,
} from 'aws-cdk-lib/aws-ec2';
import { Role } from 'aws-cdk-lib/aws-iam';
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
     * Role
     */
    const executionRole = Role.fromRoleArn(
      this,
      `ecsTaskExecutionRole`,
      `arn:aws:iam::${Stack.of(this).account}:role/ecsTaskExecutionRole`
    )

    /**
     * VPC
     */
    const vpc = new Vpc(this, `${context.name}-vpc`, {
      vpcName: `${context.name}-vpc`,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      natGateways: 0,
      subnetConfiguration: [],
      ipAddresses: IpAddresses.cidr("10.0.0.0/16")
    });

    /**
     * Security Group
     */
    const albSg = new SecurityGroup(this, `${context.name}-alb-sg`, {
      vpc,
      allowAllOutbound: true,
      securityGroupName: `${context.name}-alb-sg`,
    });
    albSg.addIngressRule(Peer.anyIpv4(), Port.tcp(80));
    const apiServiceSg = new SecurityGroup(this, `${context.name}-api-service-sg`, {
      vpc,
      allowAllOutbound: true,
      securityGroupName: `${context.name}-api-server-sg`,
    });
    apiServiceSg.connections.allowFrom(albSg, Port.tcp(3000), 'Allow alb access')

    /**
     * Gateway
     */
    const igw = new CfnInternetGateway(this, `${context.name}-igw`, {
      tags: [
        {
          key: "Name",
          value: `${context.name}-igw`,
        },
      ],
    })
    const igwAttachment = new CfnVPCGatewayAttachment(this, `${context.name}-igw-attachment`, {
      vpcId: vpc.vpcId,
      internetGatewayId: igw.ref,
    })

    /**
     * Subnet
     */
    const appPublicSubnet = new Subnet(this, `${context.name}-app-public-subnet`, {
      vpcId: vpc.vpcId,
      cidrBlock: "10.0.0.0/24",
      availabilityZone: "ap-northeast-1a",
      mapPublicIpOnLaunch: true,
    })
    appPublicSubnet.addDefaultInternetRoute(igw.ref, igwAttachment)
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
