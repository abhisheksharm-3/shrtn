import { useState, type FormEvent, type ChangeEvent } from 'react';
import { Check, Copy, AlertCircle, ChevronDown, ChevronUp } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card, CardContent } from '@/components/ui/Card';
import { Spinner } from '@/components/ui/Spinner';
import { QrCodeDisplay } from '@/components/features/QrCodeDisplay';
import { LinkPreviewCard } from '@/components/features/LinkPreviewCard';
import { useShortener } from '@/hooks/useShortener';
import { Link } from 'react-router-dom';

function UrlShortener() {
    const [longUrl, setLongUrl] = useState('');
    const [customCode, setCustomCode] = useState('');
    const [isCopied, setIsCopied] = useState(false);
    const [isCustomCodeVisible, setIsCustomCodeVisible] = useState(false);

    const { url, optimisticUrl, error, isPending, shorten, reset } = useShortener();

    const displayUrl = url || optimisticUrl;
    const shortUrl = displayUrl ? `${window.location.origin}/${displayUrl.shortCode}` : null;
    const originalUrl = url?.originalUrl || null;

    const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        if (!longUrl.trim()) return;

        shorten({
            originalUrl: longUrl.trim(),
            customCode: customCode.trim() || undefined,
        });
    };

    const handleCopy = async () => {
        if (!shortUrl) return;
        await navigator.clipboard.writeText(shortUrl);
        setIsCopied(true);
        setTimeout(() => setIsCopied(false), 2000);
    };

    const handleReset = () => {
        setLongUrl('');
        setCustomCode('');
        setIsCopied(false);
        reset();
    };

    const handleToggleCustomCode = () => {
        setIsCustomCodeVisible(!isCustomCodeVisible);
    };

    return (
        <Card className="w-full max-w-lg">
            <CardContent className="p-6">
                {!shortUrl ? (
                    <>
                        <form onSubmit={handleSubmit} className="space-y-4">
                            <div className="space-y-3">
                                <Input
                                    placeholder="Enter your URL"
                                    value={longUrl}
                                    onChange={(e: ChangeEvent<HTMLInputElement>) => setLongUrl(e.target.value)}
                                    disabled={isPending}
                                    required
                                />

                                <button
                                    type="button"
                                    onClick={handleToggleCustomCode}
                                    className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground transition-colors"
                                >
                                    {isCustomCodeVisible ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
                                    Custom short code
                                </button>

                                {isCustomCodeVisible && (
                                    <Input
                                        placeholder="Custom code (optional)"
                                        value={customCode}
                                        onChange={(e: ChangeEvent<HTMLInputElement>) => setCustomCode(e.target.value)}
                                        disabled={isPending}
                                    />
                                )}
                            </div>

                            <Button type="submit" className="w-full" disabled={isPending || !longUrl.trim()}>
                                {isPending && <Spinner size="sm" className="mr-2" />}
                                Shorten
                            </Button>
                        </form>

                        {error && (
                            <div className="mt-4 flex items-center gap-2 text-destructive text-sm">
                                <AlertCircle size={14} />
                                <span>{error}</span>
                            </div>
                        )}
                    </>
                ) : (
                    <div className="space-y-6">
                        <div className="space-y-3">
                            <div className="flex items-center gap-2 p-3 bg-secondary border border-border">
                                <a
                                    href={shortUrl}
                                    target="_blank"
                                    rel="noopener noreferrer"
                                    className="flex-1 text-sm font-mono truncate hover:underline"
                                >
                                    {shortUrl}
                                </a>
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    onClick={handleCopy}
                                    aria-label={isCopied ? 'Copied' : 'Copy to clipboard'}
                                >
                                    {isCopied ? (
                                        <Check size={16} className="text-success" />
                                    ) : (
                                        <Copy size={16} />
                                    )}
                                </Button>
                            </div>
                        </div>

                        <div className="grid gap-4 sm:grid-cols-2">
                            <div>
                                <h4 className="text-xs font-medium text-muted-foreground mb-3 uppercase tracking-wide">
                                    QR Code
                                </h4>
                                <QrCodeDisplay url={shortUrl} size={160} />
                            </div>

                            {originalUrl && (
                                <div>
                                    <h4 className="text-xs font-medium text-muted-foreground mb-3 uppercase tracking-wide">
                                        Destination
                                    </h4>
                                    <LinkPreviewCard url={originalUrl} />
                                </div>
                            )}
                        </div>

                        <Button
                            type="button"
                            variant="outline"
                            onClick={handleReset}
                            className="w-full"
                        >
                            Shorten another URL
                        </Button>
                    </div>
                )}

                <p className="mt-6 text-xs text-muted-foreground text-center">
                    By using this service, you agree to our{' '}
                    <Link to="/terms" className="underline underline-offset-2 hover:text-foreground">
                        Terms
                    </Link>{' '}
                    and{' '}
                    <Link to="/privacy" className="underline underline-offset-2 hover:text-foreground">
                        Privacy Policy
                    </Link>
                </p>
            </CardContent>
        </Card>
    );
}

export { UrlShortener };
