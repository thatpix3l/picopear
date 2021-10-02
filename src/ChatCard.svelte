<script>
    import { onMount } from 'svelte';
    
    // Add feature to get max scroll amount in Element types
    (function(elmProto){
        if ('scrollTopMax' in elmProto) {
            return;
        }
        
        if ('scrollLeftMax' in elmProto) {
            return;
        }

        Object.defineProperties(elmProto, {
            'scrollTopMax': {
                get: function scrollTopMax() {
                  return this.scrollHeight - this.clientHeight;
                }
            },
            'scrollLeftMax': {
                get: function scrollLeftMax() {
                  return this.scrollWidth - this.clientWidth;
                }
            }

        });
    })(Element.prototype);
    
    // Add feature to parse a String and return milliseconds as an int
    function toMillisecond(string_with_time) {
        if(!(string_with_time.endsWith("s") || string_with_time.endsWith("ms"))) {
            throw "Does not end with 's' or 'ms'!";
        }
        let isMilliseconds = string_with_time.indexOf('ms') !== -1;

        if(isMilliseconds) {
            string_with_time = parseInt(string_with_time.replace('ms', ''));

        } else {
            string_with_time = parseInt(parseFloat(string_with_time.replace('s', '')) * 1000);

        }
        
        return string_with_time;

    }
    
    // Chat message object
	export let chatMsg; 
    
    // How long a card is allowed to be visible
    export let visibilityDurationMs = 5000;
    
    // Chat card css width
    export let cardWidth = 350;
    
    // Object of keys of colors
    export let colorPalettes;

    // Default chat color if not existing
    chatMsg.msgBodyFg ||= "#333";
    
    // Randomly pick color palette from array of color palettes
    let colorPalette = colorPalettes[Math.floor(Math.random() * colorPalettes.length)];
    
    // Chat vard visibility
    export let isCardVisible = true;
    
    // Reference to chat card. Updated during onMount
    let chatCardElem;
    
    // Func to scroll msg body smoothly, infinitely
    function autoScroll(elem, delay) {
                
        // Scroll every something-milliseconds
        setInterval(() => {

            // Get max scrollable X
            const scrollLeftMax = elem.scrollLeftMax;

            // Get scroll position
            let scrollPosition = parseInt(elem.scrollLeft);
            
            // Scroll if not at end
            if(scrollPosition < scrollLeftMax - 1) {
                elem.scrollBy({
                    left: 200,
                    behavior: 'smooth'
                });

            } else {
                elem.scroll({
                    left: 0,
                    behavior: 'smooth'
                });

            }

        }, delay);

    }
    
    function autoHide(elem, delay) {
        let animDisappearDuration = getComputedStyle(elem).getPropertyValue("--anim--disappear__duration");
        animDisappearDuration = toMillisecond(animDisappearDuration);
  
        setTimeout(() => {
            isCardAllowed = false;
        }, delay);

    }
    

    // Beginning card lifecycle
    onMount(() => {
        let msgBodyElem = chatCardElem.querySelector(".messageBody");
        let disappearDurationMs = toMillisecond(window.getComputedStyle(chatCardElem).getPropertyValue("--anim--disappear__duration"));

        chatCardElem.style.setProperty("--visibility__duration", `${visibilityDurationMs}ms`);
        

        // Autoscroll
        autoScroll(msgBodyElem, 1000);

        // Unrender and remove layout of card after x milliseconds,
        // not including disappearing animation duration
        autoHide(chatCardElem, visibilityDurationMs - disappearDurationMs);
    });
    
</script>

<style>
    .card {
        --anim--appear__duration: .4s;
        --anim--disappear__duration: .4s;

        max-width: var(--card__width);
        max-height: var(--card__height);
        margin-bottom: 10px;
        overflow: hidden;
    }

    .usernameBody {
        font-size: .9rem;
        padding: 2px;
    }
    
    .messageBody {
        padding: 5px;
        overflow-x: auto;
    }
    
    .messageBody, ::-webkit-scrollbar {
        /* Hide scrollbars, but still usable with mouse */
        -ms-overflow-style: none;  /* IE and Edge */
        scrollbar-width: none;  /* Firefox */
    }
    
    .messageBody > * {
        font-size: 1.25rem;
        white-space: nowrap;
    }
    
    .appear {
        animation-name: fade-in, swing-up;
        animation-duration: var(--anim--appear__duration);
    }
    
    .disappear {
        animation-name: shrink-height, fade-out;
        animation-duration: var(--anim--disappear__duration);
        animation-timing-function: ease-out;
        animation-fill-mode: forwards;
    }
    
    @keyframes swing-up {
        0% {
            transform-origin: top left;
            transform: translateY(20px) rotate(3deg);
        }
        
        50% {
            transform: translateY(0px);
        }
        
        100% {
            transform: rotate(0);
        }
    }
    
    @keyframes fade-in {
        0% {
            opacity: 0;
        }

        100% {
            opacity: 1;
        }
    }
    
    @keyframes fade-out {
        0% {
            opacity: 1;
        }

        100% {
            opacity: 0;
        }
    }
    
    @keyframes shrink-height {
        100% {
            max-height: 0px;
            margin: 0;
            padding: 0;
        }
    }
    
</style>

<div id="{chatMsg.id}" class="card {isCardVisible === true ? 'appear' : 'disappear'}" bind:this={chatCardElem}>
    <div class="usernameBody" style="background: {colorPalette.usernameBg}">
        <p style="color: {colorPalette.usernameFg}"><b>{chatMsg.username}</b></p>
    </div>
    <div class="messageBody" style="background: {colorPalette.msgBodyBg}">
        <p style="color: {colorPalette.msgBodyFg}">{chatMsg.msgText}</p>
    </div>
</div>