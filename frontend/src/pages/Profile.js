import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getProfile, getUserPosts, updateProfile } from '../api';
import { useAuth } from '../AuthContext';
import { useTheme } from '../ThemeContext';

export default function Profile() {
  const { username } = useParams();
  const { user, setUser } = useAuth();
  const { theme } = useTheme();
  const [profile, setProfile] = useState(null);
  const [posts, setPosts] = useState([]);
  const [editing, setEditing] = useState(false);
  const [form, setForm] = useState({ username: '', bio: '', avatar: '' });
  const [error, setError] = useState('');
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    setVisible(false);
    getProfile(username).then(res => {
      setProfile(res.data);
      setForm({ username: res.data.username, bio: res.data.bio || '', avatar: res.data.avatar || '' });
      getUserPosts(res.data.id).then(r => { setPosts(r.data); setTimeout(() => setVisible(true), 50); });
    });
  }, [username]);

  const handleSave = async () => {
    try { const res = await updateProfile(form); setProfile(res.data); setUser(res.data); setEditing(false); setError(''); }
    catch (e) { setError(e.response?.data?.error || 'Failed to update'); }
  };

  if (!profile) return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: theme.bg }}><div style={{ width: 32, height: 32, border: `2px solid ${theme.border}`, borderTopColor: theme.accent, borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} /></div>;
  const isOwner = user?.id === profile.id;

  return (
    <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg, transition: 'background 0.3s ease' }}>
      <div style={{ borderBottom: `1px solid ${theme.border}`, padding: '64px 0 48px', opacity: visible ? 1 : 0, transform: visible ? 'translateY(0)' : 'translateY(20px)', transition: 'all 0.6s ease' }}>
        <div style={{ maxWidth: 800, margin: '0 auto', padding: '0 48px', display: 'flex', gap: 40, alignItems: 'flex-start' }}>
          <div style={{ flexShrink: 0 }}>
            {profile.avatar
              ? <img src={profile.avatar} alt="" style={{ width: 88, height: 88, borderRadius: '50%', objectFit: 'cover', border: `2px solid ${theme.accentBorder}` }} />
              : <div style={{ width: 88, height: 88, borderRadius: '50%', background: 'linear-gradient(135deg, #ff6b35, #ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 36, fontWeight: 700, color: '#fff' }}>{profile.username[0].toUpperCase()}</div>
            }
          </div>
          <div style={{ flex: 1, paddingTop: 4 }}>
            {editing ? (
              <div style={{ display: 'flex', flexDirection: 'column', gap: 12, maxWidth: 440 }}>
                {error && <div style={{ background: theme.accentLight, border: `1px solid ${theme.accentBorder}`, color: theme.accent, borderRadius: 6, padding: '10px 14px', fontSize: 13 }}>{error}</div>}
                {[['Username', 'username', 'text', 'username'], ['Avatar URL', 'avatar', 'text', 'https://...']].map(([label, key, type, ph]) => (
                  <input key={key} style={{ background: theme.inputBg, border: `1px solid ${theme.inputBorder}`, borderRadius: 7, padding: '10px 14px', color: theme.text, fontSize: 14, fontFamily: "'DM Sans', sans-serif", outline: 'none' }}
                    type={type} placeholder={ph} value={form[key]} onChange={e => setForm({ ...form, [key]: e.target.value })} />
                ))}
                <textarea style={{ background: theme.inputBg, border: `1px solid ${theme.inputBorder}`, borderRadius: 7, padding: '10px 14px', color: theme.text, fontSize: 14, fontFamily: "'DM Sans', sans-serif", outline: 'none', resize: 'vertical' }}
                  placeholder="Short bio..." value={form.bio} rows={2} onChange={e => setForm({ ...form, bio: e.target.value })} />
                <div style={{ display: 'flex', gap: 10 }}>
                  <button onClick={handleSave} style={{ background: theme.accent, border: 'none', color: '#fff', borderRadius: 6, padding: '8px 20px', cursor: 'pointer', fontSize: 13, fontWeight: 600, fontFamily: "'DM Sans', sans-serif" }}>Save changes</button>
                  <button onClick={() => { setEditing(false); setError(''); }} style={{ background: 'none', border: `1px solid ${theme.borderStrong}`, color: theme.textSecondary, borderRadius: 6, padding: '8px 16px', cursor: 'pointer', fontSize: 13, fontFamily: "'DM Sans', sans-serif" }}>Cancel</button>
                </div>
              </div>
            ) : (
              <>
                <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 36, fontWeight: 700, color: theme.text, marginBottom: 10, letterSpacing: -0.5 }}>{profile.username}</h1>
                {profile.bio && <p style={{ color: theme.textSecondary, fontSize: 16, lineHeight: 1.6, marginBottom: 16, fontWeight: 300, maxWidth: 480 }}>{profile.bio}</p>}
                <div style={{ display: 'flex', alignItems: 'center', gap: 24 }}>
                  <span style={{ fontSize: 13, color: theme.textMuted }}>{posts.length} {posts.length === 1 ? 'story' : 'stories'}</span>
                  {isOwner && <button onClick={() => setEditing(true)} style={{ background: 'none', border: `1px solid ${theme.borderStrong}`, color: theme.textSecondary, borderRadius: 6, padding: '6px 16px', cursor: 'pointer', fontSize: 12, letterSpacing: 0.5, fontFamily: "'DM Sans', sans-serif" }}>Edit profile</button>}
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      <div style={{ maxWidth: 800, margin: '0 auto', padding: '48px 48px 80px' }}>
        <div style={{ fontSize: 10, letterSpacing: 4, color: theme.textMuted, marginBottom: 28, fontWeight: 600 }}>STORIES</div>
        {posts.length === 0 && (
          <div style={{ padding: '40px 0' }}>
            <p style={{ color: theme.textMuted, fontSize: 15, marginBottom: 16 }}>No stories written yet.</p>
            {isOwner && <Link to="/new-post" style={{ color: theme.accent, textDecoration: 'none', fontSize: 14, fontWeight: 600 }}>Write your first story →</Link>}
          </div>
        )}
        {posts.map((post, i) => (
          <Link key={post.id} to={`/posts/${post.slug}`}
            style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', padding: '28px 0', borderBottom: `1px solid ${theme.border}`, textDecoration: 'none', gap: 24, opacity: visible ? 1 : 0, transform: visible ? 'translateY(0)' : 'translateY(20px)', transition: `all 0.5s ease ${0.1 + i * 0.06}s` }}>
            <div style={{ flex: 1 }}>
              {post.tags && <span style={{ fontSize: 10, letterSpacing: 2, color: theme.accent, fontWeight: 600, textTransform: 'uppercase', display: 'block', marginBottom: 10 }}>{post.tags.split(',')[0].trim()}</span>}
              <h3 style={{ fontFamily: "'Playfair Display', serif", fontSize: 22, fontWeight: 700, color: theme.text, lineHeight: 1.3, marginBottom: 8, letterSpacing: -0.3 }}>{post.title}</h3>
              <p style={{ color: theme.textSecondary, fontSize: 14, lineHeight: 1.6, fontWeight: 300 }}>{post.body.substring(0, 100)}...</p>
            </div>
            <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-end', gap: 12, flexShrink: 0, paddingTop: 2 }}>
              <span style={{ color: theme.textMuted, fontSize: 12, whiteSpace: 'nowrap' }}>{new Date(post.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}</span>
              <span style={{ color: theme.accent, opacity: 0.5, fontSize: 18 }}>→</span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}
