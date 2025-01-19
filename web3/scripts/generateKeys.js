const { MerkleTree } = require('merkletreejs');
const { ethers } = require('ethers');

async function generateMerkleTree(authorizedAddresses) {
    // Ensure we have at least two addresses for proper proof generation
    const defaultAddresses = [
        "0xf4061f4486122C19cd2d3521B04e884890d4b08e",  // Your address
        "0x2345678901234567890123456789012345678901",  // Another address
        "0x3456789012345678901234567890123456789012"   // Third address
    ];

    // Use provided addresses or default ones
    const addresses = authorizedAddresses.length >= 2 ? authorizedAddresses : defaultAddresses;

    console.log("Creating Merkle tree with addresses:", addresses);

    const leaves = addresses.map(addr => 
        ethers.keccak256(ethers.solidityPacked(['address'], [addr]))
    );

    const tree = new MerkleTree(leaves, ethers.keccak256, { sortPairs: true });
    const root = tree.getHexRoot();

    const proofs = {};
    addresses.forEach(addr => {
        const leaf = ethers.keccak256(ethers.solidityPacked(['address'], [addr]));
        proofs[addr] = tree.getHexProof(leaf);
    });

    // Log the tree structure for debugging
    console.log("Merkle Tree Structure:");
    console.log(tree.toString());
    console.log("Root:", root);

    return {
        root,
        proofs
    };
}

async function testMatchingLogic() {
    const address = "0x1234567890123456789012345678901234567890";
    
    // Generate tree and proof using generateKeys logic
    const { root, proofs } = await generateMerkleTree([
        address,
        "0x2345678901234567890123456789012345678901"
    ]);
    
    console.log("\nVerification Test:");
    console.log("1. Address:", address);
    console.log("2. Root:", root);
    console.log("3. Proof:", proofs[address]);
    
    // Simulate Verifier.sol's verification logic
    let computedHash = ethers.keccak256(ethers.solidityPacked(['address'], [address]));
    console.log("4. Initial hash:", computedHash);
    
    for (let i = 0; i < proofs[address].length; i++) {
        const proofElement = proofs[address][i];
        // Sort just like Verifier.sol does
        if (computedHash < proofElement) {
            computedHash = ethers.keccak256(ethers.concat([computedHash, proofElement]));
        } else {
            computedHash = ethers.keccak256(ethers.concat([proofElement, computedHash]));
        }
    }
    
    console.log("5. Final computed hash:", computedHash);
    console.log("6. Matches root?", computedHash === root);
}

// Example usage
async function main() {
    await testMatchingLogic();
}

// ensure that main only runs when this file is executed directly
if (require.main === module) {
    main()
        .then(() => process.exit(0))
        .catch(error => {
            console.error(error);
            process.exit(1);
        });
}

// Export the function directly
module.exports = generateMerkleTree;
