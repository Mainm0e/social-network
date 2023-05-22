import LoginPage from "./LoginPage/LoginPage";
import RegisterPage from "./RegisterPage/RegisterPage";
import './App.css';
function App() {
  // function handle that check url and return the page
  const getPage = () => {
    const page = window.location.pathname;
    if (page === '/login') {
      return <LoginPage />;
    } else if (page === '/register') {
      return <RegisterPage />;
    } else {
      return <LoginPage />;
    }
  };


  return (
    <div className="App">
      {getPage()}
    </div>
  );
};

export default App;
