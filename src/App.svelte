<script>
	import ChatCard from "./ChatCard.svelte";

	let config = {};
	let chatMessages = new Map();
	
	let hiddenChatId;
	$: {
		chatMessages.delete(hiddenChatId);
		chatMessages = chatMessages;

	}
	
	function removeOldestChat() {
		chatMessages.delete([...chatMessages][0][0]);
		
	}

	function removeExcessChat(maxVisible) {
		if(chatMessages.size >= maxVisible) {
			removeOldestChat();
			
		}

	}
	
	function updateChat(chatMsgObj, addToTop, maxVisible) {
		removeExcessChat(maxVisible);
		addChatMsg(chatMsgObj, addToTop);
		chatMessages = chatMessages;

	}
	
	function connectStreamChat(websocketServerUrl) {
		let relaySocket = new WebSocket(websocketServerUrl);

		relaySocket.onopen = () => {
			relaySocket.send("I am the client, can you hear me?");
		}
		
		relaySocket.onmessage = (event) => {
			try {
				let responseObj = JSON.parse(event.data);

				if(!("username" in responseObj && "msgText" in responseObj)) {
					throw 'Required properties are missing for chatMsg object!';
				}
				
				updateChat(
					responseObj,
					config.cards.addToTop || false,
					config.cards.max || 20
				);

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
	
	// Add new chatMsg object to chatMsgs
	function addChatMsg(chatMsgObj) {
		chatMessages.set(getRandId(), chatMsgObj);

	}

	// Loading JSON config
	async function loadConfig()  {
		const response = await fetch("/config.json");
 		config = await response.json();
		return true;
	}
	
	// Generate a random id string
	function getRandId() {
		return Math.random().toString(36).substr(2, 9);

	}
	
	// Uh, main?
	function main() {
		connectStreamChat(config.serverWebSocketUrl || `ws://${window.location.hostname}:8888`);

	}
	
	// Load config.json, run main when finished
	let loadConfigPromise = loadConfig();
	loadConfigPromise.then(isConfigLoaded => main());
	
	</script>

<div class="container" style="height: 100%; width: 100%">
	{#await loadConfigPromise then isConfigLoaded}
		{#each (config.cards.addToTop ? [...chatMessages].reverse() : [...chatMessages]) as chatMsgArr (chatMsgArr[0])}
			<ChatCard
				{chatMsgArr}
				bind:hiddenChatId={hiddenChatId}
				removeAfter={config.cards.removeAfter}
			/>
		{:else}
			<p>No chat messages are available to display</p>
		{/each}

	{/await}
</div>