<script>
    import { resolve } from '$app/paths';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/firebase';
    import { createUserWithEmailAndPassword } from 'firebase/auth';
    import { PUBLIC_API_URL } from '$env/static/public';

    // Form state
    let fullName = '';
    let email = '';
    let password = '';
    let showPassword = false;

    // Interaction states
    let isLoading = false;
    let isSuccess = false;
    let errorMessage = '';

    // Scroll tracking
    let scrollProgress = 0;

    const togglePassword = () => {
        showPassword = !showPassword;
    };

    const handleSignup = async () => {
        errorMessage = '';
        isLoading = true;

        try {
            // 1. Create the user in Firebase Auth
            const credential = await createUserWithEmailAndPassword(auth, email, password);

            // 2. Get the ID token to authenticate with our backend
            const idToken = await credential.user.getIdToken();

            // 3. Register the user in our own database via the Go API
            const res = await fetch(`${PUBLIC_API_URL}/api/register`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${idToken}`
                },
                body: JSON.stringify({
                    email,
                    display_name: fullName
                })
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || `Registration failed (${res.status})`);
            }

            // Redirect to profile page after successful registration
            goto(resolve('/profile'));
        } catch (err) {
            if (err instanceof Error) {
                errorMessage = err.message;
            } else {
                errorMessage = 'Something went wrong. Please try again.';
            }
            console.error('Signup error:', err);
        } finally {
            isLoading = false;
        }
    };

    // Atmospheric scroll progress
    const handleScroll = () => {
        const winScroll = document.body.scrollTop || document.documentElement.scrollTop;
        const height = document.documentElement.scrollHeight - document.documentElement.clientHeight;
        
        if (winScroll > 10) {
            scrollProgress = (winScroll / height) * 100;
        } else {
            scrollProgress = 0;
        }
    };
</script>

<svelte:window on:scroll={handleScroll} />

<svelte:head>
    <title>Join /Blog - Create Your Account</title>
    <link href="https://fonts.googleapis.com/css2?family=EB+Garamond:ital,wght@0,400..800;1,400..800&family=Inter:wght@100..900&family=Source+Serif+4:ital,opsz,wght@0,8..60,200..900;1,8..60,200..900&display=swap" rel="stylesheet" />
    <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet" />
</svelte:head>

<div 
    class="fixed top-0 left-0 h-`2px` bg-primary z-`100` transition-all duration-300 ease-out" 
    style="width: {scrollProgress}%; opacity: {scrollProgress > 0 ? 1 : 0};"
></div>

<div class="flex flex-col min-h-screen bg-surface text-on-surface antialiased">
    <header class="w-full py-8 px-gutter flex justify-center">
        <a class="transition-transform duration-200 hover:scale-105" href={resolve('/')}>
            <img alt="/Blog Logo" class="h-10 w-auto opacity-90" src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAPAAAABQCAYAAAAnSfh8AAAQAElEQVR4AexdC3hcxXU+Z+7qYRsZG2Mby9Ku8Eu7Suw6SNqVA9Q8bENCiuOQFAqlobR8ENKU0DThI04TmjQECA00MYWmpISUkjTAF+OkQDDGIY0fu7JMwam0kmWwVpJt4sQGFDCWdu/JmV3t+mr33rvSap/a2e/OzsyZM69/9r8zd+bOrAD1UQgoBEoWAUXgkm06VXCFAIAisPoVKARKGAFF4BJuPFV0hUA5E1i1vkKg5BFQBC75JlQVKGcEFIHLufVV3UseAUXgkm9CVYFyRkARuDxbX9V6iiCgCDxFGlJVozwRUAQuz3ZXtZ4iCCgCT5GGVNUoTwQUgcuz3cu51lOq7orAU6o5VWXKDQFF4HJrcVXfKYWAIvCUak5VmXJDQBG43Fpc1XdKITBBAk+puqvKKARKHgFF4JJvQlWBckZAEbicW1/VveQRUAQu+SZUFShnBBSBx936SlEhUHwIFJjAhJevHehZv26QpLl87eB3ig8iVSKFQPEiUFACr193aBUiLo3Dgwg/dLl9V7g8PsqKcXvfdLq9/ZxmF9vtLrd3m8vj/Vdno/cTtcua3fF8x2NzvL2u8ZTL7X1tPOkpHYVANhAoKIG5AhvYRC8COPL087W7+oL+p4ZhuIEi+hoiugmInosqZPKFeDoi1gGCm+0WQLwIAD+FAn9coTmY1L5BV2PrNTCOzzt6xWo9Qn9MBFvHoa5UFAJ5QaCgBGaCro/XEgl+BsDfAHC46+W+UE/7tlAw8G99wcCHxkManehzb58Mz9JHwksiYWxj/9Uc79sEdICTNL0QoRaEeMzp8W2eP3/FDFOlUeFvu3cM9fcE/jcShuu53PqouEwsVc1iRaBgBF5/8eEmNAyfmRSbrUBiIv7IKswoP/5ax1v9vR0HBvbv9vcHAz8MBf236EPaco5/l1Ev2Y0A66tmV/93stzMP9jrH2D5QTbqUggUHIGCERgwYhw+D1UdWmg5NBVAGRNmYGDXCSby7UDwz3Zo883kMn7GfcBOxxB2zOBWToVAwRAoGIEJ4dTwGejnT3TisBUKYYHvWoWNV87P1n/Pw+nONPo31za2NKbRAQQcSqejwhUC+UCgIAS+fO3RWgBsgdEPgbAcPksVB0ZOSnuyhofS30uXhgPF7el0VHg5IVDcdRWFKB7C8AbkbiyWNw2/q1fzBFbMZ/YdHolNbpmFTUQmhL49nT4X62PpdFS4QqBYECgIgY3DZwJ86YUXzngrH4CcHI7ICSj7rBBq5i9eMc9eSYUqBIoDgbwTeM2aY6cj0Op49blvfTruzrV9pPflo+PJQ1Rq08ajp3QUAoVGIO8Eni7e+wgAVgJ/iGeVCCp/ws68XPVLmhenzYhg6DCvQ6fVy5GC07l8tmtpm8fl9q6ad7Zvfo6ySSTb0LBylsxL5indiQDlKAkEcktgEwgQ9I+eEtOeLVvnHjrlz60LK7Rz0+XAw3t/Op1shy/0+NqcHt/9TKTXcMb0Y+CgTkDcOa0ajrD8GJutLo/3r+rqVk16ZCBvYk6P926Xx/cM2yGaVnVc5iXzlG6X2/c6yx90Nvo+XuduucDO1Dd+gCcjs42GSm8iCOSVwJ9ookp+5r0kXsB8Dp9lnjxBdaW07QwC3GcXns0w2cM63b5HHAC7ON9bmEhnj6b/DA9O7gGiF1k+m80aAHxY1OiDLo/3n6CpKTqCgQl8Zi9qPp2JeY+ocPQi4Bc46ofYrudR0D4CuFWP0AYiehQQGlh+Ewp4QkNtu51BrOTRFKhPARHIK4FP1g6uRYCaRH1Jy9/wudErl60+nMjbxME/5m/3dfmfMQnKuqjO41taXU3tPBt/nTFx0uGzXIbLQl2B2/qCgYt1Ha6PhzN2swFwo5Nqtk1koq2usfXSmVVaEAE/D2M//qHh8PmhLv/9/T2BzaFg4DrG4KmxKspXzAjklcCImBg+891+/9PbFnTmA5zaRc1OIfBJu7x0gntDQf8tdjrZCnMubVkkgLYhYL0xTSbP1lC3/1+Msv5u/yPcGz9klCHAeVUV1bsWLvP+kVFu5q5bfM4SIcTjAHgWJH1GIuHr5OunRjFi+E6jP+7m9hrgnvoFowHQD8fDlV0YBPJIYEJCSAy5mMw5n30+Y4lvZr3He2tFpePXDK+LTepF8BsA+kh/0J/cO6XqTkRirYuoif/AJPJKdZ30b0k72RDoKa94Mn6LNAHPyqFxsn7cX1vbPF2rqHgKAbjnjktjNt8sdhzq6QjGfKe++7o69hLwM/gpUcyFOEuPRDZyb702YYLtP40Fqu9CIZA3Akf3/gIYe4GsDp8F4sUuj+8f4sbp9m6pqYC3BOC3gNd2xwJMR/hH+lAE6GPvHT+xqK8r8D9jw3Pnczb6/hYQV5vl8M6IvstM3t+1h29A9GZyGCIumFmpfSVZHvdrNY6rAGFF3D/GRr19jN/oIdxh9Eo33wRO04T4kbwpSL8yxYGAyGMxjJsXont/s5y3fL79KqcZNfzj/hN2p1zc8xwiwLsBadNAV+Anb7zx6jspSjkSRH/8AswJR/R68nDWWAwCeNXoj7sJ8TP1FrPBPBG1Lq6XYuvCcoMIY3QkRV8KEM921DhukE5ligMBka9i8DPUqc0Lhr2/WcufQA4H/fxD7wWglN4qng8i1HJvch+S+DX30vK0jk21jek3MMTjT8Z2zBTXcN4pw1mZJgGke540JRWn50CsuEKmkWwQ6JJkWdxPAt6Iu1NsoZvmJfUYv0Q7Sr8yhUUgLwSeyN7fTOHQgf69r8vfxs9nS3lIPPsd3TEzrMNFPFS+DaLPuakpI6I8rePTFUIL8hLL3XObmk5L1cqmBE2JJnPgsti/TkpouSMLAROjG5mWNHL9FgBngcVHkG65+wt1wfMCFhGBmq1CEnLlyBsCeSHwRPb+Zqvm8gSNwW7/9lBX4J6T4eH3c69su+MJeW10ul6zr87dvDxbZRibzgUOQLR+kYTIctQg00EAy6E+36RWNzS1GucX+AkBZsp4VoYARqzCQESsy4Jom65lmiogJwjkicBwZbz0PKyz3fsb18umLd+B5l55Az/bbbVNF6FBgPaia2mbx1Yvg0DnshNtCGDdwyNY9rAyO57DtwxHRKHr0Cr14gYjGj9KxH0mNuJsE2lURIBzog6zL6JXzMRKVhgEck5gOXwGxESvRmn2/uYSBl2Iq5jEtq9uIuKZpOkPZ7schJA4fdMsbbLrEWUEwrC0LA3iGEKGev3dlrocwPnNYcvqsgzjevRYRVLy/COQcwIbh888jE279zeXEAx07jqGSI+kywMRPyjfXkqnN5FwgWR7aB4/p1sPaWVGSJq0LI1Oyc+7EZ44fNRKnxvekqQChGUY6PhjqzSVHADyDAK3Y25z5Dt2YtaS8rj316pWJ0mMq3fVUCReOrFKayJyItTs9Hl4bUtgjm/bViRS00eM3A5Eb5vlSwRjhtxGHSa+1+g3uP2hbr/tG20GXeXMAwK2P4rJ5p98dA4S5Pztq3RlPhLcfZAA9qfT49FCdp+DkawnhrgwfKOzHSIj6g5Ws7wEUV9yYF9Xx+Ew4UfNSMyjjHVm68fzzvbN57CUmxcB9J4kvCo5D+UvLAI5JTAajs7hOz5flVl9+2oS0FmvgY4mSgjZ3ov7u9GkzS3CiHlATIogKsHmE9bR9PxrOROv67QGTJbSECvvSE6yuppSZNxw+yIjcKG8+SXrK39hEcgpgZkEieEz92h53ftrCysB2oZHA9Fy2SYaPMEvfSQsXzSxjMWz81WWgRzAJJrOluU1qB+zTL+/p719dCnt69wOiZEAItzgdHu/X9/oW1/v9v6Z0+P7TwS8aUwmRHeGgv4Vo+dhjwlSnsIjkDMCr8ng6Jz8wUHmGxsMBWBCDRq8k3YOHNjby72g5euLAGhLYASoBosPEQxCb6/tyZ2jS2lfOvEeugno/+NJIeInhYDNAvFxBPhzKef0Ooj0r0IYm/qCgY1SpkxxIpAzAhuPzolWPY97f6P5WXzVLm6t5x9tnUVwQkw6vpTwZM9hudeWSZU8i5yUK1mGM/FsX1KJJ1Tvbr2uuop2AuEsJujtBLAuQpEL44ZG4H2RITGde9yWULD9K337d3fF4yq7OBHIGYHRcHQOEeV+7y8C/x7Tg+yoxKvTaXFCx+HEuz9IpzfR8MjIyEOMhW4Wj0m4wEwel3GZzoy7x9o0PDIcvmesbKyvoWHlLB4qPydQPIKAcyNhaGOC3hXq8m8dCO75Rdzw2nGn/CeLsbGVr5gRyAmBU47OQZzU7LNDw7RnQaFOtpNAshEaGhqqeSb8L6XbziDRl0OhfcftdDIJk8NoRHzONC6iLYERYYlpPIDvHXqtI2QRBi5P8wJ9WlU7IkY3NnBPv0U9z1qhVXrynBA45egcgMnNPhOkJTAkvYlk1hT6tLnfYL1Gs7C4jJ//tvQFA5vifjObSWD5SiT3lLbvCnP6X2adcHK67PdYbc6PLfegyRCaRoZh5G6Oa3kRad9EMJAf+RZmqa0CSg2BnBCY7/anjs4BmPTeXx2wIR2wCHChlc48XtuMzbCKz1rpxOS0mZ//DDPnMWnKN6IJmWJaXPc57BJsTC9Ov4PXZXk2ODW4psKxNlUKgOD4OJh+8Gt2R+DWNrY0cnmuMUZFnqhyeXyfNsqUu3QRsPyhZV4lQl4+SrwIgFnY+8udxpVpy4O42unxPljX6FshdZ1LWxbVNbZeys9+36+uhl75w5VyM8M94nHuGW+UGx7Mwo2y2kXNTibgYqMs2e1yt7Yly4z+UDDwj0CQ+kqnoDuST5yU50QDituM8aWbdHq4r8v/Nem2MprAlRZhm5we35DT7d3hdPt+xeZVl9v3usvjPe7y+MgVsw863d59Lrd3p9Pte7Le03pz3bJVCy3SU+ICIZB1AicfnUNEmydaN/ncVresxStnTZ1u788Q0fpkCUPiyGuYmoBXXPwjRId2QBPiWY77SQQwH/ISdXP57oiMwAruGb9rSCrFeWbjuTWuRu9FFZWOn3KatrgRivt4XfWSuqZVZ6QkFBNQX9D/1zwU/2bMG/tGwPe59NOer1va5ps/f8UMp6dlDcyYvgURTp2/TDDEN5svhroDN8ZiWX8LwrlWociYIOIHEeFcNsu5m2+AxP5hlCMMFyK+HxBXIcIVAsQDmqYPMPF/YDXUB/XJOwK2P8QMS5PYXE4AQ3b/+2uWvsvtuwLAcUjTNH901hTxMjO9ccuIjnKP2Q1Au9ls1gnuZfuGSBjb+FnXLXvD8UzqzBAjL4HAbWB1xpShQAjg5XXV5zSKWJ5swep6qCvwBYBwM99Enmd/7OKRhOag3dVnTPs9graV0zoP5IeJy+V+EDDcyDebb7BIZ2N7scJ+W4UMArk819ZUaTsbGlZKkmeQQmoUJckcAZF5VIuYBHOYuL+Uhomz6YlOtDz5wSwF7pme4qEhZs0EA/MkUXl4vIrNhv6g//NsPzywf/eE/oGB45wz8TIFKs3qM/s0TAAAA1pJREFUaJT1dXXsDQUDlyDqC7hnvZExu4+xe4x1nmX34wT6/bpOn3pXDNVyGW5m/XRH73DU2NUfDPyczE6YjAVn/I2ATTyz/Z2ME1ARs4ZA1gn89NaF1295fuHqqNla98WslXSKJ3Sws/1IKOj/Lt9s/o7XZ6/lm8WH2X1NqKv91v7uwENHOzt/nwkECJFrIdp7ZxLbNs5VZy35gOUQ3TamCswaAlkncNZKphLKCgJ93MMD6Jdyb86PEVlJMpoIAjgqtcrVUY/6KhgComA5q4zzhkAkokV4WH7qOXs0Zx5eH+DeeTsP3XcA0Mvslhsi+pjs4/obVk5mIRt1TQKByUZVBJ4sgkUcf+HS1pVOj/eXclIMET8ji6oTbRzRI24eomOoK7CE5xwu4qH7efx8fQ67PSxv4KH7PLaj8xB44uRsJv+tMm6yQaTxEj05qvJnCQFF4CwBWWzJyKUezYFbEPD8eNmI4K7+YODOQ917xj2cPnjw/97kZ/L7QaeL4+kkbKGrzQ4JMArjUAQuDO45z3VmleNcJu+YP0/jJbAXM80YTw7vHROX6JW+zj0vj5EpT94RUATOO+T5yZAiesr+YATKeNKJqivlX9ckCq/rcAeoT8ERKGUCFxy8Yi6AGD7KE1OQNMTFjbEXZSZW8jp3ywUEmNiyyM/Rm/p7ApsnlorSzgUCisC5QLUI0jx48OB7EQpfSQBjD3hHeNLl9v7C5W79nNPta66zeN1Tylnvaqfb+6iG2nZEiM04EzzAz9HRCbEiqGbZF0EReAr/BAaCHfvEiTeW8wTUl4hoIFFVxNWA4l5E2KOR/jsm6W+ZrDtHTbfT4zsm5YD4X4j4F7ys9DbfCB4LR2glz1T/TSId5Sg4AorABW+C3BZA9sR93YGvh4KBetLpT3m990Fe7+W1X0qckskkncNkXcVmUbQ0RLLXfiZK/Aicz8tKp/NM9LWDPQH1typRgIrnSxG4eNpiIiXJSDfUHXiC13tv5l6U134DZ8bXeg32WUzWRia7l2WXRYnf4/9VRpmpSHlBQBE4LzCrTBQCuUFAETg3uKpUFQJ5QUAROC8wq0wUArlBQBE4N7iqVHOHgErZgIAisAEM5VQIlBoCisCl1mKqvAoBAwKKwAYwlFMhUGoI/AEAAP//Rm8aaAAAAAZJREFUAwCrvs37t/Ok5wAAAABJRU5ErkJggg==" />
        </a>
    </header>

    <main class="grow flex items-center justify-center px-gutter py-12">
        <section class="w-full max-w-`440px` flex flex-col items-center">
            
            <div class="bg-surface-container-lowest p-8 md:p-12 rounded-lg w-full text-center shadow-[0px_4px_20px_rgba(0,0,0,0.04)] hover:scale-[1.005] transition-transform duration-200 ease-in-out">
                <header class="mb-element-gap">
                    <h1 class="font-headline-md text-headline-md text-primary mb-2">Join /Blog.</h1>
                    <p class="font-body-md text-body-md text-secondary opacity-80">Start sharing your stories with the world today.</p>
                </header>
                
                <form class="space-y-6 text-left" on:submit|preventDefault={handleSignup}>
                    <div class="group">
                        <label class="font-label-md text-label-md text-secondary mb-2 block" for="full_name">Full Name</label>
                        <input 
                            bind:value={fullName}
                            class="w-full bg-transparent border-0 border-b border-outline-variant focus:ring-0 focus:border-primary transition-colors py-3 px-0 font-body-md text-body-md text-on-surface placeholder:text-outline-variant" 
                            id="full_name" 
                            name="full_name" 
                            placeholder="Enter Name" 
                            required 
                            type="text"
                        />
                    </div>
                    
                    <div class="group">
                        <label class="font-label-md text-label-md text-secondary mb-2 block" for="email">Email address</label>
                        <input 
                            bind:value={email}
                            class="w-full bg-transparent border-0 border-b border-outline-variant focus:ring-0 focus:border-primary transition-colors py-3 px-0 font-body-md text-body-md text-on-surface placeholder:text-outline-variant" 
                            id="email" 
                            name="email" 
                            placeholder="Enter Email" 
                            required 
                            type="email"
                        />
                    </div>
                    
                    <div class="group">
                        <label class="font-label-md text-label-md text-secondary mb-2 block" for="password">Password</label>
                        <div class="relative">
                            <input 
                                bind:value={password}
                                class="w-full bg-transparent border-0 border-b border-outline-variant focus:ring-0 focus:border-primary transition-colors py-3 px-0 font-body-md text-body-md text-on-surface placeholder:text-outline-variant" 
                                id="password" 
                                name="password" 
                                placeholder="••••••••" 
                                required 
                                type={showPassword ? "text" : "password"}
                            />
                            <button 
                                class="absolute right-0 top-1/2 -translate-y-1/2 text-secondary opacity-50 hover:opacity-100 transition-opacity" 
                                type="button"
                                on:click={togglePassword}
                            >
                                <span class="material-symbols-outlined text-sm" style="font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;">
                                    {showPassword ? 'visibility_off' : 'visibility'}
                                </span>
                            </button>
                        </div>
                    </div>

                    {#if errorMessage}
                        <div class="rounded-md bg-red-100 border border-red-300 text-red-700 px-4 py-3 text-sm">
                            {errorMessage}
                        </div>
                    {/if}
                    
                    <button 
                        class="w-full font-ui-button text-ui-button py-4 rounded-lg hover:opacity-90 transition-all active:scale-[0.98] mt-4 flex justify-center items-center gap-2 {isSuccess ? 'bg-surface-container-highest text-primary' : 'bg-primary text-on-primary'}" 
                        type="submit"
                        disabled={isLoading}
                    >
                        {#if isLoading}
                            <span class="animate-spin material-symbols-outlined" style="font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;">progress_activity</span> 
                            <span>Processing...</span>
                        {:else if isSuccess}
                            <span class="material-symbols-outlined" style="font-variation-settings: 'FILL' 0, 'wght' 400, 'GRAD' 0, 'opsz' 24;">check_circle</span> 
                            <span>Account Created</span>
                        {:else}
                            <span>Create Account</span>
                        {/if}
                    </button>
                </form>

                <div class="mt-8 pt-8 border-t border-surface-container-high">
                    <p class="font-label-md text-label-md text-secondary mb-4">
                        Already have an account? 
                        <a class="text-primary font-bold hover:underline" href={resolve('/auth/login')}>Log In</a>
                    </p>
                    <p class="font-label-md text-[12px] text-outline opacity-70 leading-relaxed max-w-[320px] mx-auto">
                        By creating an account, you agree to /Blog's 
                        <a class="underline hover:text-primary" href="#terms">Terms of Service</a> 
                        and 
                        <a class="underline hover:text-primary" href="#privacy">Privacy Policy</a>.
                    </p>
                </div>
            </div>
        </section>
    </main>

    <footer class="w-full py-8 px-gutter flex flex-col md:flex-row justify-center items-center gap-4 border-t border-surface-container-high mt-auto">
        <p class="font-label-md text-label-md text-secondary opacity-60">© 2024 /Blog Editorial. All rights reserved.</p>
        <div class="flex gap-6">
            <a class="font-label-md text-label-md text-secondary hover:text-primary transition-colors" href="#privacy">Privacy</a>
            <a class="font-label-md text-label-md text-secondary hover:text-primary transition-colors" href="#terms">Terms</a>
            <a class="font-label-md text-label-md text-secondary hover:text-primary transition-colors" href="#archive">Archive</a>
            <a class="font-label-md text-label-md text-secondary hover:text-primary transition-colors" href="#contact">Contact</a>
        </div>
    </footer>
</div>