<script lang="ts">
    import { ValidateMnemonic } from "../../wailsjs/go/main/App";

    export let seedPhraseList: string[] = [];
    export let onConfirm: () => void = () => {};
    let blocksIndexToConfirm: number[] = getRandomNumbers(seedPhraseList.length);
    

    function getRandomNumbers(seedLength: number) : number[] {
        const numbers: Set<number> = new Set();
        while (numbers.size < seedLength / 3) {
            numbers.add(Math.floor(Math.random() * seedLength));
        }

        return Array.from(numbers);
    }

    function checkSeedInputs() : void {
        const inputs = document.querySelectorAll('.seed-phrase-block') as NodeListOf<HTMLInputElement>;
        if(!inputs) {
            console.error('Seed phrase inputs not found');
            return;
        }

        const confirmButton = document.getElementById('confirm-seed-button') as HTMLButtonElement;
        if (!confirmButton) {
            console.error('Confirm button not found');
            return;
        }

        const validationLabel = document.getElementById('seed-phrase-validation') as HTMLParagraphElement;
        if (!validationLabel) {
            console.error('Validation label not found');
            return;
        }

        const allFilled = Array.from(inputs).every(input => input.value.trim() !== '');
        if(allFilled) {
            const seedPhrase: string = Array.from(inputs).map(input => input.value.trim()).join(' ');
            ValidateMnemonic(seedPhrase).then(isValid => {
                if(isValid) {
                    confirmButton.disabled = false;
                    validationLabel.style.display = 'none';
                } else {
                    validationLabel.style.display = 'block';
                    validationLabel.textContent = 'Invalid seed phrase';
                }
            });
        }
    } 
</script>

<div class="header-content">
    <h4>Confirm your Secret Recovery Phrase</h4>
</div>
<div class="seed-words-container">
    <div class="seed-words-container">
        { #each seedPhraseList as word, index }
            <input type="text" on:input={checkSeedInputs} class="seed-phrase-block" value={blocksIndexToConfirm.includes(index) ? '' : word} readonly={!blocksIndexToConfirm.includes(index) || null} />
        { /each }
    </div>
</div>
<p id="seed-phrase-validation" class="validation-label" ></p>
<div class="confirm-seed-button-container">
    <button id="confirm-seed-button" on:click={onConfirm} disabled=true>Restore wallet</button>
</div> 

<style>
h4 {
    margin-top: 20px;
    color: #333; 
    text-align: center; 
    padding: 0 10px;
}

.seed-words-container {
    display: grid;
    grid-template-columns: repeat(3, minmax(120px, 1fr)); 
    column-gap: 50px; 
    row-gap: 15px; 
    justify-items: center;
    align-items: space-between;
    margin-bottom: 20px;
    width: 100%;
    max-width: 480px; 
}

.seed-phrase-block {
    width: 65%; 
    font-size: 15px;
    background-color: #fff; 
    border: 1px solid #ccc;
    border-radius: 7px;
    height: 40px;
    font-family: "Nunito";
    padding: 0 35px 0 10px;
    color: #333; 
}

.seed-phrase-block::placeholder {
    color: #999; 
}


.confirm-seed-button-container {
    display: flex;
    justify-content: center;
    width: 100%;
}

#confirm-seed-button {
    padding: 12px 20px;
    margin-top: 20px;
    background-color: #0066cc;
    max-width: 300px;
    font-family: "Nunito";
    font-weight: bold;
    color: #fff; 
    border-radius: 20px;
    border: none;
    width: 100%;
    cursor: pointer;
    transition: background-color 0.3s ease;
}


#confirm-seed-button:hover:not(:disabled) {
    background-color: #004c99; 
}

#confirm-seed-button:disabled {
    background-color: #ccc; 
    color: #666;
    cursor: not-allowed;
}

.validation-label {
    margin-top: -8px;
    background-color: #ffdddd; 
    border-left: 4px solid #d9534f;
    font-size: 0.9rem;
    color: #a94442; 
    padding: 8px 12px;
    border-radius: 5px;
    display: none;
    width: fit-content;
    max-width: 100%;
    transition: opacity 0.3s ease-in-out;
}

@media (min-width: 1500px) {
    .header-content  {
        margin-top: 1.2%;
    }

    .header-content h4 {
        font-size: 1.4rem;
    }

    .seed-words-container {
        grid-template-columns: repeat(3, minmax(58%, 1fr));
        column-gap: 20%;
        row-gap: 35%;
        height: 35%;
        margin-top: 0.7%;
        margin-right: 5%;
    }

    h4 {
        font-size: 1.3rem;
    }

    #confirm-seed-button {
        margin-top: -0.8%;
        width: 10%;
    }
}

@media (min-width: 900px) and (max-width: 1500px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(50%, 1fr)); 
        column-gap: 20%;
        row-gap: 20%;
        margin-left: 3%;
        margin-top: 1%;
    }

    h4 {
        font-size: 1.2rem;
    }

    #confirm-seed-button {
        width: 15%;
        margin-top: 12%;
    }
}

@media (min-width: 600px) and (max-width: 900px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(50%, 1fr)); 
        column-gap: 20%;
        row-gap: 20%;
        margin-left: 3%;
        margin-top: 1%;
    }

    h4 {
        font-size: 1.2rem;
    }

    #confirm-seed-button {
        margin-top: 15%;
        width: 25%;
    }
}

@media (max-width: 600px)  {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(45%, 1fr));
        column-gap: 25%;
        margin-left: 12%;
    }

    h4 {
        font-size: 1rem;
    }

    #confirm-seed-button {
        width: 25%;
        margin-top: 16%;
        
    }
}
</style>