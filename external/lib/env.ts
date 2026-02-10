const getEnv = (key: string, defaultValue?: string): string => {
  const value = process.env[key] || defaultValue;
  if (value === undefined) {
    throw new Error(`❌ Ошибка конфигурации: Переменная окружения "${key}" не найдена.`);
  }
  return value;
};

export const ENV = {
  NODE_ENV: getEnv('NODE_ENV', 'development'),
  
  TELEGRAM_BOT_TOKEN: getEnv('TELEGRAM_BOT_TOKEN'),
  TELEGRAM_BOT_NAME: getEnv('TELEGRAM_BOT_NAME'),
  SERVER_API_URL: getEnv("SERVER_API_URL"),
  
  get isDev() {
    return this.NODE_ENV === 'development';
  },
};