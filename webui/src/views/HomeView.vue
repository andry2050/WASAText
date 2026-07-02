<template>
	<div class="container mt-3" style="max-width: 800px;">
		
		<div class="d-flex justify-content-between align-items-center pb-2 mb-3 border-bottom">
			<h1 class="h3 mb-0">Le mie Chat</h1>
			
			<div>
				<button class="btn btn-primary me-2" @click="openChatModal">
					👤 Nuova Chat
				</button>
				<button class="btn btn-success" @click="openGroupModal">
					👥 Nuovo Gruppo
				</button>
			</div>
		</div>

		<div v-if="errormsg" class="alert alert-danger">{{ errormsg }}</div>

		<div v-if="loading" class="text-center text-muted">Caricamento conversazioni...</div>

		<div v-else-if="conversations.length === 0" class="text-center text-muted mt-5">
			Nessuna conversazione trovata. Cerca un utente o crea un gruppo per iniziare a chattare.
		</div>
		
		<div v-else class="list-group">
			<div 
				v-for="chat in conversations" 
				:key="chat.id" 
				class="list-group-item list-group-item-action d-flex align-items-center" 
				style="cursor: pointer;" 
				@click="$router.push('/chat/' + chat.id)"
			>
				<img v-if="chat.photo_url" :src="$axios.defaults.baseURL + chat.photo_url" class="rounded-circle me-3" style="width: 50px; height: 50px; object-fit: cover;">

				<img v-if="chat.photo_url" :src="getPhotoUrl(chat.photo_url)" class="rounded-circle me-3" style="width: 50px; height: 50px; object-fit: cover;">
				
				<div v-else class="me-3 fs-3">
					{{ chat.type === 'group' ? '👥' : '👤' }}
				</div>

				<div>
					<strong class="mb-1 d-block">{{ chat.name }}</strong> 
					<small class="text-muted">{{ chat.last_message_preview || 'Nessun messaggio' }}</small>
				</div>
			</div>
		</div>

		<div v-if="showGroupModal" class="modal-overlay d-flex justify-content-center align-items-center">
			<div class="bg-white p-4 rounded shadow-lg w-100" style="max-width: 500px;">
				<h3 class="h5 mb-3">Crea un Nuovo Gruppo</h3>
				
				<div class="mb-3">
					<label class="form-label fw-bold">Nome del Gruppo</label>
					<input v-model="newGroupName" type="text" class="form-control" placeholder="Es. Compagni di Università">
				</div>

				<div class="mb-3 border-top pt-3">
					<label class="form-label fw-bold">Aggiungi Membri</label>
					<div class="input-group mb-2">
						<input v-model="searchQuery" type="text" class="form-control" placeholder="Cerca username..." @keyup.enter="searchUsers">
						<button class="btn btn-outline-primary" @click="searchUsers">Cerca</button>
					</div>
					
					<ul v-if="searchResults.length > 0" class="list-group mb-2" style="max-height: 150px; overflow-y: auto;">
						<li v-for="user in searchResults" :key="user.id" class="list-group-item d-flex justify-content-between align-items-center p-2">
							{{ user.username }}
							<button 
								class="btn btn-sm btn-outline-success" 
								@click="addMemberToSelection(user)"
								:disabled="isUserSelected(user.id)"
							>
								{{ isUserSelected(user.id) ? 'Aggiunto' : '+ Aggiungi' }}
							</button>
						</li>
					</ul>
				</div>

				<div v-if="selectedMembers.length > 0" class="mb-3 p-2 bg-light rounded">
					<span class="d-block small fw-bold mb-1">Membri selezionati:</span>
					<div class="d-flex flex-wrap gap-1">
						<span v-for="member in selectedMembers" :key="member.id" class="badge bg-primary text-white d-flex align-items-center p-2">
							{{ member.username }}
							<button class="btn-close btn-close-white ms-2" style="font-size: 0.5rem;" @click="removeMemberFromSelection(member.id)"></button>
						</span>
					</div>
				</div>

				<div class="d-flex justify-content-end gap-2 mt-4">
					<button class="btn btn-secondary" @click="closeGroupModal">Annulla</button>
					<button class="btn btn-success" @click="createGroup" :disabled="!newGroupName.trim() || creatingGroup">
						{{ creatingGroup ? 'Creazione...' : 'Crea Gruppo' }}
					</button>
				</div>
			</div>
		</div>

		<div v-if="showChatModal" class="modal-overlay d-flex justify-content-center align-items-center">
			<div class="bg-white p-4 rounded shadow-lg w-100" style="max-width: 500px;">
				<h3 class="h5 mb-3">Cerca un utente</h3>
				<div class="input-group mb-3">
					<input v-model="searchQuery" type="text" class="form-control" placeholder="Cerca username..." @keyup.enter="searchUsers">
					<button class="btn btn-primary" @click="searchUsers">Cerca</button>
				</div>
				
				<ul class="list-group mb-3" style="max-height: 200px; overflow-y: auto;">
					<li v-for="user in searchResults" :key="user.id" class="list-group-item list-group-item-action d-flex justify-content-between align-items-center" style="cursor: pointer;" @click="$router.push('/chat/' + user.id)">
						{{ user.username }}
						<span class="text-primary small">💬 Scrivi</span>
					</li>
				</ul>
				
				<div class="d-flex justify-content-end">
					<button class="btn btn-secondary" @click="showChatModal = false">Chiudi</button>
				</div>
			</div>
		</div>

	</div>
