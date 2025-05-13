import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';

function Home() {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        setTimeout(() => {
            setIsVisible(true);
        }, 100); // Стартуем анимацию через 100мс
    }, []);

    return (
        <div className="min-h-screen flex flex-col justify-center items-center p-0 m-0">
            <div className="text-center mb-12">
                <h1 className="text-5xl font-montserrat font-medium text-white mb-4 hover:text-gray-500 transition duration-300">
                    Добро пожаловать в КИУ-Parking!
                </h1>
            </div>

            {/* 2 ряда по 2 одинаковых блоков */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 w-full px-4 sm:px-6 lg:px-8 mb-12">
                <div
                    className={`${
                        isVisible ? 'transform translate-x-0 opacity-100' : 'transform translate-x-20 opacity-0'
                    } w-full bg-[#363636] text-white p-8 rounded-lg shadow-lg transition-all duration-500`}
                >
                    <p className="font-montserrat font-medium text-3xl mb-4 hover:text-gray-500 transition duration-300">
                        Сервис аренды машиномест.
                    </p>
                    <p className="font-montserrat font-medium text-gray-300 transition duration-300">
                        Для сотрудников и студентов КИУ
                    </p>
                </div>

                <div
                    className={`${
                        isVisible ? 'transform translate-x-0 opacity-100' : 'transform translate-x-20 opacity-0'
                    } w-full bg-[#795238] text-white p-8 rounded-lg shadow-lg transition-all duration-500`}
                >
                    <p className="font-montserrat font-medium text-3xl mb-4 hover:text-gray-500 transition duration-300">
                        Удобная парковка.
                    </p>
                    <p className="font-montserrat font-medium text-gray-300 transition duration-300">
                        Парковка на территории универитета
                    </p>
                </div>

                <div
                    className={`${
                        isVisible ? 'transform translate-x-0 opacity-100' : 'transform translate-x-20 opacity-0'
                    } w-full bg-[#795238] text-white p-8 rounded-lg shadow-lg transition-all duration-500`}
                >
                    <p className="font-montserrat font-medium text-3xl mb-4 hover:text-gray-500 transition duration-300">
                        Охраняемая территория.
                    </p>
                    <p className="font-montserrat font-medium text-gray-300 transition duration-300">
                        Въезд на парковку через шлагбаум.
                    </p>
                    <p className="font-montserrat font-medium text-gray-300 transition duration-300">
                        Работает служба охраны.
                    </p>
                </div>

                <div
                    className={`${
                        isVisible ? 'transform translate-x-0 opacity-100' : 'transform translate-x-20 opacity-0'
                    } w-full bg-[#363636] text-white p-8 rounded-lg shadow-lg transition-all duration-500`}
                >
                    <p className="font-montserrat font-medium text-3xl mb-8 hover:text-gray-500 transition duration-300">
                        Быстрая аренда машиноместа.
                    </p>
                    <p className="font-montserrat font-medium text-gray-300 transition duration-300">
                        Удобный интерфейс для бронирования.
                    </p>
                    <p className="font-montserrat font-medium text-gray-300 transition duration-300">
                        Оплата картой, по СБП и QR.
                    </p>
                </div>
            </div>

            <div className="flex space-x-4 mb-8">
                <Link to="/map">
                    <button className="font-montserrat font-medium py-3 px-6 bg-[#9E7758] text-white rounded-lg hover:bg-[#11120e] transition-transform duration-300 hover:scale-105 focus:outline-none">
                        Как добраться?
                    </button>
                </Link>
                <Link to="/booking">
                    <button className="font-montserrat font-medium py-3 px-6 bg-[#9E7758] text-white rounded-lg hover:bg-[#11120e] transition-transform duration-300 hover:scale-105 focus:outline-none">
                        Забронировать место
                    </button>
                </Link>
                <Link to="/login">
                    <button className="font-montserrat font-medium py-3 px-6 bg-[#9E7758] text-white rounded-lg hover:bg-[#11120e] transition-transform duration-300 hover:scale-105 focus:outline-none">
                        Войти / Зарегистрироваться
                    </button>
                </Link>
            </div>
        </div>
    );
}

export default Home;