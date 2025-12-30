import { useQuery } from '@tanstack/react-query';
import { fetchLinkPreview } from '@/api/client';

export function useLinkPreview(url: string | null) {
    return useQuery({
        queryKey: ['linkPreview', url],
        queryFn: () => fetchLinkPreview(url!),
        enabled: !!url,
        staleTime: 1000 * 60 * 30,
        retry: 1,
    });
}
