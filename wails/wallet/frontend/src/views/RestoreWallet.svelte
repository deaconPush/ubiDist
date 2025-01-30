<script>
    import { ValidateMnemonic, RestoreWallet } from '../../wailsjs/go/main/App'    
    import SeedRecovery from '../components/SeedRecovery.svelte';
    import CreatePassword from '../components/CreatePassword.svelte';
    
    let seedPhrase = ''
    let seedPhraseBlocks = 12;
    let showSeedRecovery = true;

    function restoreWallet() {
        const password = document.getElementById('wallet-password').value;
        RestoreWallet(password, seedPhrase).then((err) => {
            if (err) {
                console.log("Error restoring wallet: ", err);
            }
        })
    }
    function confirmRecoveryPhrase() {
        const inputs = document.querySelectorAll('.seed-phrase-block');
        seedPhrase = Array.from(inputs).map(input => input.value.trim()).join(' ');
        ValidateMnemonic(seedPhrase).then(isValid => {
            if(isValid) {
                showSeedRecovery = false;
            } else {
                
            }
        })
    }
</script>


<main>
    {#if showSeedRecovery}
        <SeedRecovery 
        seedPhraseBlocks={seedPhraseBlocks}
        onConfirm={confirmRecoveryPhrase}
        />
    {:else}
        <CreatePassword 
            handleClick={restoreWallet}
            walletLabel="Restore Wallet"
        />
    {/if}
    
</main>

<style>
main {
    display: flex;
    flex-direction: column;
    justify-content: flex-start; 
    align-items: center;
    height: 100vh;
    padding: 20px;
    color: #222; 
    background-color: #f4f4f9; 
    font-family: "Nunito", sans-serif;
    flex-grow: 1;
}
</style>
