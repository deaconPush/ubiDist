<script>
  let passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{12,}$/;
  function isPasswordValid(password) {
    return passwordRegex.test(password);
  }

  function checkPasswordStrength() {
    const password = document.getElementById('wallet-password').value;
    const passwordStrengthContainer = document.getElementById('password-strength');
    const isPasswordStrong = isPasswordValid(password);
    passwordStrengthContainer.style.display = 'block';
    if (isPasswordStrong) {
      passwordStrengthContainer.textContent = 'Password strength: Strong';
     } else {
        passwordStrengthContainer.textContent = 'Password strength: Weak';
      }
    }

  function checkPasswordsMatch() {
    const password = document.getElementById('wallet-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;
    const isPasswordStrong = isPasswordValid(password);
    const createWalletButton = document.getElementById('create-wallet-button');
    const passwordMatchContainer = document.getElementById('password-match');
    passwordMatchContainer.style.display = 'block';
    if (password !== confirmPassword) {
      passwordMatchContainer.textContent = 'Passwords do not match';
    } else {
      passwordMatchContainer.textContent = '';
      if (isPasswordStrong) {
        createWalletButton.disabled = false;
        createWalletButton.style.backgroundColor = 'black';
      }
    }
  } 
  
  </script>
  
  <main>
    <div class="header-content">
      <h2>Create Wallet</h2>
      <h3>Create Password</h3>
      <p>This password will unlock your wallet on this application.</p>
    </div>
  
    <form class="form-container">
      <div class="input-group">
        <label for="password">New password (16 characters min, at least one character should be in lowercase, one in uppercase and one number)</label>
        <input type="password" id="wallet-password" on:input={checkPasswordStrength} />
        <p id="password-strength" class="validation-label" ></p>
      </div>
      <div class="input-group">
        <label for="confirm-password">Confirm password</label>
        <input type="password" id="confirm-password" on:input={checkPasswordsMatch} />
        <p id="password-match" class="validation-label" ></p>
      </div>
      <button id="create-wallet-button" type="button" disabled=true>Create a new wallet</button>
    </form>
  </main>
  
  <style>
  main {
    color: greenyellow;
    font-family: "Nunito", sans-serif;
    height: 100vh;
    padding: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;
  }
  
  .header-content {
    text-align: center;
    margin-bottom: 20px; /* Adds space between the header and form */
  }
  
  h2 {
    font-size: 2em;
    margin-bottom: 10px;
  }
  
  h3 {
    font-size: 1.5em;
    margin-bottom: 10px;
  }
  
  p {
    font-size: 1.3em;
    font-weight: bold;
    margin-bottom: 20px;
  }
  
  .form-container {
    display: flex;
    flex-direction: column;
    gap: 20px; /* Adds space between input groups */
    width: 100%;
    max-width: 400px;
  }
  
  .input-group {
    display: flex;
    flex-direction: column;
    gap: 5px; /* Space between label and input */
  }
  
  label {
    font-size: 1em;
    font-weight: bold;
  }
  
  input {
    padding: 10px;
    font-size: 1em;
    border: 1px solid #ccc;
    border-radius: 5px;
    width: 100%;
  }
  
  #create-wallet-button {
    text-align: center; 
    justify-content: center;
    color: white;
    background-color: black;
    font-family: "Nunito";
    font-weight: bold;
    width: 100%;
    height: 40px;
    border-radius: 30px;
    cursor: pointer;
    border: none;
    transition: background-color 0.3s;
  }

  #create-wallet-button:disabled {
    background-color: lightslategray;
  }
  
  #create-wallet-button:hover {
    background-color: #333; /* Slightly lighter black on hover */
  }

  .validation-label {
    padding-top: 10px;
    font-size: 0.9em;
    display: none;
    color: red;
    margin-top: -10px;
  }
  </style>
  