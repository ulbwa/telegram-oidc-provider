import { redirect } from "next/navigation";
import { LoginRequestData } from "@/lib/types";
import { ENV } from "@/lib/env";
import { LoginClient } from "./login-client";

export default async function LoginPage({
  searchParams,
}: {
  searchParams: Promise<{ login_challenge?: string }>;
}) {
  console.log("▶ [LoginPage] Start: Initialization");

  const { login_challenge } = await searchParams;
  console.log(`ℹ [LoginPage] Params received. Challenge: ${login_challenge || "MISSING"}`);

  if (!login_challenge) {
    console.warn("⚠ [LoginPage] Warning: No login_challenge found.");
    
    if (ENV.isDev) {
      console.log("↪ [LoginPage] Dev mode detected. Redirecting to Google fallback...");
      redirect("https://google.com?complete=no_login_challenge");
    } else {
      console.log("⏹ [LoginPage] Prod mode. Rendering error UI (Missing Challenge).");
      return (
        <div className="flex min-h-screen items-center justify-center bg-background p-10 text-center text-destructive">
          Ошибка авторизации: неверный запрос.
        </div>
      );
    }
  }

  const serverApiUrl = ENV.SERVER_API_URL;
  console.log(`ℹ [LoginPage] Config: API URL is ${serverApiUrl}`);
  
  let data: LoginRequestData;

  try {
    const fetchUrl = `${serverApiUrl}/login?login_challenge=${login_challenge}`;
    console.log(`⏳ [LoginPage] Fetching data from: ${fetchUrl}`);

    const res = await fetch(fetchUrl, { cache: "no-store" });
    console.log(`ℹ [LoginPage] API Response Status: ${res.status} ${res.statusText}`);

    if (!res.ok) {
      console.error("✖ [LoginPage] Error: API response was not OK.");
      throw new Error(`API returned status: ${res.status} (${res.statusText})`);
    }

    data = await res.json();
    console.log("✅ [LoginPage] Data successfully parsed:", JSON.stringify(data, null, 2));

  } catch (e: any) {
    console.error("💥 [LoginPage] Exception caught:", e);

    if (ENV.isDev) {
      console.log("⏹ [LoginPage] Rendering Dev Error UI with Stack Trace.");
      return (
        <div className="flex min-h-screen flex-col gap-4 bg-destructive/5 p-8 text-destructive overflow-auto">
          <h1 className="text-2xl font-bold">
            Ошибка при загрузке данных
          </h1>

          <div className="rounded-lg border border-destructive/20 bg-card p-4 shadow-sm text-foreground">
            <h3 className="font-semibold text-destructive">Error Message:</h3>
            <pre className="mt-2 whitespace-pre-wrap text-sm text-destructive/80">
              {e.message || JSON.stringify(e)}
            </pre>
          </div>

          <div className="rounded-lg border border-destructive/20 bg-card p-4 shadow-sm text-foreground">
            <h3 className="font-semibold text-destructive">Target API URL:</h3>
            <code className="text-sm bg-muted px-2 py-1 rounded">
              {serverApiUrl}/login
            </code>
          </div>

          {e.stack && (
            <div className="rounded-lg border border-destructive/20 bg-card p-4 shadow-sm text-foreground">
              <h3 className="font-semibold text-destructive">Stack Trace:</h3>
              <pre className="mt-2 overflow-x-auto text-xs text-muted-foreground">
                {e.stack}
              </pre>
            </div>
          )}
        </div>
      );
    } else {
      console.log("⏹ [LoginPage] Rendering Public Error UI (Service Unavailable).");
      return (
        <div className="flex min-h-screen items-center justify-center bg-background p-4">
          <div className="text-center">
            <h1 className="text-xl font-semibold text-foreground">
              Сервис временно недоступен
            </h1>
            <p className="mt-2 text-muted-foreground">
              Попробуйте повторить попытку позже.
            </p>
          </div>
        </div>
      );
    }
  }

  if (data.skip) {
    console.log("↪ [LoginPage] Skip flag is TRUE. Redirecting back to Hydra...");
    redirect(`${serverApiUrl}/login?login_challenge=${login_challenge}`);
  }

  console.log("🎨 [LoginPage] Rendering LoginClient component.");
  
  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <LoginClient
        data={data}
        challenge={login_challenge}
        serverApiUrl={serverApiUrl!}
      />
    </div>
  );
}