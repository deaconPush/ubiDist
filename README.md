# UBI DEFI
This repository is intended to consolidate my knowledge and proof of concepts (POCs) in blockchain technology.

## Table of contents
- [Hardhat](#hardhat)
- [Wails](#wails)


## Hardhat
This folder contains basic Hardhat scripts for deploying contracts on the Hardhat local network. Detailed instructions for running a Hardhat node or deploying smart contracts to the local network can be found in the hardhat folder's README.

## Wails
This repository contains applications built using the [Wails](https://wails.io/) framework, which enables the development of cross-platform desktop applications with Go. Below is an overview of the folder structure and included projects.

### Applications

#### 1. **Hello World**
Located in the `helloworld` folder, this project is a basic Wails application that follows the official [Wails documentation tutorials](https://wails.io/docs/gettingstarted). It is designed to help you get started with Wails and understand its core concepts.

#### 2. **Dogs API**
The `dogs-api` project is another tutorial-based example that demonstrates how to create a Wails application that interacts with APIs. This project builds upon the foundational knowledge from the `helloworld` project and showcases more advanced features.

#### 3. **Wallet**
The `wallet` project is a desktop application designed for managing transactions. It allows users to:
- **Receive funds**
- **Transfer funds**
- **Interact with different networks**

> **Note:**  
> Currently, only the **Hardhat local network** is supported for development and testing purposes.