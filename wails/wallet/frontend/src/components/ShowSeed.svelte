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
        { #each seedPhraseList as word }
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
    margin-bottom: 10px;
}

h4 {
    font-size: 2rem;
    margin-bottom: 5px;
    font-weight: bold;
    color: #0066cc;
}

h3 {
    font-size: 1rem;
    margin-bottom: 5px;
}


.seed-words-container {
    display: grid;
    grid-template-columns: repeat(3, minmax(120px, 1fr)); 
    column-gap: 50px; 
    row-gap: 15px; 
    justify-items: center;
    align-items: center;
    margin-top: 10px;
    margin-bottom: 16px;
    width: 100%;
    max-width: 480px; 
}

.seed-words-container .hidden{
    background-color: rgba(48, 23, 23, 0.7); 
    box-shadow: 0px 6px 25px rgba(24, 22, 22, 0.5);
    filter: blur(5px);
    backdrop-filter: blur(5px);
    position: relative;    
}

.seed-phrase-block {
    width: 75%; 
    font-size: 15px;
    background-color: #fff; 
    border: 1px solid #ccc;
    border-radius: 7px;
    height: 40px;
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
    gap: 2px; 
    width: 100%;
    margin-left: -230px;
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
    height: 20px;
    width: 20px;
}
.next-button-container {
    margin-top: -15px;
    display: flex;
    justify-content: center;
    width: 100%;
}

#next-button {
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

@media (min-width: 1025px) {
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

    #next-button {
        width: 7%;
    }
    #eye-icon {
        width: 2%;
        height: 100%;
    }
}

@media (min-width: 768px) and (max-width: 1024px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(54%, 1fr)); 
        column-gap: 15%;
        row-gap: 15%;
    }

    .seed-phrase-block {
        width: 50%;
        justify-content: center;
        align-items: center;
    }

    h3 {
        font-size: 1.1rem;
    }

    h4 {
        font-size: 1.3rem;
    }

    #next-button {
        width: 13%;
    }

    .toggle-seedphrase-container{
        margin-left: -15%;
        margin-top: 12%;
    }

    #eye-icon{
        height: 100%;
        width: 3%;
    }
}

@media (min-width: 480px) and (max-width: 768px) {
    .seed-words-container {
        grid-template-columns: repeat(3, minmax(50%, 1fr));
        column-gap: 15%;
        margin-right: -4%;
    }

    h3 {
        font-size: 1.1rem;
    }

    h4 {
        font-size: 1.3rem;
    }

    #next-button {
        width: 22%;
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