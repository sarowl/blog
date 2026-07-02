<script>
	import { onMount } from 'svelte';
	import Navbar from '$lib/components/Navbar.svelte';
	import { getPosts } from '$lib/api';
    import {resolve} from '$app/paths';

	/**
	 * @typedef {Object} Post
	 * @property {string} id
	 * @property {string} user_id
 	 * @property {string} author_name
 	 * @property {string} title
 	 * @property {string} body
 	 * @property {string} created_at
	 */

	/** @type {Post[]} */
	let posts = $state([]);
    let loading = $state(true);
    let error = $state('');

	onMount(async () => {
	try {
		const data = await getPosts();
		posts = Array.isArray(data) ? data : [];
	} catch (err) {
		console.log('caught error', err);
	} finally {
		loading = false;
	}
});

	/**
	 * Formats an ISO date string into a readable date.
	 * @param {string} isoDate
	 * @returns {string}
	 */
	function formatDate(isoDate) {
		return new Date(isoDate).toLocaleDateString('en-US', {
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}
</script>

<svelte:head>
	<title>Feed</title>
</svelte:head>

<Navbar />

<main class="feed">
	<h1 class="feed__heading">Feed</h1>

	{#if loading}
		<p class="feed__empty">Loading posts...</p>
	{:else if error}
		<p class="feed__empty">{error}</p>
	{:else if posts.length === 0}
		<p class="feed__empty">No public posts yet.</p>
	{:else}
		<ul class="feed__list">
			{#each posts as post (post.id)}
				<li class="card">
					<div class="card__meta">
						<span class="card__author">{post.author_name || post.user_id}</span>
						<span class="card__dot" aria-hidden="true">&bull;</span>
						<time class="card__date" datetime={post.created_at}>
							{formatDate(post.created_at)}
						</time>
					</div>
					<h2 class="card__title">
                        <a href={resolve(`/posts/${post.id}`)}>{post.title}</a>
                    </h2>
                    
					{#if post.body.length > 150}
                        <p class="card__content">{post.body.substring(0, 150)}...</p>
                        <a href={resolve(`/posts/${post.id}`)} class="card__read-more">Read more...</a>
                    {:else}
                        <p class="card__content">{post.body}</p>
                    {/if}
				</li>
			{/each}
		</ul>
	{/if}
</main>

<style>
	.feed {
		max-width: 640px;
		margin: 0 auto;
		padding: 2rem 1rem;
	}

	.feed__heading {
		font-size: 1.75rem;
		font-weight: 700;
		margin-bottom: 1.5rem;
	}

	.feed__empty {
		color: #6b7280;
		text-align: center;
		padding: 2rem 0;
	}

	.feed__list {
		list-style: none;
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.card {
		border: 1px solid #e5e7eb;
		border-radius: 0.75rem;
		padding: 1.25rem 1.5rem;
		background-color: #ffffff;
	}

	.card__meta {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		font-size: 0.875rem;
		color: #6b7280;
		margin-bottom: 0.5rem;
	}

	.card__author {
		font-weight: 600;
		color: #111827;
	}

	.card__dot {
		color: #d1d5db;
	}

	.card__title {
		font-size: 1.25rem;
		font-weight: 700;
		margin: 0 0 0.5rem;
		color: #111827;
	}

	.card__content {
		font-size: 0.95rem;
		line-height: 1.5;
		color: #374151;
		margin: 0;
	}
</style>