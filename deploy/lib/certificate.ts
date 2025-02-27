import {
  aws_route53 as route53,
  aws_certificatemanager as acm,
  CfnOutput,
  Stack,
  StackProps,
} from "aws-cdk-lib";
import { ICertificate } from "aws-cdk-lib/aws-certificatemanager";
import { Construct } from "constructs";

export interface CertificateStackProps extends StackProps {
  readonly zoneId: string;
  readonly zoneName: string;
  readonly domainName: string;
}

export class Certificate extends Stack {
  public readonly certificate: ICertificate;

  constructor(scope: Construct, id: string, props: CertificateStackProps) {
    super(scope, id, props);

    const zone = route53.HostedZone.fromHostedZoneAttributes(
      this,
      "BONextHostedZone",
      {
        hostedZoneId: props.zoneId,
        zoneName: props.zoneName,
      },
    );

    this.certificate = new acm.Certificate(this, "BONextCertificate", {
      domainName: `*.${props.zoneName}`,
      validation: acm.CertificateValidation.fromDns(zone),
      subjectAlternativeNames: [props.domainName],
    });

    new CfnOutput(this, "BONextCertificateArn", {
      value: this.certificate.certificateArn,
    });
  }
}
