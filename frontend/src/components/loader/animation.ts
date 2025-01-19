import { gsap } from "gsap";
/* eslint-disable @typescript-eslint/no-explicit-any */
export const introAnimation = (wordGroupsRef: any) => {
  const tl = gsap.timeline();
  tl.to(wordGroupsRef.current, {
    yPercent: -80,
    duration: 5,
    ease: "power3.inOut",
  });

  return tl;
};
/* eslint-disable @typescript-eslint/no-explicit-any */
export const collapseWords = (wordGroupsRef: any) => {
  const tl = gsap.timeline();
  tl.to(wordGroupsRef.current, {
    "clip-path": "polygon(0% 50%, 100% 50%, 100% 50%, 0% 50%)",
    duration: 3,
    ease: "expo.inOut",
  });

  return tl;
};
