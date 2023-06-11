import React, { useState } from "react";
import defaultAvatar from "../avatar/default-avatar.png";
import "./info3.css";

const Info3 = ({selectedOption, onChange, registerStatus}) => {
    const type = "info3"

    // Upload avatar
  const [selectedFile, setSelectedFile] = useState(null);
  const [previewSource, setPreviewSource] = useState(defaultAvatar);

  const handleFileChange = (event) => {
    console.log(event.target.files);
    const file = event.target.files[0];
    setSelectedFile(file);
    previewFile(file);
    onChange(event); // Call the onChange prop
  };

  const previewFile = (file) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onloadend = () => {
      setPreviewSource(reader.result);
      // Do something with the file data, such as storing it in state or sending it to a server
      const fileData = reader.result;
      setAvatar(fileData);
    };
  };
  
  const [nickname, setNickname] = useState("");
  const [aboutme, setAboutme] = useState("");
  const [avatar, setAvatar] = useState("");
  // collect data from form
  const handleInputChanges = (event) => {
    const { name, value } = event.target;
    switch (name) {
        case "nickname":
            setNickname(value);
            break;
        case "aboutme":
            setAboutme(value);
            break;
        default:
            break;
    }         
    };

    React.useEffect(() => {
        onChange({
            type,
            nickname,
            aboutme,
            avatar,
        });
    }, [type, nickname, aboutme, avatar, onChange]);



  return (
    <form>
        <div className={`user-info3_${selectedOption !== type ? "hide" :"false"}`}>
            <div className='avatar-input-container' >
                <label>Avatar:</label>
                <div className='avatar-container'>
                {previewSource && (
                    <img src={previewSource} alt='Preview'/>
                )}
                <input type='file' onChange={handleFileChange} />
                </div>
            </div>
            <div className='input-container'>
                <label>Nickname:</label>
                <input type='text' name='nickname' onChange={handleInputChanges} />
            </div>
            <div className='input-container'>
                <label>About Me:</label>
                <input type='text' name='aboutme' onChange={handleInputChanges} />
            </div>
        </div>
    </form>
  )

}

export default Info3