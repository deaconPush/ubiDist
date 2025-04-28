import { writable } from "svelte/store";
import type { Asset } from "./types";

export const currentView = writable("Wallet Recovery");
export const assets = writable<Asset[]>([]);