"use client";

import { Button } from "@/components/ui/button";
import { useEffect, useRef, useState } from "react";

// --- Логика отправки данных (Login) ---
async function sendAuthData(uid: string, user: any) {
  if (!uid || uid === '-1') return; // Защита

  try {
    const res = await fetch(`/api/oidc/interaction/${uid}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(user),
    });

    if (res.ok) {
      window.location.href = res.url;
    } else {
      console.error("Login failed");
      alert("Ошибка авторизации. Сессия истекла.");
    }
  } catch (e) {
    console.error(e);
  }
}

// 1. КНОПКА СТАРТА (Эмуляция начала входа)
export function StartFlowButton() {
  const handleStart = () => {
    const params = new URLSearchParams({
      client_id: 'kratos-client',
      response_type: 'code',
      scope: 'openid profile email',
      redirect_uri: 'http://localhost:3000/api/oidc/debug',
    });

    window.location.href = `/api/oidc/auth?${params.toString()}`;
  };

  return (
    <Button 
      onClick={handleStart}
      className="w-full h-12 rounded-xl bg-tg-blue text-base font-semibold text-white shadow-lg shadow-tg-blue/30 transition-all hover:bg-tg-blue/90 hover:-translate-y-0.5 active:scale-[0.98]"
    >
      Войти или создать аккаунт
    </Button>
  );
}

// 2. REAL TELEGRAM WIDGET
export function RealTelegramWidget({ botName, uid }: { botName: string, uid: string }) {
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current?.innerHTML) return;
    // @ts-ignore
    window.onTelegramAuth = (user) => sendAuthData(uid, user);

    const script = document.createElement("script");
    script.src = "https://telegram.org/js/telegram-widget.js?22";
    script.setAttribute("data-telegram-login", botName);
    script.setAttribute("data-size", "large");
    script.setAttribute("data-radius", "10");
    script.setAttribute("data-request-access", "write");
    script.setAttribute("data-userpic", "false");
    script.setAttribute("data-onauth", "onTelegramAuth(user)");
    script.async = true;
    containerRef.current?.appendChild(script);
  }, [botName, uid]);

  return <div ref={containerRef} className="flex justify-center" />;
}

// 3. MOCK BUTTON (Андрей Пупкин)
export function MockLoginButton({ targetUid }: { targetUid: string }) {
  const [loading, setLoading] = useState(false);

  const handleMock = async () => {
    setLoading(true);
    const mockUser = {
      id: 123456789,
      first_name: "Андрей",
      last_name: "Пупкин",
      username: "pupkin_dev",
      photo_url: "",
      auth_date: Math.floor(Date.now() / 1000),
      hash: "mock",
      mock: true
    };
    await sendAuthData(targetUid, mockUser);
    setLoading(false);
  };

  return (
    <Button
      variant="ghost"
      onClick={handleMock}
      disabled={loading}
      className="w-full text-tg-text-secondary hover:text-tg-blue hover:bg-tg-blue/10"
    >
      {loading ? "Вход..." : "🛠 Mock: Андрей Пупкин"}
    </Button>
  );
}