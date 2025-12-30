import { ExternalLink } from 'lucide-react';
import { useLinkPreview } from '@/hooks/useLinkPreview';
import { Spinner } from '@/components/ui/Spinner';

function LinkPreviewCard({ url }: { url: string }) {
    const { data: preview, isLoading, isError } = useLinkPreview(url);

    if (isLoading) {
        return (
            <div className="border border-border p-4 flex items-center justify-center">
                <Spinner size="sm" />
                <span className="ml-2 text-sm text-muted-foreground">Loading preview...</span>
            </div>
        );
    }

    if (isError || !preview) {
        return null;
    }

    const hasContent = preview.title || preview.description || preview.image;
    if (!hasContent) {
        return null;
    }

    return (
        <a
            href={url}
            target="_blank"
            rel="noopener noreferrer"
            className="block border border-border hover:border-foreground/20 transition-colors overflow-hidden group"
        >
            {preview.image && (
                <div className="aspect-video bg-secondary overflow-hidden">
                    <img
                        src={preview.image}
                        alt=""
                        className="w-full h-full object-cover"
                        onError={(e) => {
                            e.currentTarget.style.display = 'none';
                        }}
                    />
                </div>
            )}
            <div className="p-3 space-y-1">
                <div className="flex items-start justify-between gap-2">
                    <div className="flex items-center gap-2 min-w-0">
                        {preview.favicon && (
                            <img
                                src={preview.favicon}
                                alt=""
                                className="size-4 shrink-0"
                                onError={(e) => {
                                    e.currentTarget.style.display = 'none';
                                }}
                            />
                        )}
                        <span className="text-xs text-muted-foreground truncate">
                            {preview.siteName || new URL(url).hostname}
                        </span>
                    </div>
                    <ExternalLink size={12} className="shrink-0 text-muted-foreground group-hover:text-foreground transition-colors" />
                </div>
                {preview.title && (
                    <h3 className="text-sm font-medium line-clamp-2">{preview.title}</h3>
                )}
                {preview.description && (
                    <p className="text-xs text-muted-foreground line-clamp-2">{preview.description}</p>
                )}
            </div>
        </a>
    );
}

export { LinkPreviewCard };
