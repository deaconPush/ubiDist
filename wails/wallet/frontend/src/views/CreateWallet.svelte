<script lang="ts">
  import { CreateWallet } from '../../wailsjs/go/main/App'    
  import ProgressBar from '../components/ProgressBar.svelte';
  import CreatePassword from '../components/CreatePassword.svelte';
  import ShowSeed from '../components/ShowSeed.svelte';
  import ConfirmSeed from '../components/ConfirmSeed.svelte';
  import { currentView, availableTokens } from '../stores';
  
  
  let currentStep: number = 0;
  let seedPhraseList: string[] = [];
  const steps: string[] = ["Create Wallet", "Backup Seed Phrase", "Confirm Seed Phrase"];

  function nextStep(): void {
   if (currentStep < steps.length - 1) {
     currentStep += 1;
   }
  }

  function passwordConfirmed(): void {
    const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
    
    if (!passwordInput) {
      console.error('Password input not found');
      return;
    }

    const password: string = passwordInput.value;
    CreateWallet($availableTokens, password)
    .then((data) => {
      seedPhraseList = data.split(' ');
      nextStep();
    })
    .catch((error) => {
      console.error(error);
    });
  }

  function SeedSaved(): void {
    nextStep();
  }

  function SeedConfirmed(): void {
    currentView.set('Home');
  }


  </script>
  <main>
    <ProgressBar 
    steps={steps} 
    currentStep={currentStep} 
    />
    {#if currentStep === 0}
      <CreatePassword
        handleClick={passwordConfirmed}
        walletLabel="Create Password for Wallet"
      />
    {:else if currentStep === 1}
      <ShowSeed
        seedPhraseList={seedPhraseList}
        onConfirm={SeedSaved}
      />  
    {:else if currentStep === 2}  
      <ConfirmSeed
        seedPhraseList={seedPhraseList}
        onConfirm={SeedConfirmed}
      />
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