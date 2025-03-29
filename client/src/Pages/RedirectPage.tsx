import { useParams, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'

export default function RedirectPage() {
    const { shortCode } = useParams()
    const navigate = useNavigate()
    const [error, setError] = useState<string>("")
    const [, setIsLoading] = useState(true)
    
    useEffect(() => {
        // Try the server-side redirect approach directly
        if (shortCode) {
            const backendUrl = import.meta.env.VITE_API_URL || '';
            
            // This uses a server-side redirect which bypasses CORS issues
            window.location.replace(`${backendUrl}/api/${shortCode}/redirect`);
            
            // Set a timeout to show an error if the redirect doesn't happen
            const timeoutId = setTimeout(() => {
                setError("Redirect timeout - the server didn't respond");
                setIsLoading(false);
            }, 5000);
            
            // Clean up the timeout if the component unmounts
            return () => clearTimeout(timeoutId);
        } else {
            setError("Invalid URL code");
            setIsLoading(false);
        }
    }, [shortCode]);

    // Handler for manual redirect attempt - same approach that's already working
    const tryRedirectAgain = () => {
        if (!shortCode) return;
        
        setIsLoading(true);
        setError("");
        
        const backendUrl = import.meta.env.VITE_API_URL || '';
        window.location.replace(`${backendUrl}/api/${shortCode}/redirect`);
    }

    return (
        <div className="flex items-center justify-center h-screen bg-gray-50">
            <div className="w-full max-w-sm p-6 bg-white rounded shadow-sm">
                {error ? (
                    <div className="text-center">
                        <div className="mb-4 text-red-500">●</div>
                        <h1 className="mb-2 text-xl font-medium text-gray-800">{error}</h1>
                        <p className="mb-4 text-sm text-gray-600">
                            Short code: <code className="px-1 bg-gray-100 rounded">{shortCode}</code>
                        </p>
                        <div className="flex flex-col space-y-2 sm:flex-row sm:space-y-0 sm:space-x-2 justify-center">
                            <button
                                onClick={() => navigate('/')}
                                className="px-4 py-2 text-sm text-white bg-gray-800 rounded hover:bg-gray-700"
                            >
                                Back to Home
                            </button>
                            <button
                                onClick={tryRedirectAgain}
                                className="px-4 py-2 text-sm text-gray-800 bg-gray-200 rounded hover:bg-gray-300"
                            >
                                Try Again
                            </button>
                        </div>
                    </div>
                ) : (
                    <div className="text-center">
                        <div className="w-6 h-6 mx-auto mb-4 border-2 border-gray-800 border-t-transparent rounded-full animate-spin"></div>
                        <p className="text-sm text-gray-600">Redirecting...</p>
                    </div>
                )}
            </div>
        </div>
    )
}