<script>
  import { CreateWallet } from '../../wailsjs/go/main/App'    
  import ProgressBar from '../components/ProgressBar.svelte';
  import CreatePassword from '../components/CreatePassword.svelte';
  

  let currentStep = 0;
  const steps = ["Create Wallet", "Backup Wallet", "Confirm Seed Phrase"];

  function nextStep() {
   if (currentStep < steps.length - 1) {
     currentStep += 1;
   }
  }

  function createWallet() {
    const password = document.getElementById('wallet-password').value;
    CreateWallet(password).then((data, error) => {
      if (error) {
        console.error(error);
        return;
      }
      console.log(data);
      nextStep();
    })
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
      <p>Backup Wallet</p>
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