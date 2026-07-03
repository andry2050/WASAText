<template>
	<div class="chat-container d-flex flex-column" style="height: 85vh;">
		
		<div class="d-flex justify-content-between align-items-center pb-2 border-bottom">
			<div class="d-flex align-items-center">
				<button class="btn btn-sm btn-outline-secondary me-3" @click="$router.push('/')">⬅ Indietro</button>
				
				<div class="rounded-circle me-2 bg-secondary text-white d-flex justify-content-center align-items-center" style="width: 40px; height: 40px; overflow: hidden;">
					<img v-if="chatInfo && chatInfo.photo_url" :src="getPhotoUrl(chatInfo.photo_url)" style="width: 100%; height: 100%; object-fit: cover;">
					<span v-else class="fs-4">{{ chatInfo && chatInfo.type === 'group' ? '👥' : '👤' }}</span>
				</div>
				<h2 class="h4 mb-0">{{ chatInfo ? chatInfo.name : 'Caricamento...' }}</h2>

                <button class="btn btn-sm btn-info text-white ms-3" @click="openInfoModal" title="Impostazioni">ℹ️ Info</button>
			</div>
		</div>

		<div class="messages-area flex-grow-1 overflow-auto my-3 p-2 bg-light border rounded">
			<div v-if="loading" class="text-center text-muted mt-3">Caricamento messaggi...</div>
			<div v-else-if="messages.length === 0" class="text-center text-muted mt-3">
				Nessun messaggio. Scrivi qualcosa o invia una foto per iniziare
			</div>
			
			<div class="message-list d-flex flex-column">
				<div 
					v-for="msg in messages" 
					:key="msg.id" 
					class="mb-3 p-2 border rounded w-75 position-relative pb-4 shadow-sm" 
					:style="msg.sender.id === myUserId ? 'background-color: #d1ecf1; border-color: #bee5eb;' : 'background-color: #ffffff;'"
					:class="msg.sender.id === myUserId ? 'align-self-end' : 'align-self-start'"
				>
					<div class="d-flex justify-content-between align-items-center border-bottom pb-1 mb-1" :style="msg.sender.id === myUserId ? 'border-bottom-color: rgba(0,0,0,0.1) !important;' : ''">
						<strong>{{ msg.sender.username }}</strong>
						
						<div class="d-flex align-items-center">
							<button @click="replyTo(msg)" class="btn btn-sm p-0 border-0 bg-transparent me-2" title="Rispondi">↩️</button>
							<button @click="openForwardModal(msg.id)" class="btn btn-sm p-0 border-0 bg-transparent me-2" title="Inoltra">↪️</button>
							<button @click="toggleEmojiPicker(msg.id)" class="btn btn-sm p-0 border-0 bg-transparent me-2" title="Reagisci">😀</button>

							<small class="text-muted ms-1">
								{{ new Date(msg.timestamp).toLocaleString(undefined, { month: 'short', day: 'numeric', hour: '2-digit', minute:'2-digit' }) }}
							</small>
							
							<button v-if="msg.sender.username === currentUsername" @click="deleteMessage(msg.id)" class="btn btn-sm text-danger p-0 border-0 bg-transparent ms-2" title="Elimina">🗑️</button>
						</div>
					</div>
					
					<div v-if="msg.content && isReply(msg.content)">
						<div class="py-1 px-2 mb-1 d-flex flex-column" style="background-color: rgba(255, 255, 255, 0.55); border-left: 4px solid #0056b3; border-radius: 4px;">
							<strong class="small" style="color: #0056b3;">{{ getReplyUsername(msg.content) }}</strong>
							<span class="text-muted small text-truncate">{{ getReplySnippet(msg.content) }}</span>
						</div>
						<div class="text-break mt-1">{{ getActualMessage(msg.content) }}</div>
					</div>
					<div v-else-if="!msg.is_photo && (!msg.content || !msg.content.includes('/uploads/'))" class="text-break">
						{{ msg.content }}
					</div>

					<div v-if="msg.is_photo || (msg.content && msg.content.includes('/uploads/'))" class="mt-2 text-center">
						<img :src="getPhotoUrl(msg.content)" alt="Foto" class="img-fluid rounded shadow-sm" style="max-height: 250px; object-fit: contain;">
					</div>

					<div class="position-absolute bottom-0 end-0 p-1 me-1" style="font-size: 0.85rem;">
						<span v-if="msg.sender.id === myUserId" :class="msg.status === 'read' ? 'text-primary' : 'text-muted'" style="font-weight: bold;">
							{{ msg.status === 'read' ? '✓✓' : '✓' }}
						</span>
					</div>

                    <div v-if="activeEmojiPickerMsgId === msg.id" class="mt-2 p-1 bg-white border rounded shadow-sm d-flex gap-1 justify-content-start">
						<button v-for="emoji in availableEmojis" :key="emoji" @click="addReaction(msg.id, emoji)" class="btn btn-light btn-sm fs-5 p-1 border-0">{{ emoji }}</button>
					</div>

					<div v-if="msg.reactions && msg.reactions.length > 0" class="mt-2 d-flex flex-wrap gap-1">
						<span 
							v-for="reaction in msg.reactions" 
							:key="reaction.id" 
							class="badge border text-dark p-1 d-flex align-items-center bg-light reaction-badge"
							:style="reaction.user.username === currentUsername ? 'cursor: pointer;' : ''"
							@click="reaction.user.username === currentUsername ? removeReaction(msg.id, reaction.id) : null"
							:title="'Aggiunta da ' + reaction.user.username"
						>
							<span class="fs-6">{{ reaction.emoji }}</span>
						</span>
					</div>
				</div>
			</div>
		</div>

		<div v-if="selectedFile" class="text-success small mb-1 px-2 d-flex justify-content-between">
			<span>📎 Foto selezionata: <strong>{{ selectedFile.name }}</strong></span>
			<button class="btn btn-sm btn-link text-danger p-0 text-decoration-none" @click="removeFile">❌ Rimuovi</button>
		</div>

		<div v-if="replyingToMsg" class="py-1 px-2 mb-1 d-flex justify-content-between align-items-center" style="background-color: #e3f2fd; border-left: 4px solid #0056b3; border-radius: 4px;">
			<span class="small text-truncate">
				<strong style="color: #0056b3;">{{ replyingToMsg.sender.username }}</strong><br>
				<span class="text-muted">{{ replyingToMsg.is_photo ? '📷 Foto' : replyingToMsg.content }}</span>
			</span>
			<button class="btn-close btn-sm" style="font-size: 0.5rem;" @click="replyingToMsg = null"></button>
		</div>

		<div class="input-group mt-auto">
			<input type="file" ref="fileInput" class="d-none" accept="image/*" @change="onFileSelected" />
			<button class="btn btn-outline-secondary" @click="triggerFileInput" title="Allega foto">📷</button>
			<input 
				v-model="newMessage" 
				ref="msgInput"
				type="text" 
				class="form-control" 
				placeholder="Scrivi un messaggio..." 
				@keyup.enter="sendMessage"
				:disabled="selectedFile !== null" 
			/>
			<button class="btn btn-primary" @click="sendMessage" :disabled="(!newMessage.trim() && !selectedFile) || sending">
				{{ sending ? '...' : 'Invia' }}
			</button>
		</div>

	</div>
    <div v-if="showForwardModal" class="modal-overlay d-flex justify-content-center align-items-center">
        <div class="bg-white p-4 rounded shadow-lg w-75" style="max-width: 400px;">
            <h3 class="h5 mb-3">Inoltra a...</h3>
            
            <div v-if="myConversations.length === 0" class="text-muted mb-3">
                Non hai altre chat.
            </div>
            
            <ul v-else class="list-group mb-3 overflow-auto" style="max-height: 250px;">
                <li 
                    v-for="chat in myConversations" 
                    :key="chat.id" 
                    class="list-group-item list-group-item-action" 
                    style="cursor: pointer;" 
                    @click="confirmForward(chat.id)"
                >
                    {{ chat.name }}
                </li>
            </ul>
            
            <button class="btn btn-secondary w-100" @click="closeForwardModal">Annulla</button>
        </div>
    </div>
    <div v-if="showInfoModal" class="modal-overlay d-flex justify-content-center align-items-center">
        <div class="bg-white p-4 rounded shadow-lg w-100" style="max-width: 450px; max-height: 90vh; overflow-y: auto;">
            <div class="d-flex justify-content-between align-items-center mb-3">
                <h3 class="h5 mb-0">{{ chatInfo.type === 'group' ? 'Impostazioni Gruppo' : 'Info Profilo' }}</h3>
                <button class="btn-close" @click="closeInfoModal"></button>
            </div>
            
            <template v-if="chatInfo.type === 'group'">
                <div class="mb-4">
                    <label class="form-label fw-bold small">Cambia Nome</label>
                    <div class="input-group">
                        <input v-model="newGroupName" type="text" class="form-control form-control-sm" placeholder="Nuovo nome...">
                        <button class="btn btn-sm btn-primary" @click="updateGroupName" :disabled="!newGroupName">Salva</button>
                    </div>
                </div>

                <div class="mb-4">
                    <label class="form-label fw-bold small">Cambia Foto</label>
                    <input type="file" ref="groupPhotoInput" class="d-none" accept="image/*" @change="onGroupPhotoSelected" />
                    <div class="d-flex gap-2">
                        <button class="btn btn-sm btn-outline-secondary w-100" @click="$refs.groupPhotoInput.click()">Scegli Foto...</button>
                        <button class="btn btn-sm btn-success w-100" v-if="selectedGroupPhoto" @click="updateGroupPhoto">Carica Foto</button>
                    </div>
                </div>

                <div class="mb-4 border-top pt-3">
                    <label class="form-label fw-bold small">Aggiungi Membro</label>
                    <div class="input-group mb-2">
                        <input v-model="searchUsername" type="text" class="form-control form-control-sm" placeholder="Username esatto...">
                        <button class="btn btn-sm btn-outline-primary" @click="addMemberToGroup" :disabled="!searchUsername">Aggiungi</button>
                    </div>
                </div>

                <div class="border-top pt-3 text-center">
                    <button class="btn btn-danger w-100" @click="leaveGroup">
                        🚪 Abbandona Gruppo
                    </button>
                </div>
            </template>

            <template v-else>
                <div class="text-center py-4">
                    <div class="rounded-circle mx-auto mb-3 bg-secondary text-white d-flex justify-content-center align-items-center" style="width: 120px; height: 120px; overflow: hidden;">
                        <img v-if="chatInfo.photo_url" :src="getPhotoUrl(chatInfo.photo_url)" style="width: 100%; height: 100%; object-fit: cover;">
                        <span v-else class="fs-1">👤</span>
                    </div>
                    <h4 class="mb-1">{{ chatInfo.name }}</h4>
                    <p class="text-muted small">Utente WASAText</p>
                </div>
            </template>
        </div>
    </div>
