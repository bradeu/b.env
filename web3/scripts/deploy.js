const hre = require("hardhat");
const generateMerkleTree = require("./generateKeys");
const fs = require('fs');

async function main() {
    console.log("Deploying contracts...");

    // IMPORTANT: Save these exact addresses
    const authorizedAddresses = [
        "0xf4061f4486122C19cd2d3521B04e884890d4b08e",  // Your address
        "0x2345678901234567890123456789012345678901",
        "0x3456789012345678901234567890123456789012"
    ];

    const { root, proofs } = await generateMerkleTree(authorizedAddresses);
    console.log("Merkle Root:", root);
    console.log("Proofs:", JSON.stringify(proofs, null, 2));

    // Save the deployment data
    const deployData = {
        root,
        proofs,
        authorizedAddresses
    };
    
    fs.writeFileSync(
        './scripts/deployData.json', 
        JSON.stringify(deployData, null, 2)
    );

    // Deploy Verifier with the root
    const Verifier = await hre.ethers.getContractFactory("Verifier");
    const verifier = await Verifier.deploy(root);
    await verifier.waitForDeployment();
    const verifierAddress = await verifier.getAddress();
    console.log("Verifier deployed to:", verifierAddress);

    // Deploy SecretStorage
    const SecretStorage = await hre.ethers.getContractFactory("SecretStorage");
    const secretStorage = await SecretStorage.deploy(verifierAddress);
    await secretStorage.waitForDeployment();
    const secretStorageAddress = await secretStorage.getAddress();
    console.log("SecretStorage deployed to:", secretStorageAddress);

    // Save addresses
    deployData.addresses = {
        verifier: verifierAddress,
        secretStorage: secretStorageAddress
    };
    
    fs.writeFileSync(
        './scripts/deployData.json', 
        JSON.stringify(deployData, null, 2)
    );
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});