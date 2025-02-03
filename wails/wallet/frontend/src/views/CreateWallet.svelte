<script lang="ts">
  import { CreateWallet } from '../../wailsjs/go/main/App'    
  import ProgressBar from '../components/ProgressBar.svelte';
  import CreatePassword from '../components/CreatePassword.svelte';
  import ShowSeed from '../components/ShowSeed.svelte';
  
  let currentStep: number = 0;
  let seedPhraseList: string[] = [];
  const steps: string[] = ["Create Wallet", "Backup Seed Phrase", "Confirm Seed Phrase"];
  function nextStep(): void {
   if (currentStep < steps.length - 1) {
     currentStep += 1;
   }
  }

  function createWallet(): void {
    const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
    
    if (!passwordInput) {
      console.error('Password input not found');
      return;
    }

    const password: string = passwordInput.value;
    CreateWallet(password)
    .then((data) => {
      seedPhraseList = data.split(' ');
      nextStep();
    })
    .catch((error) => {
      console.error(error);
    });
  }

  function confirmSeedPhrase(): void {
    nextStep();
  }
  </script>
  <main>
    <ProgressBar 
    steps={steps} 
    currentStep={currentStep} 
    />
    {#if currentStep === 0}
      <CreatePassword
        handleClick={createWallet}
        walletLabel="Create Wallet"
      />
    {:else if currentStep === 1}
      <ShowSeed
        seedPhraseList={seedPhraseList}
        onConfirm={confirmSeedPhrase}
      />  
    {:else if currentStep === 2}
      <p>Confirm Seed Phrase</p>
    {/if}
  </main>

<style>
  main {
    font-family: "Nunito", sans-serif;
    color: #333; 
    background-color: #f9f9f9; 
    height: 100vh;
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

</style>