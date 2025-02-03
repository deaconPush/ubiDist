<script lang="ts">
    import { ValidateMnemonic, RestoreWallet } from '../../wailsjs/go/main/App'    
    import SeedRecovery from '../components/SeedRecovery.svelte';
    import CreatePassword from '../components/CreatePassword.svelte';
    
    let seedPhrase: string = ''
    let seedPhraseBlocks: number = 12;
    let showSeedRecovery: boolean = true;

    function restoreWallet(): void {
        const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
        if (!passwordInput) {
            console.error('Password input not found');
            return;
        }

        const password = passwordInput.value;
        RestoreWallet(password, seedPhrase)
        .then(() => {
            console.log("Wallet restored");
        })
        .catch((error) => {
            console.error(error);
        });
    }
    function confirmRecoveryPhrase(): void {
        const inputs = document.querySelectorAll('.seed-phrase-block') as NodeListOf<HTMLInputElement>;
        if (!inputs) {
            console.error('Seed phrase inputs not found');
            return;
        }
        
        seedPhrase = Array.from(inputs).map(input => input.value.trim()).join(' ');
        ValidateMnemonic(seedPhrase).then(isValid => {
            if(isValid) {
                showSeedRecovery = false;
            } else {
                console.error('Invalid seed phrase');
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
