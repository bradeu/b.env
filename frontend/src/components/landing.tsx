import {
  Center,
  Container,
  // Grid,
  // GridItem,
  // GridItemProps,
  Heading,
} from "@chakra-ui/react";
import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { useEffect, useRef } from "react";

gsap.registerPlugin(useGSAP);

export default function Landing() {
  const titleRef = useRef<HTMLHeadingElement>(null);
  // useGSAP(() => {
  //   gsap.timeline({ paused: true }).from(".title", {
  //     y: 0,
  //     duration: 1,
  //     ease: "power4.inOut",
  //   });
  // }, []);
  useEffect(() => {
    gsap.to(titleRef.current, {
      scale: 1,
      duration: 0.5,
      ease: "power4.inOut",
    });
  }, []);
  // const itemProps: GridItemProps = {
  //   w: "100%",
  //   bg: "black",
  //   width: { base: "100%", md: "10em" },
  //   height: "fit-content",
  //   my: 5,
  //   background: "#212121",
  //   p: { base: 8, md: 10 },
  //   rounded: "2xl",
  // };
  return (
    <Center
      as={Container}
      position={"relative"}
      height={"100%"}
      flexDir={"column"}
    >
      <Heading
        ref={titleRef}
        fontSize={"5xl"}
        color={"white"}
        fontWeight={"initial"}
        scale={0}
      >
        Safe-to-use .env solution
      </Heading>
      {/* <Grid
        templateColumns={"repeat(5, 1fr)"}
        position={"relative"}
        w={"100%"}
        color={"white"}
      >
        <GridItem colSpan={3} {...itemProps}>
          1
        </GridItem>
        <GridItem colSpan={2} {...itemProps}>
          2
        </GridItem>
        <GridItem colSpan={2} {...itemProps}>
          3
        </GridItem>
        <GridItem colSpan={3} {...itemProps}>
          4
        </GridItem>
      </Grid> */}
    </Center>
  );
}
