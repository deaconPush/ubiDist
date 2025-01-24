<script>
  import Home from "./views/Home.svelte";
  import RestoreWallet from "./views/RestoreWallet.svelte";
  import { currentView } from "./stores.js";

  $: view = $currentView;
</script>

<main>
  {#if view === "home"}
    <Home />
  {:else if view === "RestoreWallet"}
    <RestoreWallet />
  {/if}
</main>


