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
                    validationLabel.style.visibility = 'hidden';
                } else {
                    confirmButton.disabled = true;
                    validationLabel.style.visibility = 'visible';
                    validationLabel.textContent = 'Invalid seed phrase';
                }
            });
        } else {
            validationLabel.style.visibility = 'hidden';
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
    margin-top: 3.5vh;
    color: #333; 
    text-align: center; 
    padding: 0 2%;
}

.seed-words-container {
    display: grid;
    justify-items: center;
    align-items: space-between;
}

.seed-phrase-block {
    width: 65%; 
    font-size: 1rem;
    background-color: #fff; 
    border: 1px solid #ccc;
    border-radius: 1.5vh;
    height: 6vh;
    font-family: "Nunito";
    padding: 0 20% 0 4%;
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
    background-color: #0066cc;
    font-family: "Nunito";
    font-weight: bold;
    color: #fff; 
    border-radius: 5vh;
    border: none;
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
    background-color: #ffdddd; 
    border-left: 1vh solid #d9534f;
    font-size: 0.9rem;
    color: #a94442; 
    padding: 0.8% 0.8%;
    border-radius: 1vh;
    visibility: hidden;
    transition: opacity 0.3s ease-in-out;
}

@media (min-width: 1500px) {
    .header-content  {
        margin-top: -0.5%;
        width: 100%;
    }

    .header-content h4 {
        font-size: 1.2rem;
    }

    .seed-words-container {
        grid-template-columns: repeat(3, minmax(28%, 1fr));
        column-gap: 10%;
        width: 100%;
        row-gap: 35%;
        height: 35%;
        margin-top: 0.2%;
        margin-left: 55%;
        margin-bottom: 1%;
    }

    h4 {
        font-size: 1.3rem;
    }

    .validation-label{
        margin-top: 3%
        
    }

    .confirm-seed-button-container{
        margin-top: 3%;
    }

    #confirm-seed-button {
        width: 8%;
        height: 6vh
    }

    #confirm-seed-button:disabled {
        margin-top: -2%;
    }
}

@media (min-width: 900px) and (max-width: 1500px) {
    .header-content  {
        margin-top: 1.2%;
    }
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(45%, 1fr)); 
        column-gap: 20%;
        row-gap: 20%;
        margin-left: 17%;
        margin-bottom: 20%;
    }



    h4 {
        font-size: 1.2rem;
    }

    .validation-label{
        margin-top: -7%
    }


    #confirm-seed-button {
        width: 16%;
        height: 5vh;
    }

    #confirm-seed-button:disabled {
        margin-top: 1%;
    }

    #confirm-seed-button:enabled {
        margin-top: -4%;    
    }

}

@media (min-width: 600px) and (max-width: 900px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(50%, 1fr)); 
        column-gap: 12%;
        row-gap: 35%;
        height: 8vh;
        margin-top: 1%;
        margin-left: 10%;
        margin-bottom: 20%;
    }

    .seed-phrase-block  {
        width: 50%;
    }

    h4 {
        font-size: 1.2rem;
    }   

    .validation-label{
        margin-top: -1%;
        height: 2.5vh;
        margin-bottom: 3%;
    }

    #confirm-seed-button {
        width: 17%;
        height: 5vh;
    }

}

@media (max-width: 600px)  {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(52%, 1fr));
        column-gap: 10%;
        row-gap: 35%;
        margin-left: 8%;
        margin-bottom: 38%;
    }

    h4 {
        font-size: 1rem;
    }   

    .validation-label{
        margin-top: -21%;
        height: 4.5vh;
    }

    #confirm-seed-button {
        width: 20%;
        height: 7vh;
    }


    
}
</style>