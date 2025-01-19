const { expect } = require("chai");
const { ethers } = require("hardhat");
const { MerkleTree } = require('merkletreejs');

describe("ZKP Secret Storage System", function () {
  let verifier;
  let secretStorage;
  let owner;
  let authorizedUser;
  let unauthorizedUser;
  let merkleTree;
  let merkleRoot;

  before(async function () {
    // Get signers
    [owner, authorizedUser, unauthorizedUser] = await ethers.getSigners();

    // Create Merkle tree with authorized addresses
    const authorizedAddresses = [
      authorizedUser.address
    ];

    // Generate leaves
    const leaves = authorizedAddresses.map(addr => 
      ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(['address'], [addr]))
    );

    // Create Merkle tree
    merkleTree = new MerkleTree(leaves, ethers.utils.keccak256, { sortPairs: true });
    merkleRoot = merkleTree.getHexRoot();

    // Deploy contracts
    const Verifier = await ethers.getContractFactory("Verifier");
    verifier = await Verifier.deploy(merkleRoot);
    await verifier.deployed();

    const SecretStorage = await ethers.getContractFactory("SecretStorage");
    secretStorage = await SecretStorage.deploy(verifier.address);
    await secretStorage.deployed();
  });

  it("Should allow authorized user to store and retrieve API key", async function () {
    const proof = merkleTree.getHexProof(
      ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(['address'], [authorizedUser.address]))
    );

    // Store API key
    const encryptedApiKey = "encrypted_test_key";
    await secretStorage.connect(authorizedUser).storeEncryptedApiKey(
      authorizedUser.address,
      encryptedApiKey,
      proof
    );

    // Retrieve API key
    const retrievedKey = await secretStorage.connect(authorizedUser).getEncryptedApiKeyForAddress(
      authorizedUser.address,
      proof
    );

    expect(retrievedKey).to.equal(encryptedApiKey);
  });

  it("Should prevent unauthorized user from storing API key", async function () {
    // Generate invalid proof for unauthorized user
    const proof = merkleTree.getHexProof(
      ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(['address'], [unauthorizedUser.address]))
    );

    // Attempt to store API key
    await expect(
      secretStorage.connect(unauthorizedUser).storeEncryptedApiKey(
        unauthorizedUser.address,
        "encrypted_test_key",
        proof
      )
    ).to.be.revertedWith("Not authorized to store API keys");
  });

  it("Should verify Merkle proof correctly", async function () {
    const validProof = merkleTree.getHexProof(
      ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(['address'], [authorizedUser.address]))
    );

    const invalidProof = merkleTree.getHexProof(
      ethers.utils.keccak256(ethers.utils.defaultAbiCoder.encode(['address'], [unauthorizedUser.address]))
    );

    // Check valid proof
    const isValidProof = await verifier.verify(authorizedUser.address, validProof);
    expect(isValidProof).to.be.true;

    // Check invalid proof
    const isInvalidProof = await verifier.verify(unauthorizedUser.address, invalidProof);
    expect(isInvalidProof).to.be.false;
  });
}); 