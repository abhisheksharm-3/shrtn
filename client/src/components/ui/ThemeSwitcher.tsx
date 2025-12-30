import { Sun, Moon, Monitor } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { useTheme } from '@/components/providers/ThemeProvider';

const THEME_OPTIONS = [
    { value: 'light', icon: Sun, label: 'Light' },
    { value: 'dark', icon: Moon, label: 'Dark' },
    { value: 'system', icon: Monitor, label: 'System' },
] as const;

function ThemeSwitcher() {
    const { theme, setTheme } = useTheme();

    const handleCycle = () => {
        const currentIndex = THEME_OPTIONS.findIndex((opt) => opt.value === theme);
        const nextIndex = (currentIndex + 1) % THEME_OPTIONS.length;
        setTheme(THEME_OPTIONS[nextIndex].value);
    };

    const currentOption = THEME_OPTIONS.find((opt) => opt.value === theme) || THEME_OPTIONS[2];
    const Icon = currentOption.icon;

    return (
        <Button
            variant="ghost"
            size="icon"
            onClick={handleCycle}
            aria-label={`Current theme: ${currentOption.label}. Click to change.`}
            title={`Theme: ${currentOption.label}`}
        >
            <Icon size={16} />
        </Button>
    );
}

export { ThemeSwitcher };
