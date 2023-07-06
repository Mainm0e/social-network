const GroupHeader = ({ groupid, type, state, refreshComponent }) => {
  return (
    <div className="group-header">
      <div className="group-header-title">
        <h1>Group</h1>
      </div>
    </div>
  );
};

export default GroupHeader;

//group dummy data
// const group = {
//     id: 1,
//     name: "group1",
//     description: "this is group1",
//     createrId: 1,
//     createrName: "user1",
//     followers: 10,
