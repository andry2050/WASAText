<template>
	<div class="chat-container d-flex flex-column" style="height: 85vh;">
		
		<div class="d-flex justify-content-between align-items-center pb-2 border-bottom">
			<div class="d-flex align-items-center">
				<button class="btn btn-sm btn-outline-secondary me-3" @click="$router.push('/')">
					⬅ Indietro
				</button>
				
				<img v-if="chatInfo && chatInfo.photo_url" :src="getPhotoUrl(chatInfo.photo_url)" class="rounded-circle me-2" style="width: 40px; height: 40px; object-fit: cover;">
				<h2 class="h4 mb-0">{{ chatInfo ? chatInfo.name : 'Caricamento...' }}</h2>

                <button class="btn btn-sm btn-info text-white ms-3" @click="openInfoModal" title="Impostazioni">
					ℹ️ Info
				</button>
			</div>
		</div>

		<div class="messages-area flex-grow-1 overflow-auto my-3 p-2 bg-light border rounded">
			<div v-if="loading" class="text-center text-muted mt-3">Caricamento messaggi...</div>
			<div v-else-if="messages.length === 0" class="text-center text-muted mt-3">
				Nessun messaggio. Scrivi qualcosa o invia una foto per iniziare
			</div>
			
			<div v-else class="message-list d-flex flex-column">
				<div 
					v-for="msg in messages" 
					:key="msg.id" 
					class="mb-3 p-2 border rounded"
					:class="msg.sender === currentUsername ? 'bg-success bg-opacity-10 align-self-end w-75' : 'bg-white align-self-start w-75'"
				>
					<div class="d-flex justify-content-between align-items-center border-bottom pb-1 mb-1">
						<strong>{{ msg.sender }}</strong>
						
						<div>
							<small class="text-muted me-2">{{ new Date(msg.timestamp).toLocaleString() }}</small>
							
							<small v-if="msg.sender_id === myUserId" :class="msg.status === 'read' ? 'text-primary' : 'text-muted'">
									{{ msg.status === 'read' ? '✓✓' : '✓' }}
							</small>
							<button 
								@click="toggleEmojiPicker(msg.id)" 
								class="btn btn-sm text-warning p-0 border-0 bg-transparent me-2" 
								title="Reagisci"
							>
								😀
							</button>

							<button 
								@click="openForwardModal(msg.id)" 
								class="btn btn-sm text-primary p-0 border-0 bg-transparent me-2" 
								title="Inoltra"
							>
								↪️
							</button>

							<button 
								v-if="msg.sender === currentUsername" 
								@click="deleteMessage(msg.id)" 
								class="btn btn-sm text-danger p-0 border-0 bg-transparent" 
								title="Elimina"
							>
								🗑️
							</button>
						</div>
					</div>
					
					<div v-if="!msg.is_photo">
						{{ msg.content }}
					</div>
					<div v-else class="text-center mt-2">
						<img :src="getPhotoUrl(msg.content)" alt="Foto" class="img-fluid rounded shadow-sm" style="max-height: 250px; object-fit: contain;">
					</div>

					<div v-if="activeEmojiPickerMsgId === msg.id" class="mt-2 p-1 bg-white border rounded shadow-sm d-flex gap-1 justify-content-start">
						<button 
							v-for="emoji in availableEmojis" 
							:key="emoji" 
							@click="addReaction(msg.id, emoji)"
							class="btn btn-light btn-sm fs-5 p-1 border-0"
						>
							{{ emoji }}
						</button>
					</div>

					<div v-if="msg.reactions && msg.reactions.length > 0" class="mt-2 d-flex flex-wrap gap-1">
						<span 
							v-for="reaction in msg.reactions" 
							:key="reaction.reactionid" 
							class="badge border text-dark p-1 d-flex align-items-center"
							:class="reaction.user.name === currentUsername ? 'bg-primary bg-opacity-25' : 'bg-light'"
							:style="reaction.user.name === currentUsername ? 'cursor: pointer;' : ''"
							@click="reaction.user.name === currentUsername ? removeReaction(msg.id, reaction.reactionid) : null"
							:title="'Inserita da ' + reaction.user.name + (reaction.user.name === currentUsername ? ' (Clicca per rimuovere)' : '')"
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

		<div class="input-group mt-auto">
			
			<input type="file" ref="fileInput" class="d-none" accept="image/*" @change="onFileSelected" />
			
			<button class="btn btn-outline-secondary" @click="triggerFileInput" title="Allega foto">
				📷
			</button>
			
			<input 
				v-model="newMessage" 
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
                <h3 class="h5 mb-0">Impostazioni Gruppo</h3>
                <button class="btn-close" @click="closeInfoModal"></button>
            </div>
            
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
                <small v-if="selectedGroupPhoto" class="text-success d-block mt-1">📎 {{ selectedGroupPhoto.name }}</small>
            </div>

            <div class="mb-4 border-top pt-3">
                <label class="form-label fw-bold small">Aggiungi Membro (cerca username)</label>
                <div class="input-group mb-2">
                    <input v-model="searchUsername" type="text" class="form-control form-control-sm" placeholder="Username esatto...">
                    <button class="btn btn-sm btn-outline-primary" @click="addMemberToGroup" :disabled="!searchUsername">Aggiungi</button>
                </div>
            </div>

            <div class="border-top pt-3 text-center">
                <p class="text-muted small mb-2">Se esci, non potrai più leggere o inviare messaggi qui.</p>
                <button class="btn btn-danger w-100" @click="leaveGroup">
                    🚪 Abbandona Gruppo
                </button>
            </div>
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
		}
	},
	methods: {
		async loadMessages() {
			this.loading = true;
			try {
				let response = await this.$axios.get(`/conversations/${this.conversationId}`);
				this.messages = response.data.messages.reverse(); 
			} catch (e) {
				console.error("Errore caricamento messaggi:", e);
				alert("Impossibile caricare i messaggi.");
			}
			this.loading = false;
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
			const baseURL = this.$axios.defaults.baseURL || "";
			return baseURL + path;
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
			const text = this.newMessage.trim();
			if (!text && !this.selectedFile) return;

			this.sending = true;
			try {
				const formData = new FormData();
				if (this.selectedFile) {
					formData.append("file", this.selectedFile);
				} else {
					formData.append("text", text);
				}

				await this.$axios.post(`/conversations/${this.conversationId}/messages`, formData);
				this.newMessage = ""; 
				this.removeFile();
				this.loadMessages(); 
			} catch (e) {
				alert("Impossibile inviare il messaggio.");
			}
			this.sending = false;
		},

		async loadChat() {
			try {
				// Segna come letti
				await this.$axios.put(`/conversations/${this.conversationId}/read`);
				
				// Cerca il nome e la foto della chat dalla lista generale
				let res = await this.$axios.get("/conversations");
				let conv = res.data.find(c => c.id === this.conversationId);
				if (conv) {
					this.chatInfo = conv;
				}
			} catch (e) {
				console.error("Errore caricamento info chat");
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

        async addReaction(messageId, emoji) {
			try {
				await this.$axios.post(`/messages/${messageId}/comments`, {
					emoji: emoji
				});
				this.activeEmojiPickerMsgId = null; 
				this.loadMessages(); 
			} catch (e) {
				console.error("Errore aggiunta reazione:", e);
				alert("Impossibile aggiungere la reazione.");
			}
		},

        async removeReaction(messageId, reactionId) {
			try {
				await this.$axios.delete(`/messages/${messageId}/comments/${reactionId}`);
				this.loadMessages(); // Ricarichiamo per vederla sparire
			} catch (e) {
				console.error("Errore rimozione reazione:", e);
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
		this.loadMessages();
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