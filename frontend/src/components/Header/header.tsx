import React from "react";
import "./header.css";

function Header() {
  return (
    <div>
      <div className="head">
        <span className="intro" style={{ width: "50%", float: "left" }}>
          <b>Hey, I'm Peter</b>
        </span>

        <div style={{ width: "50%", float: "right", textAlign: "right" }}>
          <span className="github">
            <a href="https://github.com/PeteMango" target="_blank">
              [github]
            </a>{" "}
          </span>{" "}
          <span className="linkedin">
            <a href="https://www.linkedin.com/in/p25wang/" target="_blank">
              [linkedin]
            </a>{" "}
          </span>
          <span className="spotify">
            <a href="https://open.spotify.com/user/whcpeterwang?si=ceb6ac4ad731465b" target="_blank">
              [spotify]
            </a>{" "}
          </span>
          <span className="mail">
            <a href="mailto:p25wang@uwaterloo.ca" target="_blank">
              [@]
            </a>
          </span>
        </div>
      </div>
      <hr />
    </div>
  );
}

export default Header;
