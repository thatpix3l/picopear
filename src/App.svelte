<script>
	import ChatCard from "./ChatCard.svelte";
	

	let chatMsgs = [];
	let config = {};
	
	function autoClearExcessChat(delay, maxVisible) {
		setInterval(() => {
			if(chatMsgs.length > maxVisible) {
				chatMsgs.splice(0, chatMsgs.length - maxVisible);
				chatMsgs = chatMsgs;
			}
		}, delay);
	}
	
	function connectStreamChat(websocketServerUrl) {
			let relaySocket = new WebSocket(websocketServerUrl);

		relaySocket.onopen = () => {
			relaySocket.send("I am the client, can you hear me?");
		}
		
		relaySocket.onmessage = (event) => {
			try {
				let responseObj = JSON.parse(event.data);

				if("username" in responseObj && "msgText" in responseObj) {
					addChatMsg(responseObj, config.cards.addToTop || false);

				} else {
					throw 'Required properties are missing for chatMsg object!';

				}

			} catch(e) {
				console.log(e);

			}
			
		}

		relaySocket.onclose = (event) => {
			console.log("Socket closed. Reconnect will be attempted in 1 second.", event.reason);
			setTimeout(() => {
				connectStreamChat(websocketServerUrl);
			}, 1000);
		}
	}
	
	// Add new chatMsg object to end of chatMsgs
	function addChatMsg(chatMsgObj, addToTop){
		chatMsgObj.id = getRandId();

		if(addToTop) {
			chatMsgs.unshift(chatMsgObj);
			chatMsgs = chatMsgs;

		} else {
			chatMsgs[chatMsgs.length] = chatMsgObj;

		}

	}

	// Loading JSON config
	async function loadConfig()  {
		const response = await fetch("/config.json");
 		config = await response.json();
		return true;
	}
	
	// Generate a random id
	function getRandId() {
		return Math.random().toString(36).substr(2, 9);

	}
	
	// Uh, main?
	function main() {
		autoClearExcessChat(3000, 30);
		connectStreamChat(config.streamChatUrl || `ws://${window.location.hostname}:8888`);
		console.log(config)

	}
	
	let loadConfigPromise = loadConfig();
	loadConfigPromise.then(isConfigLoaded => main());
	
	</script>

<div class="container">
	{#await loadConfigPromise then isConfigLoaded}
		{#if chatMsgs.length === 0}
			<p>No chat messages are available to display: {chatMsgs}</p>
		{:else}
			{#each chatMsgs as chatMsg (chatMsg.id)}
				<ChatCard {chatMsg} colorPalettes={config.colorPalettes}/>
			{/each}
		{/if}
	{/await}
</div>