<script lang="ts">
    import logo from '../assets/images/pear-logo.png';
    import { currentView } from "../stores";
    import { RecoverWallet, WalletExists } from "../../wailsjs/go/main/App";
    
    let walletExists = false;
  
    function checkWallet(): void {
      WalletExists()
      .then((exists) => {
        walletExists = exists;
      })
      .catch((err) => {
        console.error(err);
      });
    }

    checkWallet();

    function createWallet(): void {
      currentView.set("Create Wallet");
    }
  
    function importWallet(): void {
      currentView.set("Restore Wallet");
    }

    function checkInput(): void {
      const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
      const unlockButton = document.getElementById('unlock-wallet-button') as HTMLButtonElement;
      if (!passwordInput || !unlockButton) {
        console.error('Password input or unlock button not found');
        return;
      }

      const password: string = passwordInput.value;
      unlockButton.disabled = password.length === 0;
    }

    function recoverWallet(event: Event) : void {
      event.preventDefault(); 
      const passwordInput = document.getElementById('wallet-password') as HTMLInputElement;
      if (!passwordInput) {
        console.error('Password input not found');
        return;
      }
      
      const password: string = passwordInput.value;
      RecoverWallet(password)
      .then(() => {
        currentView.set("Home");
      })
      .catch((err) => {
        console.error(err);
        alert("Invalid password");
      });

    }

  

  </script>
  
  <main>
    <img alt="Wails logo" id="logo" src="{logo}">
    <h3>Pear wallet</h3>
    {#if walletExists === false}
      <div class="wallet-buttons">
        <button  id="create-wallet-button" on:click={createWallet}>Create a new wallet</button>
        <button id="import-existing-button" on:click={importWallet}>Import an existing wallet</button>
      </div>
    {:else}
    <form class="form-container">
      <div class="input-group">
        <label for="password">Password</label>
        <input 
        type="password" 
        id="wallet-password"
        on:input={checkInput}
         />
      </div>
      <button 
      id="unlock-wallet-button"
      disabled=true
      on:click={recoverWallet}>
      Unlock Wallet
    </button>
    </form>    
    {/if}
  </main>
  
  <style>
    #logo {
      background-color: #A6B057;
      display: block;
      margin: auto;
      margin-top: -8%;
      padding: 20% 0 0;
      background-position: center;
    }

    h3 {
      font-size: 1.3rem;
    }

    .wallet-buttons {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 0.9vh;
    }

    #create-wallet-button {
      text-align: center;
      justify-content: center;
      width: 18%;
      color: #FDE024;
      background-color: black;
      font-family: "Nunito";
      font-style: normal;
      font-weight: bold;
      height: 4vh;
      border-radius: 3vh;
      font-size: 1rem;
    }

    #import-existing-button {
      text-align: center;
      justify-content: center;
      color: #FDE024;
      background-color: darkslategrey;
      font-family: "Nunito";
      font-style: normal;
      font-weight: bold;
      width: 17%;
      height: 4vh;
      border-radius: 3vh;
      font-size: 1rem;
    }

    .form-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      align-content: center;
      background-color: #A6B057;
    }

    .input-group {
        display: flex;
        flex-direction: column;
      }

      .input-group label {
        align-self: flex-start;
        font-size: 0.9rem;
        color: #2F2F2F;
      }

    input[type="password"] {
      border: none;
      border-bottom: 0.3vh solid #3D4A27;
      outline: none;
      background-color: transparent;
      color: #2F2F2F;
      transition: border-color 0.3s ease;
    }

    input[type="password"]::placeholder {
      color: #5C5C5C;
    }

    input[type="password"]:focus {
      border-bottom-color: #fff;
    }

    #unlock-wallet-button {
      background-color: black;
      color: #FDE024;
      border: none;
      font-family: "Nunito";
      font-style: normal;
      font-weight: bold;
      width: 50%;
      height: 4vh;
      border-radius: 3vh;
      font-size: 1.1rem;
    }

    #unlock-wallet-button:disabled {
      background-color: #ccc;
      color: #666;
      cursor: not-allowed;
    }

    @media (min-width: 950px) {
      #logo {
        width: 50%;
        height: 80vh;
        padding-bottom: 5%;
        margin-bottom: -10%;
        margin-top: -19%;
      }

      h3 {
        font-size: 1.3rem;
      }

      .wallet-buttons {
        gap: 1.5vh;
      }

      #create-wallet-button {
        width: 14%;
        font-size: 1rem;
      }

      #import-existing-button {
        width: 14%;
        font-size: 1rem;
      }
      .form-container {
        padding: 4%;
        width: 90%;
        margin-top: 0;
        gap: 3vh;
      }

      .input-group {
        margin-top: -3%;
        gap: 1vh;
      }

      input[type="password"] {
        width: 700px;
      }

      #unlock-wallet-button {
        width: 43%;
        height: 4vh;
      }
    }

    @media (min-width: 600px) and (max-width: 950px) {
      #logo {
        width: 55%;
        margin-bottom: -7%;
      }

      h3 {
        font-size: 1.2rem;
      }

      .wallet-buttons {
        gap: 1.5vh;
      }

      #create-wallet-button {
        width: 31%;
        font-size: 0.9rem;
        height: 6vh;
      }

      #import-existing-button {
        width: 31%;
        height: 6vh;
        font-size: 0.9rem;
      }

      .form-container{
        padding: 4%;
        width: 90%;
        margin-top: 0;
        gap: 3vh;
      }

      input[type="password"] {
        width: 450px;
      }

      #unlock-wallet-button {
        width: 55%;
        font-size: 1.1rem;
      }
    }

    @media (min-width: 300px) and (max-width: 600px) {
      #logo {
        width: 55%;
      }

      h3 {
        font-size: 1.1rem;
      }

      .wallet-buttons {
        gap: 2.5vh;
      }

      #create-wallet-button {
        width: 40%;
        height: 7vh;
      }

      #import-existing-button {
        width: 40%;
        height: 7vh;
      }

      .form-container{
        padding: 4%;
        width: 90%;
        margin-top: -5%;
        gap: 4vh;
      }

      input[type="password"] {
        width: 290px;
      }

      #unlock-wallet-button {
        width: 58%;
        height: 6vh;
        font-size: 1rem;
      }
    }
  </style>
  