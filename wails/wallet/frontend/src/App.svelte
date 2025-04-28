<script>
  import WalletRecovery from "./views/WalletRecovery.svelte";
  import RestoreWallet from "./views/RestoreWallet.svelte";
  import CreateWallet  from "./views/CreateWallet.svelte";
  import Home from "./views/Home.svelte";
  import Send from "./views/Send.svelte";
  import { currentView } from "./stores";
  
  $: view = $currentView;
</script>

<main>
  {#if view === "Wallet Recovery"}
    <WalletRecovery />
  {:else if view === "RestoreWallet"}
    <RestoreWallet />
  {:else if view === "CreateWallet"}
    <CreateWallet />
  {:else if view === "Home"}
    <Home />
  {:else if view === "Send"}
    <Send />
  {/if}
</main>
