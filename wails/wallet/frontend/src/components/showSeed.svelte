<script lang="ts">
    export let seedPhraseBlocks: number = 0;
    export let seedList: string[] = [];
    export let onConfirm: () => void = () => {};
</script>

<div class="seed-words-container">
    { #each Array(seedPhraseBlocks) as _, i }
        <input readonly type="text" class="seed-phrase-block" value={seedList[i]} />
    { /each }
</div>
<div class="next-button-container">
    <button id="confirm-recovery-button" on:click={onConfirm}>Next</button>
</div>