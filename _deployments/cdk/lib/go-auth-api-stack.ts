import * as cdk from 'aws-cdk-lib';
import { SecretValue, Stack, Tags } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import {
  aws_ec2 as ec2,
  aws_ecr as ecr,
  aws_ecs as ecs,
  aws_elasticloadbalancingv2 as elb,
  aws_iam as iam,
  aws_logs as logs,
  aws_route53 as route53,
  aws_route53_targets as targets,
  aws_ssm as ssm,
  aws_rds as rds,
} from 'aws-cdk-lib';
import { SecurityGroup } from 'aws-cdk-lib/aws-ec2';

interface dns {
  zoneID: string
  zoneName: string
  domain: string
}
interface certificate {
  ref: string
}
interface loadBalancer {
  listener: {
    certificate: certificate
  }
}
interface StageContext {
  name: string
  dns: dns
  loadBalancer: loadBalancer
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
     * Get paramters from SSM
     */
    const dbPort = ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/db/port`) as unknown as number
    const dbUser = ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/db/user`)
    const dbName = ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/db/database`)

    /**
     * Role
     */
    const executionRole = new iam.Role(this, `ecsTaskExecutionRole`, {
      roleName: `${context.name}-ecsTaskExecutionRole`,
      description: "ECS execution role",
      assumedBy: new iam.ServicePrincipal('ecs-tasks.amazonaws.com'),
      managedPolicies: [
        iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AmazonECSTaskExecutionRolePolicy"),
      ],
    })
    executionRole.addToPolicy(new iam.PolicyStatement({
      effect: iam.Effect.ALLOW,
      actions: [
        "ssm:GetParameters",
        "secretsmanager:GetSecretValue",
        "kms:Decrypt",
      ],
      resources: [
        `arn:aws:ssm:${this.region}:${this.account}:parameter/*`,
        `arn:aws:secretsmanager:${this.region}:${this.account}:secret:*`,
        `arn:aws:kms:${this.region}:${this.account}:key/*`,
      ],
    }))

    /**
     * VPC
     */
    const vpc = new ec2.Vpc(this, `${context.name}-vpc`, {
      vpcName: `${context.name}-vpc`,
      enableDnsHostnames: true,
      enableDnsSupport: true,
      natGateways: 0,
      subnetConfiguration: [],
      ipAddresses: ec2.IpAddresses.cidr("10.0.0.0/16"),
    });

    /**
     * Security Group
     */
    const albSg = new ec2.SecurityGroup(this, `${context.name}-alb-sg`, {
      vpc,
      allowAllOutbound: true,
      securityGroupName: `${context.name}-alb-sg`,
    });
    albSg.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(80));
    const apiServiceSg = new ec2.SecurityGroup(this, `${context.name}-service-sg`, {
      vpc,
      allowAllOutbound: true,
      securityGroupName: `${context.name}-service-sg`,
    });
    apiServiceSg.connections.allowFrom(albSg, ec2.Port.tcp(3000), 'Allow alb access')
    const bastionSg = new SecurityGroup(this, `${context.name}-bastion-sg`, {
      vpc,
      securityGroupName: `${context.name}-bastion-sg`,
      description: "security group for bastion",
    })
    bastionSg.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(22));
    const dbSg = new SecurityGroup(this, `${context.name}-db-sg`, {
      vpc,
      securityGroupName: `${context.name}-db-sg`,
      description: "security group for auth db",
    })
    dbSg.connections.allowFrom(apiServiceSg, ec2.Port.tcp(dbPort), 'Allow db access from api')
    dbSg.connections.allowFrom(bastionSg, ec2.Port.tcp(dbPort), 'Allow db access from bastion')

    /**
     * Gateway
     */
    const igw = new ec2.CfnInternetGateway(this, `${context.name}-igw`, {
      tags: [
        {
          key: "Name",
          value: `${context.name}-igw`,
        },
      ],
    })
    const igwAttachment = new ec2.CfnVPCGatewayAttachment(this, `${context.name}-igw-attachment`, {
      vpcId: vpc.vpcId,
      internetGatewayId: igw.ref,
    })

    /**
     * Subnet
     */
    const appAvailabilityZones = ["ap-northeast-1a", "ap-northeast-1c"]
    const apPublicSubnets = appAvailabilityZones.map((az, i) => {
      const azSuffix = az.replace(/^.*-/, "")
      const appPublicSubnet = new ec2.Subnet(this, `${context.name}-app-public-subnet-${azSuffix}`, {
        vpcId: vpc.vpcId,
        cidrBlock: `10.0.${i + 1}.0/24`,
        availabilityZone: az,
        mapPublicIpOnLaunch: true,
      })
      appPublicSubnet.addDefaultInternetRoute(igw.ref, igwAttachment)
      Tags.of(appPublicSubnet).add('Name', `${context.name}-app-public-subnet-${azSuffix}`)
      return appPublicSubnet
    })
    const apPrivateSubnets = appAvailabilityZones.map((az, i) => {
      const azSuffix = az.replace(/^.*-/, "")
      const apPrivateSubnet = new ec2.Subnet(this, `${context.name}-app-private-subnet-${azSuffix}`, {
        vpcId: vpc.vpcId,
        cidrBlock: `10.0.${i + 11}.0/24`,
        availabilityZone: az,
        mapPublicIpOnLaunch: false,
      })
      apPrivateSubnet.addDefaultInternetRoute(igw.ref, igwAttachment)
      Tags.of(apPrivateSubnet).add('Name', `${context.name}-app-private-subnet-${azSuffix}`)
      return apPrivateSubnet
    })

    /**
     * VPC Endpoint
     */
    const ecrDockerVpce = new ec2.InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-ecr-dkr`,
      {
        service: ec2.InterfaceVpcEndpointAwsService.ECR_DOCKER,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        subnets: {
          subnets: apPrivateSubnets,
        },
      }
    );
    const ecrApiVpce = new ec2.InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-ecr-api`,
      {
        service: ec2.InterfaceVpcEndpointAwsService.ECR,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        subnets: {
          subnets: apPrivateSubnets,
        },
      }
    );
    const ssmVpce = new ec2.InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-ssm`,
      {
        service: ec2.InterfaceVpcEndpointAwsService.SSM,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        subnets: {
          subnets: apPrivateSubnets,
        },
      }
    );
    const secretMngVpce = new ec2.InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-secret-mng`,
      {
        service: ec2.InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        subnets: {
          subnets: apPrivateSubnets,
        },
      }
    );
    const logsVpce = new ec2.InterfaceVpcEndpoint(
      this,
      `${context.name}-vpce-logs`,
      {
        service: ec2.InterfaceVpcEndpointAwsService.CLOUDWATCH_LOGS,
        vpc: vpc,
        open: true,
        privateDnsEnabled: true,
        subnets: {
          subnets: apPrivateSubnets,
        },
      }
    );
    const s3Vpce = new ec2.GatewayVpcEndpoint(this, `${context.name}-vpce-s3`, {
      service: ec2.GatewayVpcEndpointAwsService.S3,
      vpc: vpc,
      subnets: [
        {
          subnets: apPrivateSubnets,
        },
      ],
    });

    const rdb = new rds.DatabaseInstance(this, `${context.name}-db`, {
      instanceIdentifier: `${context.name}-db`,
      databaseName: ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/db/database`),
      vpc: vpc,
      availabilityZone: "ap-northeast-1a",
      vpcSubnets: vpc.selectSubnets({
        subnets: apPrivateSubnets,
      }),
      securityGroups: [
        dbSg,
      ],
      engine: rds.DatabaseInstanceEngine.MYSQL,
      instanceType: ec2.InstanceType.of(ec2.InstanceClass.T3, ec2.InstanceSize.MICRO),
      allocatedStorage: 20,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
      storageType: rds.StorageType.GP2,
      port: dbPort,
      storageEncrypted: true,
      credentials: rds.Credentials.fromPassword(
        ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/db/user`),
        SecretValue.ssmSecure(`/${env}/go-auth-api/db/pass`),
      )
    })

    /**
     * EC2
     */
    // bastion
    const bastionKey = new ec2.CfnKeyPair(this, `${context.name}-bastion-key`, {
      keyName: `${context.name}-bastion-key`,
    })
    bastionKey.applyRemovalPolicy(cdk.RemovalPolicy.DESTROY)
    new ec2.Instance(this, `${context.name}-bastion`, {
      vpc,
      instanceName: `${context.name}-bastion`,
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T3,
        ec2.InstanceSize.MICRO,
      ),
      machineImage: ec2.MachineImage.latestAmazonLinux({
        generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2,
      }),
      securityGroup: bastionSg,
      vpcSubnets: {
        subnets: [
          apPublicSubnets[0],
        ],
      },
      keyName: bastionKey.keyName,
    })

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
    const logGroup = new logs.LogGroup(this, `${context.name}-log`, {
      logGroupName: '/aws/cdk/ecs-alb-fargate-service/web',
      retention: logs.RetentionDays.ONE_DAY,
      removalPolicy: cdk.RemovalPolicy.DESTROY,
    })

    /**
     * ALB
     */
    const alb = new elb.ApplicationLoadBalancer(this, `${context.name}-alb`, {
      vpc,
      internetFacing: true,
      securityGroup: albSg,
      loadBalancerName: `${context.name}-alb`,
      vpcSubnets: {
        subnets: apPublicSubnets,
      },
    })
    const httpListener = alb.addListener(`${context.name}-alb-http-listener`, {
      open: true,
      protocol: elb.ApplicationProtocol.HTTP,
      port: 80,
      defaultAction: elb.ListenerAction.redirect({
        protocol: "HTTPS",
        port: "443",
      }),
    })
    const httpsListener = alb.addListener(`${context.name}-alb-https-listener`, {
      open: true,
      protocol: elb.ApplicationProtocol.HTTPS,
      port: 443,
      certificates: [
        elb.ListenerCertificate.fromArn(
          `arn:aws:acm:ap-northeast-1:${Stack.of(this).account}:certificate/${context.loadBalancer.listener.certificate.ref}`,
        )
      ],
    })
    const targetGroup = new elb.ApplicationTargetGroup(this, `${context.name}-alb-target-group`, {
      vpc,
      protocol: elb.ApplicationProtocol.HTTP,
      targetType: elb.TargetType.IP,
      targetGroupName: `${context.name}-alb-target-group`,
      port: 3000,
      healthCheck: {
        path: "/health",
        interval: cdk.Duration.seconds(60),
        healthyHttpCodes: "200",
      },
    })
    httpsListener.addTargetGroups(`${context.name}-alb-https-target-group`, {
      targetGroups: [
        targetGroup,
      ],
    })

    /**
     * DNS
     */
    const hostedZone = route53.HostedZone.fromHostedZoneAttributes(
      this,
      `${context.name}-hosted-zone`,
      {
        hostedZoneId: ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/dns/zoneID`),
        zoneName: ssm.StringParameter.valueForStringParameter(this, `/${env}/go-auth-api/dns/zoneName`),
      },
    )
    const aliasRecord = new route53.ARecord(this, `${context.name}-alias-record`, {
      target: route53.RecordTarget.fromAlias(new targets.LoadBalancerTarget(alb)),
      zone: hostedZone,
      recordName: "auth",
    })

    /**
     * ECS
     */
    const cluster = new ecs.Cluster(this, `${context.name}-cluster`, {
      clusterName: `${context.name}-cluster`,
      vpc: vpc,
    })
    const taskDefinition = new ecs.TaskDefinition(this, `${context.name}-task`, {
      compatibility: ecs.Compatibility.FARGATE,
      cpu: "256",
      memoryMiB: "512",
      family: `${context.name}-task`,
      executionRole: executionRole,
    })
    taskDefinition.addContainer(`${context.name}-container`, {
      containerName: `${context.name}`,
      image: ecs.ContainerImage.fromEcrRepository(
        ecr.Repository.fromRepositoryName(this, `${context.name}-repo`, image.repositoryName), revision
      ),
      command: ["cmd/auth-api-server"],
      portMappings: [
        {
          name: "http-port-mapping",
          hostPort: 3000,
          containerPort: 3000,
          appProtocol: ecs.AppProtocol.http,
          protocol: ecs.Protocol.TCP,
        },
      ],
      logging: new ecs.AwsLogDriver({
        streamPrefix: 'ecs',
        logGroup: logGroup,
      }),
      environment: {
        "MYSQL_HOST": rdb.instanceEndpoint.hostname,
        "MYSQL_PORT": `${dbPort}`,
        "MYSQL_USER": dbUser,
        "MYSQL_DATABASE": dbName,
      },
      secrets: {
        'MYSQL_PASSWORD': ecs.Secret.fromSsmParameter(
          ssm.StringParameter.fromSecureStringParameterAttributes(this, 'ParameterRDBCredential', {
            parameterName: `/${env}/go-auth-api/db/pass`,
          })
        ),
      }
    })
    const service = new ecs.FargateService(this, `${context.name}-service`, {
      serviceName: `${context.name}-service`,
      cluster,
      taskDefinition,
      vpcSubnets: {
        subnets: apPrivateSubnets,
      },
      securityGroups: [
        apiServiceSg,
      ],
      desiredCount: 1,
    })
    service.attachToApplicationTargetGroup(targetGroup)
  }
}
