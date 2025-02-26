import { Construct } from "constructs";
import {
  Distribution,
  HeadersFrameOption,
  HeadersReferrerPolicy,
  OriginAccessIdentity,
  ResponseHeadersPolicy,
  SecurityPolicyProtocol,
} from "aws-cdk-lib/aws-cloudfront";
import { ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import {
  CfnOutput,
  Duration,
  RemovalPolicy,
  Stack,
  StackProps,
} from "aws-cdk-lib";
import { ARecord, HostedZone, RecordTarget } from "aws-cdk-lib/aws-route53";
import { Bucket, BucketEncryption } from "aws-cdk-lib/aws-s3";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import { S3BucketOrigin } from "aws-cdk-lib/aws-cloudfront-origins";
import {
  BucketDeployment,
  CacheControl,
  Source,
} from "aws-cdk-lib/aws-s3-deployment";
import { CloudFrontTarget } from "aws-cdk-lib/aws-route53-targets";
import { StringParameter } from "aws-cdk-lib/aws-ssm";

export interface CognitoStackProps extends StackProps {
  readonly fullDomainName: string;
  readonly subdomainName: string;
  readonly zoneId: string;
  readonly zoneName: string;
  readonly certificate: ICertificate;
}

export class SinglePageApp extends Stack {
  constructor(scope: Construct, id: string, props: CognitoStackProps) {
    super(scope, id, props);

    const zone = HostedZone.fromHostedZoneAttributes(this, `${id}-HostedZone`, {
      hostedZoneId: props.zoneId,
      zoneName: props.zoneName,
    });
    const bucketName = `${id}-bucket-` + props.zoneId.toLowerCase();
    const bucket = new Bucket(this, `{id}-Bucket`, {
      bucketName: bucketName,
      publicReadAccess: false,
      removalPolicy: RemovalPolicy.DESTROY,
      encryption: BucketEncryption.S3_MANAGED,
      enforceSSL: true,
      websiteIndexDocument: "index.html",
      websiteErrorDocument: "index.html",
    });

    const cloudfrontOAI = new OriginAccessIdentity(this, `${id}-OAI`, {
      comment: `OAI for ${bucket.bucketArn} app bucket`,
    });

    const cloudfrontS3Access = new PolicyStatement();

    cloudfrontS3Access.addActions("s3:GetBucket*");
    cloudfrontS3Access.addActions("s3:GetObject*");
    cloudfrontS3Access.addResources(bucket.bucketArn);
    cloudfrontS3Access.addResources(`${bucket.bucketArn}/*`);
    cloudfrontS3Access.addCanonicalUserPrincipal(
      cloudfrontOAI.cloudFrontOriginAccessIdentityS3CanonicalUserId,
    );

    bucket.addToResourcePolicy(cloudfrontS3Access);

    const responseHeadersPolicy = new ResponseHeadersPolicy(
      this,
      `${id}-origin-request`,
      {
        responseHeadersPolicyName: `security-headers`,
        securityHeadersBehavior: {
          strictTransportSecurity: {
            accessControlMaxAge: Duration.days(365),
            includeSubdomains: true,
            preload: true,
            override: true,
          },
          contentSecurityPolicy: {
            contentSecurityPolicy: `frame-ancestors 'none';`,
            override: true,
          },
          contentTypeOptions: {
            override: true,
          },
          frameOptions: {
            frameOption: HeadersFrameOption.DENY,
            override: true,
          },
          xssProtection: {
            modeBlock: true,
            protection: true,
            override: true,
          },
          referrerPolicy: {
            referrerPolicy: HeadersReferrerPolicy.SAME_ORIGIN,
            override: true,
          },
        },
        customHeadersBehavior: {
          customHeaders: [
            {
              header: "Permissions-Policy",
              value:
                "accelerometer=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()",
              override: true,
            },
          ],
        },
      },
    );

    // Rewrites URLs to append .html if the URL doesn't end with a file extension (needed because output files are named with .html)
    // const urlRewriteFunction = new cloudfront.Function(this, 'UrlRewriteFunction', {
    //   code: cloudfront.FunctionCode.fromInline(`
    //     function handler(event) {
    //       var request = event.request;
    //       var uri = request.uri;

    //       // Redirect root path to /task-management
    //       if (uri === '/') {
    //         uri = '/task-management'; // Set the URI to /task-management
    //       }

    //       // Check if URI doesn't end with a file extension
    //       if (!uri.includes('.')) {
    //         uri += '.html'; // Append .html to the URI
    //       }

    //       request.uri = uri; // Update the request URI

    //       return request;
    //     }
    //   `),
    // });

    const distribution = new Distribution(this, `${id}-distribution`, {
      domainNames: [props.fullDomainName],
      certificate: props.certificate,
      minimumProtocolVersion: SecurityPolicyProtocol.TLS_V1_2_2021,
      defaultRootObject: "index.html",
      errorResponses: [
        {
          httpStatus: 403,
          responseHttpStatus: 403,
          responsePagePath: "/index.html",
          ttl: Duration.days(1),
        },
        {
          httpStatus: 404,
          responseHttpStatus: 404,
          responsePagePath: "/index.html",
          ttl: Duration.days(1),
        },
      ],
      defaultBehavior: {
        origin: S3BucketOrigin.withOriginAccessControl(bucket),
        compress: true,
        responseHeadersPolicy,
        // functionAssociations: [
        //   {
        //     eventType: cloudfront.FunctionEventType.VIEWER_REQUEST,
        //     function: urlRewriteFunction,
        //   },
        // ],
      },
    });

    // Upload the files in the /out/_next folder with a Cache-Control value of 1 year (static assets)
    new BucketDeployment(this, `${id}-deploy-with-invalidation`, {
      sources: [Source.asset("../out/_next")],
      destinationBucket: bucket,
      destinationKeyPrefix: "_next",
      distribution: distribution,
      distributionPaths: ["/*"],
      cacheControl: [CacheControl.fromString("max-age=31536000")],
      prune: false,
    });

    // Upload the files in the /out folder with a no-cache header (HTML files)
    new BucketDeployment(this, `${id}-deploy-with-invalidation-no-cache`, {
      sources: [
        Source.asset("../out", {
          exclude: ["_next/**/*"], // Exclude all files in the _next directory
        }),
      ],
      destinationBucket: bucket,
      destinationKeyPrefix: "",
      distribution: distribution,
      distributionPaths: ["/*"],
      cacheControl: [CacheControl.noCache()],
      prune: false,
    });

    new ARecord(this, `${id}-SiteAliasRecord`, {
      recordName: props.subdomainName,
      target: RecordTarget.fromAlias(new CloudFrontTarget(distribution)),
      zone,
    });

    const distributionIdParam = new StringParameter(
      this,
      `${id}-distribution-id-param`,
      {
        parameterName: `${id}-distribution-id-param`,
        description: `${id}- cloudfront distribution id`,
        stringValue: distribution.distributionId,
      },
    );

    const distributionDomainNameParam = new StringParameter(
      this,
      `${id}-distribution-domain-name-param`,
      {
        parameterName: `${id}-distribution-domain-name-param`,
        description: `${id}-cloudfront distribution domain name`,
        stringValue: distribution.distributionDomainName,
      },
    );

    new CfnOutput(this, `${id}-Distribution-Id-Param`, {
      value: distributionIdParam.stringValue,
    });

    new CfnOutput(this, `${id}-Distribution-DomainName-Param`, {
      value: distributionDomainNameParam.stringValue,
    });
  }
}
