import { Link } from 'react-router-dom';

function Footer() {
    return (
        <footer className="bg-[#11120e] fixed bottom-0 left-0 w-full z-10 shadow-xl backdrop-blur-xl">
            <div className="mx-auto max-w-screen-xl px-4 sm:px-6 lg:px-8 py-4">
                <div className="flex items-center justify-between">
                    <div className="flex-1 text-white text-sm">
                        <p>&copy; 2025 КИУ-Parking. Все права защищены.</p>
                    </div>
                </div>
            </div>
        </footer>
    );
}

export default Footer;