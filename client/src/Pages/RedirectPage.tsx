import { useParams, useNavigate } from 'react-router-dom'
import { useEffect, useState } from 'react'

export default function RedirectPage() {
    const { shortCode } = useParams()
    const navigate = useNavigate()
    const [error, setError] = useState<string>("")
    const [, setIsLoading] = useState(true)
    
    useEffect(() => {
        async function fetchOriginalUrl() {
            try {
                const backendUrl = import.meta.env.VITE_API_URL || '';
                const response = await fetch(`${backendUrl}/api/${shortCode}`)
                
                if (!response.ok) {
                    if (response.status === 404) {
                        setError('URL not found')
                        setIsLoading(false)
                        return
                    }
                    throw new Error('Failed to fetch URL')
                }
                
                const data = await response.json()
                window.location.href = data.OriginalURL
            } catch (err) {
                console.error('Error:', err)
                setError('Something went wrong')
                setIsLoading(false)
            }
        }
        
        fetchOriginalUrl()
    }, [shortCode])

    return (
        <div className="flex items-center justify-center h-screen bg-gray-50">
            <div className="w-full max-w-sm p-6 bg-white rounded shadow-sm">
                {error ? (
                    <div className="text-center">
                        <div className="mb-4 text-red-500">●</div>
                        <h1 className="mb-2 text-xl font-medium text-gray-800">{error}</h1>
                        <p className="mb-4 text-sm text-gray-600">
                            Invalid code: <code className="px-1 bg-gray-100 rounded">{shortCode}</code>
                        </p>
                        <button
                            onClick={() => navigate('/')}
                            className="px-4 py-2 text-sm text-white bg-gray-800 rounded hover:bg-gray-700"
                        >
                            Back to Home
                        </button>
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
