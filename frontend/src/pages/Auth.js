import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { login, register } from '../api';
import { useAuth } from '../AuthContext';
import { useTheme } from '../ThemeContext';

function AuthLayout({ title, subtitle, children, switchText, switchLink, switchLabel }) {
  const { theme, toggle, mode } = useTheme();
  return (
    <div style={{ display: 'flex', minHeight: '100vh', background: theme.bg, transition: 'background 0.3s ease' }}>
      <div style={{ width: '45%', background: theme.isDark ? 'linear-gradient(160deg, #0f0f18 0%, #1a0a04 100%)' : 'linear-gradient(160deg, #f5ede4 0%, #fde8d8 100%)', position: 'relative', overflow: 'hidden', borderRight: `1px solid ${theme.border}`, display: 'flex', alignItems: 'center' }}>
        <div style={{ padding: '60px 56px', width: '100%', position: 'relative', zIndex: 2 }}>
          <Link to="/" style={{ display: 'flex', alignItems: 'center', gap: 8, textDecoration: 'none', marginBottom: 80 }}>
            <span style={{ color: theme.accent, fontSize: 8 }}>●</span>
            <span style={{ fontFamily: "'Playfair Display', serif", fontSize: 16, fontWeight: 700, letterSpacing: 5, color: theme.text }}>INKWELL</span>
          </Link>
          <div style={{ maxWidth: 320 }}>
            <div style={{ fontFamily: "'Playfair Display', serif", fontStyle: 'italic', fontSize: 26, lineHeight: 1.5, color: theme.isDark ? 'rgba(255,255,255,0.75)' : theme.text, marginBottom: 20 }}>
              "The scariest moment is always just before you start."
            </div>
            <div style={{ fontSize: 12, letterSpacing: 2, color: theme.accent, fontWeight: 500, opacity: 0.8 }}>— Stephen King</div>
          </div>
        </div>
        <div style={{ position: 'absolute', bottom: -60, right: -60, width: 300, height: 300, borderRadius: '50%', background: `radial-gradient(circle, ${theme.accentLight} 0%, transparent 70%)` }} />
      </div>

      <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '60px 48px', position: 'relative' }}>
        <button onClick={toggle} style={{ position: 'absolute', top: 24, right: 32, background: theme.inputBg, border: `1px solid ${theme.inputBorder}`, borderRadius: 20, padding: '5px 12px', cursor: 'pointer', fontSize: 14, color: theme.textSecondary, fontFamily: "'DM Sans', sans-serif", display: 'flex', alignItems: 'center', gap: 6 }}>
          {mode === 'dark' ? '☀️' : '🌙'} <span style={{ fontSize: 11 }}>{mode === 'dark' ? 'Light' : 'Dark'}</span>
        </button>
        <div style={{ width: '100%', maxWidth: 400, animation: 'fadeUp 0.6s ease both' }}>
          <div style={{ marginBottom: 40 }}>
            <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 38, fontWeight: 700, color: theme.text, marginBottom: 8, letterSpacing: -0.5 }}>{title}</h1>
            <p style={{ color: theme.textMuted, fontSize: 15 }}>{subtitle}</p>
          </div>
          {children}
          <div style={{ marginTop: 28, display: 'flex', gap: 8, justifyContent: 'center', fontSize: 14 }}>
            <span style={{ color: theme.textMuted }}>{switchText}</span>
            <Link to={switchLink} style={{ color: theme.accent, textDecoration: 'none', fontWeight: 600 }}>{switchLabel}</Link>
          </div>
        </div>
      </div>
    </div>
  );
}

