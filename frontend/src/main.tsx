import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { BrowserRouter, Route, Routes } from 'react-router'
import UploadDetail from './pages/UploadDetail.tsx'
import Video from './pages/Video.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path='/' element={<App />} />
        <Route path='/upload-detail/:id' element={<UploadDetail />} />
        <Route path='/video' element={<Video />} />
      </Routes>
    </BrowserRouter>,
  </StrictMode>,
)
