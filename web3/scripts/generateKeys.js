const { MerkleTree } = require('merkletreejs');
const { ethers } = require('ethers');

async function generateMerkleTree(authorizedAddresses) {
    // Convert addresses to leaves using ethers v6 syntax
    const leaves = authorizedAddresses.map(addr => 
        ethers.keccak256(ethers.solidityPacked(['address'], [addr]))
    );

    // Create tree
    const tree = new MerkleTree(leaves, ethers.keccak256, { sortPairs: true });
    const root = tree.getHexRoot();

    // Generate proofs for each address
    const proofs = {};
    authorizedAddresses.forEach(addr => {
        const leaf = ethers.keccak256(ethers.solidityPacked(['address'], [addr]));
        proofs[addr] = tree.getHexProof(leaf);
    });

    return {
        root,
        proofs
    };
}

async function verifyStep(address, proofArray) {
    // Create the leaf
    const leaf = ethers.keccak256(
        ethers.solidityPacked(['address'], [address])
    );
    console.log('Leaf:', leaf);

    if (!proofArray || proofArray.length === 0) {
        console.log('No proof elements (single-node tree)');
        console.log('Root should equal leaf:', leaf);
        return;
    }

    // Get the first proof element
    const proof = proofArray[0];
    console.log('Proof Element:', proof);
    
    // Sort the hashes
    const pair = [leaf, proof].sort();
    
    // Hash the concatenated values
    const calculatedRoot = ethers.keccak256(
        ethers.concat(pair)
    );
    
    console.log('Calculated Root:', calculatedRoot);
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

export default generateMerkleTree