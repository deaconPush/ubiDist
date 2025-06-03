export type Asset = {
    balance: number;
    symbol: string;
    name: string;
    logoPath: string;
    accountIndex: number;
}

export type Transaction = {
    sender: string;
    recipient: string;
    status: string;
    value: string;
    token: string; 
    createdAt: string; 
}