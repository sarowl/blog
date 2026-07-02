<script>
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import Navbar from '$lib/components/Navbar.svelte';
	import { auth } from '$lib/firebase';
	import { onAuthStateChanged } from 'firebase/auth';
	import { createPost, getPost, updatePost } from '$lib/api';

	let isCheckingAuth = $state(true);
	let isLoadingPost = $state(false);

	let title = $state('');
	let contentText = $state('');
	/** @type {HTMLTextAreaElement | undefined} */
	let titleEl = $state();
	let zenMode = $state(false);

	let postId = $state(/** @type {string | null} */ (null));
	let isSaving = $state(false);
	let saveStatus = $state('Not saved yet');
	let saveError = $state('');
	/** @type {ReturnType<typeof setTimeout> | undefined} */
	let saveTimeout;

	onMount(() => {
		const unsubscribe = onAuthStateChanged(auth, async (firebaseUser) => {
			if (!firebaseUser) {
				goto(resolve('/auth/login'));
				return;
			}

			const id = $page.params.id;
			if (id) {
				postId = id;
				isLoadingPost = true;
				try {
					const existingPost = await getPost(id);
					title = existingPost.title || '';
					contentText = existingPost.body || '';
					saveStatus = 'Loaded existing post';
				} catch (err) {
					saveError = err instanceof Error ? err.message : 'Could not load the post.';
					saveStatus = 'Could not load post';
				} finally {
					isLoadingPost = false;
				}
			}

			isCheckingAuth = false;
		});

		return unsubscribe;
	});

	function handleTitleResize() {
		if (titleEl) {
			titleEl.style.height = 'auto';
			titleEl.style.height = titleEl.scrollHeight + 'px';
		}
		triggerAutoSave();
	}

	/** @param {string} status */
	async function persist(status) {
		if (!title.trim() || !contentText.trim()) {
			saveStatus = 'Add a title and some content first';
			return;
		}
		isSaving = true;
		saveStatus = status === 'public' ? 'Publishing...' : 'Saving...';
		saveError = '';

		try {
			const payload = { title: title.trim(), body: contentText.trim(), status };
			let result;

			if (!postId) {
				const draftResult = await createPost({
					title: payload.title,
					body: payload.body,
					status: 'draft'
				});
				postId = draftResult.id;
				result = status === 'public'
					? await updatePost(draftResult.id, { ...payload, status: 'public' })
					: draftResult;
			} else {
				result = await updatePost(postId, payload);
			}

			postId = result.id;
			const time = new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
			saveStatus = status === 'public' ? `Published at ${time}` : `Draft saved at ${time}`;
		} catch (/** @type {any} */ err) {
			saveError = err.message;
			saveStatus = 'Save failed';
		} finally {
			isSaving = false;
		}
	}

	function triggerAutoSave() {
		clearTimeout(saveTimeout);
		saveTimeout = setTimeout(() => {
			persist('draft');
		}, 1500);
	}

	function handleSaveDraft() {
		clearTimeout(saveTimeout);
		persist('draft');
	}

	function handlePublish() {
		clearTimeout(saveTimeout);
		persist('public');
		goto(resolve('/feed'));
	}
</script>

<svelte:head>
	<title>{postId ? 'Edit post' : 'Write'} - /Blog Editorial</title>
	<link href="https://fonts.googleapis.com/css2?family=EB+Garamond:ital,wght@0,400..800;1,400..800&family=Inter:wght@300;400;500;600&family=Source+Serif+4:ital,opsz,wght@0,8..60,200..900;1,8..60,200..900&display=swap" rel="stylesheet">
	<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet">
</svelte:head>

