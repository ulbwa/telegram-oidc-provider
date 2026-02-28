import { ENV } from "@/lib/env";
import { NextResponse } from "next/server";

// Типы для наглядности
interface User {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
}

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const challenge = searchParams.get("login_challenge") || "";

  console.log(`⚡ [MockAPI] Request received. Challenge: "${challenge}"`);

  // Базовая структура бота и клиента
  const baseData = {
    client: { name: "MedIncident" },
    bot: { name: "SosiskaBot", username: ENV.TELEGRAM_BOT_NAME, url: "https://t.me/my_auth_bot" },
  };

  const mockUser: User = {
    id: 123456789,
    first_name: "Андрей",
    last_name: "Тестовый",
    username: "test_dev",
    photo_url: "https://github.com/shadcn.png",
  };

  // --- СЦЕНАРИЙ 1: ERROR (Тест экрана ошибки) ---
  // ?login_challenge=error
  if (challenge.includes("error")) {
    console.log("❌ [MockAPI] Simulating 500 Server Error");
    return new NextResponse(
      JSON.stringify({ error: "Simulated Internal Server Error" }),
      { status: 500, statusText: "Simulated Crash" }
    );
  }

  // --- СЦЕНАРИЙ 2: SKIP (Мгновенный редирект) ---
  // ?login_challenge=skip
  // Пользователь уже вошел и дал права ранее
  if (challenge.includes("skip")) {
    console.log("⏩ [MockAPI] Scenario: SKIP (Redirect back immediately)");
    return NextResponse.json({
      auth: true,
      skip: true,
      user: mockUser,
      ...baseData,
    });
  }

  // --- СЦЕНАРИЙ 3: AUTH (Подтверждение входа) ---
  // ?login_challenge=auth
  // Пользователь известен (есть кука), но нужно нажать "Да, войти"
  if (challenge.includes("auth")) {
    console.log("👤 [MockAPI] Scenario: AUTH RECOGNIZED (Show confirm card)");
    return NextResponse.json({
      auth: true,
      skip: false,
      user: mockUser,
      ...baseData,
    });
  }

  // --- СЦЕНАРИЙ 4: WIDGET (Новый пользователь) ---
  // ?login_challenge=widget (или любой другой)
  // Пользователь неизвестен, показываем виджет Telegram
  console.log("🎨 [MockAPI] Scenario: WIDGET (Default state)");
  return NextResponse.json({
    auth: false,
    skip: false,
    ...baseData,
  });
}