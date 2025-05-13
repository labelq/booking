import '../App.css'
import { Link } from "react-router-dom";

function Header() {
    return (
        <div>
            <header className="bg-[#11120e] fixed top-0 left-0 w-full z-10 shadow-xl backdrop-blur-xl">
                <div className="mx-auto max-w-screen-xl px-4 sm:px-6 lg:px-8">
                    <div className="flex h-16 items-center justify-between">
                        <div className="flex-1 md:flex md:items-center md:gap-12">
                            <a>
                                <span className="sr-only">Home</span>
                                <img
                                    className="h-18"
                                    src="/kiu-logo-menu.png"
                                    alt="Home Logo"
                                />
                            </a>
                            <nav aria-label="Global" className="hidden md:block">
                                <ul className="flex items-center gap-6 text-sm">
                                    <li>
                                        <Link
                                            className="text-gray-500 transition hover:text-gray-500/75 dark:text-white dark:hover:text-white/75"
                                            to="/"
                                        >
                                            КИУ-Parking
                                        </Link>
                                    </li>
                                </ul>
                            </nav>
                        </div>

                        <div className="md:flex md:items-center md:gap-12">
                            <nav aria-label="Global" className="hidden md:block">
                                <ul className="flex items-center gap-6 text-sm">
                                    <li>
                                        <Link
                                            className="text-gray-500 transition hover:text-gray-500/75 dark:text-white dark:hover:text-white/75"
                                            to="/map"
                                        >
                                            Как добраться?
                                        </Link>
                                    </li>
                                    <li>
                                        <Link
                                            className="text-gray-500 transition hover:text-gray-500/75 dark:text-white dark:hover:text-white/75"
                                            to="/booking"
                                        >
                                            Бронировать
                                        </Link>
                                    </li>
                                </ul>
                            </nav>

                            <div className="flex items-center gap-4">
                                <div className="sm:flex sm:gap-4">
                                    <Link
                                        className="py-2.5 px-5 bg-[#646560] text-white rounded-lg hover:bg-[#3e3f3a] transition-transform duration-300 hover:scale-105 focus:outline-none"
                                        to="/login"
                                    >
                                        Вход
                                    </Link>

                                    <div className="hidden sm:flex">
                                        <Link
                                            className="py-2.5 px-5 bg-[#646560] text-white rounded-lg hover:bg-[#3e3f3a] transition-transform duration-300 hover:scale-105 focus:outline-none"
                                            to="/register"
                                        >
                                            Регистрация
                                        </Link>
                                    </div>
                                </div>

                                <div className="block md:hidden">
                                    <button
                                        className="rounded-sm bg-gray-100 p-2 text-gray-600 transition hover:text-gray-600/75 dark:bg-gray-800 dark:text-white dark:hover:text-white/75"
                                    >
                                        <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            className="size-5"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            stroke="currentColor"
                                            stroke-width="2"
                                        >
                                            <path stroke-linecap="round" stroke-linejoin="round" d="M4 6h16M4 12h16M4 18h16"/>
                                        </svg>
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </header>
        </div>
    )
}

export default Header;