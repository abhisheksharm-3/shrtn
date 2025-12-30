import { Link } from 'react-router-dom';
import { Navbar } from '@/components/layout/Navbar';
import { ArrowLeft } from 'lucide-react';

function TermsPage() {
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

                    <h1 className="text-3xl font-bold mb-8">Terms of Service</h1>

                    <div className="prose prose-zinc dark:prose-invert max-w-none space-y-6 text-muted-foreground">
                        <p className="text-foreground font-medium">
                            Last updated: December 2024
                        </p>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">1. Acceptance of Terms</h2>
                            <p>
                                By accessing and using shrtn ("the Service"), you accept and agree to be bound by these Terms of Service. If you do not agree to these terms, please do not use the Service.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">2. Description of Service</h2>
                            <p>
                                shrtn provides a URL shortening service that allows users to create shortened versions of long URLs. The Service includes the creation of short links, QR code generation, and link analytics.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">3. Acceptable Use</h2>
                            <p>You agree not to use the Service to:</p>
                            <ul className="list-disc pl-6 space-y-1">
                                <li>Shorten URLs that lead to malicious, harmful, or illegal content</li>
                                <li>Distribute spam, phishing attempts, or malware</li>
                                <li>Violate any applicable laws or regulations</li>
                                <li>Infringe on the intellectual property rights of others</li>
                                <li>Attempt to circumvent rate limits or abuse the Service</li>
                            </ul>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">4. Content Responsibility</h2>
                            <p>
                                You are solely responsible for the content of the URLs you shorten. We reserve the right to disable any shortened URL that violates these terms without prior notice.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">5. Service Availability</h2>
                            <p>
                                We strive to maintain high availability but do not guarantee uninterrupted access to the Service. We may modify, suspend, or discontinue the Service at any time without notice.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">6. Limitation of Liability</h2>
                            <p>
                                The Service is provided "as is" without warranties of any kind. We shall not be liable for any indirect, incidental, or consequential damages arising from your use of the Service.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">7. Changes to Terms</h2>
                            <p>
                                We reserve the right to modify these terms at any time. Continued use of the Service after changes constitutes acceptance of the new terms.
                            </p>
                        </section>

                        <section className="space-y-3">
                            <h2 className="text-xl font-semibold text-foreground">8. Contact</h2>
                            <p>
                                For questions about these Terms, please open an issue on our GitHub repository.
                            </p>
                        </section>
                    </div>
                </article>
            </main>
        </div>
    );
}

export default TermsPage;
