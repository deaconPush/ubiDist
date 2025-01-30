<script>
      export let handleClick = () => {};
      export let walletLabel = '';

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
      const walletButton = document.getElementById('wallet-button');
      const passwordMatchContainer = document.getElementById('password-match');
      passwordMatchContainer.style.display = 'block';
      if (password !== confirmPassword) {
        passwordMatchContainer.textContent = 'Passwords do not match';
      } else {
        passwordMatchContainer.textContent = '';
        if (isPasswordStrong) {
          walletButton.disabled = false;
          walletButton.style.backgroundColor = 'black';
        }
      }
    }   
    </script>
    
<div class="header-content">
  <h2>Create Password</h2>
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
  <button id="wallet-button" type="button" disabled=true on:click={handleClick}>{walletLabel}</button>
</form>
    
    <style>
    .header-content {
      text-align: center;
      margin-bottom: 20px;
    }
  
    h2 {
      font-size: 2rem;
      color: #0066cc; 
      margin-bottom: 10px;
    }
  
    p {
      font-size: 1rem;
      color: #555; 
      margin-bottom: 20px;
    }
  
    /* Form styles */
    .form-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      align-content: center;
      gap: 20px;
      width: 100%;
      max-width: 400px;
      background-color: #fff;
      padding: 20px;
      border-radius: 10px;
      box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1); 
    }
  
    .input-group {
      display: flex;
      flex-direction: column;
      gap: 2px;
      align-items: center;
    }
  
    label {
      font-size: 1rem;
      font-weight: bold;
      color: #222;
    }
  
    input {
      padding: 4px;
      font-size: 1rem;
      border: 1px solid #ccc;
      border-radius: 5px;
      max-width: 190px;
      color: #333; 
      background-color: #fefefe; 
    }
  
    input:focus {
      outline: none;
      border-color: #0066cc; 
      box-shadow: 0 0 3px #0066cc;
    }
  
    #wallet-button {
      text-align: center;
      color: white;
      background-color: #0066cc;
      font-family: "Nunito", sans-serif;
      font-weight: bold;
      width: 100%;
      height: 25px;
      border-radius: 25px;
      cursor: pointer;
      border: none;
      transition: background-color 0.3s ease;
    }
  
    #wallet-button:hover:not(:disabled) {
      background-color: #004c99; 
    }
  
    #wallet-button:disabled {
      background-color: #ccc; 
      color: #666;
      cursor: not-allowed;
    }
  
    /* Validation messages */
    .validation-label {
      font-size: 0.9rem;
      color: red; 
      display: none; 
    }
  
    /* Responsive styles */
    @media (max-width: 768px) {
      main {
        padding: 15px;
      }
  
      h2 {
        font-size: 1.8rem;
      }
  
      h3 {
        font-size: 1.3rem;
      }
  
      .form-container {
        padding: 15px;
      }
  
      #wallet-button {
        height: 40px;
      }
    }
  
    @media (max-width: 480px) {
      h2 {
        font-size: 1.5rem;
      }
  
      h3 {
        font-size: 1.2rem;
      }
  
      label,
      input,
      #wallet-button {
        font-size: 0.9rem;
      }
  
      .form-container {
        gap: 15px;
      }
    }
    </style>
    