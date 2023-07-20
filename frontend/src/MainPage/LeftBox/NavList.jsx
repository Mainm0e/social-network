import React,{useState} from "react";
import "./LeftBox.css";
const NavList = ({ type,links }) => {
    const [showList, setShowList] = useState(false);
    const toggleList = () => {
        setShowList(!showList);
        };
    const List = () => {
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
           <List />
        </div>
    );
}
export default NavList;
