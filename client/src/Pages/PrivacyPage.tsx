import { Link } from 'react-router-dom';
import { Navbar } from '@/components/layout/Navbar';
import { ArrowLeft } from 'lucide-react';

function PrivacyPage() {
    return (
        <div className="min-h-screen bg-background">
            <Navbar />

            <main className="pt-14">
                <article className="max-w-2xl mx-auto px-4 py-12">
                    <Link
                        to="/"
                        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground mb-8"
                    >
                        <ArrowLeft size={14} />
                        Back to home
                    </Link>

                    <h1 className="text-3xl font-bold mb-8">Privacy Policy</h1>

                    <div className="prose prose-zinc dark:prose-invert max-w-none space-y-6 text-muted-foreground">
                        <p className="text-foreground font-medium">
                            Last updated: December 2024
                        </p>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">1. Information We Collect</h2>
                            <p>When you use shrtn, we collect minimal information necessary to provide the Service:</p>
                            <ul className="list-disc pl-6 space-y-1">
                                <li><strong>URLs:</strong> The original URLs you submit for shortening</li>
                                <li><strong>Click Analytics:</strong> Timestamp, referrer, and anonymized IP address when someone clicks a shortened link</li>
                                <li><strong>Technical Data:</strong> Browser type and device information for analytics</li>
                            </ul>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">2. How We Use Your Information</h2>
                            <p>We use the collected information to:</p>
                            <ul className="list-disc pl-6 space-y-1">
                                <li>Provide and maintain the URL shortening service</li>
                                <li>Generate click statistics and analytics</li>
                                <li>Improve the Service and user experience</li>
                                <li>Detect and prevent abuse or fraudulent activity</li>
                            </ul>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">3. Data Storage</h2>
                            <p>
                                Your data is stored securely using industry-standard practices. Shortened URLs and their associated data are retained until explicitly deleted or the Service is discontinued.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">4. Data Sharing</h2>
                            <p>
                                We do not sell, trade, or otherwise transfer your information to third parties. We may share anonymized, aggregated statistics that cannot be used to identify individuals.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">5. Cookies</h2>
                            <p>
                                We use minimal cookies to store your theme preference (light/dark mode). We do not use tracking cookies or third-party analytics services.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">6. Your Rights</h2>
                            <p>You have the right to:</p>
                            <ul className="list-disc pl-6 space-y-1">
                                <li>Request access to your data</li>
                                <li>Request deletion of your shortened URLs</li>
                                <li>Opt out of analytics collection</li>
                            </ul>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">7. Security</h2>
                            <p>
                                We implement appropriate security measures to protect your data, including encryption in transit (HTTPS), rate limiting, and secure authentication for API access.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">8. Changes to This Policy</h2>
                            <p>
                                We may update this Privacy Policy from time to time. We will notify users of any material changes by updating the "Last updated" date.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">9. Contact</h2>
                            <p>
                                For privacy-related questions or requests, please open an issue on our GitHub repository.
                            </p>
                        </section>
                    </div>
                </article>
            </main>
        </div>
    );
}

export default PrivacyPage;
