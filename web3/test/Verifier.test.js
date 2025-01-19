const { expect } = require("chai");
const { ethers } = require("hardhat");
const { MerkleTree } = require('merkletreejs');

describe("Verifier", function () {
  let verifier;
  let owner;
  let user1;
  let user2;
  let merkleTree;
  let merkleRoot;

  beforeEach(async function () {
    // Get signers
    [owner, user1, user2] = await ethers.getSigners();

    // Create Merkle tree with test addresses
    const authorizedAddresses = [user1.address];
    
    // Generate leaves - using the same encoding as the contract
    const leaves = authorizedAddresses.map(addr => 
      ethers.keccak256(ethers.solidityPacked(['address'], [addr]))
    );

    // Create Merkle tree
    merkleTree = new MerkleTree(leaves, ethers.keccak256, { sortPairs: true });
    merkleRoot = merkleTree.getHexRoot();

    // Deploy contract
    const Verifier = await ethers.getContractFactory("Verifier");
    verifier = await Verifier.deploy(merkleRoot);
  });

  describe("Verification", function () {
    it("Should verify valid proof for authorized address", async function () {
      const leaf = ethers.keccak256(ethers.solidityPacked(['address'], [user1.address]));
      const proof = merkleTree.getHexProof(leaf);

      const isValid = await verifier.verify(user1.address, proof);
      expect(isValid).to.be.true;
    });

    it("Should reject invalid proof for unauthorized address", async function () {
      const leaf = ethers.keccak256(ethers.solidityPacked(['address'], [user2.address]));
      const proof = merkleTree.getHexProof(leaf);

      const isValid = await verifier.verify(user2.address, proof);
      expect(isValid).to.be.false;
    });

    it("Should reject empty proof", async function () {
      const isValid = await verifier.verify(user1.address, []);
      expect(isValid).to.be.false;
    });
  });

  describe("Merkle Root Management", function () {
    it("Should allow owner to update merkle root", async function () {
      const newAuthorizedAddresses = [user2.address];
      const newLeaves = newAuthorizedAddresses.map(addr => 
        ethers.keccak256(ethers.solidityPacked(['address'], [addr]))
      );
      const newTree = new MerkleTree(newLeaves, ethers.keccak256, { sortPairs: true });
      const newRoot = newTree.getHexRoot();

      await expect(verifier.connect(owner).updateMerkleRoot(newRoot))
        .to.emit(verifier, "MerkleRootUpdated")
        .withArgs(merkleRoot, newRoot);

      expect(await verifier.merkleRoot()).to.equal(newRoot);
    });

    it("Should prevent non-owner from updating merkle root", async function () {
      const newRoot = ethers.keccak256("0x1234");
      
      await expect(
        verifier.connect(user1).updateMerkleRoot(newRoot)
      ).to.be.revertedWithCustomError(verifier, "OwnableUnauthorizedAccount");
    });
  });
}); 