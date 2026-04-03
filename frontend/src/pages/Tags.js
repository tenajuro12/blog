import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { getTags, getPostsByTag } from '../api';
import { useTheme } from '../ThemeContext';

const READING_TIME = (body) => Math.max(1, Math.ceil(body.split(' ').length / 200));

export function TagsPage() {
  const { theme } = useTheme();
  const [tags, setTags] = useState([]);

  useEffect(() => {
    getTags().then(r => setTags((r.data || []).sort((a, b) => b.count - a.count)));
  }, []);

  return (
    <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg }}>
      <div style={{ maxWidth: 900, margin: '0 auto', padding: '56px 48px 80px' }}>
        <div style={{ marginBottom: 48 }}>
          <div style={{ fontSize: 11, letterSpacing: 4, color: theme.accent, marginBottom: 16, fontWeight: 600 }}>EXPLORE</div>
          <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 40, fontWeight: 700, color: theme.text, letterSpacing: -0.5 }}>All Tags</h1>
        </div>
        {tags.length === 0 && <p style={{ color: theme.textMuted }}>No tags yet.</p>}
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 12 }}>
          {tags.map(tag => (
            <Link key={tag.name} to={`/tags/${tag.name}`} style={{
              textDecoration: 'none',
              background: theme.bgCard,
              border: `1px solid ${theme.border}`,
              borderRadius: 10,
              padding: '14px 22px',
              display: 'flex', flexDirection: 'column', gap: 4,
              boxShadow: theme.isDark ? 'none' : '0 2px 8px rgba(0,0,0,0.05)',
              transition: 'all 0.2s',
            }}>
              <span style={{ fontSize: 15, fontWeight: 600, color: theme.text }}>#{tag.name}</span>
              <span style={{ fontSize: 12, color: theme.textMuted }}>{tag.count} {tag.count === 1 ? 'story' : 'stories'}</span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}

export function TagPage() {
  const { name } = useParams();
  const { theme } = useTheme();
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getPostsByTag(name).then(r => setPosts(r.data || [])).finally(() => setLoading(false));
  }, [name]);

  if (loading) return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: theme.bg }}>
      <div style={{ width: 32, height: 32, border: `2px solid ${theme.border}`, borderTopColor: theme.accent, borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} />
    </div>
  );

  return (
    <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg }}>
      <div style={{ maxWidth: 900, margin: '0 auto', padding: '56px 48px 80px' }}>
        <div style={{ marginBottom: 40 }}>
          <Link to="/tags" style={{ fontSize: 13, color: theme.textMuted, textDecoration: 'none' }}>← All tags</Link>
          <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 40, fontWeight: 700, color: theme.text, marginTop: 20, letterSpacing: -0.5 }}>
            #{name}
          </h1>
          <p style={{ color: theme.textMuted, fontSize: 15, marginTop: 8 }}>{posts.length} {posts.length === 1 ? 'story' : 'stories'}</p>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column' }}>
          {posts.map(post => (
            <Link key={post.id} to={`/posts/${post.slug}`} style={{
              display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start',
              padding: '28px 0', borderBottom: `1px solid ${theme.border}`,
              textDecoration: 'none', gap: 24,
            }}>
              <div style={{ flex: 1 }}>
                <h3 style={{ fontFamily: "'Playfair Display', serif", fontSize: 22, fontWeight: 700, color: theme.text, lineHeight: 1.3, marginBottom: 8 }}>{post.title}</h3>
                <p style={{ color: theme.textSecondary, fontSize: 14, lineHeight: 1.6 }}>{post.body.substring(0, 120)}...</p>
                <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginTop: 12 }}>
                  <div style={{ width: 22, height: 22, borderRadius: '50%', background: 'linear-gradient(135deg,#ff6b35,#ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 9, fontWeight: 700, color: '#fff' }}>
                    {post.author?.username[0].toUpperCase()}
                  </div>
                  <span style={{ color: theme.textMuted, fontSize: 12 }}>{post.author?.username}</span>
                  <span style={{ color: theme.textMuted, fontSize: 12 }}>·</span>
                  <span style={{ color: theme.textMuted, fontSize: 12 }}>{READING_TIME(post.body)} min read</span>
                </div>
              </div>
              <span style={{ color: theme.accent, opacity: 0.5, fontSize: 18, paddingTop: 4 }}>→</span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}
