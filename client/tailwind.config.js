module.exports = {
    content: [
        "./src/**/*.{html,js,jsx,ts,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                primary: "#1e3a8a",
                secondary: "#9333ea",
            },
            fontFamily: {
                sans: ['system-ui', 'Avenir', 'Helvetica', 'Arial', 'sans-serif'],
                'montserrat': ['Montserrat', 'sans-serif'],
            },
        },
    },
    plugins: [],
}