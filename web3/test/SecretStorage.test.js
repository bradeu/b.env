const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("SecretStorage", function () {
  let secretStorage;
  let owner;
  let addr1;
  let addr2;

  beforeEach(async function () {
    [owner, addr1, addr2] = await ethers.getSigners();
    
    const SecretStorage = await ethers.getContractFactory("SecretStorage");
    secretStorage = await SecretStorage.deploy();
  });

  describe("Storing Secrets", function () {
    it("Should store a secret successfully", async function () {
      const secretId = ethers.keccak256(ethers.toUtf8Bytes("test-secret-id"));
      const encryptedContent = "encrypted-content";

      await secretStorage.storeSecret(secretId, encryptedContent);
      
      expect(await secretStorage.doesSecretExist(secretId)).to.equal(true);
    });

    it("Should not allow storing a secret with existing ID", async function () {
      const secretId = ethers.keccak256(ethers.toUtf8Bytes("test-secret-id"));
      const encryptedContent = "encrypted-content";

      await secretStorage.storeSecret(secretId, encryptedContent);
      
      await expect(
        secretStorage.storeSecret(secretId, "new-content")
      ).to.be.revertedWith("Secret ID already exists");
    });
  });

  describe("Retrieving Secrets", function () {
    it("Should allow owner to retrieve their secret", async function () {
      const secretId = ethers.keccak256(ethers.toUtf8Bytes("test-secret-id"));
      const encryptedContent = "encrypted-content";

      await secretStorage.connect(addr1).storeSecret(secretId, encryptedContent);
      
      const retrievedSecret = await secretStorage.connect(addr1).getSecret(secretId);
      expect(retrievedSecret).to.equal(encryptedContent);
    });

    it("Should not allow non-owners to retrieve secret", async function () {
      const secretId = ethers.keccak256(ethers.toUtf8Bytes("test-secret-id"));
      const encryptedContent = "encrypted-content";

      await secretStorage.connect(addr1).storeSecret(secretId, encryptedContent);
      
      await expect(
        secretStorage.connect(addr2).getSecret(secretId)
      ).to.be.revertedWith("Only the owner can retrieve the secret");
    });
  });
}); 