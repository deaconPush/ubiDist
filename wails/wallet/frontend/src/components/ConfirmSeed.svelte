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
    margin-top: -10px;
    background-color: #ef9f9e;
    font-size: 0.9rem;
    color: black; 
    display: none; 
}

@media (max-width: 1024px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(100px, 1fr));
        column-gap: 20px; 
        max-width: 340px;
    }

    #confirm-seed-button {
        max-width: 250px;
    }
}

@media (max-width: 768px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(80px, 1fr)); 
        column-gap: 40px;
        max-width: 300px;
    }

    h4 {
        font-size: 1.2rem;
        margin-top: 10px;
    }

    #confirm-seed-button {
        max-width: 200px;
    }
}

@media (max-width: 480px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(80px, 1fr));
        column-gap: 40px;
    }

    h4 {
        font-size: 1rem;
    }

    #confirm-seed-button {
        max-width: 100%; 
        padding: 10px 15px;
    }
}
</style>