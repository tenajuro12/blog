import axios from 'axios';

const api = axios.create({ baseURL: '/api' });

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

// Auth
export const register    = (data) => api.post('/auth/register', data);
export const login       = (data) => api.post('/auth/login', data);
export const getMe       = ()     => api.get('/auth/me');

// Posts
export const getPosts    = (tag)  => api.get('/posts', { params: tag ? { tag } : {} });
export const getPost     = (slug) => api.get(`/posts/${slug}`);
export const createPost  = (data) => api.post('/posts', data);
export const updatePost  = (slug, data) => api.put(`/posts/${slug}`, data);
export const deletePost  = (slug) => api.delete(`/posts/${slug}`);

// Comments
export const getComments    = (slug)     => api.get(`/posts/${slug}/comments`);
export const createComment  = (slug, d)  => api.post(`/posts/${slug}/comments`, d);
export const deleteComment  = (slug, id) => api.delete(`/posts/${slug}/comments/${id}`);

// Likes
export const likePost   = (slug) => api.post(`/posts/${slug}/like`);
export const unlikePost = (slug) => api.delete(`/posts/${slug}/like`);
export const getLikes   = (slug) => api.get(`/posts/${slug}/likes`);

// Bookmarks
export const bookmarkPost   = (slug) => api.post(`/posts/${slug}/bookmark`);
export const unbookmarkPost = (slug) => api.delete(`/posts/${slug}/bookmark`);
export const getBookmarks   = ()     => api.get('/bookmarks');

// Follows
export const followUser   = (id) => api.post(`/users/${id}/follow`);
export const unfollowUser = (id) => api.delete(`/users/${id}/follow`);
export const getFollowers = (id) => api.get(`/users/${id}/followers`);
export const getFollowing = (id) => api.get(`/users/${id}/following`);
export const getFeed      = ()   => api.get('/feed');

// Tags
export const getTags        = ()    => api.get('/tags');
export const getPostsByTag  = (tag) => api.get(`/tags/${tag}`);

// Profile
export const getProfile   = (username) => api.get(`/profiles/${username}`);
export const updateProfile = (data)    => api.put('/profile', data);
export const getUserPosts  = (id)      => api.get(`/profiles/${id}/posts`);
