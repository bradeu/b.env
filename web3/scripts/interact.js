const hre = require("hardhat");
const deployData = require('./deployData.json');

async function main() {
    console.log("Loading deployment data...");
    
    const SECRET_STORAGE_ADDRESS = deployData.addresses.secretStorage;
    const senderAddress = deployData.authorizedAddresses[0];  // Your address
    const proof = deployData.proofs[senderAddress];

    console.log("Using address:", senderAddress);
    console.log("Using proof:", proof);
    console.log("Contract address:", SECRET_STORAGE_ADDRESS);

    // Get contract instance
    const SecretStorage = await hre.ethers.getContractFactory("SecretStorage");
    const secretStorage = SecretStorage.attach(SECRET_STORAGE_ADDRESS);

    try {
        console.log("Storing API key...");
        const tx = await secretStorage.storeEncryptedApiKey(
            senderAddress,
            "test-encrypted-api-key",
            proof
        );
        await tx.wait();
        console.log("API Key stored successfully!");

        console.log("Retrieving API key...");
        const apiKey = await secretStorage.getEncryptedApiKeyForAddress(
            senderAddress,
            proof
        );
        console.log("Retrieved API Key:", apiKey);

    } catch (error) {
        console.error("Transaction failed!");
        console.error("Error:", error.message);
        
        // Get the verifier address and root for debugging
        const verifierAddress = await secretStorage.verifier();
        console.log("\nDebugging info:");
        console.log("Verifier address:", verifierAddress);
        console.log("Expected root:", deployData.root);
    }
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
}); 