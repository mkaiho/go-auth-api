import * as cdk from 'aws-cdk-lib';
import { Stack, Tags } from 'aws-cdk-lib';
import {
  CfnInternetGateway,
  CfnVPCGatewayAttachment,
  GatewayVpcEndpoint,
  GatewayVpcEndpointAwsService,
  InterfaceVpcEndpoint,
  InterfaceVpcEndpointAwsService,
  IpAddresses,
  Peer,
  Port,
  SecurityGroup,
  Subnet,
  Vpc,
} from 'aws-cdk-lib/aws-ec2';
import { Repository } from 'aws-cdk-lib/aws-ecr';
import { AppProtocol, Cluster, Compatibility, ContainerImage, TaskDefinition, Protocol, FargateService, AwsLogDriver } from 'aws-cdk-lib/aws-ecs';
import { Role } from 'aws-cdk-lib/aws-iam';
import { LogGroup, RetentionDays } from 'aws-cdk-lib/aws-logs';
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
    const apiServiceSg = new SecurityGroup(this, `${context.name}-service-sg`, {
      vpc,
      allowAllOutbound: true,
      securityGroupName: `${context.name}-service-sg`,
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
     * VPC Endpoint
     */
    const ecrDockerVpce = new InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-ecr-dkr`,
      {
        service: InterfaceVpcEndpointAwsService.ECR_DOCKER,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        securityGroups: [apiServiceSg],
        subnets: {
          subnets: [apPrivateSubnet],
        },
      }
    );
    const ecrApiVpce = new InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-ecr-api`,
      {
        service: InterfaceVpcEndpointAwsService.ECR,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        securityGroups: [apiServiceSg],
        subnets: {
          subnets: [apPrivateSubnet],
        },
      }
    );
    const logsVpce = new InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-logs`,
      {
        service: InterfaceVpcEndpointAwsService.CLOUDWATCH_LOGS,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        securityGroups: [apiServiceSg],
        subnets: {
          subnets: [apPrivateSubnet],
        },
      }
    );
    const s3Vpce = new GatewayVpcEndpoint(this, `${context.name}-vpce-s3`, {
      service: GatewayVpcEndpointAwsService.S3,
      vpc: vpc,
      subnets: [
        {
          subnets: [apPrivateSubnet],
        },
      ],
    });

    /**
     * ECR
     */
    const image = this.synthesizer.addDockerImageAsset({
      sourceHash: revision,
      directoryName: `${__dirname}/../../../`,
    });

    /**
     * Log group
     */
    const logGroup = new LogGroup(this, `${context.name}-log`, {
      logGroupName: '/aws/cdk/ecs-alb-fargate-service/web',
      retention: RetentionDays.ONE_DAY,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    })

    /**
     * ECS
     */
    const cluster = new Cluster(this, `${context.name}-cluster`, {
      clusterName: `${context.name}-cluster`,
      vpc: vpc,
    })
    const taskDefinition = new TaskDefinition(this, `${context.name}-task`, {
      compatibility: Compatibility.FARGATE,
      cpu: "256",
      memoryMiB: "512",
      family: `${context.name}-task`,
      executionRole: executionRole,
    })
    taskDefinition.addContainer(`${context.name}-container`, {
      containerName: `${context.name}`,
      image: ContainerImage.fromEcrRepository(
        Repository.fromRepositoryName(this, `${context.name}-repo`, image.repositoryName), revision
      ),
      command: ["cmd/auth-api-server"],
      portMappings: [
        {
          name: "http-port-mapping",
          hostPort: 3000,
          containerPort: 3000,
          appProtocol: AppProtocol.http,
          protocol: Protocol.TCP,
        },
      ],
      logging: new AwsLogDriver({
        streamPrefix: 'ecs',
        logGroup: logGroup,
      }),
    })
    const service = new FargateService(this, `${context.name}-service`, {
      serviceName: `${context.name}-service`,
      cluster,
      taskDefinition,
      vpcSubnets: {
        subnets: [
          apPrivateSubnet,
        ],
      },
      securityGroups: [
        apiServiceSg,
      ],
    })
  }
}
