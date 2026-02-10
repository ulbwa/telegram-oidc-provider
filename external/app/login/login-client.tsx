"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Loader2 } from "lucide-react";
import { LoginRequestData } from "@/lib/types";
import WebApp from "@twa-dev/sdk";

interface LoginClientProps {
  data: LoginRequestData;
  challenge: string;
  serverApiUrl: string;
}

export function LoginClient({
  data,
  challenge,
  serverApiUrl,
}: LoginClientProps) {
  const [isCheckingTwa, setIsCheckingTwa] = useState(true);

  // 1. Проверка Telegram Web App
  useEffect(() => {
    const tg = WebApp;
    if (tg && tg.initData) {
      const initData = tg.initData;
      console.log("Telegram initData: ", String(initData));
      const redirectUrl = `${serverApiUrl}/login/mini-app?login_challenge=${challenge}&${initData}`;
      window.location.href = redirectUrl;
    } else {
      setIsCheckingTwa(false);
    }
  }, [challenge, serverApiUrl]);

  // 2. Инъекция скрипта виджета (Только если auth = false)
  useEffect(() => {
    if (isCheckingTwa) return;
    if (data.auth) return; // Если юзер известен, скрипт не грузим

    // Защита от дублей
    if (document.getElementById("telegram-widget-script")) return;

    const script = document.createElement("script");
    script.id = "telegram-widget-script";
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.setAttribute("data-telegram-login", data.bot.username);
    script.setAttribute("data-size", "large");
    script.setAttribute("data-radius", "12");
    script.setAttribute("data-request-access", "write");
    script.setAttribute("data-userpic", "false");
    script.setAttribute(
      "data-auth-url",
      `${serverApiUrl}/login/widget?login_challenge=${challenge}`,
    );

    document.getElementById("telegram-widget-container")?.appendChild(script);
  }, [isCheckingTwa, data.auth, data.bot.username, serverApiUrl, challenge]);

  // --- LOADING ---
  if (isCheckingTwa) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <Loader2 className="h-10 w-10 animate-spin text-primary" />
      </div>
    );
  }

  // ==========================================
  // VIEW 1: СУЩЕСТВУЮЩИЙ АККАУНТ (Старый дизайн)
  // ==========================================
  if (data.auth && data.user) {
    return (
      <Card className="w-full max-w-md border-none shadow-xl bg-card/80 backdrop-blur-xl">
        <CardHeader className="text-center pb-2">
          <CardTitle className="text-xl">
            Здравствуйте, {data.user.first_name}!
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center gap-6">
          <Avatar className="h-24 w-24 border-4 border-background shadow-sm">
            <AvatarImage src={data.user.photo_url} />
            <AvatarFallback className="text-lg">
              {data.user.first_name[0]}
            </AvatarFallback>
          </Avatar>

          <p className="text-center text-muted-foreground">
            Приложение{" "}
            <strong className="text-foreground">{data.client.name}</strong>{" "}
            запрашивает доступ.
            <br />
            Использовать этот аккаунт?
          </p>

          <div className="grid w-full grid-cols-2 gap-4">
            <Button
              variant="outline"
              asChild
              className="w-full h-11 border-border hover:bg-muted"
            >
              <a
                href={`${serverApiUrl}/login/reject?login_challenge=${challenge}`}
              >
                Нет, другой
              </a>
            </Button>
            <Button
              asChild
              className="w-full h-11 font-semibold bg-primary text-primary-foreground hover:bg-primary/90"
            >
              <a href={`${serverApiUrl}/login?login_challenge=${challenge}`}>
                Да, войти
              </a>
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  // ==========================================
  // VIEW 2: ВХОД ЧЕРЕЗ ВИДЖЕТ (Утка + DevTools)
  // ==========================================
  const mockLoginUrl = `${serverApiUrl}/login?login_challenge=${challenge}skip`;

  return (
    <div className="w-full max-w-[400px] text-center">
      {/* УТКА */}
      <div className="mx-auto mb-6 h-[150px] w-[150px] overflow-hidden rounded-2xl bg-muted/20">
        <iframe
          src="https://tenor.com/embed/4989962420582851395"
          className="h-full w-full border-0 pointer-events-none"
          allowFullScreen
          title="Telegram Sticker"
        />
      </div>

      <h1 className="mb-2 text-2xl font-bold text-foreground">
        Войти через Telegram
      </h1>

      <p className="mb-8 text-base leading-relaxed text-muted-foreground">
        Используйте ваш аккаунт Telegram для быстрого и безопасного входа в
        <span className="font-semibold text-primary"> {data.client.name}</span>.
      </p>

      {/* КОНТЕЙНЕР ДЛЯ ВИДЖЕТА */}
      <div className="min-h-[50px] flex justify-center mb-6">
        <div id="telegram-widget-container" />
      </div>

      {/* РАЗДЕЛИТЕЛЬ DEV TOOLS */}
      <div className="space-y-6">
        <div className="relative py-2">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t border-border"></span>
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-background px-2 text-muted-foreground">
              Dev Tools
            </span>
          </div>
        </div>

        {/* MOCK КНОПКА */}
        <div>
          <Button
            variant="outline"
            className="w-full border-dashed border-primary/50 text-primary hover:bg-primary/5 h-12 rounded-xl"
            asChild
          >
            <a href={mockLoginUrl}>Александр Пупкин</a>
          </Button>
        </div>
      </div>
    </div>
  );
}
