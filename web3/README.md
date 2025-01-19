# Web3 Secret Sharing

A blockchain-based secret sharing platform that allows users to store and retrieve encrypted secrets using MetaMask authentication.

## Smart Contract Features

- Store encrypted secrets on the blockchain
- Retrieve secrets (only by the owner)
- Access control using MetaMask authentication
- Secure storage with encryption

## Project Structure

```
web3-secret-sharing/
├── contracts/               # Smart contracts
│   └── SecretStorage.sol    # The main contract for storing encrypted secrets
├── scripts/                 # Deployment scripts
│   └── deploy.js           # Script to deploy the contract
├── test/                   # Test cases
│   └── SecretStorage.test.js
└── hardhat.config.js       # Hardhat configuration
```

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Create a `.env` file:
   ```bash
   cp .env.example .env
   ```
   Fill in your environment variables.

3. Compile contracts:
   ```bash
   npx hardhat compile
   ```

4. Run tests:
   ```bash
   npx hardhat test
   ```

5. Deploy to network:
   ```bash
   npx hardhat run scripts/deploy.js --network <network-name>
   ```

## Testing

Run the test suite:
```bash
npx hardhat test
```

## Security

- All secrets are encrypted before being stored on the blockchain
- Only the owner of a secret can retrieve it
- Uses OpenZeppelin's security contracts
- Implements reentrancy protection