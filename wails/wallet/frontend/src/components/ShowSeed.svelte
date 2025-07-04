<script lang="ts">
    import display from '../assets/images/display.png';
    import hide from '../assets/images/hide.png';

    export let seedPhraseList: string[] = [];
    export let onConfirm: () => void = () => {};
    let showSeedPhrase: boolean = false;

    function toggleSeedPhraseVisibility(): void {
        showSeedPhrase = !showSeedPhrase;
    }
</script>

<div class="header-content">
    <h4>Write down your Secret Recovery Phrase</h4>
    <h3>Write down your Secret Recovery Phrase and save it in a safe place with protected access.</h3>
</div>
<div class="seed-words-container">
    <div class="seed-words-container {showSeedPhrase ? "" : "hidden"}">
        { #each seedPhraseList as word, index (index) }
            <input readonly type="text" class="seed-phrase-block {showSeedPhrase ?  "": "hidden" }" value={word} />
        { /each }
    </div>
</div>
<div class="toggle-seedphrase-container" >
    <img 
    src={showSeedPhrase ? hide : display} 
    on:click={toggleSeedPhraseVisibility}
    id="eye-icon"
    alt="eye toggle icon"
    />
    <span
    id="toggle-seed-phrase-span"
    on:click={toggleSeedPhraseVisibility}
    >
    {showSeedPhrase ? 'Hide seed phrase' : "Reveal seed phrase"}
    </span>
</div>
<div class="next-button-container">
    <button id="next-button" on:click={onConfirm}>Next</button>
</div> 


<style>
.header-content {
    text-align: center;
    margin-bottom: 2vh;
}

h4 {
    font-size: 2rem;
    margin-bottom:1vh;
    font-weight: bold;
    color: #0066cc;
}

h3 {
    font-size: 1rem;
    margin-bottom: -1vh;
}


.seed-words-container {
    display: grid;  
    justify-items: center;
    align-items: center;
    margin-top: 2vh;
    margin-bottom: 1.5vh;
    width: 55%;
    margin-left: 6.5%;
}

.seed-words-container .hidden{
    background-color: rgba(48, 23, 23, 0.7); 
    box-shadow: 0px 6px 25px rgba(24, 22, 22, 0.5);
    filter: blur(5px);
    backdrop-filter: blur(5px);
    position: relative;    
}

.seed-phrase-block {
    width: 65%; 
    font-size: 1.1rem;
    background-color: #fff; 
    border: 1px solid #ccc;
    border-radius: 1vh;
    height: 4vh;
    font-family: "Nunito";
    padding: 0 35px 0 10px;
    color: #333; 
}

.seed-phrase-block .hidden {
    box-shadow: 0px 3px 10px rgba(0, 0, 0, 0.2);
    position: relative;
    z-index: 2; 
    background-color: #f5f5f5;
    opacity: 0.3;
    color: #333
}

.toggle-seedphrase-container {
    display: flex;
    align-items: center; 
    justify-content: center;
    gap: 0.5vh; 
    width: 100%;
    margin-left: -30%;
    cursor: pointer;
}

.toggle-seedphrase-container span {
    background-color: transparent;
    border: none;
    color: #0066cc;
    font-family: "Nunito";
    font-weight: bold;
    cursor: pointer;
    font-size: 0.9rem;
    margin: 0;
}

.toggle-seedphrase-container img {
    cursor: 'pointer';
    height: 100%;
    width: 5vh;
}
.next-button-container {
    margin-top: 2%;
    display: flex;
    justify-content: center;
    width: 100%;
}

#next-button {
    padding: 0.5% 0.5%;
    background-color: #0066cc;
    font-family: "Nunito";
    font-weight: bold;
    color: #fff;
    border-radius: 4vh;
    border: none;
    width: 100%;
    cursor: pointer;
    transition: background-color 0.3s ease;
}

@media (min-width: 1500px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(55%, 1fr));
        column-gap: 15%; 
        row-gap: 20%;
    }

    h3 {
        font-size: 1.1rem;
    }

    h4 {
        font-size: 1.3rem;
    }

    .toggle-seedphrase-container{
        margin-left: -15%;
        margin-top: 6%;
    }

    .toggle-seedphrase-container span {
        font-size: 1.1rem;
    }

    #next-button {
        width: 7%;
        font-size: 1rem;
    }
    #eye-icon {
        width: 2%;
        height: 100%;
    }
}

@media (min-width: 768px) and (max-width: 1500px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(65%, 1fr)); 
        column-gap: 30%;
        row-gap: 18%;
        margin-left: -8%;
    }

    .seed-phrase-block {
        width: 65%;
        justify-content: center;
        align-items: center;
    }

    h3 {
        font-size: 1.1rem;
    }

    h4 {
        font-size: 1.3rem;
    }

    .next-button-container{
        margin-top: 2.5%;
    }

    #next-button {
        width: 13%;
        height: 3vh;
        font-size: 1rem; 
    }

    .toggle-seedphrase-container{
        margin-left: -15%;
        margin-top: 11.5%;
    }

    .toggle-seedphrase-container span {
        font-size: 1.1rem;
    }

    #eye-icon{
        height: 100%;
        width: 3%;
    }


}

@media (min-width: 480px) and (max-width: 768px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(65%, 1fr));
        column-gap: 70%;
        row-gap: 20%;
        margin-left: -20%;
    }

    .seed-phrase-block {
        width: 100%;
        height: 6;
    }

    h3 {
        font-size: 1.1rem;
    }

    h4 {
        font-size: 1.3rem;
    }

    .next-button-container{
        margin-top: 3%;
    }

    .toggle-seedphrase-container{
        margin-top: 13%;
    }

    #next-button {
        width: 22%;
        height: 5vh;
    }

    #eye-icon{
        width: 4%;
        height: 100%;
    }

    .toggle-seedphrase-container span{
        font-size: 1rem;
     }
}
</style>