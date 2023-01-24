/** @type {import('tailwindcss').Config} */
module.exports = {
	darkMode: "class",
	content: [
		"./src/**/*.{html,jsx,js,svelte}",
	],
	theme: {
		extend: {
			fontFamily: {
				abel: ["Abel", "sans-serif"]
			},
		}
	},
	plugins: [
		require("daisyui"),
	],
	daisyui: {
		themes: ["corporate"]
	}
};
