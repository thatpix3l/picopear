<script>
	import ChatCard from "./ChatCard.svelte";

	let config = {};
	let chatMsgMap = new Map();

	function removeExcessChat( maxVisible, isAddToTop) {
		if(chatMsgMap.size > maxVisible) {
			if(isAddToTop) {
				chatMsgMap.delete([...chatMsgMap][chatMsgMap.size - 1][0]);

				
			} else {
				chatMsgMap.delete([...chatMsgMap][0][0]);
				
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
		let id = getRandId();

		if(isAddToTop) {
			chatMsgMap = new Map([[id, chatMsgObj], ...chatMsgMap]);

		} else {
			chatMsgMap = new Map([...chatMsgMap, [id, chatMsgObj]]);

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
		{#each [...chatMsgMap] as chatMsgArr (chatMsgArr[0])}
			<ChatCard
				bind:chatMsgMap={chatMsgMap}
				{chatMsgArr}
				removeAfter={config.cards.removeAfter}
			/>
		{:else}
			<p>No chat messages are available to display</p>
		{/each}
	{/await}
</div>