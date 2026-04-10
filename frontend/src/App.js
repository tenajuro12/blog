import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './AuthContext';
import { ThemeProvider, useTheme } from './ThemeContext';
import Navbar from './components/Navbar';
import Home from './pages/Home';
import { Login, Register } from './pages/Auth';
import PostDetail from './pages/PostDetail';
import { NewPost, EditPost } from './pages/PostForm';
import Profile from './pages/Profile';
import { TagsPage, TagPage } from './pages/Tags';
import Bookmarks from './pages/Bookmarks';
import Feed from './pages/Feed';

function ProtectedRoute({ children }) {
  const { user, loading } = useAuth();
  const { theme } = useTheme();
  if (loading) return (
    <div style={{ display:'flex', alignItems:'center', justifyContent:'center', height:'100vh', background: theme.bg }}>
      <div style={{ width:32, height:32, border:`2px solid ${theme.border}`, borderTopColor: theme.accent, borderRadius:'50%', animation:'spin 0.8s linear infinite' }} />
    </div>
  );
  return user ? children : <Navigate to="/login" />;
}

function AppInner() {
  const { theme } = useTheme();
  return (
    <div style={{ minHeight:'100vh', background: theme.bg, transition:'background 0.3s ease' }}>
      <Navbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/posts/:slug" element={<PostDetail />} />
        <Route path="/profile/:username" element={<Profile />} />
        <Route path="/tags" element={<TagsPage />} />
        <Route path="/tags/:name" element={<TagPage />} />
        <Route path="/new-post" element={<ProtectedRoute><NewPost /></ProtectedRoute>} />
        <Route path="/edit/:slug" element={<ProtectedRoute><EditPost /></ProtectedRoute>} />
        <Route path="/bookmarks" element={<ProtectedRoute><Bookmarks /></ProtectedRoute>} />
        <Route path="/feed" element={<ProtectedRoute><Feed /></ProtectedRoute>} />
      </Routes>
    </div>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <BrowserRouter>
          <AppInner />
        </BrowserRouter>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
