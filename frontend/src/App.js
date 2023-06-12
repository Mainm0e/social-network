import LoginPage from "./LoginPage/LoginPage";
import RegisterPage from "./RegisterPage/RegisterPage";
import MainPage from "./MainPage/MainPage";
import ErrorPage from "./ErrorPage/ErrorPage";
import MyComponent from "./test_componen";
import './App.css';
function App() {
  // function handle that check url and return the page
  const getPage = () => {
    const page = window.location.pathname;
    if (page === '/login') {
      return <LoginPage />;
    } else if (page === '/register') {
      return <RegisterPage />;
    } else if (page === '/') {
      return <MainPage />;
    } else {
      return <ErrorPage /> ;
    }
  };


  return (
    <div className="App">
      <MyComponent/>
    </div>
  );
};

export default App;
