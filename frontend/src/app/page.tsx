"use client";

import ConnectWalletButton from "@/components/wallet-btn";

import { Box, Center } from "@chakra-ui/react";
import gsap from "gsap";
import { useGSAP } from "@gsap/react";
import Loader from "@/components/loader";
import { useState } from "react";
import Landing from "@/components/landing";

gsap.registerPlugin(useGSAP);

export default function Home() {
  const [loaderFinished, setLoaderFinished] = useState(false);
  const [timeline, setTimeline] = useState<GSAPTimeline>();

  useGSAP(() => {
    const tl = gsap.timeline({
      onComplete: () => setLoaderFinished(true),
    });
    setTimeline(tl);
  });

  return (
    <Center
      position={"relative"}
      flexDir={"column"}
      height={"100vh"}
      color={"#212121"}
      fontSize={"7xl"}
    >
      {loaderFinished ? (
        <Box position={"relatives"}>
          <Landing />
          <ConnectWalletButton />
        </Box>
      ) : (
        <Loader timeline={timeline} />
      )}
    </Center>
  );
}
