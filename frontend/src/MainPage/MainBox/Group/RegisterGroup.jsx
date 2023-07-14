import React from 'react';
import { fetchData } from '../../../tools/fetchData';
import './Group.css';
import { getCookie, getUserId } from '../../../tools/cookie';


const RegisterGroup = (user) => {

    const sendRequest = async (e) => {
        e.preventDefault();
        const method = "POST";
        const type = "createGroup";
        /* name: e.target.name.value, description: e.target.description.value  */
        const payload = { 
            sessionId: getCookie("sessionId"),
            creatorProfile: {
                            userId: getUserId("userId"),
                            firstName: user.firstName,
                            lastName:user.lastName,
                            avatar:user.avatar
                            },
            groupId: null,
            title: e.target.name.value,
            description: e.target.description.value,
            date:null
        };
       fetchData(method, type, payload).then((data) => {
                window.location.href = "/group";
       });
    }

    /* registerGroup form */

    return (
        <div className="registerGroup">
            <h1>Register Group</h1>
            <form onSubmit={sendRequest} className='register-box'>
                <label htmlFor="name">Name</label>
                <input type="text" name="name" id="name" />
                <label htmlFor="description">Description</label>
                <textarea name="description" id="description" />
                <input type="submit" value="Register" />
            </form>
        </div>
    );
}

export default RegisterGroup;