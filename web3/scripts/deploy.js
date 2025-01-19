const hre = require("hardhat");
const { generateMerkleTree } = require("./generateKeys");

async function main() {
    console.log("Deploying contracts...");

    // Create array of authorized addresses
    const authorizedAddresses = [
        "0xf4061f4486122C19cd2d3521B04e884890d4b08e"
    ];

    // Get merkle tree data
    const { root, proofs } = await generateMerkleTree(authorizedAddresses);
    console.log("Merkle Root:", root);
    console.log("Proofs:", proofs);

    // First deploy Verifier with just the root
    const Verifier = await hre.ethers.getContractFactory("Verifier");
    const verifier = await Verifier.deploy(root);  // Pass only the root
    await verifier.waitForDeployment();
    const verifierAddress = await verifier.getAddress();
    console.log("Verifier deployed to:", verifierAddress);

    // Then deploy SecretStorage with Verifier's address
    const SecretStorage = await hre.ethers.getContractFactory("SecretStorage");
    const secretStorage = await SecretStorage.deploy(verifierAddress);  // Pass verifier's address
    await secretStorage.waitForDeployment();
    console.log("SecretStorage deployed to:", await secretStorage.getAddress());

    // Save proofs for later use
    console.log("\nSave these proofs for later use:");
    console.log(JSON.stringify(proofs, null, 2));
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});