// telegram-oidc-provider/src/services/account.service.ts
import { UserStore } from '@/services/user.store';
import { Account } from 'oidc-provider';

export const findAccount = async (ctx: any, id: string): Promise<Account | undefined> => {
    return {
      accountId: id,
      async claims(use, scope) {
        const user = UserStore.get(id);
        return {
          sub: id,
          email: `tg_${id}@telegram.placeholder`,
          email_verified: true,
          telegram_data:
            user ?
              {
                id: user.id,
                username: user.username,
                first_name: user.first_name,
                last_name: user.last_name,
                photo_url: user.photo_url,
                language_code: "ru",
              }
            : {},
          name: user ? [user.first_name, user.last_name].filter(Boolean).join(" "): "Telegram User",
          picture: user?.photo_url,
        };
      },
    };
  };