# UBI DEFI
This repository is intended to consolidate my knowledge and proof of concepts (POCs) in blockchain technology.

## Table of contents
- [Installation](#installation)
- [Running hardhat scripts](#running-hardhat-scripts)
- [Hardhat local deployment](#hardhat-local-deployment)


## Installation

```
nvm use
yarn install
```


## Running hardhat scripts

```
yarn hardhat compile
yarn hardhat run <script-name>
```

## Hardhat local deployment

* Run the local network
```
yarn hardhat node
```

* Deploy the module into the local network
```
yarn hardhat ignition deploy ignition/modules/Rocket.js --network localhost
```
 
 * Run a script to interact with the deployed contract
```
 yarn hardhat run scripts/<script-file> --network localhost
 ```