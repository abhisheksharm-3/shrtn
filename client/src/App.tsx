import { Routes, Route } from 'react-router-dom'
import HomePage from './Pages/HomePage'
import RedirectPage from './Pages/RedirectPage'
import NotFoundPage from './Pages/NotFoundPage'

function App() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/:shortCode" element={<RedirectPage />} />
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}

export default App