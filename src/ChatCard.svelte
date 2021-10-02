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
    
    // Object of keys of colors
    export let colorPalettes;

    // If card should be visible or not
    let isCardAllowed = true;

    // Randomly pick color palette from array of color palettes
    let colorPalette = colorPalettes[Math.floor(Math.random() * colorPalettes.length)];
    
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
        // Autoscroll ChatCard message text body
        let msgBodyElem = chatCardElem.querySelector(".messageBody");
        autoScroll(msgBodyElem, 1000);

        // Set a CSS var for computed height when a card is set or implied to have auto
        // Handy for CSS transitions that try to transition height
        let computedHeight = parseInt(window.getComputedStyle(chatCardElem).height);
        document.documentElement.style.setProperty("--card__computed-height", `${computedHeight}px`);

        // Unrender and remove layout of card after x milliseconds,
        // not including disappearing animation duration
        // This is up to the maintainer of the CSS used to stylize the cards
        autoHide(chatCardElem, visibilityDurationMs - disappearDurationMs);
    });
    
</script>

<div id="{chatMsg.id}" class="card {isCardAllowed === true ? 'appear' : 'disappear'}" bind:this={chatCardElem}>
    <div class="usernameBody" style="background: {colorPalette.usernameBg}">
        <p style="color: {colorPalette.usernameFg}"><b>{chatMsg.username}</b></p>
    </div>
    <div class="messageBody" style="background: {colorPalette.msgBodyBg}">
        <p style="color: {colorPalette.msgBodyFg}">{chatMsg.msgText}</p>
    </div>
</div>