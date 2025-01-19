"use client";

import { Center, Heading } from "@chakra-ui/react";

import { useGSAP } from "@gsap/react";
import gsap from "gsap";
import { ScrollTrigger } from "gsap/ScrollTrigger";
import { useEffect } from "react";

gsap.registerPlugin(useGSAP, ScrollTrigger);

export default function Navbar() {
  useEffect(() => {
    const showAnim = gsap
      .from(".navbar", {
        yPercent: -100,
        paused: true,
        duration: 0.2,
      })
      .progress(1);

    ScrollTrigger.create({
      start: "top top",
      end: "max",
      onUpdate: (self) => {
        if (self.direction === -1) {
          showAnim.play();
        } else {
          showAnim.reverse();
        }
      },
    });
  });

  return (
    <Center
      className="navbar"
      position={"fixed"}
      top={0}
      width={"100%"}
      height={"4rem"}
      zIndex={999}
    >
      <Heading fontSize={"2xl"} mt={5}>
        {"<"}b.env{" />"}
      </Heading>
    </Center>
  );
}
