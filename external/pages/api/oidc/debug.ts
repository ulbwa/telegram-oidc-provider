import { NextApiRequest, NextApiResponse } from 'next';

export default function handler(req: NextApiRequest, res: NextApiResponse) {
  res.status(200).json({ 
    status: "success",
    message: "Authorization Code получен! OIDC Flow завершен успешно.",
    query_params: req.query,
    note: "Теперь клиентское приложение должно обменять этот 'code' на токены."
  });
}