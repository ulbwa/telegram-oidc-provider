const getEnv = (key: string, defaultValue?: string): string => {
  const value = process.env[key] || defaultValue;
  if (value === undefined) {
    throw new Error(`❌ Ошибка конфигурации: Переменная окружения "${key}" не найдена.`);
  }
  return value;
};

export const ENV = {
  NODE_ENV: getEnv('NODE_ENV', 'development'),
  ISSUER_URL: getEnv('ISSUER_URL', 'http://localhost:3000'), 
  
  TELEGRAM_BOT_TOKEN: getEnv('TELEGRAM_BOT_TOKEN'),
  TELEGRAM_BOT_NAME: getEnv('TELEGRAM_BOT_NAME'),

  CLIENT_ID: getEnv('CLIENT_ID'),
  CLIENT_SECRET: getEnv('CLIENT_SECRET'),
  REDIRECT_URI: getEnv('REDIRECT_URI'),
  
  COOKIE_SECRET: getEnv('COOKIE_SECRET', 'super-secret-secret-must-be-very-long'),
  
  get isDev() {
    return this.NODE_ENV === 'development';
  },
};