const ethers = require('ethers');  
const crypto = require('crypto');

function generatePrivateKey() {
    const id = crypto.randomBytes(32).toString('hex');
    return "0x"+id;
}

function generateWallet(privateKey) {
    return new ethers.Wallet(privateKey);
}

module.exports = {
    generatePrivateKey,
    generateWallet
}