import { useParams, useNavigate } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { Button } from '@/components/ui/Button';
import { Card, CardContent } from '@/components/ui/Card';
import { Spinner } from '@/components/ui/Spinner';
import { getRedirectUrl } from '@/api/client';

const REDIRECT_TIMEOUT_MS = 5000;

function RedirectPage() {
    const { shortCode } = useParams<{ shortCode: string }>();
    const navigate = useNavigate();
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!shortCode) {
            setError('Invalid short code');
            return;
        }

        const redirectUrl = getRedirectUrl(shortCode);
        window.location.replace(redirectUrl);

        const timeoutId = setTimeout(() => {
            setError('Redirect timed out. The server may be unavailable.');
        }, REDIRECT_TIMEOUT_MS);

        return () => clearTimeout(timeoutId);
    }, [shortCode]);

    const handleGoHome = () => {
        navigate('/');
    };

    const handleRetry = () => {
        if (!shortCode) return;
        setError(null);
        const redirectUrl = getRedirectUrl(shortCode);
        window.location.replace(redirectUrl);
    };

    if (error) {
        return (
            <div className="min-h-screen bg-background flex items-center justify-center p-4">
                <Card className="w-full max-w-sm">
                    <CardContent className="p-6 text-center space-y-4">
                        <div className="text-destructive text-2xl">×</div>
                        <h1 className="text-lg font-medium">{error}</h1>
                        <p className="text-sm text-muted-foreground">
                            Short code: <code className="px-1 py-0.5 bg-secondary">{shortCode}</code>
                        </p>
                        <div className="flex gap-2">
                            <Button variant="outline" className="flex-1" onClick={handleGoHome}>
                                Go home
                            </Button>
                            <Button className="flex-1" onClick={handleRetry}>
                                Retry
                            </Button>
                        </div>
                    </CardContent>
                </Card>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-background flex items-center justify-center p-4">
            <Card className="w-full max-w-sm">
                <CardContent className="p-6 text-center space-y-4">
                    <Spinner className="mx-auto" />
                    <p className="text-sm text-muted-foreground">Redirecting...</p>
                </CardContent>
            </Card>
        </div>
    );
}

export default RedirectPage;
