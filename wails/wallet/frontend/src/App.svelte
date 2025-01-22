<script>
  import { GetBalance } from "../wailsjs/go/main/App.js";
  import logo from './assets/images/logo-universal.png'

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
  <h3>Wallet App</h3>
  <div>
    <h4 id="wallet-balance">Initial Balance: {initialBalance}</h4>
  </div>
</main>

<style>

  #logo {
    display: block;
    width: 50%;
    height: 50%;
    margin: auto;
    padding: 10% 0 0;
    background-position: center;
    background-repeat: no-repeat;
    background-size: 100% 100%;
    background-origin: content-box;
  }

  #wallet-balance {
    text-align: center;
  }

</style>
