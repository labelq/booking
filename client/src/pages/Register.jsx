import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Link } from "react-router-dom";

function Register() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [message, setMessage] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();

        // Проверка на пустые поля
        if (email === "" || password === "" || confirmPassword === "") {
            setMessage("Заполните все поля");
            return;
        }

        // Проверка на совпадение паролей
        if (password !== confirmPassword) {
            setMessage("Пароли не совпадают");
            return;
        }

        try {
            // Отправляем данные на сервер Go для регистрации
            const response = await fetch("http://localhost:8080/api/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ email, password }),
            });

            if (response.ok) {
                // Если регистрация успешна, получаем токен и сохраняем его
                const data = await response.json();
                localStorage.setItem("authToken", data.token); // Сохраняем токен

                // Публикуем событие для обновления состояния авторизации
                window.dispatchEvent(new Event('auth-change'));

                // Перенаправляем на главную страницу вместо страницы бронирования
                navigate("/");
            } else {
                const errorData = await response.json();
                setMessage(errorData.message || "Ошибка при регистрации");
            }
        } catch (error) {
            setMessage("Ошибка сервера, попробуйте позже");
            console.error("Registration error:", error);
        }
    };

    return (
        <div className="min-h-screen flex flex-col justify-center items-center text-white p-4">
            <h1 className="font-montserrat font-semibold text-3xl mb-8">Регистрация</h1>

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

                {/* Подтверждение пароля */}
                <div className="relative mb-6">
                    <input
                        type="password"
                        value={confirmPassword}
                        id="confirmPassword"
                        placeholder=" "
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758]  dark:border-gray-600 dark:text-white dark:focus:ring-[#646560] dark:focus:border-[#646560]"
                        onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                    <label
                        htmlFor="confirmPassword"
                        className="absolute left-4 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-500 peer-focus:text-sm peer-focus:text-[#9E7758] peer-focus:transform peer-focus:-translate-y-5 transition-all duration-200"
                    >
                        Подтвердите пароль
                    </label>
                </div>

                {/* Сообщение об ошибке */}
                {message && <p className="text-red-500 mt-4">{message}</p>}

                {/* Кнопка регистрации */}
                <button
                    type="submit"
                    className="w-full py-3 px-6 bg-[#9E7758] text-white font-semibold rounded-lg hover:bg-[#6E5A42] transition-all duration-300"
                >
                    Зарегистрироваться
                </button>
            </form>

            {/* Ссылка на авторизацию */}
            <div className="mt-4 text-center">
                <Link
                    to="/login"
                    className="text-[#9E7758] hover:text-[#6E5A42] font-medium"
                >
                    Уже есть аккаунт? Войти
                </Link>
            </div>
        </div>
    );
}

export default Register;