import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useNavigate } from 'react-router-dom';

function Booking() {
    const [carNumber, setCarNumber] = useState("");
    const [hours, setHours] = useState(1);
    const [totalPrice, setTotalPrice] = useState(100);
    const [parkingSpot, setParkingSpot] = useState("");
    const [message, setMessage] = useState("");
    const [occupiedSpots, setOccupiedSpots] = useState([]);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        const token = localStorage.getItem('authToken');
        if (!token) {
            navigate('/login');
            return;
        }
        fetchOccupiedSpots();
    }, [navigate]);

    const fetchOccupiedSpots = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/bookings', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                }
            });
            if (response.ok) {
                const data = await response.json();
                setOccupiedSpots(data.occupiedSpots || []);
            }
        } catch (error) {
            console.error('Error fetching occupied spots:', error);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        setMessage("");

        if (!parkingSpot || !carNumber) {
            setMessage("Пожалуйста, заполните все поля");
            setLoading(false);
            return;
        }

        try {
            const response = await fetch('http://localhost:8080/api/booking', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                },
                body: JSON.stringify({
                    parkingSpot: parseInt(parkingSpot),
                    carNumber,
                    hours: parseInt(hours)
                })
            });

            const data = await response.json();

            if (response.ok) {
                const endTime = new Date(data.endTime);
                setMessage(`
                    Бронирование успешно!
                    Место: ${parkingSpot}
                    Время окончания: ${endTime.toLocaleString()}
                    Номер машины: ${carNumber}
                `);
                await fetchOccupiedSpots();
                setCarNumber("");
                setParkingSpot("");
                setHours(1);
                setTotalPrice(100);
            } else {
                setMessage(data.message || "Ошибка при бронировании");
            }
        } catch (error) {
            console.error('Error:', error);
            setMessage("Произошла ошибка при бронировании");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex flex-col justify-center items-center text-white p-4">
            <h1 className="font-montserrat font-semibold text-3xl mb-8">Бронирование парковки</h1>

            {/* Фото парковки */}
            <div className="mb-8">
                <img
                    src="/0DBFF439-B104-40AD-B43F-22FD452EE2AB.JPEG"
                    alt="Парковка"
                    className="w-full max-w-lg rounded-lg shadow-lg"
                />
            </div>

            {/* Сообщение об успехе/ошибке */}
            {message && (
                <div className={`mb-4 p-4 rounded-lg ${message.includes('успешно') ? 'bg-green-600' : 'bg-red-600'}`}>
                    <pre className="whitespace-pre-line">{message}</pre>
                </div>
            )}

            {/* Форма бронирования */}
            <form onSubmit={handleSubmit} className="w-full max-w-sm">
                {/* Номер машины */}
                <div className="relative mb-6">
                    <input
                        type="text"
                        value={carNumber}
                        id="carNumber"
                        placeholder=" "
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758]"
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
                        className="peer w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758]"
                        onChange={(e) => {
                            const value = parseInt(e.target.value);
                            setHours(value);
                            setTotalPrice(value * 100);
                        }}
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
                        className="w-full px-4 py-3 bg-[#3e3f3a] text-white rounded-lg border-2 border-[#9E7758] focus:outline-none focus:ring-2 focus:ring-[#9E7758] focus:border-[#9E7758]"
                    >
                        <option value="">Выберите парковочное место</option>
                        {[...Array(16)].map((_, index) => {
                            const spotNumber = index + 1;
                            const isOccupied = occupiedSpots.includes(spotNumber);
                            return (
                                <option
                                    key={spotNumber}
                                    value={spotNumber}
                                    disabled={isOccupied}
                                >
                                    Место №{spotNumber} {isOccupied ? '(занято)' : ''}
                                </option>
                            );
                        })}
                    </select>
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
                    disabled={loading}
                    className={`w-full py-3 px-6 bg-[#9E7758] text-white font-semibold rounded-lg 
                        ${loading ? 'opacity-50 cursor-not-allowed' : 'hover:bg-[#6E5A42]'} 
                        transition-all duration-300`}
                >
                    {loading ? 'Бронирование...' : 'Подтвердить бронирование'}
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