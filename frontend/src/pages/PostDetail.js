import React, { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { getPost, deletePost, createComment, deleteComment, likePost, unlikePost, bookmarkPost, unbookmarkPost, followUser, unfollowUser } from '../api';
import { useAuth } from '../AuthContext';
import { useTheme } from '../ThemeContext';

const READING_TIME = (body) => Math.max(1, Math.ceil(body.split(' ').length / 200));

export default function PostDetail() {
  const { slug } = useParams();
  const { user } = useAuth();
  const { theme } = useTheme();
  const navigate = useNavigate();
  const [post, setPost] = useState(null);
  const [commentBody, setCommentBody] = useState('');
  const [loading, setLoading] = useState(true);
  const [visible, setVisible] = useState(false);
  const [liked, setLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [bookmarked, setBookmarked] = useState(false);
  const [following, setFollowing] = useState(false);

  useEffect(() => {
    getPost(slug).then(res => {
      setPost(res.data);
      setLiked(res.data.liked || false);
      setLikeCount(res.data.like_count || 0);
      setBookmarked(res.data.bookmarked || false);
      setTimeout(() => setVisible(true), 50);
    })

        .catch(() => navigate('/')).finally(() => setLoading(false));
  }, [slug, navigate]);

  const handleLike = async () => {
    if (!user) return;
    if (liked) { await unlikePost(slug); setLiked(false); setLikeCount(c => c - 1); }
    else        { await likePost(slug);  setLiked(true);  setLikeCount(c => c + 1); }
  };

  const handleBookmark = async () => {
    if (!user) return;
    if (bookmarked) { await unbookmarkPost(slug); setBookmarked(false); }
    else            { await bookmarkPost(slug);   setBookmarked(true);  }
  };

  const handleFollow = async () => {
    if (!user || !post) return;
    if (following) { await unfollowUser(post.author_id); setFollowing(false); }
    else           { await followUser(post.author_id);   setFollowing(true);  }
  };

  const handleDelete = async () => {
    if (!window.confirm('Delete this post?')) return;
    await deletePost(slug); navigate('/');
  };

  const handleComment = async () => {
    if (!commentBody.trim()) return;
    const res = await createComment(slug, { body: commentBody });
    setPost(p => ({ ...p, comments: [...(p.comments || []), res.data] }));
    setCommentBody('');
  };

  const handleDeleteComment = async (id) => {
    await deleteComment(slug, id);
    setPost(p => ({ ...p, comments: p.comments.filter(c => c.id !== id) }));
  };

  if (loading) return <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh', background: theme.bg }}><div style={{ width: 32, height: 32, border: `2px solid ${theme.border}`, borderTopColor: theme.accent, borderRadius: '50%', animation: 'spin 0.8s linear infinite' }} /></div>;
  if (!post) return null;
  const isAuthor = user?.id === post.author_id;

  return (
      <div style={{ paddingTop: 64, minHeight: '100vh', background: theme.bg, transition: 'background 0.3s ease' }}>
        <div style={{ maxWidth: 720, margin: '0 auto', padding: '64px 48px 40px', opacity: visible ? 1 : 0, transform: visible ? 'translateY(0)' : 'translateY(24px)', transition: 'all 0.6s ease' }}>
          <div style={{ marginBottom: 32 }}>
            <Link to="/" style={{ fontSize: 13, color: theme.textMuted, textDecoration: 'none' }}>← All stories</Link>
          </div>
          {post.tags && (
              <div style={{ display: 'flex', gap: 8, marginBottom: 24 }}>
                {post.tags.split(',').map(t => (
                    <span key={t} style={{ fontSize: 10, letterSpacing: 2, color: theme.accent, fontWeight: 600, textTransform: 'uppercase', background: theme.accentLight, borderRadius: 4, padding: '4px 10px' }}>{t.trim()}</span>
                ))}
              </div>
          )}
          <h1 style={{ fontFamily: "'Playfair Display', serif", fontSize: 'clamp(36px, 5vw, 56px)', fontWeight: 700, color: theme.text, lineHeight: 1.15, letterSpacing: -0.5, marginBottom: 36 }}>{post.title}</h1>
          <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
            <Link to={`/profile/${post.author?.username}`} style={{ display: 'flex', alignItems: 'center', gap: 14, textDecoration: 'none' }}>
              <div style={{ width: 44, height: 44, borderRadius: '50%', background: 'linear-gradient(135deg, #ff6b35, #ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 16, fontWeight: 700, color: '#fff' }}>{post.author?.username[0].toUpperCase()}</div>
              <div>
                <div style={{ color: theme.text, fontSize: 15, fontWeight: 600, marginBottom: 3 }}>{post.author?.username}</div>
                <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
                  <span style={{ color: theme.textMuted, fontSize: 13 }}>{new Date(post.created_at).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}</span>
                  <span style={{ color: theme.border }}>·</span>
                  <span style={{ color: theme.textMuted, fontSize: 13 }}>{READING_TIME(post.body)} min read</span>
                </div>
              </div>
            </Link>
            {isAuthor && (
                <div style={{ display: 'flex', gap: 10 }}>
                  <button onClick={() => navigate(`/edit/${slug}`)} style={{ background: theme.inputBg, border: `1px solid ${theme.inputBorder}`, color: theme.textSecondary, borderRadius: 6, padding: '7px 16px', cursor: 'pointer', fontSize: 13, fontFamily: "'DM Sans', sans-serif" }}>Edit</button>
                  <button onClick={handleDelete} style={{ background: theme.accentLight, border: `1px solid ${theme.accentBorder}`, color: theme.accent, borderRadius: 6, padding: '7px 16px', cursor: 'pointer', fontSize: 13, fontFamily: "'DM Sans', sans-serif" }}>Delete</button>
                </div>
            )}
          </div>

          {/* Like / Bookmark / Follow action bar */}
          <div style={{ display: 'flex', alignItems: 'center', gap: 16, marginTop: 24, paddingTop: 20, borderTop: `1px solid ${theme.border}` }}>
            {user && !isAuthor && (
                <button onClick={handleFollow} style={{
                  background: following ? theme.accentLight : 'none',
                  border: `1px solid ${following ? theme.accent : theme.borderStrong}`,
                  color: following ? theme.accent : theme.textSecondary,
                  borderRadius: 20, padding: '6px 18px', cursor: 'pointer',
                  fontSize: 13, fontWeight: 500, fontFamily: "'DM Sans', sans-serif",
                  transition: 'all 0.2s',
                }}>
                  {following ? '✓ Following' : '+ Follow'}
                </button>
            )}
            <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: 16 }}>
              <button onClick={handleLike} style={{
                background: 'none', border: 'none',
                cursor: user ? 'pointer' : 'default',
                display: 'flex', alignItems: 'center', gap: 6,
                color: liked ? theme.accent : theme.textMuted,
                fontSize: 14, fontFamily: "'DM Sans', sans-serif", padding: 0,
                transition: 'color 0.2s',
              }}>
                <span style={{ fontSize: 20 }}>{liked ? '♥' : '♡'}</span>
                <span>{likeCount}</span>
              </button>
              <button onClick={handleBookmark} style={{
                background: 'none', border: 'none',
                cursor: user ? 'pointer' : 'default',
                color: bookmarked ? theme.accent : theme.textMuted,
                fontSize: 18, fontFamily: "'DM Sans', sans-serif", padding: 0,
                transition: 'color 0.2s',
              }} title={bookmarked ? 'Remove bookmark' : 'Bookmark'}>
                {bookmarked ? '🔖' : '⬜'}
              </button>
            </div>
          </div>
        </div>

        <div style={{ maxWidth: 720, margin: '0 auto', padding: '0 48px 80px' }}>
          <div style={{ height: 1, background: `linear-gradient(to right, transparent, ${theme.borderStrong}, transparent)`, marginBottom: 48 }} />
          <div style={{ fontSize: 19, lineHeight: 1.85, color: theme.textBody, whiteSpace: 'pre-wrap', fontWeight: 300, letterSpacing: 0.2, opacity: visible ? 1 : 0, transition: 'opacity 0.8s ease 0.2s' }}>{post.body}</div>
        </div>

        <div style={{ borderTop: `1px solid ${theme.border}`, background: theme.isDark ? 'rgba(255,255,255,0.015)' : theme.bgSecondary }}>
          <div style={{ maxWidth: 720, margin: '0 auto', padding: '56px 48px 80px' }}>
            <h3 style={{ fontFamily: "'Playfair Display', serif", fontSize: 24, fontWeight: 700, color: theme.text, marginBottom: 36, display: 'flex', alignItems: 'center', gap: 12 }}>
              <span style={{ color: theme.accent, fontSize: 16 }}>✦</span>
              {(post.comments || []).length} {(post.comments || []).length === 1 ? 'Response' : 'Responses'}
            </h3>

            {user ? (
                <div style={{ background: theme.bgCard, border: `1px solid ${theme.border}`, borderRadius: 12, padding: 24, marginBottom: 40, boxShadow: theme.isDark ? 'none' : '0 2px 12px rgba(0,0,0,0.05)' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 10, marginBottom: 16 }}>
                    <div style={{ width: 28, height: 28, borderRadius: '50%', background: 'linear-gradient(135deg, #ff6b35, #ff9f6b)', display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 11, fontWeight: 700, color: '#fff' }}>{user.username[0].toUpperCase()}</div>
                    <span style={{ color: theme.textSecondary, fontSize: 13, fontWeight: 500 }}>{user.username}</span>
                  </div>
                  <textarea value={commentBody} onChange={e => setCommentBody(e.target.value)} placeholder="What are your thoughts?" style={{ width: '100%', background: 'transparent', border: 'none', outline: 'none', color: theme.text, fontSize: 15, lineHeight: 1.6, resize: 'none', fontFamily: "'DM Sans', sans-serif", fontWeight: 300 }} rows={4} />
                  <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 12, paddingTop: 12, borderTop: `1px solid ${theme.border}` }}>
                    <button onClick={handleComment} style={{ background: theme.accent, border: 'none', color: '#fff', borderRadius: 6, padding: '8px 20px', cursor: 'pointer', fontSize: 13, fontWeight: 600, fontFamily: "'DM Sans', sans-serif' " }} disabled={!commentBody.trim()}>Respond</button>
                  </div>
                </div>
            ) : (
                <div style={{ marginBottom: 40, paddingBottom: 20, borderBottom: `1px solid ${theme.border}` }}>
                  <Link to="/login" style={{ color: theme.accent, textDecoration: 'none', fontSize: 14, fontWeight: 600 }}>Sign in to respond →</Link>
                </div>
            )}

            <div style={{ display: 'flex', flexDirection: 'column', gap: 32 }}>
              {(post.comments || []).map(c => (
                  <div key={c.id} style={{ paddingBottom: 32, borderBottom: `1px solid ${theme.border}` }}>
                    <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
                      <div style={{ width: 32, height: 32, borderRadius: '50%', background: theme.inputBg, border: `1px solid ${theme.border}`, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 12, fontWeight: 600, color: theme.textSecondary }}>{c.author?.username[0].toUpperCase()}</div>
                      <div style={{ flex: 1, display: 'flex', alignItems: 'center', gap: 10 }}>
                        <Link to={`/profile/${c.author?.username}`} style={{ color: theme.textSecondary, fontSize: 13, fontWeight: 600, textDecoration: 'none' }}>{c.author?.username}</Link>
                        <span style={{ color: theme.textMuted, fontSize: 12 }}>{new Date(c.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}</span>
                      </div>
                      {user?.id === c.author_id && <button onClick={() => handleDeleteComment(c.id)} style={{ background: 'none', border: 'none', color: theme.textMuted, cursor: 'pointer', fontSize: 18, padding: '0 4px' }}>×</button>}
                    </div>
                    <p style={{ color: theme.textSecondary, fontSize: 15, lineHeight: 1.7, margin: 0, fontWeight: 300 }}>{c.body}</p>
                  </div>
              ))}
            </div>
          </div>
        </div>
      </div>
  );
}