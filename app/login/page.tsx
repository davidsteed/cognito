"use client";
import AcmeLogo from "@/app/ui/acme-logo";
import { Suspense } from "react";
import { Authenticator } from "@aws-amplify/ui-react";
import { Amplify } from "aws-amplify";
import "@aws-amplify/ui-react/styles.css";

Amplify.configure({
  Auth: {
    Cognito: {
      userPoolId: process.env.AWS_COGNITO_POOL_ID!,
      userPoolClientId: process.env.AWS_COGNITO_APP_CLIENT_ID!,
    },
  },
});
export default function LoginPage() {
  return (
    <main className="flex items-center justify-center md:h-screen">
      <div className="relative mx-auto flex w-full max-w-[512px] flex-col space-y-2.5 p-4 md:-mt-32">
        <div className="flex h-20 w-full items-end rounded-lg bg-blue-500 p-3 md:h-36">
          <div className="w-32 text-white md:w-36">
            <AcmeLogo />
          </div>
        </div>
        <Suspense>
          <Authenticator hideSignUp>
            {({ signOut, user }) => (
              <main>
                <h1>Hello {user?.username}</h1>
                <button onClick={signOut}>Sign out</button>
              </main>
            )}
          </Authenticator>
        </Suspense>
      </div>
    </main>
  );
}
