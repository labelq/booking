import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

function Admin() {
    const [bookings, setBookings] = useState([]);
    const [blockedSpots, setBlockedSpots] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState('');
    const navigate = useNavigate();
    const [users, setUsers] = useState([]);

    useEffect(() => {
        const token = localStorage.getItem('authToken');
        if (!token) {
            navigate('/login');
            return;
        }
        fetchBookings();
        fetchBlockedSpots();
        fetchUsers();
    }, [navigate]);

    const fetchBookings = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/admin/bookings', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                }
            });
            if (response.ok) {
                const data = await response.json();
                setBookings(data.bookings || []); // Изменить строку 32
            } else {
                setBookings([]); // Добавить эту строку после строки 32
            }
        } catch (error) {
            console.error('Error fetching bookings:', error);
            setBookings([]); // Добавить эту строку перед строкой 36
        }
    };

    const fetchBlockedSpots = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/admin/blocked-spots', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                }
            });

            if (!response.ok) {
                const errorText = await response.text();
                console.error('Failed to fetch blocked spots:', response.status, errorText);
                setBlockedSpots([]); // Устанавливаем пустой массив при ошибке
                return;
            }

            const data = await response.json();
            setBlockedSpots(data.blockedSpots || []); // Используем пустой массив если data.blockedSpots равен null
        } catch (error) {
            console.error('Error fetching blocked spots:', error);
            setBlockedSpots([]); // Устанавливаем пустой массив при ошибке
        }
    };

    const cancelBooking = async (bookingId) => {
        try {
            setLoading(true);
            const response = await fetch(`http://localhost:8080/api/admin/bookings/${bookingId}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                }
            });

            if (response.ok) {
                setMessage('Бронирование успешно отменено');
                fetchBookings();
            } else {
                const data = await response.json();
                setMessage(data.message || 'Ошибка при отмене бронирования');
            }
        } catch (error) {
            console.error('Error:', error);
            setMessage('Произошла ошибка при отмене бронирования');
        } finally {
            setLoading(false);
        }
    };

    const toggleSpotBlock = async (spotNumber) => {
        try {
            setLoading(true);
            const response = await fetch('http://localhost:8080/api/admin/spots/toggle-block', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                },
                body: JSON.stringify({ spotNumber })
            });

            if (response.ok) {
                fetchBlockedSpots();
                setMessage(`Статус парковочного места ${spotNumber} изменен`);
            } else {
                const data = await response.json();
                setMessage(data.message || 'Ошибка при изменении статуса места');
            }
        } catch (error) {
            console.error('Error:', error);
            setMessage('Произошла ошибка при изменении статуса места');
        } finally {
            setLoading(false);
        }
    };

    const fetchUsers = async () => {
        try {
            const response = await fetch('http://localhost:8080/api/admin/users', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                }
            });
            if (response.ok) {
                const data = await response.json();
                setUsers(data.users || []);
            } else {
                setUsers([]);
            }
        } catch (error) {
            console.error('Error fetching users:', error);
            setUsers([]);
        }
    };

    const toggleUserRole = async (userId, currentRole) => {
        try {
            setLoading(true);
            const newRole = currentRole === 'admin' ? 'user' : 'admin';

            const response = await fetch(`http://localhost:8080/api/admin/users/${userId}/role`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('authToken')}`
                },
                body: JSON.stringify({ role: newRole })
            });

            if (response.ok) {
                setMessage(`Роль пользователя успешно изменена на ${newRole}`);
                fetchUsers();
            } else {
                const data = await response.json();
                setMessage(data.message || 'Ошибка при изменении роли пользователя');
            }
        } catch (error) {
            console.error('Error:', error);
            setMessage('Произошла ошибка при изменении роли пользователя');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen p-8 bg-[#2d2d2d] text-white">
            <h1 className="text-3xl font-montserrat font-semibold mb-8">Панель администратора</h1>

            {/* Сообщение об успехе/ошибке */}
            {message && (
                <div className={`mb-6 p-4 rounded-lg ${message.includes('успешно') ? 'bg-green-600' : 'bg-red-600'}`}>
                    {message}
                </div>
            )}

            {/* Секция управления пользователями */}
            <div className="mb-12">
                <h2 className="text-2xl font-montserrat font-medium mb-4">Управление пользователями</h2>
                <div className="overflow-x-auto">
                    <table className="w-full min-w-[600px] table-auto">
                        <thead>
                        <tr className="bg-[#363636]">
                            <th className="px-4 py-2 text-left">ID</th>
                            <th className="px-4 py-2 text-left">Email</th>
                            <th className="px-4 py-2 text-left">Роль</th>
                            <th className="px-4 py-2 text-left">Действия</th>
                        </tr>
                        </thead>
                        <tbody>
                        {(users || []).map((user) => (
                            <tr key={user.id} className="border-b border-[#363636]">
                                <td className="px-4 py-2">{user.id}</td>
                                <td className="px-4 py-2">{user.email}</td>
                                <td className="px-4 py-2">
                                    <span className={`px-2 py-1 rounded-full text-sm ${
                                        user.account_type === 'admin' ? 'bg-purple-600' : 'bg-blue-600'
                                    }`}>
                                        {user.account_type}
                                    </span>
                                </td>
                                <td className="px-4 py-2">
                                    <button
                                        onClick={() => toggleUserRole(user.id, user.account_type)}
                                        disabled={loading}
                                        className={`px-3 py-1 bg-[#9E7758] rounded-lg hover:bg-[#6E5A42] 
                                            ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
                                    >
                                        {user.account_type === 'admin' ? 'Сделать пользователем' : 'Сделать админом'}
                                    </button>
                                </td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                    {users.length === 0 && (
                        <p className="text-gray-400 mt-4">Нет пользователей</p>
                    )}
                </div>
            </div>

            {/* Секция активных бронирований */}
            <div className="mb-12">
                <h2 className="text-2xl font-montserrat font-medium mb-4">Активные бронирования</h2>
                <div className="grid gap-4">
                    {bookings.map((booking) => (
                        <div key={booking.id} className="bg-[#363636] p-4 rounded-lg">
                            <p>Место: {booking.parking_spot}</p>
                            <p>Номер машины: {booking.car_number}</p>
                            <p>Время начала: {new Date(booking.reserved_at).toLocaleString()}</p>
                            <p>Длительность: {booking.hours} ч.</p>
                            <button
                                onClick={() => cancelBooking(booking.id)}
                                disabled={loading}
                                className={`mt-2 px-4 py-2 bg-red-600 rounded-lg hover:bg-red-700 
                                    ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
                            >
                                Отменить бронирование
                            </button>
                        </div>
                    ))}
                    {bookings.length === 0 && (
                        <p className="text-gray-400">Нет активных бронирований</p>
                    )}
                </div>
            </div>

            {/* Секция управления парковочными местами */}
            <div>
                <h2 className="text-2xl font-montserrat font-medium mb-4">Управление парковочными местами</h2>
                <div className="grid grid-cols-4 gap-4">
                    {[...Array(16)].map((_, index) => {
                        const spotNumber = index + 1;
                        const isBlocked = (blockedSpots || []).includes(spotNumber);
                        return (
                            <button
                                key={spotNumber}
                                onClick={() => toggleSpotBlock(spotNumber)}
                                disabled={loading}
                                className={`p-4 rounded-lg ${
                                    isBlocked
                                        ? 'bg-red-600 hover:bg-red-700'
                                        : 'bg-[#9E7758] hover:bg-[#6E5A42]'
                                } ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
                            >
                                Место {spotNumber}
                                <br />
                                {isBlocked ? 'Заблокировано' : 'Активно'}
                            </button>
                        );
                    })}
                </div>
            </div>
        </div>
    );
}

export default Admin;