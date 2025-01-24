<script>
    let seedPhraseBlocks = 12;
</script>

<main>
    <h3>Restore wallet</h3>
    <h4>Enter the Secret Recovery Phrase that you were given when you created your wallet.</h4>
    <div class="seed-words-container">
        { #each Array(seedPhraseBlocks) as _, i }
            <input type="text" class="seed-phrase-block" />
        { /each }
    </div>
    <div class="confirm-recovery-button-container">
        <button id="confirm-recovery-button">Confirm Recovery Phrase</button>
    </div>
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;  /* Start from the top */
        align-items: center;
        height: 100vh;
        padding: 20px;
        color: greenyellow;
        font-family: "Nunito", sans-serif;
        flex-grow: 1;  /* Allow the main element to grow */
    }

    h3, h4 {
        margin-top: 20px;  /* Add some space from the top */
    }

    .seed-words-container {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); /* Flexible grid */
        gap: 30px;
        justify-items: center;
        align-items: center;
        font-family: "Nunito";
        margin-bottom: 20px;
        width: 100%;
        max-width: 600px; /* Max width for the grid container */
    }

    .seed-phrase-block {
        width: 100%; /* Full width of its container */
        font-size: 15px;
        background-color: darkslategrey;
        border: 1px solid #ccc;
        border-radius: 7px;
        height: 40px;
        font-family: "Nunito";
    }

    .seed-phrase-block::placeholder {
        color: greenyellow;
        font-family: "Nunito";
    }

    .confirm-recovery-button-container {
        display: flex;
        justify-content: center; /* Horizontally center the button */
        width: 100%; /* Ensure the container spans the full width */
    }

    #confirm-recovery-button {
        padding: 12px 20px;
        margin-top: 40px;
        background-color: black;
        max-width: 300px;
        font-family: "Nunito";
        font-weight: bold;
        color: greenyellow;
        border-radius: 20px;
        border: none;
        width: 100%; /* Make the button responsive */
    }

    @media (max-width: 768px) {
        .seed-words-container {
            grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); /* More compact on smaller screens */
            gap: 20px;
        }

        #confirm-recovery-button {
            max-width: 250px; /* Adjust button size */
        }
    }

    @media (max-width: 480px) {
        .seed-words-container {
            grid-template-columns: 1fr; /* Stack the blocks in a single column on very small screens */
            gap: 15px;
        }

        #confirm-recovery-button {
            max-width: 200px; /* Further adjust button size */
        }

        h3 {
            font-size: 1.5rem; /* Adjust title size on small screens */
        }

        h4 {
            font-size: 1rem; /* Adjust subtitle size on small screens */
        }
    }
</style>
