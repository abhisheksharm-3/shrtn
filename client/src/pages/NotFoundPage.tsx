import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { Card, CardContent } from '@/components/ui/Card';

function NotFoundPage() {
    return (
        <div className="min-h-screen bg-background flex items-center justify-center p-4">
            <Card className="w-full max-w-sm">
                <CardContent className="p-6 text-center space-y-4">
                    <div className="text-4xl font-bold">404</div>
                    <h1 className="text-lg font-medium">Page not found</h1>
                    <p className="text-sm text-muted-foreground">
                        The page you're looking for doesn't exist or has been moved.
                    </p>
                    <Button asChild className="w-full">
                        <Link to="/">Go home</Link>
                    </Button>
                </CardContent>
            </Card>
        </div>
    );
}

export default NotFoundPage;
