import { auth } from './firebase';
import { PUBLIC_API_URL } from '$env/static/public';

/**
 * @param {string} path
 * @param {RequestInit} [options]
 */
async function authedFetch(path, options = {}) {
	const ready = auth.authStateReady();
	const timeout = new Promise((_, reject) =>
		setTimeout(() => reject(new Error('authStateReady timed out')), 5000)
	);
	await Promise.race([ready, timeout]);

	const user = auth.currentUser;
	if (!user) throw new Error('Not authenticated');

	const token = await user.getIdToken();

	const res = await fetch(`${PUBLIC_API_URL}${path}`, {
		...options,
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${token}`,
			...(options.headers || {})
		}
	});

if (!res.ok) {
	const text = await res.text();
	throw new Error(text || `Request failed with status ${res.status}`);
}

if (res.status === 204) return null;

const data = await res.json();
return data;
}

/**
 * @param {{ title: string, body: string, status: string }} post
 */
export function createPost({ title, body, status }) {
	return authedFetch('/api/posts', {
		method: 'POST',
		body: JSON.stringify({ title, body, status })
	});
}

/**
 * @param {string} id
 * @param {{ title: string, body: string, status: string }} post
 */
export function updatePost(id, { title, body, status }) {
	return authedFetch(`/api/posts/${id}`, {
		method: 'PUT',
		body: JSON.stringify({ title, body, status })
	});
}

export function getMe() {
	return authedFetch('/api/me', { method: 'GET' });
}

export function getPosts() {
	return authedFetch('/api/posts', { method: 'GET' });
}

export function getMyPosts() {
	return authedFetch('/api/me/posts', { method: 'GET' });
}

/**
 * @param {string} id
 */
export function getPost(id) {
	return authedFetch(`/api/posts/${id}`, { method: 'GET' });
}
/**
 * @param {string} id
 */
export function deleteMyPost(id) {
	return authedFetch(`/api/posts/${id}`, { method: 'DELETE' });
}