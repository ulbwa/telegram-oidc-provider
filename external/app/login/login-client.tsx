"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Loader2, Database, Terminal } from "lucide-react"; // Добавил иконки
import { LoginRequestData } from "@/lib/types";

interface LoginClientProps {
  data: LoginRequestData;
  challenge: string;
  serverApiUrl: string;
}

export function LoginClient({ data, challenge, serverApiUrl }: LoginClientProps) {
  const [isCheckingTwa, setIsCheckingTwa] = useState(true);
  // Сохраняем initData, чтобы показать её на экране
  const [twaInitData, setTwaInitData] = useState<string | null>(null);

  // 1. Проверка Telegram Web App (с динамическим импортом)
  useEffect(() => {
    import("@twa-dev/sdk").then((sdk) => {
      const tg = sdk.default;
      
      if (tg && tg.initData) {
        const initData = tg.initData;
        setTwaInitData(initData); // <--- Сохраняем данные для отображения

        console.log("✅ [LoginClient] TWA detected. InitData saved.");
        
        // ВАЖНО: Если ты хочешь успеть прочитать данные, можно закомментировать авто-редирект ниже
        // const redirectUrl = `${serverApiUrl}/login/mini-app?login_challenge=${challenge}&${initData}`;
        // window.location.href = redirectUrl;
        
        // Для теста просто выключаем лоадер, чтобы увидеть UI с данными
        setIsCheckingTwa(false); 
      } else {
        setIsCheckingTwa(false);
      }
    }).catch((err) => {
      console.error("Failed to load TWA SDK:", err);
      setIsCheckingTwa(false);
    });
  }, [challenge, serverApiUrl]);

  // 2. Инъекция скрипта виджета (Только если auth = false и TWA не найден)
  useEffect(() => {
    if (isCheckingTwa) return;
    if (data.auth) return;
    // Если мы нашли TWA данные, возможно, виджет нам и не нужен (зависит от логики), 
    // но оставим его для fallback сценариев
    
    if (document.getElementById("telegram-widget-script")) return;

    const script = document.createElement("script");
    script.id = "telegram-widget-script";
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.setAttribute("data-telegram-login", data.bot.username);
    script.setAttribute("data-size", "large");
    script.setAttribute("data-radius", "12");
    script.setAttribute("data-request-access", "write");
    script.setAttribute("data-userpic", "false");
    script.setAttribute("data-auth-url", `${serverApiUrl}/login/widget?login_challenge=${challenge}`);
    
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

  // Вспомогательный компонент для отображения initData
  const TwaDebugBlock = () => {
    if (!twaInitData) return null;

    // Пытаемся распарсить строку для красивого отображения
    let parsedData = twaInitData;
    try {
        const params = new URLSearchParams(twaInitData);
        const obj: any = {};
        params.forEach((value, key) => {
            // Пробуем распарсить JSON внутри параметров (например, user)
            try { obj[key] = JSON.parse(value); } catch { obj[key] = value; }
        });
        parsedData = JSON.stringify(obj, null, 2);
    } catch (e) {}

    return (
      <div className="mt-4 w-full animate-in fade-in slide-in-from-bottom-2">
         <div className="rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-3 text-left">
            <div className="mb-2 flex items-center gap-2 text-xs font-bold text-yellow-600 dark:text-yellow-400">
                <Database className="h-3 w-3" />
                <span>TWA Init Data Detected</span>
                <span>{twaInitData}</span>
            </div>
            <pre className="max-h-[150px] overflow-auto rounded bg-background/50 p-2 text-[10px] font-mono leading-tight text-foreground/80 scrollbar-thin scrollbar-thumb-border">
                {parsedData}
            </pre>
            <div className="mt-2">
                 <Button 
                    size="sm" 
                    variant="secondary" 
                    className="h-7 w-full text-xs"
                    onClick={() => {
                        const url = `${serverApiUrl}/login/mini-app?login_challenge=${challenge}&${twaInitData}`;
                        window.location.href = url;
                    }}
                 >
                    Force Login with TWA Data
                 </Button>
            </div>
         </div>
      </div>
    );
  };

  // ==========================================
  // VIEW 1: СУЩЕСТВУЮЩИЙ АККАУНТ (Строгая карточка)
  // ==========================================
  if (data.auth && data.user) {
    return (
      <div className="flex min-h-screen flex-col items-center justify-center bg-background p-4">
        <Card className="w-full max-w-md border-none shadow-xl bg-card/80 backdrop-blur-xl animate-in fade-in zoom-in duration-300">
            <CardHeader className="text-center pb-2">
            <CardTitle className="text-xl">Здравствуйте, {data.user.first_name}!</CardTitle>
            </CardHeader>
            <CardContent className="flex flex-col items-center gap-6">
            <Avatar className="h-24 w-24 border-4 border-background shadow-sm">
                <AvatarImage src={data.user.photo_url} />
                <AvatarFallback className="text-lg">{data.user.first_name[0]}</AvatarFallback>
            </Avatar>
            
            <p className="text-center text-muted-foreground">
                Приложение <strong className="text-foreground">{data.client.name}</strong> запрашивает доступ.
                <br />
                Использовать этот аккаунт?
            </p>

            <div className="grid w-full grid-cols-2 gap-4">
                <Button variant="outline" asChild className="w-full h-11 border-border hover:bg-muted">
                <a href={`${serverApiUrl}/login/reject?login_challenge=${challenge}`}>
                    Нет, другой
                </a>
                </Button>
                <Button asChild className="w-full h-11 font-semibold bg-primary text-primary-foreground hover:bg-primary/90">
                <a href={`${serverApiUrl}/login?login_challenge=${challenge}`}>
                    Да, войти
                </a>
                </Button>
            </div>
            
            {/* Показываем данные TWA даже если юзер залогинен (для дебага) */}
            <TwaDebugBlock />

            </CardContent>
        </Card>
      </div>
    );
  }

  // ==========================================
  // VIEW 2: ВХОД ЧЕРЕЗ ВИДЖЕТ (Утка + DevTools)
  // ==========================================
  const mockLoginUrl = `${serverApiUrl}/login?login_challenge=${challenge}skip`;

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4 font-sans text-foreground">
        <div className="w-full max-w-[400px] text-center animate-in fade-in slide-in-from-bottom-4 duration-500">
            
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

        {/* --- DEV TOOLS AREA --- */}
        <div className="space-y-6">
            <div className="relative py-2">
            <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t border-border"></span>
            </div>
            <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-background px-2 text-muted-foreground">Dev Tools</span>
            </div>
            </div>

            {/* БЛОК TWA ДАННЫХ (Появляется если есть initData) */}
            <TwaDebugBlock />

            {/* MOCK КНОПКА */}
            <div>
                <Button 
                    variant="outline" 
                    className="w-full border-dashed border-primary/50 text-primary hover:bg-primary/5 h-12 rounded-xl"
                    asChild
                >
                    <a href={mockLoginUrl}>
                    <Terminal className="mr-2 h-4 w-4" />
                    Василий Пупкин
                    </a>
                </Button>
            </div>
        </div>
        </div>
    </div>
  );
}