const { mnemonicToEntropy } = require("ethereum-cryptography/bip39");
const { wordlist } = require("ethereum-cryptography/bip39/wordlists/english");
const { HDKey } = require("ethereum-cryptography/hdkey");
const { secp256k1 } = require("ethereum-cryptography/secp256k1");
const { keccak256 } = require("ethereum-cryptography/keccak");
const { bytesToHex } = require("ethereum-cryptography/utils");

async function restoreWallet(_mnemonic) {
  const entropy = mnemonicToEntropy(_mnemonic, wordlist);
  const hdRootKey = HDKey.fromMasterSeed(entropy);
  const privateKey = hdRootKey.deriveChild(0).privateKey;
  const publicKey = secp256k1.getPublicKey(privateKey);
  const address = keccak256(publicKey).slice(-20);
  return { address }
}
