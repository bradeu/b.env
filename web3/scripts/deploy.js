const hre = require("hardhat");

async function main() {
  console.log("Deploying SecretStorage contract...");

  const SecretStorage = await hre.ethers.getContractFactory("SecretStorage");
  const secretStorage = await SecretStorage.deploy();

  await secretStorage.waitForDeployment();

  console.log("SecretStorage deployed to:", await secretStorage.getAddress());
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
}); 