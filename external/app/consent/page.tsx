import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { User, Hash } from "lucide-react";
import Link from "next/link";

export default function ConsentPage() {
  const appName = "Medincident";
  const appUrl = "sosiska.com";
  const botUserName = "@SosiskaBot";

  const refuseUrl = "/login";
  const acceptUrl = "/";

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-tg-bg p-4 font-sans text-foreground">
      <div className="absolute -left-[10%] -top-[10%] h-120 w-120 rounded-full bg-tg-blue/5 blur-3xl" />
      <div className="absolute -bottom-[10%] -right-[10%] h-96 w-96 rounded-full bg-blue-400/5 blur-3xl" />

      <Card className="relative z-10 w-full max-w-sm rounded-3xl border border-black/5 bg-tg-bg/75 shadow-xl backdrop-blur-xl">
        <CardContent className="p-8">
          <div className="mb-6 text-center">
            <iframe
              src="https://tenor.com/embed/8528851418612932470"
              className="pointer-events-none h-full w-full border-0 h-24 w-24 mx-auto antialiased"
              title="Logo Animation"
            />

            <h2 className="text-2xl font-bold text-foreground">{appName}</h2>

            <p className="mt-1 text-sm text-tg-text-secondary">
              {botUserName}
              {" • "}
              <Link
                href={appUrl}
                className="text-tg-blue no-underline hover:opacity-80 transition-opacity"
                target="_blank">
                {appUrl}
              </Link>
            </p>
          </div>

          <p className="mb-6 text-center text-base leading-normal text-foreground">
            Это приложение запрашивает доступ <br /> к вашим данным:
          </p>

          <div className="mb-8 overflow-hidden rounded-xl bg-tg-bg">
            <ScopeRow
              icon={<User size={20} />}
              title="Профиль"
              desc="Имя, юзернейм и фото"
            />
            <ScopeRow
              icon={<Hash size={20} />}
              title="Telegram ID"
              desc="Уникальный идентификатор"
              isLast
            />
          </div>

          <div className="flex gap-3">
            <Button
              variant="ghost"
              className="h-12 flex-1 rounded-xl text-base font-semibold text-tg-blue hover:bg-tg-blue/10 hover:text-tg-blue"
              asChild
            >
              <Link href={refuseUrl}>Отмена</Link>
            </Button>
            <Button
              className="h-12 flex-1 rounded-xl bg-tg-blue text-base font-semibold text-white shadow-none transition-all hover:bg-tg-blue/90 active:scale-[0.98]"
              asChild
            >
              <Link href={acceptUrl}>Разрешить</Link>
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function ScopeRow({
  icon,
  title,
  desc,
  isLast = false,
}: {
  icon: React.ReactNode;
  title: string;
  desc: string;
  isLast?: boolean;
}) {
  return (
    <div
      className={`flex items-center px-4 py-3.5 ${!isLast ? "border-b border-black/5" : ""}`}
    >
      <div className="mr-4 flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-tg-blue/10 text-tg-blue">
        {icon}
      </div>
      <div>
        <div className="mb-0.5 text-base font-semibold text-foreground leading-tight">
          {title}
        </div>
        <div className="text-sm text-tg-text-secondary leading-tight">
          {desc}
        </div>
      </div>
    </div>
  );
}
