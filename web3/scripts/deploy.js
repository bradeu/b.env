const hre = require("hardhat");
import { generateMerkleTree } from "./generateKeys";

async function main() {
  console.log("Deploying SecretStorage contract...");

  const SOME_KEY = 0x10000; //this should be changed to reak frontend.
  const root = generateMerkleTree(SOME_KEY);
  const SecretStorage = await hre.ethers.getContractFactory("SecretStorage");
  const secretStorage = await SecretStorage.deploy(root);

  await secretStorage.waitForDeployment();

  console.log("SecretStorage deployed to:", await secretStorage.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
}); 