import { Navbar } from '@/components/layout/Navbar';
import { UrlShortener } from '@/components/features/UrlShortener';
import { Check } from 'lucide-react';

const FEATURES = [
    'Free forever',
    'No registration',
    'HTTPS secure',
] as const;

function HomePage() {
    return (
        <div className="min-h-screen bg-background">
            <Navbar />

            <main className="pt-14">
                <section className="min-h-[calc(100vh-3.5rem)] flex items-center justify-center px-4 py-16">
                    <div className="w-full max-w-5xl mx-auto grid lg:grid-cols-2 gap-12 lg:gap-16 items-center">
                        <div className="space-y-6">
                            <h1 className="text-4xl sm:text-5xl font-bold tracking-tight">
                                Shorten your URLs.
                                <br />
                                <span className="text-muted-foreground">Share them anywhere.</span>
                            </h1>

                            <p className="text-lg text-muted-foreground max-w-md">
                                Transform long links into clean, memorable short URLs. Track clicks and manage your links with ease.
                            </p>

                            <ul className="flex flex-wrap gap-4">
                                {FEATURES.map((feature) => (
                                    <li key={feature} className="flex items-center gap-2 text-sm">
                                        <span className="flex items-center justify-center size-5 bg-foreground text-background">
                                            <Check size={12} strokeWidth={3} />
                                        </span>
                                        <span>{feature}</span>
                                    </li>
                                ))}
                            </ul>
                        </div>

                        <div className="flex justify-center lg:justify-end">
                            <UrlShortener />
                        </div>
                    </div>
                </section>
            </main>
        </div>
    );
}

export default HomePage;
