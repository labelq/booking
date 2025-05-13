import React from 'react';
import { Link } from "react-router-dom";

function Map() {
    return (
        <div className="min-h-screen flex flex-col pt-16"> {/* Добавили pt-16 для отступа от шапки */}
            <h1 className="font-montserrat font-medium mb-2">
                Как добраться до парковки университета?
            </h1>
            <p className="font-montserrat font-medium text-xl ">
                Подъезжая к университету сверните на улицу Бурхана Шахиди.
            </p>
            <p className="font-montserrat font-medium text-xl mb-6">
                Заезд находится с левой стороны, сразу после Созведия талантов.
            </p>

            <div style={{ position: 'relative', overflow: 'hidden' }}>
                <iframe
                    src="https://yandex.ru/map-widget/v1/?l=trf%2Ctrfe&ll=49.107658%2C55.787677&z=18"
                    width="100%"  // Ширина карты на всю ширину контейнера
                    height="500"  // Высота карты
                    frameBorder="1"
                    allowFullScreen="true"
                    title="Yandex Map"
                    style={{ position: 'relative' }}
                ></iframe>
            </div>

            <div className="flex justify-center mt-6">
                <Link to="/booking">
                    <button className="font-montserrat font-medium py-3 px-6 bg-[#9E7758] text-white rounded-lg hover:bg-[#11120e] transition-transform duration-300 hover:scale-105 focus:outline-none">
                        Забронировать место
                    </button>
                </Link>
            </div>

        </div>
    );
}

export default Map;