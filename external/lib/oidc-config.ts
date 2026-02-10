import { ENV } from "@/lib/env";
import Provider, { Configuration } from "oidc-provider";
import { exportJWK, generateKeyPair } from "jose";
import { findAccount } from "./oidc-func";

export const oidcConfig: Configuration = {
  clients: [
    {
      client_id: ENV.CLIENT_ID,
      client_secret: ENV.CLIENT_SECRET,
      grant_types: ["authorization_code"],
      redirect_uris: [
        "http://localhost:3000/api/oidc/debug",
        ENV.REDIRECT_URI || "",
      ].filter(Boolean),
      response_types: ["code"],
      scope: "openid profile email",
      token_endpoint_auth_method: "client_secret_basic",
    },
  ],
  adapter: undefined,
  findAccount: findAccount,
  interactions: {
    url(ctx, interaction) {
      return `/${interaction.prompt.name}?uid=${interaction.uid}`;
    },
  },
  cookies: {
    keys: [ENV.COOKIE_SECRET],
  },
  claims: {
    openid: ["sub"],
    email: ["email", "email_verified"],
    profile: ["name", "picture", "telegram_data"],
  },
  features: {
    devInteractions: { enabled: false },
  },
  jwks: {
    keys: [],
  },
};

const globalForOidc = global as unknown as { oidcProvider: Provider };

export async function getOidcProvider() {
  if (globalForOidc.oidcProvider) {
    console.log("PROVIDER CACHED")
    return globalForOidc.oidcProvider;
  }
    console.log("PROVIDER CREATED")

  const keypair = await generateKeyPair("RS256", { extractable: true });
  const jwk = await exportJWK(keypair.privateKey);

  const config = {
    ...oidcConfig,
    jwks: { keys: [{ ...jwk, kid: "sig-rs-01", use: "sig" }] },
  };

  const provider = new Provider(`${ENV.ISSUER_URL}/api/oidc`, config);

  // 1. Ошибки сервера (500)
  provider.on('server_error', (ctx, err) => {
    console.error('🔥 OIDC SERVER ERROR:', err);
    console.error('Context:', ctx.method, ctx.url);
  });

  // 2. Ошибки авторизации (когда клиент прислал что-то не то)
  provider.on('authorization.error', (ctx, error) => {
    console.warn('⚠️ Authorization Error:', error);
    console.warn('Details:', error.error_description);
  });

  // 3. Начало взаимодействия (редирект на /login)
  provider.on('interaction.started', (ctx, prompt) => {
    console.log('🔹 Interaction interaction.started:', prompt);
    console.log('   Prompt:', prompt.name);
  });

  // 4. Завершение взаимодействия
  provider.on('discovery.error', (ctx, result) => {
    console.log('✅ Interaction discovery.error:', result);
  });

  provider.callback()

  // if (ENV.isDev) {
  //   provider.proxy = true;
  // }
  globalForOidc.oidcProvider = provider;
  return provider;
}
