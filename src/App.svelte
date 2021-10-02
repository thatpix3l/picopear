<script>
	import ChatCard from "./ChatCard.svelte";

	let config = {};
	let chatMsgArr = [];

	function removeExcessChat(
		maxVisible = config.cards.max || false,
		isAddToTop = config.cards.addToTop || true
	) {
		if(chatMsgArr.length > maxVisible) {
			if(isAddToTop) {
				chatMsgArr = chatMsgArr.slice(0, maxVisible - 1);

			} else {
				chatMsgArr = chatMsgArr.slice(-maxVisible);
				
			}

		}
	}
	
	function updateChat(chatMsgObj, isAddToTop, maxVisible) {
		addChatMsg(chatMsgObj, isAddToTop);
		removeExcessChat(maxVisible, isAddToTop);

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
	function addChatMsg(chatMsgObj, isAddToTop){
		chatMsgObj.id = getRandId();

		if(isAddToTop) {
			chatMsgArr = [chatMsgObj].concat(chatMsgArr);

		} else {
			chatMsgArr = chatMsgArr.concat([chatMsgObj]);

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
		connectStreamChat(config.serverWebSocketUrl || `ws://${window.location.hostname}:8888`);

	}
	
	// Load config.json, run main when finished
	let loadConfigPromise = loadConfig();
	loadConfigPromise.then(isConfigLoaded => main());
	
	</script>

<div class="container" style="height: 100%; width: 100%">
	{#await loadConfigPromise then isConfigLoaded}
		{#if chatMsgArr.length === 0}
			<p>No chat messages are available to display: {chatMsgArr}</p>
		{:else}
			{#each chatMsgArr as chatMsg, chatMsgIndex (chatMsg.id)}
				<ChatCard
					bind:chatMsgArr={chatMsgArr}
					{chatMsgIndex}
					{chatMsg}
					colorPalettes={config.colorPalettes}
					removeAfter={config.cards.removeAfter}
				/>
			{/each}
		{/if}
	{/await}
</div>