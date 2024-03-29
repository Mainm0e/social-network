import { useEffect } from "react";
import "./Group.css";
import {memberlist} from "../../../tools/link"
const GroupHeader = ({ group, handleRefresh }) => {
  useEffect(() => {
  }, [group]);

  return (
    <div className="group-header">
      <div className="group-header-name">
        <label htmlFor="group-name">Group:</label>
        <span id="group-name">{group.title}</span>
      </div>
      <div className="group-header-description">
        <label htmlFor="group-description">Description:</label>
        <span id="group-description">{group.description}</span>
      </div>
      <div className="group-header-followers">
        <label htmlFor="group-followers" onClick={()=> memberlist()}>Member:</label>
        <span id="group-followers">{group.noMembers}</span>
      </div>
      <div className="group-header-creater">
        <label htmlFor="group-creater">Creater:</label>
        <div className="group-header-creater-info">
          <div className="group-header-creater-avatar">
            <img src={group.creatorProfile.avatar} alt="creater-avatar" />
          </div>
          <div className="group-header-creater-name">
            <span id="group-creater-firstName">
              {group.creatorProfile.firstName}
            </span>
            <span> </span>
            <span id="group-creater-lastName">
              {group.creatorProfile.lastName}
            </span>
          </div>
        </div>
      </div>
      <div className="group-header-date">
        <label htmlFor="group-date">Date:</label>
        <span id="group-date">{group.date}</span>
      </div>
    </div>
  );
};

export default GroupHeader;