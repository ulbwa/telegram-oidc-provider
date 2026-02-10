export interface LoginRequestData {
  auth: boolean;
  skip: boolean;
  user?: {
    id: number;
    first_name: string;
    last_name?: string;
    username?: string;
    photo_url?: string;
  };
  client: {
    name: string;
  };
  bot: {
    name: string;
    username: string;
    url: string;
  };
}