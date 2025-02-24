<script lang="ts">
    import logo from '../assets/images/pear-logo.png';
    import { currentView } from "../stores";

    function createWallet(): void {
      currentView.set("CreateWallet");
    }
  
    function importWallet(): void {
      currentView.set("RestoreWallet");
    }
  </script>
  
  <main>
    <img alt="Wails logo" id="logo" src="{logo}">
    <h3>Pear wallet</h3>
    <div class="wallet-buttons">
      <button id="create-wallet-button" on:click={createWallet}>Create a new wallet</button>
      <button id="import-existing-button" on:click={importWallet}>Import an existing wallet</button>
    </div>
  </main>
  
  <style>

    #logo {
      background-color: #A6B057;
      display: block;
      width: 40%;
      height: 40%;
      margin: auto;
      margin-bottom: -4%;
      margin-top: -8%;
      padding: 10% 0 0;
      background-position: center;
      background-repeat: no-repeat;
      background-origin: content-box;
    }

    h3 {
      font-size: 1.3rem;
    }
  
    .wallet-buttons {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 10px;
    }
  
    #create-wallet-button {
      text-align: center; 
      justify-content: center;
      width: 17%;
      color: #FDE024;
      background-color:  black;
      font-family: "Nunito";
      font-style: normal;
      font-weight: bold;
      height:40px;
      border-radius: 30px;
      font-size: 1rem;
    }
  
    #import-existing-button {
      text-align: center; 
      justify-content: center;
      color: #FDE024;
      background-color:darkslategrey;
      font-family: "Nunito";
      font-style: normal;
      font-weight: bold;
      width: 17%;
      height:40px;
      border-radius: 30px;
      font-size: 1rem;
    }

    @media (min-width:769px) and (max-width: 1536px) {
      #logo {
        width: 60%;
        margin-bottom: -4%;
      }

      h3 {
        font-size: 1.3rem;
      }

      #create-wallet-button {
        width: 30%;
        font-size: 1rem;
      }

      #import-existing-button {
        width: 30%;
        font-size: 1rem;
     }
    }

    @media (min-width: 300px) and (max-width: 768px) {
      #logo {
        width: 55%;
        margin-bottom: -4%;
      }

      h3 {
        font-size: 1.2rem;
      }

      #create-wallet-button {
        width: 32%;
        font-size: 0.9rem

      }

      #import-existing-button {
        width: 32%;
        font-size: 0.9rem;
      }
    }

    @media (max-width: 350px) {
      #logo {
        width: 40%;
      }

      #create-wallet-button {
        background-color: blue;
      }
    }
  
  </style>
  