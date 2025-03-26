import { useState } from 'react';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Check, Copy, Loader2, AlertCircle } from "lucide-react";

const Hero = () => {
  const [longUrl, setLongUrl] = useState('');
  const [customCode, setCustomCode] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);
  const [showCustomCode, setShowCustomCode] = useState(false);

  const handleSubmit = async (e: { preventDefault: () => void; }) => {
    e.preventDefault();
    if (!longUrl.trim()) return;

    setIsLoading(true);
    setError('');
    setShortUrl('');
    
    try {
      const backendUrl = import.meta.env.VITE_API_URL || '';
      const response = await fetch(`${backendUrl}/api/shorten`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          OriginalURL: longUrl,
          CustomCode: customCode.trim() || undefined
        }),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Failed to shorten URL');
      }

      const data = await response.json();
      // Format the short URL with the domain and short code
      const baseUrl = window.location.origin;
      setShortUrl(`${baseUrl}/${data.ShortCode}`);
    } catch (err) {
      setError(
        err instanceof Error 
          ? err.message 
          : 'An error occurred while shortening the URL'
      );
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(shortUrl).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  };

  return (
    <div className="min-h-screen bg-gray-950 text-white flex flex-col md:flex-row items-center justify-center px-4 md:px-8 lg:px-16 py-16 relative overflow-hidden">
      {/* Background gradient */}
      <div className="absolute top-0 right-0 w-full h-full bg-gradient-to-br from-blue-500/10 to-purple-500/5 opacity-50"></div>
      
      {/* Blob */}
      <div className="absolute bottom-0 right-0 w-2/3 h-2/3 bg-blue-600/20 rounded-full filter blur-3xl -z-10 transform translate-x-1/4 translate-y-1/4"></div>
      
      {/* Left side content */}
      <div className="z-10 w-full md:w-1/2 mb-12 md:mb-0">
        <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold mb-4">
          <span className="text-blue-400">Shorten</span> Your URLs,<br />
          <span className="text-blue-400">Expand</span> Your Reach
        </h1>
        
        <p className="text-gray-300 text-base md:text-lg mb-10 max-w-lg">
          Transform long, unwieldy links into clean, memorable, and trackable 
          short URLs in seconds. Boost your click-through rates with our 
          powerful URL shortening service.
        </p>
        
        <div className="bg-gray-900/70 backdrop-blur-sm p-6 rounded-xl max-w-lg">
          <form onSubmit={handleSubmit}>
            <div className="flex flex-col gap-3 mb-4">
              <Input 
                placeholder="Enter your long URL here" 
                className="flex-1 bg-gray-800 border-gray-700 h-12"
                value={longUrl}
                onChange={(e) => setLongUrl(e.target.value)}
                disabled={isLoading}
              />
              
              {showCustomCode && (
                <Input 
                  placeholder="Custom short code (optional)" 
                  className="flex-1 bg-gray-800 border-gray-700 h-12"
                  value={customCode}
                  onChange={(e) => setCustomCode(e.target.value)}
                  disabled={isLoading}
                />
              )}
              
              <div className="flex flex-col md:flex-row gap-3">
                <Button 
                  type="button"
                  variant="outline"
                  className="h-12 border-gray-700 text-gray-300"
                  onClick={() => setShowCustomCode(!showCustomCode)}
                  disabled={isLoading}
                >
                  {showCustomCode ? 'Hide Custom Code' : 'Custom Code'}
                </Button>
                
                <Button 
                  type="submit"
                  className="bg-blue-500 hover:bg-blue-600 h-12 px-6 flex-1"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <Loader2 className="h-5 w-5 animate-spin mr-2" />
                  ) : null}
                  Shorten URL
                </Button>
              </div>
            </div>
          </form>
          
          {error && (
            <div className="flex items-center gap-2 text-red-400 mb-4">
              <AlertCircle size={16} />
              <p className="text-sm">{error}</p>
            </div>
          )}
          
          {shortUrl && (
            <div className="bg-gray-800 p-3 rounded-lg mb-4 flex items-center justify-between">
              <a 
                href={shortUrl} 
                target="_blank" 
                rel="noopener noreferrer"
                className="text-blue-400 hover:underline text-sm md:text-base truncate"
              >
                {shortUrl}
              </a>
              <Button 
                size="sm" 
                variant="ghost" 
                className="h-8 px-2 text-gray-300 hover:text-white ml-2 flex-shrink-0"
                onClick={copyToClipboard}
              >
                {copied ? (
                  <Check size={16} className="text-green-500" />
                ) : (
                  <Copy size={16} />
                )}
              </Button>
            </div>
          )}
          
          <p className="text-gray-400 text-sm">
            By using our service, you agree to our <a href="#" className="text-blue-400 hover:underline">Terms of Service</a> and <a href="#" className="text-blue-400 hover:underline">Privacy Policy</a>
          </p>
        </div>
        
        <div className="flex flex-wrap gap-6 mt-8">
          <div className="flex items-center gap-2">
            <div className="bg-green-500 rounded-full p-1">
              <Check size={16} className="text-black" />
            </div>
            <span className="text-gray-300">100% Free</span>
          </div>
          
          <div className="flex items-center gap-2">
            <div className="bg-green-500 rounded-full p-1">
              <Check size={16} className="text-black" />
            </div>
            <span className="text-gray-300">No Registration Required</span>
          </div>
          
          <div className="flex items-center gap-2">
            <div className="bg-green-500 rounded-full p-1">
              <Check size={16} className="text-black" />
            </div>
            <span className="text-gray-300">HTTPS Secure</span>
          </div>
        </div>
      </div>
      
      {/* Right side illustration */}
      <div className="z-10 w-full md:w-1/2 max-w-lg">
        <div className="bg-gray-800/60 backdrop-blur-md rounded-3xl p-6 border border-gray-700/50 shadow-2xl">
          {/* Browser window */}
          <div className="bg-gray-900 rounded-lg overflow-hidden">
            {/* Browser header */}
            <div className="bg-gray-800 px-4 py-2 flex items-center gap-1.5">
              <div className="h-3 w-3 rounded-full bg-red-500"></div>
              <div className="h-3 w-3 rounded-full bg-yellow-500"></div>
              <div className="h-3 w-3 rounded-full bg-green-500"></div>
              <div className="bg-gray-700 h-6 rounded ml-2 flex-1"></div>
            </div>
            
            {/* Browser content */}
            <div className="p-4 space-y-4">
              {/* Blue Item */}
              <div className="rounded-lg bg-gray-800 p-4">
                <div className="flex items-center gap-3">
                  <div className="h-5 w-5 rounded bg-blue-500"></div>
                  <div className="h-4 bg-gray-600 rounded flex-1"></div>
                </div>
                <div className="h-3 bg-gray-600 rounded mt-4 w-full"></div>
                <div className="h-3 bg-gray-600 rounded mt-2 w-full"></div>
                <div className="h-3 bg-gray-600 rounded mt-2 w-3/4"></div>
              </div>
              
              {/* Green Item */}
              <div className="rounded-lg bg-gray-800 p-4">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="h-5 w-5 rounded bg-green-500"></div>
                    <div className="h-4 bg-gray-600 rounded w-24"></div>
                  </div>
                  <div className="h-6 w-20 rounded bg-blue-500"></div>
                </div>
                <div className="flex mt-4 gap-3">
                  <div className="bg-blue-500 text-white text-xs py-1 px-3 rounded">
                    {shortUrl || "short.url/a1b2c"}
                  </div>
                  <div className="bg-gray-700 flex-1 rounded"></div>
                </div>
              </div>
              
              {/* Orange Item */}
              <div className="rounded-lg bg-gray-800 p-4">
                <div className="flex items-center gap-3">
                  <div className="h-5 w-5 rounded bg-orange-500"></div>
                  <div className="h-4 bg-gray-600 rounded flex-1"></div>
                </div>
                <div className="h-16 bg-gray-700 rounded mt-4 w-full"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Hero;