import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useNavigate } from 'react-router-dom';



function Booking() {
    const [carNumber, setCarNumber] = useState(""); // Номер машины
    const [hours, setHours] = useState(1); // Количество часов
    const [totalPrice, setTotalPrice] = useState(100); // Стоимость
    const [parkingSpot, setParkingSpot] = useState(null); // Номер выбранного места
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('authToken');
        if (!token) {
            navigate('/login'); // Перенаправляем на страницу входа, если нет токена
        }
    }, [navigate]);

    // Обработчик изменения количества часов
    const handleHoursChange = (e) => {
        const hours = e.target.value;
        setHours(hours);
        setTotalPrice(hours * 100); // 100р за каждый час
    };

    // Обработчик отправки формы
    const handleSubmit = (e) => {
        e.preventDefault();
        if (!parkingSpot) {
            alert("Пожалуйста, выберите место для парковки.");
            return;
        }
        // Логика для обработки бронирования
        alert(`Вы успешно забронировали место №${parkingSpot} на ${hours} час(ов) за ${totalPrice} рублей.`);
    };

    return (
        <div className="min-h-screen flex flex-col justify-center items-center text-white p-4">
            <h1 className="font-montserrat font-semibold text-3xl mb-8">Бронирование парковки</h1>

            {/* Фото парковки */}
            <div className="mb-8">
                <img
                    src="/0DBFF439-B104-40AD-B43F-22FD452EE2AB.JPEG" // Путь к изображению парковки
                    alt="Парковка"
                    className="w-full max-w-lg rounded-lg shadow-lg"
                />
            </div>

            {/* Форма бронирования */}
            <form onSubmit={handleSubmit} className="w-full max-w-sm">
                {/* Номер машины */}
                <div className="relative mb-6">
                    <input
                        type="text"
                        value={carNumber}
                        id="carNumber"
                        placeholder=" "
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758] dark:border-gray-600 dark:text-white dark:focus:ring-[#646560] dark:focus:border-[#646560]"
                        onChange={(e) => setCarNumber(e.target.value)}
                    />
                    <label
                        htmlFor="carNumber"
                        className="absolute left-4 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-500 peer-focus:text-sm peer-focus:text-[#9E7758] peer-focus:transform peer-focus:-translate-y-5 transition-all duration-200"
                    >
                        Номер машины
                    </label>
                </div>

                {/* Количество часов */}
                <div className="relative mb-6">
                    <input
                        type="number"
                        value={hours}
                        id="hours"
                        min="1"
                        placeholder=" "
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758] dark:border-gray-600 dark:text-white dark:focus:ring-[#646560] dark:focus:border-[#646560]"
                        onChange={handleHoursChange}
                    />
                    <label
                        htmlFor="hours"
                        className="absolute left-4 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-500 peer-focus:text-sm peer-focus:text-[#9E7758] peer-focus:transform peer-focus:-translate-y-5 transition-all duration-200"
                    >
                        Количество часов
                    </label>
                </div>

                {/* Выбор парковочного места */}
                <div className="relative mb-6">
                    <select
                        value={parkingSpot}
                        onChange={(e) => setParkingSpot(e.target.value)}
                        className="w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758] dark:border-gray-600 dark:text-white dark:focus:ring-[#646560] dark:focus:border-[#646560]"
                    >
                        <option value="" disabled>Выберите парковочное место</option>
                        {[...Array(16)].map((_, index) => (
                            <option key={index} value={index + 1}>
                                Место №{index + 1}
                            </option>
                        ))}
                    </select>
                    <label
                        className="absolute left-4 top-1/2 transform -translate-y-1/2 text-sm text-gray-400 peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-500 peer-focus:text-sm peer-focus:text-[#9E7758] peer-focus:transform peer-focus:-translate-y-5 transition-all duration-200"
                    >
                        Выберите место
                    </label>
                </div>

                {/* Итоговая стоимость */}
                <div className="mb-6">
                    <p className="text-lg font-montserrat font-medium text-white">
                        Итоговая стоимость: <span className="font-semibold">{totalPrice} р.</span>
                    </p>
                </div>

                {/* Кнопка бронирования */}
                <button
                    type="submit"
                    className="w-full py-3 px-6 bg-[#9E7758] text-white font-semibold rounded-lg hover:bg-[#6E5A42] transition-all duration-300"
                >
                    Подтвердить бронирование
                </button>
            </form>

            {/* Ссылка на главную страницу */}
            <div className="mt-4 text-center">
                <Link
                    to="/"
                    className="text-[#9E7758] hover:text-[#6E5A42] font-medium"
                >
                    Вернуться на главную
                </Link>
            </div>
        </div>
    );
}

export default Booking;