export function Login() {
  const [form, setForm] = useState({ email: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { signin } = useAuth();
  const { theme } = useTheme();
  const navigate = useNavigate();

  const handleSubmit = async () => {
    setLoading(true);
    try { const res = await login(form); signin(res.data.token, res.data.user); navigate('/'); }
    catch (e) { setError(e.response?.data?.error || 'Invalid credentials'); }
    finally { setLoading(false); }
  };

  return (
    <AuthLayout title="Welcome back." subtitle="Sign in to continue your story."
      switchText="Don't have an account?" switchLink="/register" switchLabel="Sign up →">
      {error && <div style={{ background: theme.accentLight, border: `1px solid ${theme.accentBorder}`, color: theme.accent, borderRadius: 8, padding: '12px 16px', fontSize: 13, marginBottom: 20 }}>{error}</div>}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
        {[['Email', 'email', 'email', 'you@example.com'], ['Password', 'password', 'password', '••••••••']].map(([label, key, type, ph]) => (
          <div key={key} style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <label style={{ fontSize: 11, letterSpacing: 1.5, color: theme.textMuted, fontWeight: 500, textTransform: 'uppercase' }}>{label}</label>
            <input style={{ background: theme.inputBg, border: `1px solid ${theme.inputBorder}`, borderRadius: 8, padding: '13px 16px', fontSize: 15, color: theme.text, outline: 'none', transition: 'border-color 0.2s' }}
              type={type} placeholder={ph} value={form[key]} onChange={e => setForm({ ...form, [key]: e.target.value })}
              onKeyDown={e => e.key === 'Enter' && handleSubmit()} />
          </div>
        ))}
        <button style={{ background: theme.accent, border: 'none', borderRadius: 8, padding: 14, fontSize: 15, fontWeight: 600, color: '#fff', cursor: 'pointer', fontFamily: "'DM Sans', sans-serif", opacity: loading ? 0.6 : 1, marginTop: 8 }} onClick={handleSubmit} disabled={loading}>
          {loading ? 'Signing in...' : 'Sign in →'}
        </button>
      </div>
    </AuthLayout>
  );
}

export function Register() {
  const [form, setForm] = useState({ username: '', email: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { signin } = useAuth();
  const { theme } = useTheme();
  const navigate = useNavigate();

  const handleSubmit = async () => {
    setLoading(true);
    try { const res = await register(form); signin(res.data.token, res.data.user); navigate('/'); }
    catch (e) { setError(e.response?.data?.error || 'Registration failed'); }
    finally { setLoading(false); }
  };

  return (
    <AuthLayout title="Start writing." subtitle="Create your account and share your ideas."
      switchText="Already have an account?" switchLink="/login" switchLabel="Sign in →">
      {error && <div style={{ background: theme.accentLight, border: `1px solid ${theme.accentBorder}`, color: theme.accent, borderRadius: 8, padding: '12px 16px', fontSize: 13, marginBottom: 20 }}>{error}</div>}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>
        {[['Username', 'username', 'text', 'yourname'], ['Email', 'email', 'email', 'you@example.com'], ['Password', 'password', 'password', '••••••••']].map(([label, key, type, ph]) => (
          <div key={key} style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
            <label style={{ fontSize: 11, letterSpacing: 1.5, color: theme.textMuted, fontWeight: 500, textTransform: 'uppercase' }}>{label}</label>
            <input style={{ background: theme.inputBg, border: `1px solid ${theme.inputBorder}`, borderRadius: 8, padding: '13px 16px', fontSize: 15, color: theme.text, outline: 'none' }}
              type={type} placeholder={ph} value={form[key]} onChange={e => setForm({ ...form, [key]: e.target.value })}
              onKeyDown={e => e.key === 'Enter' && handleSubmit()} />
          </div>
        ))}
        <button style={{ background: theme.accent, border: 'none', borderRadius: 8, padding: 14, fontSize: 15, fontWeight: 600, color: '#fff', cursor: 'pointer', fontFamily: "'DM Sans', sans-serif", opacity: loading ? 0.6 : 1, marginTop: 8 }} onClick={handleSubmit} disabled={loading}>
          {loading ? 'Creating account...' : 'Create account →'}
        </button>
      </div>
    </AuthLayout>
  );
}
