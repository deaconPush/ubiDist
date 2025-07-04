<script lang="ts">
  export let steps: string[] = [];
  export let currentStep: number = 0;
</script>

<div class="progress-container">
  {#each steps as step, index (index)}
    <div class="step {index <= currentStep ? 'active' : ''}">
      <div class="circle">{index + 1}</div>
      <p>{step}</p>
    </div>

    {#if index < steps.length - 1}
      <div class="progress-line {index < currentStep ? 'filled' : ''}"></div>
    {/if}
  {/each}
</div>

<style>
  .progress-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    width: 100%;
    height: 100px;
    max-width: 450px;
    margin-bottom: 5px;
  }

  .step {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    flex: 1;
    position: relative;
  }

  .circle {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background-color: #ccc;
    display: flex;
    justify-content: center;
    align-items: center;
    font-weight: bold;
    color: #333;
    font-size: 14px;
    transition:
      background-color 0.3s ease,
      transform 0.2s ease;
  }

  .step.active .circle {
    background-color: #007bff;
    color: white;
    transform: scale(1.1);
  }

  .progress-line {
    flex: 1;
    height: 6px;
    background-color: #ddd;
    margin: 0 5px;
    transition: background-color 0.3s;
  }

  .progress-line.filled {
    background-color: #007bff;
  }

  p {
    font-size: 0.8rem;
    font-weight: 600;
    color: #222;
    margin-top: 5px;
  }

  /* Improve focus for keyboard navigation */
  .step:focus {
    outline: 2px solid #007bff;
    outline-offset: 4px;
  }

  @media (prefers-reduced-motion: reduce) {
    .circle,
    .progress-line {
      transition: none;
    }
  }
</style>
