export type Asset = {
  balance: number;
  symbol: string;
  name: string;
  logoPath: string;
  accounts: AccountMap;
  selectedAccount: number;
};

export type AccountMap = {
  [key: number]: string;
};

export type Transaction = {
  sender: string;
  recipient: string;
  status: string;
  value: string;
  token: string;
  createdAt: string;
};
