#!/usr/bin/env node
import * as cdk from "aws-cdk-lib";
import { SinglePageApp } from "../lib/deploy-stack";
import { Certificate } from "../lib/certificate";
import { Cognito } from "../lib/cognito";
import { DemoAPI } from "../lib/api";

const app = new cdk.App();

const zoneId = app.node.tryGetContext("zoneId");
const zoneName = app.node.tryGetContext("zoneName");
const subdomain = app.node.tryGetContext("subdomain");
const fullDomain = `${subdomain}.${zoneName}`;

const certificate = new Certificate(app, "demo-certificate", {
  domainName: fullDomain,
  zoneId,
  zoneName,
  env: {
    region: "us-east-1",
  },
  crossRegionReferences: true,
});

const apiCertificate = new Certificate(app, "api-certificate", {
  domainName: `api.${zoneName}`,
  zoneId,
  zoneName,
  env: {
    region: "eu-west-1",
  },
  crossRegionReferences: true,
});

new SinglePageApp(app, "demo-website", {
  fullDomainName: fullDomain,
  subdomainName: subdomain,
  zoneId: zoneId,
  zoneName: zoneName,
  certificate: certificate.certificate,
  crossRegionReferences: true,
  env: {
    region: "eu-west-1",
  },
});

const cognito = new Cognito(app, "demo-cognito", {
  fullDomain: fullDomain,
  zoneName: zoneName,
  crossRegionReferences: true,
  env: {
    region: "eu-west-1",
  },
});

new DemoAPI(app, "demo-api", {
  fullDomain: fullDomain,
  zoneId: zoneId,
  zoneName: zoneName,
  userPoolId: cognito.userPool.userPoolId,
  apiCertificate: apiCertificate.certificate,
  crossRegionReferences: true,
  env: {
    region: "eu-west-1",
  },
});
