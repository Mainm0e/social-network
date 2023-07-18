import "./AlertBox.css";


// AlertBox component
// Props: title, message, status
function AlertBox(props) {
  const { title, message, status } = props;
  const boxclose = () => {
   document.querySelector(".alert-box").style.display = "none";
  }
  if (title === undefined || message === undefined || status === undefined) {
    return null;
  } else {
  return (
        <div 
        className="alert-box" 
        style={{display: status === false ? "block":"none" }}>
          <span className="closebtn" onClick={boxclose} >
            &times;
          </span>
          <strong>{title}</strong>
          <br/>
          <br/>
          {/* message is array */}
          {message.map((msg, index) => (
            <span key={index}>{msg}<br/></span>
          ))}
         
        </div>
  );
          }
}

export default AlertBox;
