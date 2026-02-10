// app/oidc-status/page.tsx
"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertCircle, CheckCircle, Loader2 } from "lucide-react";

export default function OidcStatusPage() {
  const [status, setStatus] = useState<"loading" | "success" | "error">("loading");
  const [config, setConfig] = useState<any>(null);
  const [errorMsg, setErrorMsg] = useState("");

  const discoveryUrl = "/api/oidc/.well-known/openid-configuration";

  const checkOidc = async () => {
    setStatus("loading");
    setConfig(null);
    setErrorMsg("");

    try {
      const res = await fetch(discoveryUrl);
      if (!res.ok) {
        throw new Error(`Ошибка ${res.status}: ${res.statusText}`);
      }
      const data = await res.json();
      setConfig(data);
      setStatus("success");
    } catch (e: any) {
      console.error(e);
      setStatus("error");
      setErrorMsg(e.message || "Не удалось соединиться с провайдером");
    }
  };

  useEffect(() => {
    checkOidc();
  }, []);

  return (
    <div className="min-h-screen bg-zinc-50 p-8 font-sans dark:bg-zinc-950 text-foreground">
      <div className="mx-auto max-w-3xl space-y-6">
        
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">OIDC Provider Status</h1>
          <Button onClick={checkOidc} variant="outline">
            Обновить
          </Button>
        </div>

        {/* Статус бар */}
        <Card className={status === "error" ? "border-red-500/50 bg-red-500/5" : ""}>
          <CardContent className="flex items-center gap-4 p-6">
            {status === "loading" && <Loader2 className="h-8 w-8 animate-spin text-blue-500" />}
            {status === "success" && <CheckCircle className="h-8 w-8 text-green-500" />}
            {status === "error" && <AlertCircle className="h-8 w-8 text-red-500" />}
            
            <div>
              <p className="font-semibold">
                {status === "loading" && "Проверка соединения..."}
                {status === "success" && "Провайдер работает корректно"}
                {status === "error" && "Ошибка подключения"}
              </p>
              <p className="text-sm text-muted-foreground">
                URL: <code className="bg-muted px-1 py-0.5 rounded">{discoveryUrl}</code>
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Блок с JSON конфигом */}
        {status === "success" && (
          <Card>
            <CardHeader>
              <CardTitle>OpenID Configuration</CardTitle>
            </CardHeader>
            <CardContent>
              <pre className="max-h-[500px] overflow-auto rounded-lg bg-zinc-900 p-4 text-xs text-zinc-50">
                {JSON.stringify(config, null, 2)}
              </pre>
            </CardContent>
          </Card>
        )}

        {/* Блок с ошибкой */}
        {status === "error" && (
          <Card>
            <CardHeader>
              <CardTitle className="text-red-500">Детали ошибки</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="mb-4 text-sm text-muted-foreground">
                Возможные причины:
              </p>
              <ul className="list-disc pl-5 text-sm space-y-1 text-muted-foreground">
                <li>Ошибка в <code>pages/api/oidc/[...oidc].ts</code> (например, синтаксис).</li>
                <li>Не установлены библиотеки <code>oidc-provider</code> или <code>jose</code>.</li>
                <li>Неверный <code>ISSUER_URL</code> в .env (не должен быть слеша в конце).</li>
                <li>Провайдер падает при старте (смотри терминал сервера).</li>
              </ul>
              <div className="mt-4 rounded bg-red-100 p-4 text-red-900 dark:bg-red-900/20 dark:text-red-200">
                {errorMsg}
              </div>
            </CardContent>
          </Card>
        )}

      </div>
    </div>
  );
}