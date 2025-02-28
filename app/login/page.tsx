"use client";
import AcmeLogo from "@/app/ui/acme-logo";
import { Suspense } from "react";
import { Authenticator } from "@aws-amplify/ui-react";
import { Amplify } from "aws-amplify";
import { Button } from "@/app/ui/button";
import { useAuthenticator } from "@aws-amplify/ui-react";
import { useGetLocation } from "@/app/generator";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "@aws-amplify/ui-react/styles.css";

//const baseURL = 'https://api.testawsreact.com';

Amplify.configure({
  Auth: {
    Cognito: {
      userPoolId: process.env.AWS_COGNITO_POOL_ID!,
      userPoolClientId: process.env.AWS_COGNITO_APP_CLIENT_ID!,
    },
  },
});

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
    },
  },
});

const App = () => {
  const { authStatus } = useAuthenticator((context) => [context.authStatus]);
  const { isLoading: isBranchDetailsLoading } = useGetLocation({});

  return (
    <div>
      <h1>
        Hello {authStatus} {isBranchDetailsLoading}
      </h1>
    </div>
  );
};

export default function LoginPage() {
  return (
    <main className="flex items-center justify-center md:h-screen">
      <div className="relative mx-auto flex w-full max-w-[512px] flex-col space-y-2.5 p-4 md:-mt-32">
        <div className="flex h-20 w-full items-end rounded-lg bg-blue-500 p-3 md:h-36">
          <div className="w-32 text-white md:w-36">
            <AcmeLogo />
          </div>
        </div>
        <QueryClientProvider client={queryClient}>
          <Suspense>
            <Authenticator hideSignUp>
              {({ signOut, user }) => (
                <main>
                  <h1>Hello {user?.username}</h1>
                  <App></App>
                  <Button onClick={signOut}>Sign out</Button>
                </main>
              )}
            </Authenticator>
          </Suspense>
        </QueryClientProvider>
      </div>
    </main>
  );
}
