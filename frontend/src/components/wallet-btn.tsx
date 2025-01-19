"use client";

import { Box, Button, Heading, Text } from "@chakra-ui/react";
import { useGSAP } from "@gsap/react";
import { ethers } from "ethers";
import gsap from "gsap";
import { useEffect, useRef, useState } from "react";
import { useAccount } from "wagmi";
import Form from "./form";

declare global {
  interface Window {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    ethereum?: any;
  }
}

gsap.registerPlugin(useGSAP);

export default function ConnectWalletButton() {
  const [address, setAddress] = useState<string>("");

  const containerRef = useRef<HTMLDivElement>(null);

  const timeline = gsap.timeline();
  useGSAP(() => {
    if (timeline) {
      timeline.to(containerRef.current, {
        clipPath: "polygon(0% 0%, 100% 0%, 100% 100%, 0% 100%)",
        duration: "1",
        ease: "power4.inOut",
      });
    }
  });

  useEffect(() => {
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

  return (
    <Box height={"100vh"}>
      {isConnected ? (
        <Box
          ref={containerRef}
          clipPath={"polygon(0% 0%, 100% 0%, 100% 0%, 0% 0%)"}
          width={{ base: "100%", md: "10em" }}
          my={5}
          background={"#212121"}
          p={{ base: 8, md: 10 }}
          rounded={"2xl"}
          border={"2px solid #333"}
          color={"white"}
        >
          {/* <Text fontSize={"md"}>Success, your address is {`${address}!`}</Text> */}
          <Text fontSize={"md"}>
            Your address is {`${address.toString().substring(0, 9)}`}... !
          </Text>
          <Heading mb={10}>Enter you API</Heading>
          <Form address={address} />
        </Box>
      ) : (
        <Box>
          <Heading>Connect to your Wallet!</Heading>
          <Button onClick={connectWallet}>Sign in with Metamask</Button>
        </Box>
      )}
    </Box>
  );
}
