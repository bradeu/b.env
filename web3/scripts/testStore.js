const hre = require("hardhat");
const deployData = require('./deployData.json');
const generateMerkleTree = require('./generateKeys');  

async function main() {
    console.log("\n=== Testing Storage Transaction ===\n");

    const [signer] = await hre.ethers.getSigners();
    const signerAddress = await signer.getAddress();
    console.log("Signer address:", signerAddress);

    const authorizedAddresses = [
        signerAddress,  // Use actual signer address
        "0x2345678901234567890123456789012345678901",
        "0x3456789012345678901234567890123456789012"
    ];

    const { root, proofs } = await generateMerkleTree(authorizedAddresses);
    
    const Verifier = await hre.ethers.getContractFactory("Verifier");
    const verifier = await Verifier.deploy(root);
    await verifier.waitForDeployment();
    console.log("New Verifier deployed with correct root");

    const SecretStorage = await hre.ethers.getContractFactory("SecretStorage");
    const secretStorage = await SecretStorage.deploy(await verifier.getAddress());
    await secretStorage.waitForDeployment();
    console.log("New SecretStorage deployed");

    const proof = proofs[signerAddress];
    
    console.log("\nTransaction Details:");
    console.log("Signer Address:", signerAddress);
    console.log("Proof:", proof);

    try {
        console.log("\nAttempting to store API key...");
        const tx = await secretStorage.storeEncryptedApiKey(
            signerAddress,  // Store for signer
            "test-api-key",
            proof,
            { gasLimit: 500000 }
        );
        
        console.log("Transaction sent:", tx.hash);
        const receipt = await tx.wait();
        console.log("Transaction successful!");
    } catch (error) {
        console.error("\nTransaction failed!");
        console.error("Error message:", error.message);
    }
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
}); 