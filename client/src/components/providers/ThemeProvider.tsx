import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';

type ThemeType = 'light' | 'dark' | 'system';

type ThemeContextType = {
    theme: ThemeType;
    resolvedTheme: 'light' | 'dark';
    setTheme: (theme: ThemeType) => void;
};

const STORAGE_KEY = 'shrtn-theme';

const ThemeContext = createContext<ThemeContextType | null>(null);

function getSystemTheme(): 'light' | 'dark' {
    if (typeof window === 'undefined') return 'light';
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

function getStoredTheme(): ThemeType {
    if (typeof window === 'undefined') return 'system';
    const stored = localStorage.getItem(STORAGE_KEY);
    if (stored === 'light' || stored === 'dark' || stored === 'system') {
        return stored;
    }
    return 'system';
}

type ThemeProviderPropsType = {
    children: ReactNode;
};

function ThemeProvider({ children }: ThemeProviderPropsType) {
    const [theme, setThemeState] = useState<ThemeType>(getStoredTheme);
    const [resolvedTheme, setResolvedTheme] = useState<'light' | 'dark'>(
        theme === 'system' ? getSystemTheme() : theme
    );

    useEffect(() => {
        const resolved = theme === 'system' ? getSystemTheme() : theme;
        setResolvedTheme(resolved);

        const root = document.documentElement;
        root.classList.remove('light', 'dark');
        root.classList.add(resolved);
    }, [theme]);

    useEffect(() => {
        if (theme !== 'system') return;

        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        const handleChange = () => {
            setResolvedTheme(getSystemTheme());
            document.documentElement.classList.remove('light', 'dark');
            document.documentElement.classList.add(getSystemTheme());
        };

        mediaQuery.addEventListener('change', handleChange);
        return () => mediaQuery.removeEventListener('change', handleChange);
    }, [theme]);

    const setTheme = (newTheme: ThemeType) => {
        localStorage.setItem(STORAGE_KEY, newTheme);
        setThemeState(newTheme);
    };

    return (
        <ThemeContext value={{ theme, resolvedTheme, setTheme }}>
            {children}
        </ThemeContext>
    );
}

function useTheme() {
    const context = useContext(ThemeContext);
    if (!context) {
        throw new Error('useTheme must be used within a ThemeProvider');
    }
    return context;
}

export { ThemeProvider, useTheme };