</template>

<script>
export default {
	data() {
    return {
			conversationId: this.$route.params.id,
			chat: null,
			chatInfo: null,
			myUserId: localStorage.getItem("token"),
			messages: [],
			newMessage: "",
			selectedFile: null,
			loading: false,
			sending: false,
			currentUsername: localStorage.getItem("username"),

			showForwardModal: false,
			forwardingMessageId: null,
			myConversations: [],

			activeEmojiPickerMsgId: null, 
			availableEmojis: ["👍", "❤️", "😂", "😮", "😢", "😡"],

			showInfoModal: false,
			newGroupName: "",
			selectedGroupPhoto: null,
			searchUsername: "",
			
			polling: null,       
			replyingToMsg: null  
		}
	},
	methods: {
		async loadMessages(showLoader = true) {
			if (showLoader) this.loading = true;
			try {
				// Segna i messaggi silensiosamente come letti
				await this.$axios.put(`/conversations/${this.conversationId}/read`).catch(()=>{});
				
				let response = await this.$axios.get(`/conversations/${this.conversationId}`);
				if (response.data && response.data.messages) {
					this.messages = response.data.messages.reverse(); 
				} else {
					this.messages = [];
				}
			} catch (e) {
				if (e.response && e.response.status === 404) { this.messages = []; }
			}
			if (showLoader) this.loading = false;
		},
		
		replyTo(msg) {
			this.replyingToMsg = msg;
			this.$refs.msgInput.focus(); // Porta il focus sulla tastiera
		},

		isReply(text) {
			return text && typeof text === 'string' && text.startsWith('[Risposta a ') && text.includes(']');
		},
		getReplyUsername(text) {
			let match = text.match(/\[Risposta a (.*?):/);
			return match ? match[1] : '';
		},
		getReplySnippet(text) {
			let match = text.match(/\[Risposta a .*?: (.*?)\]/);
			return match ? match[1] : '';
		},
		getActualMessage(text) {
			let idx = text.indexOf(']');
			let actual = text.substring(idx + 1);
			if (actual.startsWith('\n')) actual = actual.substring(1);
			return actual.trim();
		},

		triggerFileInput() { this.$refs.fileInput.click(); },
		
		onFileSelected(event) {
			const file = event.target.files[0];
			if (file) {
				this.selectedFile = file;
				this.newMessage = "";
			}
		},
		
		removeFile() {
			this.selectedFile = null;
			this.$refs.fileInput.value = "";
		},
		
		getPhotoUrl(path) {
			if (!path) return "";
			// Costruisce sempre l'URL corretto puntando al backend
			const baseUrl = this.$axios.defaults.baseURL || "http://localhost:3000";
			return path.startsWith("http") ? path : baseUrl + (path.startsWith("/") ? "" : "/") + path;
		},

        toggleEmojiPicker(messageId) {
			if (this.activeEmojiPickerMsgId === messageId) {
				this.activeEmojiPickerMsgId = null; 
			} else {
				this.activeEmojiPickerMsgId = messageId; 
			}
		},

        openInfoModal() {
			this.showInfoModal = true;
			this.newGroupName = "";
			this.searchUsername = "";
			this.selectedGroupPhoto = null;
		},
        
		closeInfoModal() {
			this.showInfoModal = false;
		},

		async sendMessage() {
			let text = this.newMessage.trim();
			if (!text && !this.selectedFile) return;

			if (this.replyingToMsg && text) {
				let snippet = this.replyingToMsg.is_photo ? "📷 Foto" : this.replyingToMsg.content;
				if (this.isReply(snippet)) snippet = this.getActualMessage(snippet);
				if (snippet.length > 50) snippet = snippet.substring(0, 50) + "...";
				text = `[Risposta a ${this.replyingToMsg.sender.username}: ${snippet}]\n${text}`;
			}

			this.sending = true;
			try {
				const formData = new FormData();
				if (this.selectedFile) formData.append("file", this.selectedFile);
				else formData.append("text", text);

				await this.$axios.post(`/conversations/${this.conversationId}/messages`, formData);
				
				this.newMessage = ""; 
				this.removeFile();
				this.replyingToMsg = null;
				this.loadMessages(); 
			} catch (e) {
				alert("Impossibile inviare il messaggio.");
			}
			this.sending = false;
		},

		async loadChat() {
			try {
				// Segna silenziosamente i messaggi come letti
				await this.$axios.put(`/conversations/${this.conversationId}/read`).catch(()=>{});
				
				let res = await this.$axios.get("/conversations");
				
				let conv = res.data.find(c => c.id === this.conversationId || (c.type === 'direct' && c.id.includes(this.conversationId)));
				
				if (conv) {
					this.chatInfo = conv;
				} else {
					// Fallback temporaneo finché non invia il primo messaggio
					this.chatInfo = { name: "Nuova Chat Privata", type: 'direct' };
				}
			} catch (e) {
				this.chatInfo = { name: "Nuova Chat Privata", type: 'direct' };
			}
		},

		async deleteMessage(messageId) {
			if (!confirm("Vuoi davvero eliminare questo messaggio?")) return;
			try {
				await this.$axios.delete(`/messages/${messageId}`);
				this.loadMessages();
			} catch (e) {
				alert("Impossibile eliminare il messaggio.");
			}
		},

		async openForwardModal(messageId) {
			this.forwardingMessageId = messageId;
			try {
				let response = await this.$axios.get("/conversations");
				this.myConversations = response.data;
				this.showForwardModal = true;
			} catch (e) {
				alert("Impossibile caricare le chat per l'inoltro.");
			}
		},

		closeForwardModal() {
			this.showForwardModal = false;
			this.forwardingMessageId = null;
		},

		async confirmForward(targetConversationId) {
			try {
				await this.$axios.post(`/conversations/${targetConversationId}/forward`, {
					message_id: this.forwardingMessageId
				});
				alert("Messaggio inoltrato con successo!");
				this.closeForwardModal();
			} catch (e) {
				alert("Errore durante l'inoltro del messaggio.");
			}
		},

        async addReaction(msgId, emoji) {
			// Trova il messaggio a cui si sta reagendo
			let msg = this.messages.find(m => m.id === msgId);
			if (!msg) return;

			// Controlla se è già inserito QUESTA STESSA emoji in questo messaggio
			let existingReaction = msg.reactions.find(r => r.user.username === this.currentUsername && r.emoji === emoji);

			if (existingReaction) {
				// Se lo è la rimuove
				await this.removeReaction(msgId, existingReaction.id);
				this.activeEmojiPickerMsgId = null; 
				return;
			}

			try {
				await this.$axios.post(`/messages/${msgId}/comments`, { emoji: emoji });
				this.activeEmojiPickerMsgId = null; 
				this.loadMessages(false); // Ricarica in background per far apparire la reazione
			} catch (e) {
				alert("Errore nell'aggiunta della reazione.");
			}
		},

        async removeReaction(msgId, reactionId) {
			try {
				await this.$axios.delete(`/messages/${msgId}/comments/${reactionId}`);
				this.loadMessages(false); // Ricarica i messaggi in background
			} catch (e) {
				alert("Impossibile rimuovere la reazione.");
			}
		},

        async updateGroupName() {
			try {
				await this.$axios.put(`/groups/${this.conversationId}/name`, { name: this.newGroupName.trim() });
				alert("Nome del gruppo aggiornato!");
				this.newGroupName = "";
				// Ricarica la chat se necessario
			} catch (e) {
				alert("Errore o permessi negati per cambiare nome.");
			}
		},

        onGroupPhotoSelected(event) {
			const file = event.target.files[0];
			if (file) this.selectedGroupPhoto = file;
		},
		async updateGroupPhoto() {
			try {
				const formData = new FormData();
				formData.append("file", this.selectedGroupPhoto);
				await this.$axios.put(`/groups/${this.conversationId}/photo`, formData);
				
				alert("Foto del gruppo aggiornata!");
				this.selectedGroupPhoto = null;
			} catch (e) {
				alert("Errore o permessi negati per cambiare foto.");
			}
		},

        async addMemberToGroup() {
			try {
				// Siccome l'API richiede l'ID utente, prima cerchiamo l'utente per nome
				let res = await this.$axios.get(`/users`, { params: { name: this.searchUsername.trim() } });
				
				if (!res.data || res.data.length === 0) {
					alert("Utente non trovato!");
					return;
				}
				
				// Prendiamo l'ID del primo utente trovato con quel nome
				const targetUserId = res.data[0].id; 

				// Aggiungiamolo al gruppo
				await this.$axios.post(`/groups/${this.conversationId}/members`, { user_id: targetUserId });
				alert(`Utente ${this.searchUsername} aggiunto al gruppo!`);
				this.searchUsername = "";
				
			} catch (e) {
				alert("Errore durante l'aggiunta. Assicurati di essere in un gruppo e non in una chat singola.");
			}
		},

        async leaveGroup() {
			if (!confirm("Sei sicuro di voler abbandonare questo gruppo?")) return;
			
			try {
				// Leggiamo il nostro ID dal localStorage (che avevi salvato al login)
				const myUserId = localStorage.getItem("token");
				
				await this.$axios.delete(`/groups/${this.conversationId}/members/${myUserId}`);
				
				alert("Hai abbandonato il gruppo.");
				this.closeInfoModal();
				
				// Ti riporto alla Home perché non sei più nel gruppo!
				this.$router.push("/");
				
			} catch (e) {
				alert("Impossibile abbandonare (forse è una chat privata e non un gruppo?).");
			}
		}
	},
	mounted() {
		this.loadChat();
		this.loadMessages(true);
		
		
		this.polling = setInterval(() => {
			this.loadMessages(false);
			this.loadChat(); 
		}, 2000);
	},
	beforeUnmount() {
		if (this.polling) clearInterval(this.polling);
	}
}
</script>

<style scoped>
.message-list {
	display: flex;
	flex-direction: column;
}
/* Disabilita l'input di testo visivamente se c'è un file selezionato */
input:disabled {
	background-color: #e9ecef;
	cursor: not-allowed;
}
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100vw;
	height: 100vh;
	background: rgba(0, 0, 0, 0.5); 
	z-index: 1050;
}
</style>