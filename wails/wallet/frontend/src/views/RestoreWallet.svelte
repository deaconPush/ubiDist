<script lang="ts">
    import { RestoreWallet } from '../../wailsjs/go/main/App'    
    import ProgressBar from '../components/ProgressBar.svelte';
    import SeedRecovery from '../components/SeedRecovery.svelte';
    import CreatePassword from '../components/CreatePassword.svelte';
    import { currentView } from '../stores';

    
    let seedPhrase: string = ''
    let seedPhraseBlocks: number = 12;
    let showSeedRecovery: boolean = true;
    let currentStep: number = 0;
    const steps: string[] = ["Seed Recovery", "Create Password"];

    function nextStep(): void {
    if (currentStep < steps.length - 1) {
        currentStep += 1;
    }
    }

    function restoreWallet(): void {
        const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
        if (!passwordInput) {
            console.error('Password input not found');
            return;
        }

        const password = passwordInput.value;
        RestoreWallet(password, seedPhrase)
        .then(() => {
            currentView.set('Home');
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
        nextStep();    
        }
</script>


<main>
    <ProgressBar 
    steps={steps} 
    currentStep={currentStep} 
    />
    {#if currentStep === 0}
        <SeedRecovery 
        seedPhraseBlocks={seedPhraseBlocks}
        onConfirm={confirmRecoveryPhrase}
        />
    {:else if currentStep === 1}
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
