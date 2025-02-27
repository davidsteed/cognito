import {
  aws_cognito as cognito,
  CfnOutput,
  Stack,
  StackProps,
} from "aws-cdk-lib";
import { Construct } from "constructs";

export interface CognitoStackProps extends StackProps {
  readonly fullDomain: string;
  readonly zoneName: string;
}

export class Cognito extends Stack {
  public readonly userPool: cognito.UserPool;
  constructor(scope: Construct, id: string, props: CognitoStackProps) {
    super(scope, id, props);

    this.userPool = new cognito.UserPool(this, id + "DevicePool", {
      userPoolName: id + "UserPool",
      signInAliases: {
        username: true,
      },
      selfSignUpEnabled: false,
      customAttributes: {
        locationid: new cognito.StringAttribute({
          minLen: 5,
          maxLen: 15,
          mutable: false,
        }),
      },
      passwordPolicy: {
        minLength: 8,
        requireDigits: true,
        requireLowercase: true,
        requireUppercase: true,
        requireSymbols: true,
      },
      // lambdaTriggers: {
      //   preAuthentication: dummyLambda,
      //   preSignUp: dummyLambda,
      // },
    });

    const client = new cognito.CfnUserPoolClient(this, id + "AppClient", {
      userPoolId: this.userPool.userPoolId,
      clientName: "SPMClient",
      supportedIdentityProviders: ["COGNITO"],
      callbackUrLs: ["https://" + props.fullDomain],
      allowedOAuthFlows: ["implicit"],
      preventUserExistenceErrors: "ENABLED",
      enableTokenRevocation: true,
      explicitAuthFlows: [
        "ALLOW_ADMIN_USER_PASSWORD_AUTH",
        "ALLOW_CUSTOM_AUTH",
        "ALLOW_USER_PASSWORD_AUTH",
        "ALLOW_USER_SRP_AUTH",
        "ALLOW_REFRESH_TOKEN_AUTH",
      ],
      allowedOAuthFlowsUserPoolClient: true,
      allowedOAuthScopes: ["openid"],
    });

    // Create Resource Server
    const cognitoAPIDomain = "auth." + props.zoneName;
    new cognito.CfnUserPoolResourceServer(this, id + "ResourceServer", {
      userPoolId: this.userPool.userPoolId,
      identifier: cognitoAPIDomain,
      name: id,
      scopes: [
        { scopeName: "deviceAccess", scopeDescription: "Access Device APIs" },
      ],
    });

    new CfnOutput(this, `${id}-UserPoolId`, {
      value: this.userPool.userPoolId,
    });
    new CfnOutput(this, `${id}-ClientId`, { value: client.ref });
  }
}
