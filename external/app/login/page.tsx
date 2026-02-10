import { Button } from "@/components/ui/button";

export default function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-tg-bg p-4">
      <div className="w-full max-w-[400px] text-center">
        
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

        <Button
          className="group relative h-auto w-full min-w-[240px] justify-start rounded-full bg-tg-blue py-3 px-6 text-white shadow-lg shadow-tg-blue/30 transition-all hover:bg-tg-blue hover:-translate-y-0.5 hover:shadow-xl hover:shadow-tg-blue/40 active:scale-[0.98]"
        >
          <div className="mr-4 flex h-9 w-9 items-center justify-center rounded-full bg-white/10 group-hover:bg-white/20 transition-colors">
            <svg viewBox="0 0 24 24" fill="currentColor" className="h-6 w-6">
              <path d="M9.78 18.65l.28-4.23 7.68-6.92c.34-.31-.07-.46-.52-.19L7.74 13.3 3.64 12c-.88-.25-.89-.86.2-1.3l15.97-6.16c.73-.33 1.43.18 1.15 1.3l-2.72 12.81c-.19.91-.74 1.13-1.5.71L12.6 16.3l-1.99 1.93c-.23.23-.42.42-.83.42z" />
            </svg>
          </div>

          <div className="text-left">
            <div className="text-[13px] font-medium opacity-80 leading-tight">
              Войти как
            </div>
            <div className="text-[15px] font-bold">
              Андрей Пупкин
            </div>
          </div>
        </Button>
      </div>
    </div>
  );
}