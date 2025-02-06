<script lang="ts">
    let balance: number = 0;
    
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
    <!-- <div class="assets-container">
        <h1>Your assets</h1>
        <div class="assets-list">
            <div class="asset">
                <img src={ethLogo} alt="eth logo">
                <h2>Ethereum</h2>
                <h3>{ethBalance}</h3>
            </div>
        </div> -->
</main>


<style>

main {
        font-family: "Nunito", sans-serif;
        color: #333; 
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
        border-radius: 8px;
        cursor: pointer;
        transition: background 0.3s;
    }

    .balance-buttons-container button:hover {
        background: #004c99;
    }

  </style>