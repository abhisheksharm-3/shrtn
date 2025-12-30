import { Routes, Route } from 'react-router-dom';
import HomePage from '@/pages/HomePage';
import RedirectPage from '@/pages/RedirectPage';
import NotFoundPage from '@/pages/NotFoundPage';
import TermsPage from '@/pages/TermsPage';
import PrivacyPage from '@/pages/PrivacyPage';

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/terms" element={<TermsPage />} />
      <Route path="/privacy" element={<PrivacyPage />} />
      <Route path="/:shortCode" element={<RedirectPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  );
}

export default App;