<script lang="ts">
    import { GetAssets } from  '../../wailsjs/go/main/App'

    const tokens: object = {
        'ETH': 'Ethereum'
    };
    let balance: number = 0;
    
    type Asset = {
        balance: number;
        symbol: string;
        name: string;
        logoPath: string;
    };
    let assets: Asset[] = [];

    function getLogoPath(symbol: string): string {
        return `src/assets/logos/${symbol}.png`;
    }

    function initAssets(): void {
        const tokenSymbols = Object.keys(tokens);
        GetAssets(tokenSymbols)
        .then((assetsData) => {
            assets = Object.keys(assetsData).map((symbol) => (
                {
                    balance: assetsData[symbol],
                    symbol: symbol,
                    name: tokens[symbol],
                    logoPath: getLogoPath(symbol)
                }
            ))
        })
        .catch((error) => {
            alert("Error fetching assets: " + error);
        });
    }

    initAssets();
</script>

<main>
    <div class="balance-container">
        <h4>Total Balance</h4>
        <h2>$ {balance}</h2>
        <div class="balance-buttons-container">
            <button id="send-crypto-button">Send</button>
            <button id="receive-crypto-button">Receive</button>
        </div>
    </div>
    <div class="assets-container">
        <h4 id="assets-title">Assets</h4>
        <div class="assets-list">
            <div class="asset">
                {#each assets as asset}
                    <img src={asset.logoPath} alt={asset.symbol} />
                    <div class="coin-description-container">
                        <h6 class="coin-description-symbol">{asset.symbol}</h6>
                        <h5 class="coin-description-name">{asset.name}</h5>
                    </div>
                    <h3 class="coin-balance">{asset.balance}</h3>
                {/each}
         </div>
    </div>
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
    }
    .coin-description-name {
        font-size: 0.9rem;
        margin: 3px;
        color: #666;
    }

    .coin-balance {
        font-size: 1.1rem;
        font-weight: bold;
        margin-left: 63%;
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
    }
  
  

  </style>