<script>
    export let onConfirm = () => {};
    export let seedPhraseBlocks = 0;

    function checkSeedInputs() {
        const inputs = document.querySelectorAll('.seed-phrase-block');
        const confirmButton = document.getElementById('confirm-recovery-button');
        const allFilled = Array.from(inputs).every(input => input.value.trim() !== '');
        confirmButton.disabled = !allFilled;
    }
</script>

<h4>Enter the Secret Recovery Phrase that you were given when you created your wallet.</h4>
<div class="seed-words-container">
    { #each Array(seedPhraseBlocks) as _, i }
        <input on:input={checkSeedInputs} type="text" class="seed-phrase-block" />
    { /each }
</div>
<div class="confirm-recovery-button-container">
    <button id="confirm-recovery-button" on:click={onConfirm} disabled=true>Confirm Recovery Phrase</button>
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
    column-gap: 40px; 
    row-gap: 15px; 
    justify-items: center;
    align-items: center;
    margin-bottom: 20px;
    width: 100%;
    max-width: 600px; 
}

.seed-phrase-block {
    width: 100%; 
    font-size: 15px;
    background-color: #fff; 
    border: 1px solid #ccc;
    border-radius: 7px;
    height: 40px;
    font-family: "Nunito";
    padding: 0 10px;
    color: #333; 
}

.seed-phrase-block::placeholder {
    color: #999; 
}

.confirm-recovery-button-container {
    display: flex;
    justify-content: center;
    width: 100%;
}

#confirm-recovery-button {
    padding: 12px 20px;
    margin-top: 40px;
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

#confirm-recovery-button:hover:not(:disabled) {
    background-color: #004c99; 
}

#confirm-recovery-button:disabled {
    background-color: #ccc; 
    color: #666;
    cursor: not-allowed;
}

@media (max-width: 1024px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(100px, 1fr));
        column-gap: 40px; 
    }

    #confirm-recovery-button {
        max-width: 250px;
    }
}

@media (max-width: 768px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(80px, 1fr)); 
        column-gap: 40px;
    }

    h3, h4 {
        font-size: 1.2rem;
        margin-top: 10px;
    }

    #confirm-recovery-button {
        max-width: 200px;
    }
}

@media (max-width: 480px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(80px, 1fr));
        column-gap: 40px;
    }

    h3 {
        font-size: 1.5rem;
    }

    h4 {
        font-size: 1rem;
    }

    #confirm-recovery-button {
        max-width: 100%; 
        padding: 10px 15px;
    }
}
</style>
