<script>
  import { GetBalance } from "../wailsjs/go/main/App.js";
  import logo from './assets/images/pear-logo.png'

  let initialBalance = 0;
  function init() {
    getBalance();
  }

  init();

  function getBalance() {
    GetBalance().then((balance, error) => {
      if (error) {
        console.log(error);
        return;
      }
      initialBalance = balance;
    })
  }
  
</script>

<main>
  <img alt="Wails logo" id="logo" src="{logo}">
  <h3>Pear wallet</h3>
  <div class="wallet-buttons">
    <button id="create-wallet-button">Create a new wallet</button>
    <button id="import-existing-button">Import an existing wallet</button>
  </div>
</main>

<style>

  #logo {
    background-color: #A6B057;
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-origin: content-box;
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
    max-width: 100px;
    color: #FDE024;
    background-color:  black;
    font-family: "Nunito";
    font-style: normal;
    font-weight: bold;
    width: 100%;
    height:40px;
    max-width: 225px;
    border-radius: 30px;
  }

  #import-existing-button {
    text-align: center; 
    justify-content: center;
    max-width: 100px;
    color: #FDE024;
    background-color:darkslategrey;
    font-family: "Nunito";
    font-style: normal;
    font-weight: bold;
    width: 100%;
    height:40px;
    max-width: 225px;
    border-radius: 30px;
  }

</style>
