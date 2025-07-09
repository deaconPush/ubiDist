FROM node:18-alpine

WORKDIR /app

COPY ./hardhat/ ./

RUN npm install

EXPOSE 8545

CMD ["npx", "hardhat", "node", "--hostname", "0.0.0.0"]
