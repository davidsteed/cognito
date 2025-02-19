#!/usr/bin/env node
import * as cdk from 'aws-cdk-lib';
import { Cognito} from '../lib/deploy-stack';
import { Certificate } from '../lib/certificate';

const app = new cdk.App();

const zoneId = app.node.tryGetContext('zoneId');
const zoneName = app.node.tryGetContext('zoneName');
const subdomain = app.node.tryGetContext("subdomain");
const fullDomain = `${subdomain}.${zoneName}`;

const certificate = new Certificate(app, 'certificate', {
  domainName: fullDomain,
  zoneId,
  zoneName,
  env: {
    region: 'us-east-1',
  },
  crossRegionReferences: true,
});


new Cognito(app, 'cognito', {
  fullDomainName: fullDomain,
  subdomainName: subdomain,
  zoneId:zoneId,
  zoneName:zoneName,
  certificate: certificate.certificate,
  crossRegionReferences: true,
  env: {
    region: 'eu-west-1',
  },
});