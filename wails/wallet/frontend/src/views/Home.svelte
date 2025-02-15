<script lang="ts">
    import ethLogo from '../assets/images/eth-logo.png';
    import { GetBalance } from '../../wailsjs/go/main/App';
    let balance: number = 0;
    let ethBalance: number = 180;
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
                <img src={ethLogo} alt="eth logo">
                <div class="coin-description-container">
                    <h6 class="coin-description-symbol">ETH</h6>
                    <h5 class="coin-description-name">Ethereum</h5>
                </div>
                <h3 class="coin-balance">{ethBalance}</h3>
            </div>
         </div>
    </div>
</main>


<style>

main {
        font-family: "Nunito", sans-serif;
        color: black;
        background-color: #f9f9f9; 
        height: 100vh;
        padding: 20px;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .balance-container h4 {
        font-weight: bold;
        margin-top: -5px;
        font-size: 1rem;
        margin-right: 300px;

    }

    .balance-container h2 {
        margin-right: 300px;
    }

    .balance-container {
        background: #fefefe;
        padding: 20px;
        border-radius: 12px;
        text-align: center;
        width: 90%;
        max-width: 400px;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    }

    .balance-buttons-container button {
        background: #0066cc;
        color: #fff;
        border: none;
        padding: 10px 15px;
        margin: 5px;
        width: 180px;
        justify-content: center;
        border-radius: 8px;
        cursor: pointer;
        transition: background 0.3s;
    }

    .balance-buttons-container button:hover {
        background: #004c99;
    }

    .assets-container {
        margin-top: 15px;
        width: 90%;
        padding: 5px 20px;
        max-width: 400px;
        margin: 25px;
        border-radius: 12px;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        background: #fefefe;
        transition: background 0.3s;
    }


    #assets-title {
        font-weight: bold;
        margin-right: 350px;
        margin-top: 5px;

    
    }

    .assets-list {
        padding: 1px;
        margin-top: -8px;
        border-radius: 12px;

    }

    .asset {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        padding: 1px;
        margin-top: 10px;
        border-bottom: 1.5px solid #ccc;
    }

    .asset img {
        width: 30px;
        height: 30px;
        margin-right: 10px;
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
        margin-left:240px;
    }

  
  

  </style>