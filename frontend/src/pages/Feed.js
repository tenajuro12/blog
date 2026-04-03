import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getFeed } from '../api';
import { useTheme } from '../ThemeContext';

const READING_TIME = (body) => Math.max(1, Math.ceil(body.split(' ').length / 200));

export default function Feed() {
  const { theme } = useTheme();
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getFeed().then(r => setPosts(r.data || [])).finally(() => setLoading(false));
  }, []);

  if (loading) return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: theme.bg }}>
      <div style={{ width: 32, height: 32, border: `2px solid ${theme.border}`, borderTopColor: theme.accent, borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
    </div>
  );

  return (
    <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg }}>
      <div style={{ maxWidth: 800, margin: '0 auto', padding: '56px 48px 80px' }}>
        <div style={{ marginBottom: 48 }}>
          <div style={{ fontSize: 11, letterSpacing: 4, color: theme.accent, marginBottom: 16, fontWeight: 600 }}>FOLLOWING</div>
          <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 40, fontWeight: 700, color: theme.text, letterSpacing: -0.5 }}>Your Feed</h1>
          <p style={{ color: theme.textMuted, marginTop: 8, fontSize: 15 }}>Stories from people you follow</p>
        </div>

        {posts.length === 0 && (
          <div style={{ textAlign: 'center', padding: '60px 0' }}>
            <div style={{ fontSize: 36, marginBottom: 16, color: theme.accent, opacity: 0.3 }}>✦</div>
            <p style={{ color: theme.textMuted, fontSize: 15, marginBottom: 20 }}>Your feed is empty. Follow some authors to see their stories here.</p>
            <Link to="/" style={{ color: theme.accent, textDecoration: 'none', fontSize: 14, fontWeight: 600 }}>Discover writers →</Link>
          </div>
        )}

        {posts.map(post => (
          <Link key={post.id} to={`/posts/${post.slug}`} style={{
            display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start',
            padding: '28px 0', borderBottom: `1px solid ${theme.border}`,
            textDecoration: 'none', gap: 24,
          }}>
            <div style={{ flex: 1 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 12 }}>
                <div style={{ width: 28, height: 28, borderRadius: '50%', background: 'linear-gradient(135deg,#ff6b35,#ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 11, fontWeight: 700, color: '#fff' }}>
                  {post.author?.username[0].toUpperCase()}
                </div>
                <span style={{ color: theme.textSecondary, fontSize: 13, fontWeight: 500 }}>{post.author?.username}</span>
                <span style={{ color: theme.textMuted, fontSize: 12 }}>· {new Date(post.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}</span>
              </div>
              {post.tags && <span style={{ fontSize: 10, letterSpacing: 2, color: theme.accent, fontWeight: 600, textTransform: 'uppercase', display: 'block', marginBottom: 8 }}>{post.tags.split(',')[0].trim()}</span>}
              <h3 style={{ fontFamily: "'Playfair Display', serif", fontSize: 22, fontWeight: 700, color: theme.text, lineHeight: 1.3, marginBottom: 8 }}>{post.title}</h3>
              <p style={{ color: theme.textSecondary, fontSize: 14, lineHeight: 1.6 }}>{post.body.substring(0, 120)}...</p>
              <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginTop: 12 }}>
                <span style={{ fontSize: 12, color: theme.textMuted }}>{READING_TIME(post.body)} min read</span>
                {post.like_count > 0 && <span style={{ fontSize: 12, color: theme.textMuted }}>♥ {post.like_count}</span>}
              </div>
            </div>
            <span style={{ color: theme.accent, opacity: 0.5, fontSize: 18, paddingTop: 4 }}>→</span>
          </Link>
        ))}
      </div>
    </div>
  );
}
