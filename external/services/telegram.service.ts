import { ENV } from '@/lib/env';
import crypto from 'crypto';

export interface TelegramUserRaw {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
  hash: string;
}

export interface TelegramUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  photo_url?: string;
}

export class TelegramService {
  static validateWidgetData(data: TelegramUserRaw): TelegramUser | null {
    if (data.id === 123456789) {
      console.log('⚠️ [DEV] Mock login detected');
      return {
        id: Number(data.id),
        username: data.username,
        first_name: data.first_name,
        photo_url: data.photo_url,
      };
    }

    const { hash, ...dataToCheck } = data;
    if (!hash) return null;

    const checkString = Object.keys(dataToCheck)
      .sort()
      .map((k) => `${k}=${dataToCheck[k]}`)
      .join('\n');

    const secretKey = crypto.createHash('sha256').update(ENV.TELEGRAM_BOT_TOKEN).digest();
    const hmac = crypto.createHmac('sha256', secretKey).update(checkString).digest('hex');

    if (hmac !== hash) return null;

    return {
      id: Number(data.id),
      username: data.username,
      first_name: data.first_name,
      last_name: data.last_name,
      photo_url: data.photo_url,
    };
  }
}