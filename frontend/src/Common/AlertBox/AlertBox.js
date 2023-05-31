import "./AlertBox.css";


// AlertBox component
// Props: title, message, status
function AlertBox(props) {
  const { title, message, status } = props;

  const boxclose = () => {
   document.querySelector(".alert-box").style.display = "none";
  }
  return (
        <div 
        className="alert-box" 
        style={{display: status === false ? "block":"none" }}>
          <span className="closebtn" onClick={boxclose} >
            &times;
          </span>
          <strong>{title}</strong>
          <br />
          <p>{message}</p>
        </div>
  );
}

export default AlertBox;
