

const Notification = ({data}) => {
    console.log("Notification",data)
    data = data ? data : dummyData
    
    return (
        <div className="notification-container">
            {data.map((item,index) => {
                return <DisplayNotification key={index} message={item.message} type={item.type}/>
            }
            )}
        </div>
    )
}

export default Notification

const DisplayNotification = ({message,type}) => {
    return (
        <div className={`notification ${type}`}>
            <p>{message}</p>
        </div>
    )
}

const dummyData = [
    {
        message: "This is a success message",
        type: "success"
    },
    {
        message: "This is a error message",
        type: "error"
    },
    {
        message: "This is a warning message",
        type: "warning"
    },
    {
        message: "This is a info message",
        type: "info"
    }
]