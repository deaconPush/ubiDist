<script>
  import { GetRandomImageUrl } from "../wailsjs/go/main/App.js";
  import { GetBreedList } from "../wailsjs/go/main/App.js";
  import { GetImageUrlsByBreed } from "../wailsjs/go/main/App.js";

  let randomImageUrl = "";
  let breeds = [];
  let photos = [];
  let selectedBreed;
  let showRandomPhoto = false;
  let showBreedPhotos = false;

  function init() {
    getBreedList();
  }

  init();

  function getRandomImageUrl() {
    showRandomPhoto = false;
    showBreedPhotos = false;
    GetRandomImageUrl().then((result) => (randomImageUrl = result));
    showRandomPhoto = true;
  }

  function getBreedList() {
    GetBreedList().then((result) => (breeds = result));
  }

  function getImageUrlsByBreed() {
    init();
    showRandomPhoto = false;
    showBreedPhotos = false;
    GetImageUrlsByBreed(selectedBreed).then((result) => (photos = result));
    showBreedPhotos = true;
  }
</script>

<h3>Dogs API</h3>
<div>
  <button class="btn" on:click={getRandomImageUrl}>
    Fetch a dog randomly
  </button>
  Click on down arrow to select a breed
  <select bind:value={selectedBreed}>
    {#each breeds as breed}
      <option value={breed}>
        {breed}
      </option>
    {/each}
  </select>
  <button class="btn" on:click={getImageUrlsByBreed}>
    Fetch by this breed
  </button>
</div>
<br />
{#if showRandomPhoto}
  <img id="random-photo" src={randomImageUrl} alt="No dog found" />
{/if}
{#if showBreedPhotos}
  {#each photos as photo}
    <img id="breed-photos" src={photo} alt="No dog found" />
  {/each}
{/if}

<style>
  #random-photo {
    width: 600px;
    height: auto;
  }

  #breed-photos {
    width: 300px;
    height: auto;
  }

  .btn:focus {
    border-width: 3px;
  }
</style>