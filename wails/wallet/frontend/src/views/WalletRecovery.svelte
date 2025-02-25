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
      background-color:  black;
      font-family: "Nunito";
      font-style: normal;
      font-weight: bold;
      height:4vh;
      border-radius: 3vh;
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
      height: 4vh;
      border-radius: 3vh;
      font-size: 1rem;
    }

    @media (min-width: 950px) and (max-width: 2048px) {
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

      .wallet-buttons{
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
    }

    @media  (min-width: 600px) and (max-width: 950px) {
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
    }
  
  </style>
  