"use client";

import ConnectWalletButton from "@/components/wallet-btn";

import { Center } from "@chakra-ui/react";

export default function Home() {
  return (
    <Center flexDir={"column"} height={"100vh"} fontSize={"7xl"}>
      <ConnectWalletButton />
    </Center>
  );
}
