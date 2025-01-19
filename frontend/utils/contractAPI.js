import { ethers } from 'ethers';
import SecretStorageABI from '../../web3/artifacts/contracts/SecretStorage.sol/SecretStorage.json';
import VerifierABI from '../../web3/artifacts/contracts/Verifier.sol/Verifier.json';

// Hardhat local network
const provider = new ethers.JsonRpcProvider("http://127.0.0.1:8545");

// Contract addresses from your local deployment
const SECRET_STORAGE_ADDRESS = "0x0B306BF915C4d645ff596e518fAf3F9669b97016";
const VERIFIER_ADDRESS = "0x9A676e781A523b5d0C0e43731313A708CB607508";

// Contract instances
const secretStorage = new ethers.Contract(
    SECRET_STORAGE_ADDRESS,
    SecretStorageABI.abi,
    provider
);

const verifier = new ethers.Contract(
    VERIFIER_ADDRESS,
    VerifierABI.abi,
    provider
);

export const API = {
    // Storage Functions
    async storeApiKey(address, encryptedKey, proof) {
        try {
            const tx = await secretStorage.storeEncryptedApiKey(address, encryptedKey, proof);
            await tx.wait();
            return {
                success: true,
                txHash: tx.hash
            };
        } catch (error) {
            return {
                success: false,
                error: error.message
            };
        }
    },

    async getApiKey(address, proof) {
        try {
            const key = await secretStorage.getEncryptedApiKeyForAddress(address, proof);
            return {
                success: true,
                key
            };
        } catch (error) {
            return {
                success: false,
                error: error.message
            };
        }
    },

    // Verifier Functions
    async verifyAddress(address, proof) {
        try {
            const isValid = await verifier.verify(address, proof);
            return {
                success: true,
                isValid
            };
        } catch (error) {
            return {
                success: false,
                error: error.message
            };
        }
    },

    async getMerkleRoot() {
        try {
            const root = await verifier.merkleRoot();
            return {
                success: true,
                root
            };
        } catch (error) {
            return {
                success: false,
                error: error.message
            };
        }
    }
};