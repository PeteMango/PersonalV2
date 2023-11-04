import React from "react";
import "./body.css";

function Body() {
  return (
    <div>
      <div className="body">
        <div>
          I'm a 2A <span className="highlight">Software Engineering</span>{" "}
          student @{" "}
          <span className="school">
            <a
              href="https://uwaterloo.ca/future-students/programs/software-engineering"
              target="_blank"
            >
              UWaterloo
            </a>
          </span>
        </div>
        <br />
        <div>
          Incoming:
          <ul className="list">
            <li>
              <span className="highlight">software developer</span> @{" "}
              <span className="windriver">
                <a href="https://www.windriver.com/" target="_blank">
                  Wind River
                </a>
              </span>
            </li>
          </ul>
        </div>
        <br />
        <div>
          Currently:
          <ul className="list">
            <li>
              <span className="highlight">
                undergraduate research assistant
              </span>{" "}
              @{" "}
              <span className="avril">
                <a
                  href="https://uwaterloo.ca/autonomous-vehicle-research-intelligence-lab/"
                  target="_blank"
                >
                  AVRIL
                </a>
              </span>
            </li>
          </ul>
        </div>
        <br />
        <div>
          Previously:
          <ul className="list">
            <li>
              <span className="highlight">systems software developer</span> @{" "}
              <span className="blackberry">
                <a
                  href="https://www.qnx.com/developers/docs/7.1/#com.qnx.doc.security.system/topic/manual/qcrypto.html"
                  target="_blank"
                >
                  Blackberry
                </a>
              </span>
            </li>
            <li>
              <span className="highlight"> autonomous software developer</span>{" "}
              @{" "}
              <span className="watonomous">
                <a href="https://www.watonomous.ca/" target="_blank">
                  WATonmous
                </a>
              </span>
            </li>
          </ul>
        </div>
        <br />
        <div>
          In my spare time, I like to{" "}
          <span className="instagram">
            <a href="https://www.instagram.com/mango._.climbs/" target="_blank">
              climb rocks
            </a>
          </span>
          ,{" "}
          <span className="go">
            <a href="https://online-go.com/user/view/665071" target="_blank">
              compete in go
            </a>
          </span>{" "}
          and{" "}
          <span className="league">
            <a href="https://www.op.gg/summoners/na/petemango" target="_blank">
              play League
            </a>
          </span>
          !
        </div>
        <br />
      </div>
    </div>
  );
}

export default Body;
