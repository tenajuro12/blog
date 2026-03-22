import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getPosts } from '../api';
import { useTheme } from '../ThemeContext';

const READING_TIME = (body) => Math.max(1, Math.ceil(body.split(' ').length / 200));

export default function Home() {
  const { theme } = useTheme();
  const [posts, setPosts] = useState([]);
  const [tag, setTag] = useState('');
  const [loading, setLoading] = useState(true);
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    getPosts(tag)
      .then((res) => setPosts(res.data || []))
      .finally(() => { setLoading(false); setTimeout(() => setVisible(true), 50); });
  }, [tag]);

  const featured = posts[0];
  const rest = posts.slice(1);
  const allTags = [...new Set(posts.flatMap(p => (p.tags || '').split(',').map(t => t.trim()).filter(Boolean)))];

  if (loading) return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: theme.bg }}>
      <div style={{ width: 32, height: 32, border: `2px solid ${theme.border}`, borderTopColor: theme.accent, borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
    </div>
  );

  return (
    <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg, transition: 'background 0.3s ease' }}>
      {/* Hero */}
      <div style={{ position: 'relative', overflow: 'hidden', padding: '100px 48px 80px', borderBottom: `1px solid ${theme.border}` }}>
        <div style={{ maxWidth: 700, position: 'relative', zIndex: 2, animation: 'fadeUp 0.8s ease both' }}>
          <div style={{ fontSize: 11, letterSpacing: 4, color: theme.accent, marginBottom: 24, fontWeight: 500, opacity: 0.8 }}>EST. 2024 · THE WRITERS' PLATFORM</div>
          <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 'clamp(56px, 8vw, 96px)', fontWeight: 900, lineHeight: 1.05, color: theme.text, marginBottom: 24, letterSpacing: -2 }}>
            Ideas worth<br />
            <em style={{ fontStyle: 'italic', color: theme.accent }}>reading.</em>
          </h1>
          <p style={{ fontSize: 18, color: theme.textSecondary, lineHeight: 1.6, maxWidth: 480, fontWeight: 300 }}>
            Discover stories, thinking, and expertise from writers on any topic.
          </p>
        </div>
        <div style={{ position: 'absolute', right: 80, top: '50%', transform: 'translateY(-50%)' }}>
          <div style={{ width: 1, height: 200, background: `linear-gradient(to bottom, transparent, ${theme.accent}40, transparent)`, margin: '0 auto 20px' }} />
          <div style={{ width: 120, height: 120, borderRadius: '50%', border: `1px solid ${theme.accent}20`, margin: '0 auto' }} />
        </div>
      </div>

      <div style={{ maxWidth: 1200, margin: '0 auto', padding: '48px 48px 80px' }}>
        {/* Tags */}
        {allTags.length > 0 && (
          <div style={{ display: 'flex', gap: 8, flexWrap: 'wrap', marginBottom: 48 }}>
            <button onClick={() => setTag('')} style={{ background: tag === '' ? theme.accent : theme.inputBg, border: `1px solid ${tag === '' ? theme.accent : theme.inputBorder}`, color: tag === '' ? '#fff' : theme.textMuted, borderRadius: 100, padding: '6px 16px', fontSize: 12, letterSpacing: 0.5, cursor: 'pointer', fontFamily: "'DM Sans', sans-serif", transition: 'all 0.2s' }}>
              All
            </button>
            {allTags.map(t => (
              <button key={t} onClick={() => setTag(t === tag ? '' : t)}
                style={{ background: tag === t ? theme.accent : theme.inputBg, border: `1px solid ${tag === t ? theme.accent : theme.inputBorder}`, color: tag === t ? '#fff' : theme.textMuted, borderRadius: 100, padding: '6px 16px', fontSize: 12, letterSpacing: 0.5, cursor: 'pointer', fontFamily: "'DM Sans', sans-serif", transition: 'all 0.2s' }}>
                {t}
              </button>
            ))}
          </div>
        )}

        {posts.length === 0 && (
          <div style={{ textAlign: 'center', padding: '80px 0' }}>
            <div style={{ fontSize: 40, marginBottom: 16, color: theme.accent, opacity: 0.4 }}>✦</div>
            <p style={{ color: theme.textMuted, fontSize: 16, marginBottom: 24 }}>No stories yet. Be the first to write.</p>
            <Link to="/register" style={{ color: theme.accent, textDecoration: 'none', fontSize: 14, fontWeight: 600 }}>Start writing →</Link>
          </div>
        )}

        {/* Featured */}
        {featured && (
          <Link to={`/posts/${featured.slug}`} style={{
            display: 'flex', position: 'relative', overflow: 'hidden',
            background: theme.isDark ? `linear-gradient(135deg, ${theme.accentLight} 0%, rgba(255,107,53,0.03) 100%)` : theme.bgCard,
            border: `1px solid ${theme.isDark ? 'rgba(255,107,53,0.15)' : theme.border}`,
            borderRadius: 16, padding: '48px 56px', marginBottom: 56,
            textDecoration: 'none', cursor: 'pointer',
            opacity: visible ? 1 : 0, transform: visible ? 'translateY(0)' : 'translateY(32px)',
            transition: 'all 0.6s ease',
            boxShadow: theme.isDark ? 'none' : '0 4px 24px rgba(0,0,0,0.06)',
          }}>
            <div style={{ flex: 1 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginBottom: 20 }}>
                <span style={{ fontSize: 10, letterSpacing: 3, color: theme.accent, fontWeight: 600 }}>FEATURED</span>
                <span style={{ fontSize: 12, color: theme.textMuted, letterSpacing: 0.5 }}>{READING_TIME(featured.body)} min read</span>
              </div>
              <h2 style={{ fontFamily: "'Playfair Display', serif", fontSize: 'clamp(28px, 4vw, 48px)', fontWeight: 700, color: theme.text, lineHeight: 1.2, marginBottom: 20, letterSpacing: -0.5 }}>
                {featured.title}
              </h2>
              <p style={{ color: theme.textSecondary, fontSize: 16, lineHeight: 1.7, marginBottom: 32, maxWidth: 600, fontWeight: 300 }}>
                {featured.body.substring(0, 200)}...
              </p>
              <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                  <div style={{ width: 28, height: 28, borderRadius: '50%', background: 'linear-gradient(135deg, #ff6b35, #ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 11, fontWeight: 700, color: '#fff' }}>
                    {featured.author?.username[0].toUpperCase()}
                  </div>
                  <span style={{ color: theme.textSecondary, fontSize: 13, fontWeight: 500 }}>{featured.author?.username}</span>
                  <span style={{ color: theme.border, fontSize: 16 }}>·</span>
                  <span style={{ color: theme.textMuted, fontSize: 13 }}>{new Date(featured.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</span>
                </div>
                <span style={{ color: theme.accent, fontSize: 14, fontWeight: 600 }}>Read story →</span>
              </div>
            </div>
          </Link>
        )}

        {/* Grid */}
        {rest.length > 0 && (
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))', gap: 24 }}>
            {rest.map((post, i) => (
              <Link key={post.id} to={`/posts/${post.slug}`}
                style={{
                  background: theme.bgCard, border: `1px solid ${theme.border}`,
                  borderRadius: 12, padding: 28, textDecoration: 'none',
                  display: 'block', boxShadow: theme.isDark ? 'none' : '0 2px 12px rgba(0,0,0,0.05)',
                  opacity: visible ? 1 : 0, transform: visible ? 'translateY(0)' : 'translateY(32px)',
                  transition: `all 0.6s ease ${0.1 + i * 0.07}s`,
                }}>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                  {post.tags && <span style={{ fontSize: 10, letterSpacing: 2, color: theme.accent, fontWeight: 600, textTransform: 'uppercase' }}>{post.tags.split(',')[0].trim()}</span>}
                  <span style={{ fontSize: 11, color: theme.textMuted }}>{READING_TIME(post.body)} min</span>
                </div>
                <h3 style={{ fontFamily: "'Playfair Display', serif", fontSize: 22, fontWeight: 700, color: theme.text, lineHeight: 1.3, marginBottom: 12, letterSpacing: -0.3 }}>
                  {post.title}
                </h3>
                <p style={{ color: theme.textSecondary, fontSize: 14, lineHeight: 1.6, marginBottom: 20, fontWeight: 300 }}>
                  {post.body.substring(0, 100)}...
                </p>
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, borderTop: `1px solid ${theme.border}`, paddingTop: 16 }}>
                  <div style={{ width: 22, height: 22, borderRadius: '50%', background: 'linear-gradient(135deg, #ff6b35, #ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 9, fontWeight: 700, color: '#fff' }}>
                    {post.author?.username[0].toUpperCase()}
                  </div>
                  <span style={{ color: theme.textSecondary, fontSize: 12, flex: 1 }}>{post.author?.username}</span>
                  <span style={{ color: theme.textMuted, fontSize: 11 }}>{new Date(post.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}</span>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
