import { ENV } from "@/lib/env";
import { RealTelegramWidget, MockLoginButton, StartFlowButton } from "./login-widget";

export default async function LoginPage({
  searchParams,
}: {
  searchParams: Promise<{ uid?: string }>;
}) {
  const { uid } = await searchParams;
  
  // Проверяем наличие UID
  const hasSession = uid && uid !== 'missing-uid' && uid.length > 5;
  
  const botName = ENV.TELEGRAM_BOT_NAME || "bot";

  return (
    <div className="flex min-h-screen items-center justify-center bg-tg-bg p-4">
      <div className="w-full max-w-[400px] text-center">
        
        {/* Анимация */}
        <div className="mx-auto mb-6 h-[150px] w-[150px] overflow-hidden rounded-2xl">
          <iframe
            src="https://tenor.com/embed/4989962420582851395"
            className="h-full w-full border-0 pointer-events-none"
            allowFullScreen
            title="Telegram Sticker"
          />
        </div>

        <h1 className="mb-2 text-2xl font-bold text-tg-text">
          Войти через Telegram
        </h1>

        <p className="mb-8 text-base leading-relaxed text-tg-text-secondary">
          Используйте ваш аккаунт Telegram для быстрого и безопасного входа.
        </p>

        <div className="space-y-4">
            {!hasSession ? (
                // СОСТОЯНИЕ 1: Нет сессии -> Кнопка старта
                // При клике: /api/oidc/auth -> редирект сюда же, но с ?uid=...
                <div className="animate-in fade-in zoom-in duration-300">
                    <StartFlowButton />
                    <p className="mt-4 text-xs text-gray-400">
                        Нажмите, чтобы начать процесс OIDC авторизации
                    </p>
                </div>
            ) : (
                // СОСТОЯНИЕ 2: Есть сессия (UID) -> Виджеты
                <div className="animate-in fade-in slide-in-from-bottom-4 duration-500 space-y-4">
                    <RealTelegramWidget botName={botName} uid={uid} />
                    
                    <div className="relative py-2">
                        <div className="absolute inset-0 flex items-center"><span className="w-full border-t border-gray-200"></span></div>
                        <div className="relative flex justify-center text-xs uppercase"><span className="bg-tg-bg px-2 text-gray-400">Dev Tools</span></div>
                    </div>
                    
                    <MockLoginButton targetUid={uid} />
                </div>
            )}
        </div>

      </div>
    </div>
  );
}