<script lang="ts">
    import { GetAssets, GetTransactions } from  '../../wailsjs/go/main/App'
    import { currentView, assets, selectedAccounts } from '../stores';
    import type { Asset, Transaction } from '../types/index';

    const tokens: object = {
        'ETH': 'Ethereum'
    };    
    let balance: number = 0;
    let assetsArray: Asset[] = [];
    let tokenAccounts = $selectedAccounts;
    let displayTransactions: boolean = false;
    let walletTransactions: Transaction[];
    let lastTransaction: Transaction;
    let lastTransactionDate: string;
    let showDropdown: boolean = false;
    let dropdownRef: HTMLDivElement;

    function getLogoPath(symbol: string): string {
        return `src/assets/logos/${symbol}.png`;
    }

    function initAssets(): void {
        GetAssets(tokenAccounts)
        .then((assetsData) => {
            assetsArray = Object.keys(assetsData).map((symbol) => (
                {
                    balance: assetsData[symbol]["balance"],
                    symbol: symbol,
                    name: tokens[symbol],
                    logoPath: getLogoPath(symbol),
                    selectedAccount: tokenAccounts[symbol],
                    accounts: assetsData[symbol]["accounts"],
                }
            ))
            assets.set(assetsArray);
        })
        .catch((error) => {
            alert("Error fetching assets: " + error);
        });
    }

    function sendCrypto(): void {
        currentView.set('Send');
    }

    function getTransactions(): void {
        GetTransactions()
        .then((transactions: Transaction[]) => {
            if(transactions !== null) {
                walletTransactions = transactions;
                lastTransaction = walletTransactions[0]
                const date = new Date(lastTransaction.createdAt).toDateString().split(' ');
                lastTransactionDate = `${date[2]} of ${date[1]} ${date[3]}`
                displayTransactions = true;
            }
        })
        .catch((err) => {
            alert("Error retrieving transactions " + err);
        })
    }

    function listTokenAccounts(): void {
        showDropdown = true;
    }

    function selectAccount(asset: Asset, key: string) {
        tokenAccounts[asset.symbol] = Number(key);
        selectedAccounts.set(tokenAccounts)
        GetAssets(tokenAccounts)
        .then((assetsData) => {
            assetsArray = Object.keys(assetsData).map((symbol) => (
                {
                    balance: assetsData[symbol]["balance"],
                    symbol: symbol,
                    name: tokens[symbol],
                    logoPath: getLogoPath(symbol),
                    selectedAccount: Number(key),
                    accounts: assetsData[symbol]["accounts"],
                }
            ))
            assets.set(assetsArray);
            showDropdown = false;
        })
        .catch((error) => {
            alert("Error fetching assets: " + error);
        });
  }

    
  function onWindowClick(event: MouseEvent): void {
        if (dropdownRef.contains(event.target as Node) == false) {
            showDropdown = false;
        }
    }

    initAssets();
    getTransactions();
</script>

