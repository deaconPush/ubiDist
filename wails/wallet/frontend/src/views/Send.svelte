<script lang="ts">
    import { assets, currentView } from '../stores';
    import type { Asset } from '../types/index';
    import { ValidateAddress,EstimateGas, SendTransaction } from '../../wailsjs/go/main/App'    
    

    $: userAssets = $assets;
    let currentAsset: Asset; 
    let currentComponent: string = "Available Assets"
    let sendTokenTitle: string;
    let sendingAddress: string;
    let confirmedTransactionAmount: number;
    let gasPrice: number;
    let showPasswordModal: boolean = false;
    


    function clickCard(asset: Asset): void {
        currentAsset = asset;
        let symbol: string = asset.symbol;
        sendTokenTitle = `Send ${symbol}`;
        currentComponent = "Validate Address"
    }

    function validateAddress(): void {
        const inputComponent = document.getElementById("address-input") as HTMLInputElement;
        const addressValidationLabel = document.getElementById("address-validation-label") as HTMLParagraphElement;
        const addressInputButton = document.getElementById("address-input-button") as HTMLButtonElement;
        if (!inputComponent || !addressValidationLabel || !addressInputButton) {
            console.error("Error retrieving html components");
        }

        const address: string = inputComponent.value;
        const token: string = sendTokenTitle.split(" ")[1]

        if(address.length !== 0){
            ValidateAddress(address, token).
            then((ok: boolean) => {
                if(!ok){
                    addressValidationLabel.style.display = "block";
                    addressValidationLabel.textContent = "Address is invalid!"
                    addressInputButton.disabled = true;
                }
                else {
                    addressValidationLabel.textContent = ""
                    addressValidationLabel.style.display = "none";
                    addressInputButton.disabled = false;
                    sendingAddress = address;

                }
            })
            return;
        }

        addressValidationLabel.textContent = ""
        addressValidationLabel.style.display = "none";
    }

    function confirmAddress(): void {
        currentComponent = "Set Token Amount"
    }

    function validateAmount(e: Event): void {
        const target = e.target as HTMLInputElement;
        const amountValidationLabel = document.getElementById("amount-validation-label") as HTMLParagraphElement;
        const continueTransactionButton =  document.getElementById("continue-transaction-button") as HTMLButtonElement;
        if(!target || !amountValidationLabel || !continueTransactionButton) {
            alert("HTML elements are not valid");
            return;
        }

        const maximumAmount: number = currentAsset.balance;
        const value: string = target.value;
        
        if (value === "") {
            target.value = "0";
            return;
        }

        if (!/^\d*\.?\d*$/.test(value)) {
            target.value = value.slice(0, -1);
            return;
        }

         if (/^0\d+/.test(value)) {
            target.value = value.replace(/^0+/, '');
            return;
        }

        const amount: number = parseFloat(value) as number;

        if (amount > maximumAmount) {
            amountValidationLabel.textContent = "You are not allowed to exceed your balance"
            amountValidationLabel.style.display = "block";
            continueTransactionButton.disabled = true;
            return;
        }

        if (amount > 0 && amount < maximumAmount){
            amountValidationLabel.textContent = ""
            amountValidationLabel.style.display = "none";
            continueTransactionButton.disabled = false;
            return;
        }

        amountValidationLabel.textContent = ""
            amountValidationLabel.style.display = "none";
            continueTransactionButton.disabled = true;
        }

        function confirmAmount(): void {
            const continueTransactionButton =  document.getElementById("continue-transaction-button") as HTMLButtonElement;
            const transactionAmountInput = document.getElementById("amount-input") as HTMLInputElement;
            if(!continueTransactionButton || !transactionAmountInput) {
                alert("HTML elements are not valid");
                return;
            }

            const transactionAmount: string = transactionAmountInput.value;
            confirmedTransactionAmount = parseFloat(transactionAmount);
            EstimateGas(currentAsset.symbol, sendingAddress, transactionAmount)
            .then((gas: string) => {
                gasPrice = parseFloat(gas);
                currentComponent = "Confirm Transaction";            
            })
            .catch((error) => alert("Error estimating gas price: " + error))
            
        }

        function confirmTransaction(): void {
            const confirmTransactionButton = document.getElementById("confirm-transaction-button") as HTMLButtonElement;

            if (!confirmTransactionButton) {
                alert("HTML element is not valid");
                return;
            }

            showPasswordModal = true;
        }

        function sendTransaction(e: Event): void {
            e.preventDefault();
            const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
            if (!passwordInput) {
                alert("HTML element is not valid");
                return;
            }

            const password: string = passwordInput.value;
            SendTransaction(currentAsset.symbol, password, sendingAddress, confirmedTransactionAmount.toString())
            .then((ok: boolean) =>{
                if (ok){
                    alert("Successful transaction YAY!")
                }
                
                currentView.set("Home");
            }  )
            .catch((error) => {
                alert("Error processing transaction: " + error)
            })

        }

    
