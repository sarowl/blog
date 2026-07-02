<!-- src/lib/components/Navbar.svelte -->
<script>
	import { resolve } from '$app/paths';
</script>

<nav class="flex items-center justify-between px-6 py-3 bg-white border-b border-gray-200">
  <a class="text-xl font-bold text-gray-800" href={resolve('/feed')}>
    /Blog
  </a>

  <div class="flex items-center space-x-6">
    <a class="flex items-center text-gray-600" href={resolve('/profile')} aria-label="Profile">
      <svg xmlns="http://www.w3.org/2000/svg" 
           fill="none" viewBox="0 0 24 24" 
           stroke-width="1.5" stroke="currentColor" 
           class="w-6 h-6">
        <path stroke-linecap="round" stroke-linejoin="round" 
              d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.5 20.25a8.25 8.25 0 1115 0v.75H4.5v-.75z" />
      </svg>
    </a>

    <a class="px-4 py-2 rounded-md bg-black text-white font-medium" href={resolve('/write')}>
      Write
    </a>
  </div>
</nav>
