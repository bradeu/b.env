import { ethers } from 'ethers';
import SecretStorageABI from '../../web3/artifacts/contracts/SecretStorage.sol/SecretStorage.json';

// Your deployed contract address
const CONTRACT_ADDRESS = "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512";  // Your actual contract address

// Create a provider (use Infura, Alchemy, or other providers)
const INFURA_ID = "c7fb4a919d8c46cca86fbbb373273d95";  // Get this from Infura
const provider = new ethers.JsonRpcProvider(`https://sepolia.infura.io/v3/${INFURA_ID}`);

// Create a read-only contract instance
const contract = new ethers.Contract(
    CONTRACT_ADDRESS,
    SecretStorageABI.abi,
    provider
);

// Public API-like functions
export const contractAPI = {
    // Read functions (no wallet needed)
    async getPublicData(address) {
        try {
            const data = await contract.getPublicData(address);
            return data;
        } catch (error) {
            console.error("Error fetching public data:", error);
            throw error;
        }
    },

    // Write functions (wallet needed)
    async storeData(data, proof) {
        try {
            // For write operations, we need a signer
            if (typeof window.ethereum === 'undefined') {
                throw new Error("Please install MetaMask!");
            }

            const provider = new ethers.BrowserProvider(window.ethereum);
            const signer = await provider.getSigner();
            const contractWithSigner = contract.connect(signer);

            const tx = await contractWithSigner.storeEncryptedApiKey(
                await signer.getAddress(),
                data,
                proof
            );
            await tx.wait();
            return tx.hash;
        } catch (error) {
            console.error("Error storing data:", error);
            throw error;
        }
    },

    // Get contract address
    getContractAddress() {
        return CONTRACT_ADDRESS;
    },

    // Get ABI
    getABI() {
        return SecretStorageABI.abi;
    }
};