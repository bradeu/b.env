import React, { useEffect, useRef } from "react";
import { words } from "./data";

import "./loader.css";
import { introAnimation, collapseWords } from "./animation";

const Loader = ({ timeline }: { timeline?: GSAPTimeline }) => {
  const loaderRef = useRef(null);
  const progressRef = useRef(null);
  const progressNumberRef = useRef(null);
  const wordGroupsRef = useRef(null);

  useEffect(() => {
    if (timeline) {
      timeline
        .add(introAnimation(wordGroupsRef))
        .add(collapseWords(loaderRef), "-=1");
    }
  }, [timeline]);

  return (
    <div className="loader__wrapper">
      <div className="loader__progressWrapper">
        <div className="loader__progress" ref={progressRef}></div>
        <span className="loader__progressNumber" ref={progressNumberRef}>
          0
        </span>
      </div>
      <div className="loader" ref={loaderRef}>
        <div className="loader__words">
          <div className="loader__overlay"></div>
          <div ref={wordGroupsRef} className="loader__wordsGroup">
            {words.map((word, index) => {
              return (
                <span key={index} className="loader__word">
                  {word}
                </span>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Loader;
