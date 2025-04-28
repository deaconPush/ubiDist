<script lang="ts">
    import { assets } from '../stores';

    $: userAssets = $assets;
    console.log(userAssets);

    function clickCard(asset) {
        alert('Card clicked for: ' + asset.symbol); 
    }
</script>


<main>
    <h3>Available assets</h3>
    <div class="assets-container">
        {#each userAssets as asset}
            {#if asset.balance > 0}
                <div class="asset-card" on:click={clickCard(asset)}
                    style="cursor: pointer; border: 1px solid #ccc; padding: 10px; margin: 5px; border-radius: 5px;">
                    <h3>{asset.symbol}</h3>
                    <h4>{asset.balance}</h4>
                    <h4>asset.balan</h4>
                    <img src={asset.logoPath} alt={asset.name} class="asset-logo"  style={
                        `width: 50px; height: 50px; border-radius: 50%;`
                    }/>
                </div>
            {/if}
        {/each}  
    </div>
</main>

<style>
main {
        font-family: "Nunito", sans-serif;
        color: black;
        background-color: #f9f9f9; 
        height: 100vh;
        padding: 5%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }
</style>

