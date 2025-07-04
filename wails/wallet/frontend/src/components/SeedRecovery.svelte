<script lang="ts">
  import display from '../assets/images/display.png';
  import hide from '../assets/images/hide.png';
  import { ValidateMnemonic } from '../../wailsjs/go/main/App';

  export let onConfirm: () => void = () => {};
  export let seedPhraseBlocks: number = 0;

  function checkSeedInputs(): void {
    const inputs = document.querySelectorAll('.seed-phrase-block') as NodeListOf<HTMLInputElement>;
    if (!inputs) {
      console.error('Seed phrase inputs not found');
      return;
    }

    const confirmButton = document.getElementById('confirm-recovery-button') as HTMLButtonElement;
    if (!confirmButton) {
      console.error('Confirm button not found');
      return;
    }

    const validationLabel = document.getElementById(
      'seed-phrase-validation'
    ) as HTMLParagraphElement;
    if (!validationLabel) {
      console.error('Validation label not found');
      return;
    }

    const allFilled = Array.from(inputs).every((input) => input.value.trim() !== '');
    if (allFilled) {
      const seedPhrase: string = Array.from(inputs)
        .map((input) => input.value.trim())
        .join(' ');
      ValidateMnemonic(seedPhrase).then((isValid) => {
        if (isValid) {
          confirmButton.disabled = false;
          validationLabel.style.display = 'none';
        } else {
          validationLabel.style.display = 'block';
          validationLabel.textContent = 'Invalid seed phrase';
        }
      });
    }
  }

  function toggleSeedBlockVisibility(event: MouseEvent): void {
    const eyeIcon = event.target as HTMLImageElement;
    if (!eyeIcon) {
      console.error('Target is not an image element');
      return;
    }

    const input = eyeIcon.previousElementSibling as HTMLInputElement;
    if (!input) {
      console.error('Previous element is not an input element');
      return;
    }

    if (input.type === 'password') {
      input.type = 'text';
      eyeIcon.src = hide;
    } else {
      input.type = 'password';
      eyeIcon.src = display;
    }
  }
</script>

<h4>Enter the Secret Recovery Phrase that you were given when you created your wallet.</h4>
<div class="seed-words-container">
  {#each Array(seedPhraseBlocks) as _, i (i)}
    <div class="input-wrapper">
      <input
        on:input={checkSeedInputs}
        type="password"
        class="seed-phrase-block"
        id={`seed-prhase-block-${i}`}
      />
      <img
        src={display}
        class="eye-icon"
        alt="eye toggle icon"
        on:click={toggleSeedBlockVisibility}
        id={`eye-icon-${i}`}
      />
    </div>
  {/each}
</div>
<p id="seed-phrase-validation" class="validation-label"></p>
<div class="confirm-recovery-button-container">
  <button id="confirm-recovery-button" on:click={onConfirm} disabled="true"
    >Confirm Recovery Phrase</button
  >
</div>

<style>
  h4 {
    margin-top: 4vh;
    color: #333;
    text-align: center;
    padding: 0 5%;
  }

  .seed-words-container {
    display: grid;
    justify-items: center;
    align-items: space-between;
    margin-bottom: 3vh;
  }

  .input-wrapper {
    position: relative;
    width: 100%;
    display: flex;
    align-items: center;
  }

  .seed-phrase-block {
    width: 65%;
    font-size: 15px;
    background-color: #fff;
    border: 1px solid #ccc;
    border-radius: 7px;
    height: 40px;
    font-family: 'Nunito';
    padding: 0 35px 0 10px;
    color: #333;
  }

  .eye-icon {
    position: absolute;
    right: -17px;
    transform: translateX(50%);
    cursor: pointer;
    z-index: 10;
    height: 20px;
    width: 20px;
  }
  .seed-phrase-block::placeholder {
    color: #999;
  }

  .confirm-recovery-button-container {
    display: flex;
    justify-content: center;
    width: 100%;
  }

  #confirm-recovery-button {
    padding: 12px 20px;
    margin-top: 15px;
    background-color: #0066cc;
    font-family: 'Nunito';
    font-weight: bold;
    color: #fff;
    border-radius: 20px;
    border: none;
    width: 25%;
    cursor: pointer;
    transition: background-color 0.3s ease;
  }

  #confirm-recovery-button:hover:not(:disabled) {
    background-color: #004c99;
  }

  #confirm-recovery-button:disabled {
    background-color: #ccc;
    color: #666;
    cursor: not-allowed;
  }

  .validation-label {
    margin-top: -8px;
    background-color: #ffdddd;
    border-left: 4px solid #d9534f;
    font-size: 0.9rem;
    color: #a94442;
    padding: 8px 12px;
    border-radius: 5px;
    display: none;
    width: fit-content;
    max-width: 100%;
    transition: opacity 0.3s ease-in-out;
  }

  @media (min-width: 1800px) {
    .seed-words-container {
      grid-template-columns: repeat(3, minmax(25%, 1fr));
      height: 35%;
      column-gap: 11%;
    }

    h4 {
      font-size: 1.3rem;
    }

    .eye-icon {
      right: -4%;
      width: 17%;
      height: 3.5vh;
    }

    .confirm-recovery-button-container {
      margin-top: -1.5vh;
    }

    #confirm-recovery-button {
      width: 12%;
      height: 6vh;
    }
  }

  @media (min-width: 769px) and (max-width: 1800px) {
    .seed-words-container {
      grid-template-columns: repeat(3, minmax(25%, 1fr));
      height: 35%;
      column-gap: 11%;
    }

    h4 {
      font-size: 1.3rem;
    }

    .eye-icon {
      right: -4%;
      width: 17%;
      height: 3.5vh;
    }

    .confirm-recovery-button-container {
      margin-top: -1.5vh;
    }

    #confirm-recovery-button {
      width: 25%;
      height: 6vh;
    }
  }

  @media (min-width: 481px) and (max-width: 768px) {
    .seed-words-container {
      grid-template-columns: repeat(3, minmax(15%, 1fr));
      column-gap: 10%;
      row-gap: 2.5vh;
      width: 80%;
    }

    h4 {
      font-size: 1.1rem;
      margin-top: 2vh;
    }

    #confirm-recovery-button {
      width: 30%;
    }

    .eye-icon {
      right: -15%;
      width: 20%;
      height: 60%;
    }
  }
</style>
