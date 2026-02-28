import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function ErrorPage() {
  const title = "Все сломалось мы все обречены!";
  const description = "Произошло что то очень плохое, я это чувствую";
  const redirectUrl = "/login";

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-tg-bg p-4 text-center">
      <div className="mb-6 h-[150px] w-[150px] overflow-hidden rounded-2xl">
        <iframe
          src="https://tenor.com/embed/8909202877275126207"
          className="h-full w-full border-0 pointer-events-none"
          title="Error Animation"
        />
      </div>

      <h1 className="mb-2 text-2xl font-bold text-tg-text">
        {title}
      </h1>

      <p className="mx-auto mb-8 max-w-xs text-base leading-relaxed text-tg-text-secondary">
        {description}
      </p>

      <Button
        variant="ghost"
        className="text-tg-blue hover:bg-tg-blue/10 hover:text-tg-blue"
        asChild
      >
        <Link href={redirectUrl}>Вернуться обратно</Link>
      </Button>
    </div>
  );
}
