"use client";

import { Box, Button, Heading, Text } from "@chakra-ui/react";
import { ethers } from "ethers";
import { useEffect, useState } from "react";
import { useAccount } from "wagmi";
import Form from "./form";
declare global {
  interface Window {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    ethereum?: any;
  }
}

export default function ConnectWalletButton() {
  const [address, setAddress] = useState<string>("");
  const [isClient, setIsClient] = useState<boolean>(false);

  useEffect(() => {
    setIsClient(true);
    const storedAddress = localStorage.getItem("address");
    if (storedAddress) {
      setAddress(storedAddress);
      console.log(address);
    }
  }, [address]);

  const connectWallet = async () => {
    if (window.ethereum == null) {
      console.log("No metamask wallet installed!");
    } else {
      try {
        const provider = new ethers.BrowserProvider(window.ethereum);
        const signer = await provider.getSigner();
        const userAddress = await signer.getAddress();

        console.log(`${signer}`);

        setAddress(userAddress);
        localStorage.setItem("address", userAddress);
        // console.log(userAddress);
      } catch {
        console.log("error");
      }
    }
  };
  const { isConnected } = useAccount();

  if (!isClient) {
    return <Text fontSize={"md"}>Loading...</Text>;
  }

  return (
    <Box>
      {isConnected ? (
        <Box
          width={"10em"}
          height={"5em"}
          my={5}
          background={"#212121"}
          p={10}
          rounded={"2xl"}
        >
          {/* <Text fontSize={"md"}>Success, your address is {`${address}!`}</Text> */}
          <Form />
        </Box>
      ) : (
        <Box>
          <Heading>Connec to your Wallet!</Heading>
          <Button onClick={connectWallet}>Sign in with Metamask</Button>
        </Box>
      )}
    </Box>
  );
}