{#if isCheckingAuth || isLoadingPost}
	<div class="min-h-screen flex items-center justify-center bg-surface text-on-surface">
		<span class="material-symbols-outlined animate-spin">progress_activity</span>
	</div>
{:else}
<div class="bg-surface text-on-surface selection:bg-primary-fixed selection:text-on-primary-fixed min-h-screen flex flex-col scroll-smooth antialiased {zenMode ? 'writing-mode-active' : ''}">
	<Navbar />

	<main class="flex-1 pb-section-gap px-margin-mobile md:px-0 flex flex-col items-center">
		<div class="w-full max-w-container-max">
			<div class="mb-6 flex items-center justify-end gap-4">
				<button onclick={handleSaveDraft} disabled={isSaving} class="font-ui-button text-ui-button text-secondary hover:text-primary transition-colors px-4 py-2">
					Save Draft
				</button>
				<button onclick={handlePublish} disabled={isSaving} class="font-ui-button text-ui-button bg-primary text-on-primary px-6 py-2 rounded-lg hover:scale-[1.02] active:scale-[0.98] transition-transform duration-200">
					Publish
				</button>
			</div>
			
			<div class="mb-12 flex items-center justify-between opacity-60">
				<span class="font-label-md text-label-md flex items-center gap-2 transition-all">
					<span class="material-symbols-outlined {isSaving ? 'animate-pulse' : ''}">
						{isSaving ? 'sync' : 'edit_note'}
					</span>
					{saveStatus}
				</span>
				<button
					class="toolbar-btn {zenMode ? 'text-primary bg-primary-fixed' : ''}"
					onclick={() => zenMode = !zenMode}
					title="Toggle Zen Mode"
				>
					<span class="material-symbols-outlined">{zenMode ? 'visibility_off' : 'visibility'}</span>
				</button>
			</div>

			{#if saveError}
				<div class="mb-8 p-4 rounded-lg bg-red-50 text-red-600 font-label-md text-label-md" role="alert">
					{saveError}
				</div>
			{/if}

			<textarea 
				bind:this={titleEl}
				bind:value={title}
				oninput={handleTitleResize}
				class="w-full bg-transparent border-none focus:ring-0 p-0 mb-8 font-display-lg text-[48px] font-semibold leading-[1.1] placeholder:text-outline-variant resize-none overflow-hidden outline-none" 
				placeholder="Title" 
				rows="1"
			></textarea>

			<div 
				id="editor"
				contenteditable="true"
				bind:textContent={contentText}
				oninput={triggerAutoSave}
				class="editor-content w-full font-body-lg text-body-lg min-h-500px text-on-surface-variant leading-relaxed focus:outline-none" 
				data-placeholder="Tell your story..." 
			></div>
		</div>
	</main>

	<footer class="w-full py-section-gap px-gutter flex flex-col md:flex-row justify-between items-center max-w-container-max mx-auto bg-surface transition-opacity duration-500 {zenMode ? 'opacity-10' : 'opacity-100'}">
		<div class="mb-4 md:mb-0">
			<span class="font-display-lg text-[24px] font-bold text-primary">/Blog</span>
			<p class="font-label-md text-label-md text-secondary mt-2">© 2024 /Blog Editorial. All rights reserved.</p>
		</div>
	</footer>
</div>
{/if}

<style>
	.material-symbols-outlined {
		font-variation-settings: 'FILL' 0, 'wght' 300, 'GRAD' 0, 'opsz' 24;
		font-size: 20px;
		display: inline-block;
		vertical-align: middle;
	}

	.editor-content[contenteditable]:empty::before {
		content: attr(data-placeholder);
		color: #727974;
		font-style: italic;
		pointer-events: none;
		display: block;
	}

	:global(::-webkit-scrollbar) {
		width: 6px;
	}
	:global(::-webkit-scrollbar-track) {
		background: transparent;
	}
	:global(::-webkit-scrollbar-thumb) {
		background: #e3e2e0;
		border-radius: 3px;
	}
	:global(::-webkit-scrollbar-thumb:hover) {
		background: #bcc7dd;
	}

	.toolbar-btn {
		padding: 0.5rem;
		border-radius: 0.5rem;
		color: #545f72;
		transition: all 0.2s;
		cursor: pointer;
	}
	.toolbar-btn:hover {
		background-color: #eeeeeb;
		color: #173329;
	}
</style>
