"use client";

import ConnectWalletButton from "@/components/wallet-btn";
// import { WalletOptions } from "@/components/wallet-options";
import { Center, Heading } from "@chakra-ui/react";



export default function Home() {
  return (
    <Center flexDir={"column"} height={"100vh"} fontSize={"7xl"}>
      <Heading>Sign In</Heading>
      <ConnectWalletButton />
    </Center>
  );
}
