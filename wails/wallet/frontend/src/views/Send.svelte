<script lang="ts">
    import { assets } from '../stores';
    import type { Asset } from '../types/index';
    import { ValidateAddress } from '../../wailsjs/go/main/App'    

    $: userAssets = $assets;
    let currentComponent: string = "Available Assets"
    let sendTokenTitle: string;


    function clickCard(asset: Asset): void {
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

                }
            })
            return;
        }

        addressValidationLabel.textContent = ""
        addressValidationLabel.style.display = "none";
    }


</script>


<main>
    {#if currentComponent === "Validate Address"}
        <h3 id="send-token-title">{sendTokenTitle}</h3>
        <div class="input-group">
            <input id="address-input" type="text" on:input={validateAddress} placeholder="Enter address" />
            <p id="address-validation-label"></p>
        </div>
        
        <button id="address-input-button" disabled on:click={() => alert('Address submitted')}>Continue</button>
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

