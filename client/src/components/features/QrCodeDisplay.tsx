import { useEffect, useRef, useState } from 'react';
import QRCode from 'qrcode';
import { Download } from 'lucide-react';
import { Button } from '@/components/ui/Button';
import { useTheme } from '@/components/providers/ThemeProvider';

function QrCodeDisplay({ url, size = 200 }: { url: string; size?: number }) {
    const canvasRef = useRef<HTMLCanvasElement>(null);
    const { resolvedTheme } = useTheme();
    const [isReady, setIsReady] = useState(false);

    useEffect(() => {
        if (!canvasRef.current || !url) return;

        const colors = {
            dark: { dark: '#fafafa', light: '#09090b' },
            light: { dark: '#09090b', light: '#ffffff' },
        };

        const colorConfig = colors[resolvedTheme];

        QRCode.toCanvas(canvasRef.current, url, {
            width: size,
            margin: 2,
            color: colorConfig,
            errorCorrectionLevel: 'H',
        }).then(() => {
            const canvas = canvasRef.current;
            if (!canvas) return;

            const ctx = canvas.getContext('2d');
            if (!ctx) return;

            const logoSize = size * 0.22;
            const logoX = (size - logoSize) / 2;
            const logoY = (size - logoSize) / 2;
            const padding = 4;

            ctx.fillStyle = colorConfig.light;
            ctx.fillRect(logoX - padding, logoY - padding, logoSize + padding * 2, logoSize + padding * 2);

            ctx.fillStyle = colorConfig.dark;
            ctx.font = `bold ${logoSize * 0.35}px 'Geist', system-ui, sans-serif`;
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.fillText('shrtn', size / 2, size / 2);

            setIsReady(true);
        });
    }, [url, size, resolvedTheme]);

    const handleDownload = () => {
        if (!canvasRef.current) return;

        const link = document.createElement('a');
        link.download = `shrtn-qr-${Date.now()}.png`;
        link.href = canvasRef.current.toDataURL('image/png');
        link.click();
    };

    return (
        <div className="flex flex-col items-center gap-3">
            <div className="border border-border p-3 bg-background">
                <canvas ref={canvasRef} className="block" />
            </div>
            {isReady && (
                <Button variant="outline" size="sm" onClick={handleDownload}>
                    <Download size={14} />
                    Download QR
                </Button>
            )}
        </div>
    );
}

export { QrCodeDisplay };
