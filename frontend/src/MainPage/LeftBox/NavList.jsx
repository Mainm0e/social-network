import React,{useState} from "react";
import "./LeftBox.css";
const NavList = ({ type,links }) => {
    const [showList, setShowList] = useState(false);
    const toggleList = () => {
        setShowList(!showList);
        };
    const list = () => {
        return (
            <>
            {Object.entries(links).map(([key, value]) => (
                <li key={key} className="link-list">
                <a href={value}>{key}</a>
              </li>
            ))}
            </>
        );
      };


    return (
        <div className="nav-list">
            <span onClick={toggleList}>{type}</span>
            {showList && <>{list()}</>}
        </div>
    );
}
export default NavList;