<svelte:window on:click={onWindowClick} />
<main>
        <div class="balance-container">
            <h4>Total Balance</h4>
            <h2>$ {balance}</h2>
            <div class="balance-buttons-container">
                <button 
                id="send-crypto-button"
                on:click={sendCrypto}
                >Send</button>
                <button id="receive-crypto-button">Receive</button>
            </div>
        </div>
        <div class="assets-container">
            <h4 id="assets-title">Assets</h4>
            <div class="assets-list">
                <div class="asset">
                    {#each assetsArray as asset, index (index)}
                        <img src={asset.logoPath} alt={asset.symbol} />
                        <div class="coin-description-container">
                            <h6 class="coin-description-symbol">{asset.symbol}</h6>
                            <h5 class="coin-description-name">{asset.name}</h5>
                        </div>
                        <div class="account-details">
                            <div class="account-dropdown" bind:this={dropdownRef}>
                                <button   on:click={() => listTokenAccounts()}
                                    class="account-dropdown-button">{asset.accounts[asset.selectedAccount].slice(0, 4) + "..." + asset.accounts[asset.selectedAccount].slice(-3)}</button>
                                    {#if showDropdown}
                                        <ul class="account-dropdown-content">
                                            {#each Object.entries(asset.accounts) as [key, account], index (index)}
                                            <li
                                            on:click={() => selectAccount(asset, key)}
                                            tabindex="0"
                                            >
                                                {account.slice(0, 6) + "..." + account.slice(-4)}
                                            </li>
                                            {/each}
                                        </ul>
                                    {/if}
                            </div>
                            <h3 class="">{asset.balance.toFixed(2)}</h3>
                        </div>
                        
                    {/each}
                </div>
            </div>
        </div>
        {#if displayTransactions == true }
        <div class="last-transaction-container">
            <h5 class="last-transaction-container-title">Last Transaction</h5>
            <div class="last-transaction-item">
                <img src="{getLogoPath(lastTransaction.token)}" alt={lastTransaction.token}/>
                <div class="last-transaction-description">
                    <h3>Withdrawal of {lastTransaction.token}</h3>
                    <h5>{lastTransactionDate}</h5>
                </div>
                
                <h4 class="last-transaction-value">{lastTransaction.value} {lastTransaction.token}</h4>
            </div>
        </div>
        {/if}
    </main>
    
<style>

main {
        font-family: "Nunito", sans-serif;
        color: black;
        background-color: #f9f9f9; 
        height: 100vh;
        padding: 5%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .balance-container h4 {
        font-weight: bold;
        margin-top: -1vh;
        font-size: 1rem;
        margin-right: 75%;

    }

    .balance-container h2 {
        margin-right: 85%;
        color: #002855;
    }

    .balance-container {
        background: #fefefe;
        padding: 5%;
        border-radius: 3vh;
        text-align: center;
        width: 73%;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    }

    .balance-buttons-container button {
        background: #0066cc;
        color: #fff;
        border: none;
        padding: 3% 2%;
        margin: 1vh;
        width: 45%;
        justify-content: center;
        border-radius: 2vh;
        cursor: pointer;
        transition: background 0.3s;
    }

    .balance-buttons-container button:hover {
        background: #004c99;
    }

    .assets-container {
        margin-top: 5vh;
        width: 75%;
        padding: 1% 4%;
        margin: 5%;
        border-radius: 3vh;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        background: #fefefe;
        transition: background 0.3s;
    }


    #assets-title {
        font-weight: bold;
        margin-right: 85%;
        margin-top: 1vh;

    
    }

    .assets-list {
        padding: 0.2%;
        margin-top: -3%;
        border-radius: 3vh;
    }

    .asset {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        padding: 0.2%;
        margin-top: 2vh;
        border-bottom: 1.5px solid #ccc;
    }

    .asset img {
        width: 10%;
        height: 6vh;
        margin-right: 3%;
    }
    
    .coin-description-container {
        margin-top: 15px;
        display: block;
    }

    .coin-description-symbol {
        position: relative;
        font-size: 0.8rem;
        margin: 0;
        font-weight: bold;
        color: #002855;
    }
    .coin-description-name {
        font-size: 0.9rem;
        margin: 3px;
        color: #666;
    }

    .account-details {
        font-size: 1.1rem;
        font-weight: bold;
        margin-left: 63%;
        margin-top: 3%;
        color: #002855;
    }

    .account-dropdown-button {
        background-color: #E6F0FF;
        color: #004085;
        font-weight: bold;
        border-radius: 3vh;
        height: 5vh;
    }

    .account-dropdown-content{
        position: absolute;
        background-color: white;
        border: 0.3vh solid #ddd;
        border-radius: 2vh;
        z-index: 1000;
        list-style: none;
    }

    .account-dropdown-content li {
        cursor: pointer;
        transition: background-color 0.2s ease;
        font-size: 0.9rem;
        line-height: 2rem;
        color: #333;
    }

    .account-dropdown-content li:hover {
        background-color: #f5f5f5;
    }


    .last-transaction-container {
        background: #fefefe;
        padding: 5%;
        border-radius: 3vh;
        text-align: center;
        width: 80%;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    }

    .last-transaction-container-title {
        font-weight: bold;
        margin-top: -3%;
        font-size: 1rem;
        margin-right: 75%;
    }

    .last-transaction-item {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-top: 5%;
    }

    .last-transaction-description {
        width: 70%;
        margin-top: -3%;
        margin-right: 35%;
    }

    .last-transaction-description h3 {
        font-size: 0.9rem;
        color: #002855;
    }

    .last-transaction-description h5 {
        font-size: 0.7rem;
        margin-top: -1%;
        margin-right: 30%;
        color: #4a4a4a;
    }


    .last-transaction-item img {
        width: 13%;
        height: 7vh;
        margin-top: -3%;
    }

    .last-transaction-value {
        margin-top: -3%;
        font-size: 1.1em;
        color: #002855;
        width: 50%;
    }

    .last-transaction-container {
            margin-top: 3%;
            height: 8vh;
            width: 74%;
        }



    @media (min-width: 1800px) {
        .balance-container {
            width: 40%;
            padding: 3%;
        }

        .assets-container {
            margin-top: 3%;
            width: 38%;
        }

        .account-dropdown-button {
            width: 350px;
        }   
    }

    
    @media (min-width: 600px) and (max-width: 1800px) {
        .balance-container {
            width: 90%;
            padding: 3%;
        }

        .assets-container {
            margin-top: 8%;
            width: 88%;
        }

        .account-details {
            margin-left: 52%;
        }

        .last-transaction-container {
            margin-top: 3%;
            height: 12vh;
            width: 87%;
        }

        .account-dropdown-button {
            width: 130px;
        }

        .account-dropdown-content {
            margin-top: -5%;
        }

        .account-dropdown-content li {
            text-align: left;
            margin-left: -14%;   
        }
    }
  

    @media (max-width: 600px) {
        .balance-container {
            width: 90%;
            padding: 3%;
        }

        .assets-container {
            margin-top: 8%;
            width: 88%;
        }

        .account-details {
            margin-left: 52%;
        }

        .last-transaction-container {
            margin-top: 3%;
            height: 16vh;
            width: 90%;
        }

        .account-dropdown-button {
            width: 115px;
        }

        .account-dropdown-content {
            margin-left: -3%;
            margin-top: -5%;
        }

        .account-dropdown-content li {
            text-align: left;
            margin-left: -15%;   
        }
    }
  
  

  </style>