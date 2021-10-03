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
    
    // Chat messages array
    export let chatMsgMap;
    
    // Array containing id and a chatMsg object
	export let chatMsgArr; 
    
    // How long a card is allowed to be visible, in milliseconds
    export let removeAfter = 5000;
    
    // Object of keys of colors
    export let colorPalettes;

    // All chat cards start off as visible
    let chatCardVisible = true;

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
    
    function autoRemove(removeAfterTime, hideDuration) {
        setTimeout(() => {
            chatCardVisible = false;

        }, removeAfterTime - hideDuration);
        
        setTimeout(() => {
            chatMsgMap.delete(chatMsgArr[0]);

        }, removeAfterTime);

    }
    

    // Beginning chat card lifecycle
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
        let disappearDuration = toMillisecond(window.getComputedStyle(chatCardElem).getPropertyValue("--disappear__duration"));
        autoRemove(removeAfter, disappearDuration);
    });
    
</script>

<div id="{chatMsgArr[0]}" class="card {chatCardVisible === true ? 'appear' : 'disappear'}" bind:this={chatCardElem}>
    <div class="usernameBody" style="background: {colorPalette.usernameBg}">
        <p style="color: {colorPalette.usernameFg}"><b>{chatMsgArr[1].username}</b></p>
    </div>
    <div class="messageBody" style="background: {colorPalette.msgBodyBg}">
        <p style="color: {colorPalette.msgBodyFg}">{chatMsgArr[1].msgText}</p>
    </div>
</div>