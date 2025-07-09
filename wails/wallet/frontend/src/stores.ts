import { writable } from 'svelte/store';
import { Asset, ProviderMap } from './types';

export const currentView = writable('Wallet Recovery');
export const assets = writable<Asset[]>([]);
export const availableTokens = writable<string[]>(['ETH']);
export const selectedAccounts = writable<Record<string, number>>({
  ETH: 0,
});
export const selectedProviders = writable<ProviderMap>({
  ETH: "http://localhost:8545",
});
