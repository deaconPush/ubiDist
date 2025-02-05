<script>
  import WalletRecovery from "./views/WalletRecovery.svelte";
  import RestoreWallet from "./views/RestoreWallet.svelte";
  import CreateWallet  from "./views/CreateWallet.svelte";
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
  {/if}
</main>


