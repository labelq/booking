import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Home from './pages/Home.jsx';
import Map from './pages/Map';
import Booking from './pages/Booking';
import Login from './pages/Login';
import Register from './pages/Register';
import Header from "./pages/Header.jsx";
import Footer from "./pages/Footer.jsx";
import Admin from './pages/Admin';

function App() {
    return (
        <Router>
            <Header />
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/map" element={<Map />} />
                <Route path="/booking" element={<Booking />} />
                <Route path="/login" element={<Login    />} />
                <Route path="/register" element={<Register />} />
                <Route path="/admin" element={<Admin />} />
            </Routes>
            <Footer />
        </Router>
    );
}

export default App;