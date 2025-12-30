import { useState } from 'react';
import { Link as LinkIcon, Menu, X, Github } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { ThemeSwitcher } from '@/components/ui/ThemeSwitcher';

function Navbar() {
    const [isMenuOpen, setIsMenuOpen] = useState(false);

    const handleToggleMenu = () => {
        setIsMenuOpen(!isMenuOpen);
    };

    const handleCloseMenu = () => {
        setIsMenuOpen(false);
    };

    return (
        <nav className="border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 fixed top-0 left-0 right-0 z-50">
            <div className="max-w-5xl mx-auto px-4 sm:px-6">
                <div className="flex items-center justify-between h-14">
                    <Link to="/" className="flex items-center gap-1.5 group">
                        <LinkIcon size={18} className="text-foreground" />
                        <span className="font-semibold text-lg tracking-tight">shrtn</span>
                    </Link>

                    <div className="hidden md:flex items-center gap-4">
                        <a
                            href="https://github.com/abhisheksharm-3/shrtn"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
                        >
                            <Github size={16} />
                            <span>GitHub</span>
                        </a>
                        <ThemeSwitcher />
                    </div>

                    <div className="flex items-center gap-2 md:hidden">
                        <ThemeSwitcher />
                        <Button
                            variant="ghost"
                            size="icon"
                            onClick={handleToggleMenu}
                            aria-expanded={isMenuOpen}
                            aria-label="Toggle menu"
                        >
                            {isMenuOpen ? <X size={18} /> : <Menu size={18} />}
                        </Button>
                    </div>
                </div>
            </div>

            {isMenuOpen && (
                <div className="md:hidden border-t border-border">
                    <div className="px-4 py-3 space-y-2">
                        <a
                            href="https://github.com/abhisheksharm-3/shrtn"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-2 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
                            onClick={handleCloseMenu}
                        >
                            <Github size={16} />
                            <span>View on GitHub</span>
                        </a>
                    </div>
                </div>
            )}
        </nav>
    );
}

export { Navbar };
