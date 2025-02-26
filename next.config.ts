import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "export",
  env: {
    AWS_COGNITO_REGION: process.env.AWS_COGNITO_REGION,
    AWS_COGNITO_POOL_ID: process.env.AWS_COGNITO_POOL_ID,
    AWS_COGNITO_APP_CLIENT_ID: process.env.AWS_COGNITO_APP_CLIENT_ID,
  },
};

module.exports = nextConfig;
