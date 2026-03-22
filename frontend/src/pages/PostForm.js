import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { createPost, updatePost, getPost } from '../api';
import { useTheme } from '../ThemeContext';

export function NewPost() {
  const [form, setForm] = useState({ title: '', body: '', tags: '' });
  const [error, setError] = useState(''); const [saving, setSaving] = useState(false);
  const navigate = useNavigate();
  const handleSubmit = async () => {
    if (!form.title.trim() || !form.body.trim()) return; setSaving(true);
    try { const res = await createPost(form); navigate(`/posts/${res.data.slug}`); }
    catch (e) { setError(e.response?.data?.error || 'Failed to publish'); setSaving(false); }
  };
  return <PostEditor title="New story" form={form} setForm={setForm} error={error} onSubmit={handleSubmit} saving={saving} />;
}

export function EditPost() {
  const { slug } = useParams();
  const [form, setForm] = useState({ title: '', body: '', tags: '' });
  const [error, setError] = useState(''); const [saving, setSaving] = useState(false);
  const navigate = useNavigate();
  useEffect(() => { getPost(slug).then(res => setForm({ title: res.data.title, body: res.data.body, tags: res.data.tags || '' })); }, [slug]);
  const handleSubmit = async () => {
    setSaving(true);
    try { const res = await updatePost(slug, form); navigate(`/posts/${res.data.slug}`); }
    catch (e) { setError(e.response?.data?.error || 'Failed to update'); setSaving(false); }
  };
  return <PostEditor title="Edit story" form={form} setForm={setForm} error={error} onSubmit={handleSubmit} saving={saving} />;
}

function PostEditor({ title, form, setForm, error, onSubmit, saving }) {
  const { theme } = useTheme();
  const wordCount = form.body.trim().split(/\s+/).filter(Boolean).length;
  const readTime = Math.max(1, Math.ceil(wordCount / 200));
  return (
    <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg, display: 'flex', flexDirection: 'column', transition: 'background 0.3s ease' }}>
      <div style={{ position: 'sticky', top: 64, zIndex: 50, display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '0 48px', height: 52, background: theme.isDark ? 'rgba(8,8,12,0.9)' : `${theme.bg}ee`, backdropFilter: 'blur(20px)', borderBottom: `1px solid ${theme.border}` }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 20 }}>
          <span style={{ fontSize: 13, color: theme.textMuted, letterSpacing: 0.5 }}>{title}</span>
          {wordCount > 0 && <span style={{ fontSize: 12, color: theme.accent, opacity: 0.6 }}>{wordCount} words · {readTime} min read</span>}
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
          {error && <span style={{ fontSize: 12, color: theme.accent }}>{error}</span>}
          <button onClick={onSubmit} disabled={saving || !form.title.trim() || !form.body.trim()} style={{ background: theme.accent, border: 'none', color: '#fff', borderRadius: 6, padding: '7px 18px', fontSize: 13, fontWeight: 600, cursor: 'pointer', fontFamily: "'DM Sans', sans-serif", opacity: (saving || !form.title.trim() || !form.body.trim()) ? 0.5 : 1 }}>
            {saving ? 'Publishing...' : 'Publish →'}
          </button>
        </div>
      </div>
      <div style={{ flex: 1, maxWidth: 760, width: '100%', margin: '0 auto', padding: '56px 48px 80px' }}>
        <input style={{ width: '100%', background: 'transparent', border: 'none', outline: 'none', fontFamily: "'Playfair Display', serif", fontSize: 'clamp(32px, 4vw, 52px)', fontWeight: 700, color: theme.text, lineHeight: 1.2, letterSpacing: -0.5, marginBottom: 28, caretColor: theme.accent }} placeholder="Your story title..." value={form.title} onChange={e => setForm({ ...form, title: e.target.value })} />
        <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 28 }}>
          <span style={{ fontSize: 14, opacity: 0.5 }}>🏷</span>
          <input style={{ background: 'transparent', border: 'none', outline: 'none', color: theme.accent, fontSize: 13, fontFamily: "'DM Mono', monospace", letterSpacing: 0.5, width: '100%', opacity: 0.7 }} placeholder="Add tags (comma separated)..." value={form.tags} onChange={e => setForm({ ...form, tags: e.target.value })} />
        </div>
        <div style={{ height: 1, background: theme.border, marginBottom: 40 }} />
        <textarea style={{ width: '100%', background: 'transparent', border: 'none', outline: 'none', fontSize: 19, lineHeight: 1.85, color: theme.textBody, fontFamily: "'DM Sans', sans-serif", fontWeight: 300, resize: 'none', minHeight: 500, letterSpacing: 0.2, caretColor: theme.accent }} placeholder="Tell your story..." value={form.body} onChange={e => setForm({ ...form, body: e.target.value })} />
      </div>
    </div>
  );
}
