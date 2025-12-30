import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useOptimistic, useTransition, useState, useCallback } from 'react';
import { shortenUrl, ApiError } from '@/api/client';
import type { UrlType, UrlInputType, ShortenerStateType } from '@/types/api';

const INITIAL_STATE: ShortenerStateType = {
    url: null,
    error: null,
    isPending: false,
};

export function useShortener() {
    const queryClient = useQueryClient();
    const [isPending, startTransition] = useTransition();
    const [state, setState] = useState<ShortenerStateType>(INITIAL_STATE);
    const [optimisticUrl, setOptimisticUrl] = useOptimistic<UrlType | null>(null);

    const mutation = useMutation({
        mutationFn: shortenUrl,
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ['urls'] });
            setState({ url: data, error: null, isPending: false });
        },
        onError: (error) => {
            const message = error instanceof ApiError
                ? error.message
                : 'Failed to shorten URL';
            setState({ url: null, error: message, isPending: false });
        },
    });

    const shorten = useCallback((input: UrlInputType) => {
        setState({ url: null, error: null, isPending: true });

        const optimisticResult: UrlType = {
            id: 'temp',
            shortCode: input.customCode || '...',
            originalUrl: input.originalUrl,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            clicks: 0,
        };

        startTransition(() => {
            setOptimisticUrl(optimisticResult);
            mutation.mutate(input);
        });
    }, [mutation, setOptimisticUrl]);

    const reset = useCallback(() => {
        setState(INITIAL_STATE);
    }, []);

    return {
        url: state.url,
        optimisticUrl,
        error: state.error,
        isPending: isPending || state.isPending,
        shorten,
        reset,
    };
}
