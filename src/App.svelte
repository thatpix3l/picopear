<script>
	import ChatCard from "./ChatCard.svelte";

	let chatMsgs = [
		{
			username: "Brian",
			msgText: "Bread and butter, best combo ever",
			id: "eeeeeee"
		},
		{
			username: "John Wick",
			msgText: "You killed my dog and stole my car. The least I could do was repay the favor tenfold.",
			id: "bbbbbbb"
		},
		{
			username: "James Bond",
			msgText: "The name is Bond. James Bond.",
			id: "ccccccc"
		},
		{
			username: "Luke Skywalker",
			msgText: "The force is strong in my family. My father has it. I have it. My sister has it.",
			id: "ooooooo"
		}
	]
	
	let cardWidth = 360;
	
	let maxChatVisible = 20

	const autoClearChat = (delay) => {
		setTimeout(() => {
			if(chatMsgs.length > maxChatVisible) {
			}
		}, delay);
	}
	
	const connectStreamChat = () => {
		let serverHostname = window.location.hostname;
		let relaySocket = new WebSocket(`ws://${serverHostname}:8888`);

		relaySocket.onopen = () => {
			relaySocket.send("I am the client, can you hear me?");
		}
		
		relaySocket.onmessage = (event) => {
			try {
				let responseObj = JSON.parse(event.data);
				console.log(responseObj);

				const hasProp = (prop) => {
					prop in responseObj;
				}

				if(["username", "msgText"].every(hasProp)) {
					addChatMsg(responseObj);

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
				connectStreamChat();
			}, 1000);
		}
	}
	
	const addChatMsg = (chatMsgObj) => {
		chatMsgObj.id = getRandId();
		chatMsgs = [...chatMsgs, chatMsgs];
	}

	// Func for loading json array of objects containing keys with color values
	const loadColorPalettes = () => {
		return fetch("/color_palettes.json").then((response) => {
			return response.json();
		});
	}
	
	const getRandId = () => {
		// For identifying unique messages in session
		return Math.random().toString(36).substr(2, 9);
	}
	
	connectStreamChat();
	
	</script>

<div class="container">
	{#await loadColorPalettes() then colorPalettes}
		{#if chatMsgs.length === 0}
			<p>No chat messages are available to display</p>
		{:else}
			{#each chatMsgs as chatMsg}
				<ChatCard {chatMsg} {cardWidth} colorPalettes={colorPalettes}/>
			{/each}
		{/if}
	{/await}
</div>