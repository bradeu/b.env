const { MerkleTree } = require('merkletreejs');
const { ethers } = require('ethers');

async function generateMerkleTree(authorizedAddresses) {
    // Convert addresses to leaves using ethers v6 syntax
    const leaves = authorizedAddresses.map(addr => 
        ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(['address'], [addr]))
    );

    // Create tree
    const tree = new MerkleTree(leaves, ethers.keccak256, { sortPairs: true });
    const root = tree.getHexRoot();

    // Generate proofs for each address
    const proofs = {};
    authorizedAddresses.forEach(addr => {
        const leaf = ethers.keccak256(ethers.AbiCoder.defaultAbiCoder().encode(['address'], [addr]));
        proofs[addr] = tree.getHexProof(leaf);
    });

    return {
        root,
        proofs
    };
}

// Example usage
async function main() {
    const authorizedAddresses = [
        "0x1234567890123456789012345678901234567890",
        "0x2345678901234567890123456789012345678901",
        // Add more authorized addresses here
    ];

    const { root, proofs } = await generateMerkleTree(authorizedAddresses);
    console.log('Merkle Root:', root);
    console.log('Proofs:', JSON.stringify(proofs, null, 2));
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
