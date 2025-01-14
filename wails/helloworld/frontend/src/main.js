import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import { GenerateWord } from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
        <div class="result" id="result"></div>
        <button class="btn" onclick="generateWord()">Click to generate a word</button>
      </div>
    </div>
`;
document.getElementById('logo').src = logo;

let resultElement = document.getElementById("result");
// Setup the greet function
window.generateWord = function () {
    try {
        GenerateWord()
        .then((result, err) => {
            if (err) {
                resultElement.innerText = "Error generating your word, sorry" ;
            }
            else {
                resultElement.innerText = result;
            }
            
        })
        .catch((err) => {
            console.error(err);
        });
        
    } catch (err) { 
        console.error(err);
    }
};