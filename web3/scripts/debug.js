const hre = require("hardhat");
const deployData = require('./deployData.json');

async function main() {
    console.log("\n=== Debug Information ===\n");

    // Get Verifier contract
    const verifier = await hre.ethers.getContractAt(
        "Verifier",
        deployData.addresses.verifier
    );

    // Get current root from contract
    const actualRoot = await verifier.merkleRoot();
    
    console.log("Root Comparison:");
    console.log("Expected root:", deployData.root);
    console.log("Actual root  :", actualRoot);
    console.log("Roots match? :", actualRoot === deployData.root);

    // Test address and proof
    const testAddress = deployData.authorizedAddresses[0];
    const proof = deployData.proofs[testAddress];

    console.log("\nTest Address:", testAddress);
    console.log("Proof:", proof);

    // Try direct verification
    try {
        const isValid = await verifier.verify(testAddress, proof);
        console.log("\nVerification result:", isValid);
    } catch (error) {
        console.log("\nVerification failed with error:", error.message);
    }

    // Get SecretStorage contract
    const secretStorage = await hre.ethers.getContractAt(
        "SecretStorage",
        deployData.addresses.secretStorage
    );

    // Check if verifier address matches
    const storedVerifierAddress = await secretStorage.verifier();
    console.log("\nVerifier Address Check:");
    console.log("Expected verifier:", deployData.addresses.verifier);
    console.log("Stored verifier :", storedVerifierAddress);
    console.log("Addresses match?:", storedVerifierAddress === deployData.addresses.verifier);
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
}); 