</script>


<main>
    {#if currentComponent === "Validate Address"}
        <h3 id="send-token-title">{sendTokenTitle}</h3>
        <div class="input-group">
            <input id="address-input" type="text" on:input={validateAddress} placeholder="Enter address" />
            <p id="address-validation-label"></p>
        </div>
        
        <button id="address-input-button" disabled on:click={confirmAddress}>Continue</button>
    {:else if currentComponent === "Available Assets" }
        <div class="assets-container">
            <h3>Available assets</h3>
            <div class="assets-list">
                    {#each userAssets as asset}
                        {#if asset.balance > 0}
                            <div class="asset" on:click={() => clickCard(asset)}>
                                <img src={asset.logoPath} alt={asset.symbol} />
                
                            <div class="coin-description-container">
                                <h6 class="coin-description-symbol">{asset.symbol}</h6>
                                <h5 class="coin-description-name">{asset.name}</h5>
                            </div>
                        <h3 class="coin-balance">{asset.balance}</h3>
                    </div>
                        {/if}
                    {/each}
                </div>
            </div>
    {:else if currentComponent === "Set Token Amount" }
        <h3 id="send-token-title">{sendTokenTitle}</h3>
        <input id="amount-input" type="text" autofocus inputmode="decimal" pattern="^\d*\.?\d*$" value="0" on:input={validateAmount}/>
        <h6 id="amount-symbol">{currentAsset.symbol}</h6>
        <h6>Available: {currentAsset.balance} {currentAsset.symbol}</h6>
        <div class="address-details">
            <h6 id="address-box-label">Address </h6>
            <h6 id="address-box">{sendingAddress.slice(0, 4) + "..." + sendingAddress.slice(-5)}</h6>
        </div>
        <button id="continue-transaction-button" disabled=true on:click={confirmAmount}>Continue</button>
        <p id="amount-validation-label"></p>
    {:else if currentComponent === "Confirm Transaction"}
        <h3 id="send-token-title">{sendTokenTitle}</h3>
        <div class="confirm-transaction-container">
            <h3>Confirm your transaction</h3>
            <h3>You are about to send {confirmedTransactionAmount} {currentAsset.symbol}</h3>
            <h4>Cost of the network: {gasPrice.toPrecision(4)} {currentAsset.symbol}</h4>
        </div>
        <button id="confirm-transaction-button" on:click={confirmTransaction}>Confirm</button>
        {#if showPasswordModal === true} 
            <div class="overlay">
                <form class="form-container">
                    <div class="input-group">
                    <label for="password">To proceed with the transaction you need to sign it with your password.</label>
                    <input type="password" id="wallet-password" />
                    <p id="password-validation-label"></p>
                    <button id="send-transaction-button" on:click={sendTransaction}>Send</button>
                </div>
            </div>
        {/if}
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

    .assets-container {
        margin-top: 3vh;
        width: 75%;
        padding: 1% 4%;
        margin: 5%;
        border-radius: 3vh;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        background: #fefefe;
        transition: background 0.3s;
    }

    .assets-list {
        padding: 0.2%;
        margin-top: -2%;
        border-radius: 3vh;
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

    .asset {
        display: flex;
        align-items: center;
        justify-content: flex-start;
        padding: 0.2%;
        margin-top: 2vh;
        border-bottom: 0.7px solid #ccc;
    }

    .asset img {
        width: 10%;
        height: 6vh;
        margin-right: 3%;
    }
    
    .input-group {
        display: flex;
        flex-direction: column;
        gap: 2vh;
        align-items: center;
    }

    #address-input {
        padding: 3%;
        border-radius: 4vh;
        border: 5px solid #ccc;
    }


    #address-validation-label {
        font-size: 0.9rem;
        color: red;
        display: none;
    }

    #address-input-button {
        margin-top: 15%;
        width: 40%;
        border-radius: 2vh;
        background-color: #007bff;
        color: white;
        cursor: pointer;
    }

    #address-input-button:disabled {
        background-color: #ccc;
        color: #666;
        border: none;   
        cursor: not-allowed;
    }



    #amount-input {
        font-size: 2rem;
        height: 20vh;
        width: 35%;
        margin-left: 20%;
        background-color: #f9f9f9;
        outline: none;
        border: none;
        box-shadow: none;   
    }

    #amount-symbol {
        margin-left: 10%;
        margin-top: -5%;
    }

    .address-details {
        display: flex;
        flex-direction: row;
        justify-content: space-evenly;
        border-bottom: 0.5px solid #ccc;
        height: 8vh;
        gap: 40%;
    }

    #address-box-label {
        align-self: flex-start;

    }

    #continue-transaction-button {
        height: 5vh;
        margin-top: 10%;
        border-radius: 2vh;
        width: 30%;
        background-color: #007bff;
        color: white;
        cursor: pointer;
    }
    
    #continue-transaction-button:disabled {
        background-color: #ccc;
        color: #666;
        border: none;   
        cursor: not-allowed;
    }


    #amount-validation-label {
        font-size: 0.9rem;
        color: red;
        display: none;
    }

    .confirm-transaction-container {
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        background: #fefefe;
        transition: background 0.3s;
        border-radius: 3vh;
        width: 65%;
    }

    #confirm-transaction-button {
        width: 60%;
        margin-top: 5%;
        border-radius: 2vh;
        font-size: 1.1rem;
        height: 5vh;
        background-color: #007bff;
        color: white;
        cursor: pointer;
    }

    .overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0, 0, 0, 0.6);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 999;
    }

    .form-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        align-content: center;
        height: 25vh;
        width: 65%;
        border-radius: 4vh;
        padding-top: 6%;
        background-color: #fff;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
    }

    #password-validation-label {
        font-size: 0.9rem;
        color: red;
        display: none;
    }

    #wallet-password {
        padding: 1.5%;
        font-size: 1rem;
        border: 1px solid #ccc;
        border-radius: 2.5vh;
        color: #333;
        background-color: #fefefe;
        height: 3vh;
        width: 60%;
    }

    #send-transaction-button {
        margin-top: 2%;
        width: 30%;
        height: 5vh;
        border-radius: 2vh;
        background-color: #007bff;
        color: white;
        cursor: pointer;
    }

    


    @media (min-width: 1800px) {

        .assets-container {
            margin-top: 3%;
            width: 38%;
        }
    }

    @media (max-width: 600px) {


        .assets-container {
            margin-top: 3%;
            width: 88%;
        }

        #address-input {
            margin-top: 8%;
            width: 300px;
            height: 3vh;

        }

        #address-input-button {
            font-size: 1.1rem;
            width: 200px;
            height: 5vh;

        }

    }
  
    
</style>

