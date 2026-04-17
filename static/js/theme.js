(function () {
    const THEME_STORAGE_KEY = "athena-theme";
    const DARK_THEME = "dark";
    const LIGHT_THEME = "light";

    function getSavedTheme() {
        try {
            return localStorage.getItem(THEME_STORAGE_KEY) === DARK_THEME ? DARK_THEME : LIGHT_THEME;
        } catch (_) {
            return LIGHT_THEME;
        }
    }

    function saveTheme(theme) {
        try {
            localStorage.setItem(THEME_STORAGE_KEY, theme);
        } catch (_) {
            // localStorage may be unavailable in some environments
        }
    }

    function applyTheme(theme) {
        const isDarkMode = theme === DARK_THEME;
        document.body.classList.toggle("dark-mode", isDarkMode);

        document.querySelectorAll(".theme-toggle").forEach((toggle) => {
            toggle.checked = isDarkMode;
        });
    }

    document.addEventListener("DOMContentLoaded", () => {
        applyTheme(getSavedTheme());

        document.querySelectorAll(".theme-toggle").forEach((toggle) => {
            toggle.addEventListener("change", (event) => {
                const selectedTheme = event.target.checked ? DARK_THEME : LIGHT_THEME;
                applyTheme(selectedTheme);
                saveTheme(selectedTheme);
            });
        });
    });
})();