</template>

<style scoped>
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

<script>
export default {
	data() {
		return {
			conversations: [],
			errormsg: null,
			loading: false,

			showChatModal: false, 
			showGroupModal: false,
			newGroupName: "",
			searchQuery: "",
			searchResults: [],
			selectedMembers: [], 
			creatingGroup: false,
		}
	},
	methods: {
		async loadConversations() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/conversations");
				this.conversations = response.data;
			} catch (e) {
				this.errormsg = "Errore nel caricamento delle chat.";
			}
			this.loading = false;
		},

		openChatModal() {
			this.showChatModal = true;
			this.searchQuery = "";
			this.searchResults = [];
		},

		openGroupModal() {
			this.showGroupModal = true;
			this.newGroupName = "";
			this.searchQuery = "";
			this.searchResults = [];
			this.selectedMembers = [];
		},
		
		closeGroupModal() {
			this.showGroupModal = false;
		},

		async searchUsers() {
			if (!this.searchQuery.trim()) return;
			try {
				let response = await this.$axios.get(`/users`, { params: { username: this.searchQuery.trim() } });
				this.searchResults = response.data;
			} catch (e) {
				console.error("Errore ricerca utenti:", e);
			}
		},

		isUserSelected(userId) {
			return this.selectedMembers.some(m => m.id === userId);
		},

		addMemberToSelection(user) {
			if (!this.isUserSelected(user.id)) {
				this.selectedMembers.push(user);
			}
		},

		removeMemberFromSelection(userId) {
			this.selectedMembers = this.selectedMembers.filter(m => m.id !== userId);
		},

		getPhotoUrl(path) {
			if (!path) return "";
			const baseUrl = this.$axios.defaults.baseURL || "http://localhost:3000";
			return path.startsWith("http") ? path : baseUrl + (path.startsWith("/") ? "" : "/") + path;
		},
		
		async createGroup() {
			if (!this.newGroupName.trim()) return;
			this.creatingGroup = true;

			try {
				// ERRORE FIXATO: Devi includere ANCHE IL TUO ID nei membri del gruppo!
				const myUserId = localStorage.getItem("token");
				const memberIds = [myUserId, ...this.selectedMembers.map(m => m.id)];

				const response = await this.$axios.post("/groups", {
					name: this.newGroupName.trim(),
					members: memberIds
				});

				const newGroupId = response.data.id;
				this.closeGroupModal();
				this.$router.push('/chat/' + newGroupId);

			} catch (e) {
				alert("Errore durante la creazione del gruppo.");
				console.error(e);
			}
			this.creatingGroup = false;
		},

		mounted() {
			this.loadConversations();
		}
	}
}
</script>