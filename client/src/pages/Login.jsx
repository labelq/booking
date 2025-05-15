import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Link } from "react-router-dom";

function Login() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [message, setMessage] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch("http://localhost:8080/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            if (response.ok) {
                const data = await response.json();
                localStorage.setItem("authToken", data.token);
                // Также сохраним информацию о пользователе
                localStorage.setItem("user", JSON.stringify(data.user));

                // Явно создаем новое событие
                const authEvent = new Event('auth-change');
                window.dispatchEvent(authEvent);

                console.log('Login successful, token saved');
                console.log('Auth change event dispatched');

                navigate("/booking");
            } else {
                setMessage("Неверный логин или пароль");
            }
        } catch (error) {
            console.error("Login error:", error);
            setMessage("Произошла ошибка при входе");
        }
    };

    return (
        <div className="min-h-screen flex flex-col justify-center items-center text-white p-4">
            <h1 className="font-montserrat font-semibold text-3xl mb-8">Авторизация</h1>

            <form onSubmit={handleSubmit} className="w-full max-w-sm">
                {/* Почта */}
                <div className="relative mb-6">
                    <input
                        type="email"
                        value={email}
                        id="email"
                        placeholder=" "
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758]  dark:border-gray-600 dark:text-white dark:focus:ring-[#646560] dark:focus:border-[#646560]"
                        onChange={(e) => setEmail(e.target.value)}
                    />
                    <label
                        htmlFor="email"
                        className="absolute left-4 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-500 peer-focus:text-sm peer-focus:text-[#9E7758] peer-focus:transform peer-focus:-translate-y-5 transition-all duration-200"
                    >
                        Почта
                    </label>
                </div>

                {/* Пароль */}
                <div className="relative mb-6">
                    <input
                        type="password"
                        value={password}
                        id="password"
                        placeholder=" "
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758]  dark:border-gray-600 dark:text-white dark:focus:ring-[#646560] dark:focus:border-[#646560]"
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <label
                        htmlFor="password"
                        className="absolute left-4 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-500 peer-focus:text-sm peer-focus:text-[#9E7758] peer-focus:transform peer-focus:-translate-y-5 transition-all duration-200"
                    >
                        Пароль
                    </label>
                </div>

                {/* Кнопка входа */}
                <button
                    type="submit"
                    className="w-full py-3 px-6 bg-[#9E7758] text-white font-semibold rounded-lg hover:bg-[#6E5A42] transition-all duration-300"
                >
                    Войти
                </button>
            </form>

            {/* Сообщение об ошибке */}
            {message && <p className="text-red-500 mt-4">{message}</p>}

            {/* Ссылка на регистрацию */}
            <div className="mt-4 text-center">
                <Link
                    to="/register"
                    className="text-[#9E7758] hover:text-[#6E5A42] font-medium"
                >
                    Еще не зарегистрированы? Зарегистрироваться
                </Link>
            </div>
        </div>
    );
}

export default Login;