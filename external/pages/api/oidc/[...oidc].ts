import { NextApiRequest, NextApiResponse } from "next";
import { getOidcProvider } from "@/lib/oidc-config";
import { TelegramService } from "@/services/telegram.service";
import { UserStore } from "@/services/user.store";

export const config = {
  api: {
    bodyParser: false,
    externalResolver: true,
  },
};

const parseBody = async (req: NextApiRequest) => {
  const buffers = [];
  for await (const chunk of req) {
    buffers.push(chunk);
  }
  const data = Buffer.concat(buffers).toString();
  return data ? JSON.parse(data) : {};
};

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse,
) {
  const provider = await getOidcProvider();
  console.log(`🔍 DEBUG REQ: ${req.method} ${req.url}`);
  // Если req.body существует, значит парсер НЕ отключился, и это причина ошибки
  if ((req as any).body) {
    console.error(
      "❌ ОШИБКА: Body Parser не отключен! Удалите папку .next и перезапустите сервер.",
    );
  }
  console.log("   - Req URL:", req.url); // Должен быть /api/oidc/.well-known/openid-configuration
  console.log("   - Provider Issuer:", provider.issuer); // Должен быть http://localhost:3000/api/oidc
  console.log("   Headers[authorization]:", req.headers["authorization"]);
  console.log("   Content-Type:", req.headers["content-type"]);

  const { oidc } = req.query;

  // ---------------------------------------------------------------------------
  // БЛОК 1: Кастомные действия (Interactions: Login / Consent)
  // ---------------------------------------------------------------------------
  if (
    req.method === "POST" &&
    Array.isArray(oidc) &&
    oidc[0] === "interaction"
  ) {
    const uid = oidc[1];
    const action = oidc[2];

    try {
      // === SCENARIO: LOGIN ===
      if (action === "login") {
        // Читаем тело вручную!
        const telegramData = await parseBody(req);

        const user = TelegramService.validateWidgetData(telegramData);

        if (!user) {
          res.status(403).json({ error: "Invalid Telegram Hash" });
          return;
        }

        const isNewUser = !UserStore.has(user.id.toString());
        if (isNewUser) {
          console.log(`✨ [REG] New User: ${user.id}`);
        }

        UserStore.set(user.id.toString(), user);

        const result = {
          login: { accountId: user.id.toString() },
        };

        await provider.interactionFinished(req, res, result, {
          mergeWithLastSubmission: false,
        });
        return;
      }

      // === SCENARIO: CONFIRM ===
      if (action === "confirm") {
        const interaction = await provider.Interaction.find(uid);
        if (!interaction) {
          res.status(404).json({ error: "Interaction not found" });
          return;
        }

        const grant = new provider.Grant({
          accountId: interaction.session?.accountId,
          clientId: interaction.params.client_id as string,
        });

        const details = interaction.prompt.details as any;
        if (details.missingOIDCScope) {
          grant.addOIDCScope(details.missingOIDCScope.join(" "));
        }

        const grantId = await grant.save();
        const result = { consent: { grantId } };

        await provider.interactionFinished(req, res, result, {
          mergeWithLastSubmission: true,
        });
        return;
      }

      // === SCENARIO: ABORT ===
      if (action === "abort") {
        const result = {
          error: "access_denied",
          error_description: "End-User aborted interaction",
        };
        await provider.interactionFinished(req, res, result, {
          mergeWithLastSubmission: false,
        });
        return;
      }
    } catch (err) {
      console.error("Interaction Error:", err);
      if (!res.headersSent) {
        res.status(500).json({ error: "Interaction processing failed" });
      }
      return;
    }
  }

  // ---------------------------------------------------------------------------
  // БЛОК 2: Стандартные OIDC запросы
  // ---------------------------------------------------------------------------
  await provider.callback()(req, res);
}
