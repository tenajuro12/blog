import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../AuthContext';
import { useTheme } from '../ThemeContext';

export default function Navbar() {
  const { user, signout } = useAuth();
  const { theme, toggle, mode } = useTheme();
  const navigate = useNavigate();
  const location = useLocation();
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 20);
    window.addEventListener('scroll', onScroll);
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  const active = (path) => location.pathname === path ? theme.text : theme.textMuted;

  return (
    <nav style={{
      position:'fixed', top:0, left:0, right:0, zIndex:100,
      display:'flex', justifyContent:'space-between', alignItems:'center',
      padding:'0 48px', height:64,
      background: scrolled ? theme.navBgScrolled : theme.navBg,
      backdropFilter:'blur(24px)',
      borderBottom:`1px solid ${theme.border}`,
      transition:'all 0.3s ease',
    }}>
      <Link to="/" style={{ display:'flex', alignItems:'center', gap:8, textDecoration:'none' }}>
        <span style={{ color:theme.accent, fontSize:8, animation:'pulse 2s infinite' }}>●</span>
        <span style={{ fontFamily:"'Playfair Display', serif", fontSize:17, fontWeight:700, letterSpacing:5, color:theme.text }}>INKWELL</span>
      </Link>

      <div style={{ display:'flex', gap:28, alignItems:'center' }}>
        <Link to="/" style={{ color:active('/'), textDecoration:'none', fontSize:13, fontWeight:500 }}>Stories</Link>
        <Link to="/tags" style={{ color:active('/tags'), textDecoration:'none', fontSize:13, fontWeight:500 }}>Tags</Link>
        {user && <>
          <Link to="/feed" style={{ color:active('/feed'), textDecoration:'none', fontSize:13, fontWeight:500 }}>Feed</Link>
          <Link to="/new-post" style={{ color:active('/new-post'), textDecoration:'none', fontSize:13, fontWeight:500 }}>Write</Link>
        </>}
      </div>

      <div style={{ display:'flex', alignItems:'center', gap:16 }}>
        <button onClick={toggle} style={{
          background:theme.inputBg, border:`1px solid ${theme.inputBorder}`,
          borderRadius:20, padding:'5px 12px', cursor:'pointer',
          fontSize:13, color:theme.textSecondary, fontFamily:"'DM Sans', sans-serif",
          display:'flex', alignItems:'center', gap:6,
        }}>
          {mode === 'dark' ? '☀️' : '🌙'}
          <span style={{ fontSize:11 }}>{mode === 'dark' ? 'Light' : 'Dark'}</span>
        </button>

        {user ? <>
          <Link to="/bookmarks" style={{ color:theme.textMuted, textDecoration:'none', fontSize:20 }} title="Bookmarks">🔖</Link>
          <Link to={`/profile/${user.username}`} style={{ display:'flex', alignItems:'center', gap:8, textDecoration:'none' }}>
            <div style={{ width:30, height:30, borderRadius:'50%', background:'linear-gradient(135deg,#ff6b35,#ff9f6b)', display:'flex', alignItems:'center', justifyContent:'center', fontSize:12, fontWeight:700, color:'#fff' }}>
              {user.username[0].toUpperCase()}
            </div>
            <span style={{ color:theme.textSecondary, fontSize:13 }}>{user.username}</span>
          </Link>
          <button onClick={() => { signout(); navigate('/'); }} style={{ background:'none', border:'none', cursor:'pointer', color:theme.textMuted, fontSize:12, fontFamily:"'DM Sans', sans-serif", padding:0 }}>
            Sign out
          </button>
        </> : <>
          <Link to="/login" style={{ color:theme.textMuted, textDecoration:'none', fontSize:13 }}>Sign in</Link>
          <Link to="/register" style={{ background:theme.accent, color:'#fff', textDecoration:'none', borderRadius:6, padding:'7px 16px', fontSize:13, fontWeight:600 }}>Get started →</Link>
        </>}
      </div>
    </nav>
  );
}
