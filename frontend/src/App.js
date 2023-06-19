import LoginPage from "./LoginPage/LoginPage";
import RegisterPage from "./RegisterPage/RegisterPage";
import MainPage from "./MainPage/MainPage";
import ErrorPage from "./ErrorPage/ErrorPage";
import './App.css';
import {getCookie} from "./tools/cookie";

function App() {
  // function handle that check url and return the page
  // check cookies and return the page
  const getPage = () => {
    const page = window.location.pathname;
    const sessionId = getCookie("sessionId");
    if (sessionId === null){
      if (page === "/"||page ==="/login"){
        return <LoginPage />
      } else if (page === "/register"){
        return <RegisterPage />
      }
    }
    if (sessionId !== null && page === "/"||sessionId !== null && page === "/register"||sessionId !== null && page === "/login"){
      return <MainPage />
    } else {
      return <ErrorPage />
    }

  };


  return (
    <div className="App">
      {getPage()}
    </div>
  );
};

export default App;
