<script>
    import { onMount } from 'svelte';
    import { page } from '$app/stores';
    import { auth } from '$lib/firebase';
    import { onAuthStateChanged } from 'firebase/auth';
    import { getPost } from '$lib/api';
	import Navbar from '$lib/components/Navbar.svelte';

    /** @type {{ id: string, title: string, body: string, created_at: string, status: string } | null} */
    let post = $state(null);
    let isLoading = $state(true);
    let loadError = $state('');

    /** @param {string | null | undefined} value */
    const formatDate = (value) => {
        if (!value) return '';
        return new Date(value).toLocaleDateString('en', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    onMount(() => {
        const unsubscribe = onAuthStateChanged(auth, async (firebaseUser) => {
            if (!firebaseUser) {
                return;
            }

            const id = $page.params.id;
            if (!id) {
                loadError = 'Post not found.';
                isLoading = false;
                return;
            }

            try {
                post = await getPost(id);
            } catch (err) {
                loadError = err instanceof Error ? err.message : 'Could not load this post.';
                console.error('Post load error:', err);
            } finally {
                isLoading = false;
            }
        });

        return unsubscribe;
    });
</script>

<svelte:head>
    <title>{post?.title || 'Post'} | /Blog</title>
</svelte:head>

<Navbar />

<div class="min-h-screen bg-background px-6 py-16 text-on-background">
    <div class="mx-auto max-w-3xl">
        {#if isLoading}
            <p>Loading post...</p>
        {:else if loadError}
            <p class="text-red-600">{loadError}</p>
        {:else if post}
            <article class="space-y-6">
                <div class="space-y-2">
                    <p class="text-sm uppercase tracking-[0.2em] text-secondary">{formatDate(post.created_at)}</p>
                    <h1 class="text-4xl font-bold text-on-surface">{post.title}</h1>
                </div>
                <div class="prose max-w-none whitespace-pre-wrap text-on-surface-variant">
                    {post.body}
                </div>
            </article>
        {:else}
            <p>Post not found.</p>
        {/if}
    </div>
</div>