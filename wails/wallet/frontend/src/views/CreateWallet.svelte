<script>
  import { CreateWallet } from '../../wailsjs/go/main/App'    
  import CreatePassword from '../components/CreatePassword.svelte';

  function createWallet() {
    const password = document.getElementById('wallet-password').value;
    CreateWallet(password).then((data, error) => {
      if (error) {
        console.error(error);
        return;
      }
      console.log(data);
    })
  }
  
  </script>
  <main>
    <CreatePassword
      handleClick={createWallet}
      walletLabel="Create Wallet"
    />
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