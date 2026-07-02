<script>
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { resolve } from '$app/paths';
    import { auth } from '$lib/firebase';
    import { onAuthStateChanged } from 'firebase/auth';
    import { getMe, getMyPosts } from '$lib/api';
    import Navbar from '$lib/components/Navbar.svelte';

    // Scroll progress tracking
    let scrollProgress = $state(0);

    /** @type {{ id: string, email: string, display_name?: string } | null} */
    let user = $state(null);
    /** @type {Array<{ id: string, title: string, body: string, created_at: string }>} */
    let posts = $state([]);
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

    /** @param {string | null | undefined} body */
    const getExcerpt = (body) => {
        if (!body) return '';
        return body.length > 180 ? `${body.slice(0, 177)}...` : body;
    };

    const handleScroll = () => {
        const winScroll = document.body.scrollTop || document.documentElement.scrollTop;
        const height = document.documentElement.scrollHeight - document.documentElement.clientHeight;

        if (winScroll > 10) {
            scrollProgress = (winScroll / height) * 100;
        } else {
            scrollProgress = 0;
        }
    };

    onMount(() => {
    const unsubscribe = onAuthStateChanged(auth, async (firebaseUser) => {
        if (!firebaseUser) {
            goto(resolve('/auth/login'));
            return;
        }

        try {
            user = await getMe();
            posts = await getMyPosts(); 
        } catch (err) {
            loadError = err instanceof Error ? err.message : 'Could not load your profile.';
            console.error('Profile load error:', err);
        } finally {
            isLoading = false;
        }
    });

    return unsubscribe;
});

</script>

<svelte:window on:scroll={handleScroll} />

<svelte:head>
    <title>{user?.display_name || 'Profile'} | /Blog</title>
    <link href="https://fonts.googleapis.com/css2?family=EB+Garamond:ital,wght@0,400..800;1,400..800&family=Inter:wght@300;400;500;600;700&family=Source+Serif+4:ital,opsz,wght@0,8..60,200..900;1,8..60,200..900&display=swap" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet">
</svelte:head>

<div
    class="fixed top-0 left-0 h-0.5 bg-primary `z-60` transition-all duration-150 ease-out"
    style="width: {scrollProgress}%; opacity: {scrollProgress > 0 ? 1 : 0};"
></div>

<div class="bg-background text-on-background font-body-md selection:bg-primary-fixed selection:text-on-primary-fixed min-h-screen antialiased">

    <Navbar />

    <main class="pt-32 pb-section-gap px-margin-mobile md:px-gutter">
        <div class="reading-column">

            {#if loadError}
                <div class="mb-element-gap p-4 rounded-lg bg-red-50 text-red-600 font-label-md text-label-md" role="alert">
                    {loadError}
                </div>
            {/if}

            <section class="flex flex-col md:flex-row items-center md:items-end gap-6 mb-section-gap">
                <div class="w-24 h-24 md:w-32 md:h-32 rounded-full bg-surface-container overflow-hidden flex items-center justify-center shadow-sm border border-outline-variant">
                    <span class="material-symbols-outlined text-outline text-5xl">person</span>
                </div>
                <div class="text-center md:text-left">
                    <h1 class="font-display-lg text-display-lg-mobile md:text-display-lg text-on-surface mb-1">
                        {isLoading ? 'Loading...' : (user?.display_name || 'Unknown')}
                    </h1>
                    {#if !isLoading && user?.email}
                        <p class="font-label-md text-label-md text-secondary">
                            {user.email}
                        </p>
                    {/if}
                </div>
            </section>

            <div class="flex items-center gap-4 mb-element-gap">
                <span class="text-label-md font-bold tracking-widest text-secondary uppercase">Recent Stories</span>
                <div class="h-px grow bg-outline-variant opacity-50"></div>
            </div>

            <div class="space-y-element-gap md:space-y-section-gap">
                {#if posts.length === 0 && !isLoading}
                    <p class="font-body-md text-body-md text-on-surface-variant">
                        No stories yet. Start writing one from the Write page.
                    </p>
                {/if}

                {#each posts as post (post.id)}
                    <article class="article-card p-6 md:p-8 rounded-2xl bg-surface-container border border-outline-variant shadow-sm hover:shadow-md transition-shadow duration-200">
                        <header class="mb-3">
                            <span class="text-label-md text-secondary font-medium">{formatDate(post.created_at)}</span>
                            <h2 class="article-title font-headline-md text-headline-md text-on-surface mt-2 transition-colors cursor-pointer leading-tight">
                                {post.title}
                            </h2>
                        </header>
                        <p class="font-body-md text-body-md text-on-surface-variant line-clamp-3 mb-4">
                            {getExcerpt(post.body)}
                        </p>
                        <a href={resolve(`/posts/${post.id}`)} class="inline-flex items-center gap-2 text-primary font-ui-button text-ui-button cursor-pointer">
                            Read full story
                            <span class="material-symbols-outlined read-more-icon">arrow_forward</span>
                        </a>
                    </article>
                {/each}
            </div>
        </div>
    </main>
</div>

<style>
    .material-symbols-outlined {
        font-variation-settings: 'FILL' 0, 'wght' 300, 'GRAD' 0, 'opsz' 24;
        display: inline-block;
        vertical-align: middle;
    }
    .reading-column {
        max-width: 720px;
        margin-left: auto;
        margin-right: auto;
    }
    
    /* Reapplied Card Interactions */
    .article-card {
        transition: transform 0.3s ease;
    }
    .article-card:hover {
        transform: translateY(-4px);
    }
    .article-card:hover .article-title {
        color: #486458; /* surface-tint */
    }
    
    .read-more-icon {
        transition: transform 0.2s ease;
    }
    .article-card:hover .read-more-icon {
        transform: translateX(4px);
    }
</style>