<template>
	<div class="profile-container container mt-4" style="max-width: 600px;">
		
		<div class="d-flex justify-content-between align-items-center pb-2 border-bottom mb-4">
			<h2 class="h3 mb-0">Il Mio Profilo</h2>
			<button class="btn btn-sm btn-outline-secondary" @click="$router.push('/')">
				⬅ Torna alla Home
			</button>
		</div>

		<div v-if="successMsg" class="alert alert-success">{{ successMsg }}</div>
		<div v-if="errorMsg" class="alert alert-danger">{{ errorMsg }}</div>

		<div class="card mb-4 shadow-sm">
			<div class="card-body">
				<h5 class="card-title">Cambia Username</h5>
				<p class="text-muted small">Il tuo username attuale è: <strong>{{ currentUsername }}</strong></p>
				
				<div class="input-group">
					<input 
						v-model="newUsername" 
						type="text" 
						class="form-control" 
						placeholder="Nuovo username (3-16 caratteri)" 
					/>
					<button 
						class="btn btn-primary" 
						@click="changeUsername" 
						:disabled="!newUsername || newUsername === currentUsername || loadingName"
					>
						{{ loadingName ? 'Salvataggio...' : 'Salva Nome' }}
					</button>
				</div>
			</div>
		</div>

		<div class="card shadow-sm">
			<div class="card-body">
				<h5 class="card-title">Cambia Foto Profilo</h5>
				<p class="text-muted small">Carica un'immagine per farti riconoscere nelle chat.</p>
				
				<input type="file" ref="photoInput" class="d-none" accept="image/*" @change="onPhotoSelected" />
				
				<div class="d-flex align-items-center gap-3">
					<button class="btn btn-outline-secondary" @click="triggerPhotoInput">
						Scegli Immagine...
					</button>
					
					<span v-if="selectedPhoto" class="text-success small">
						📎 {{ selectedFile.name }}
					</span>
				</div>

				<button 
					v-if="selectedPhoto" 
					class="btn btn-success mt-3 w-100" 
					@click="uploadPhoto"
					:disabled="loadingPhoto"
				>
					{{ loadingPhoto ? 'Caricamento in corso...' : 'Carica e Salva Foto' }}
				</button>
			</div>
		</div>

	</div>
</template>

<script>
export default {
	data() {
		return {
			currentUsername: localStorage.getItem("username"),
			newUsername: "",
			
			selectedFile: null,
			selectedPhoto: false, // Flag per sapere se c'è un file pronto

			loadingName: false,
			loadingPhoto: false,

			successMsg: "",
			errorMsg: ""
		}
	},
	methods: {
		// --- CAMBIO NOME ---
		async changeUsername() {
			this.errorMsg = "";
			this.successMsg = "";
			
			const cleanName = this.newUsername.trim();
			if (cleanName.length < 3 || cleanName.length > 16) {
				this.errorMsg = "Il nome deve essere compreso tra 3 e 16 caratteri.";
				return;
			}

			this.loadingName = true;
			try {
				await this.$axios.put("/users/me/username", {
					name: cleanName
				});
				
				// Se ha successo, aggiorniamo il localStorage e l'interfaccia!
				localStorage.setItem("username", cleanName);
				this.currentUsername = cleanName;
				this.newUsername = "";
				this.successMsg = "Username aggiornato con successo!";
				
			} catch (e) {
				// Controlliamo se il backend ci ha risposto con 409 Conflict
				if (e.response && e.response.status === 409) {
					this.errorMsg = "Questo nome utente è già in uso da qualcun altro!";
				} else {
					this.errorMsg = "Errore durante il cambio del nome.";
				}
			}
			this.loadingName = false;
		},

		// --- CAMBIO FOTO ---
		triggerPhotoInput() {
			this.$refs.photoInput.click();
		},
		
		onPhotoSelected(event) {
			const file = event.target.files[0];
			if (file) {
				this.selectedFile = file;
				this.selectedPhoto = true;
				this.successMsg = "";
				this.errorMsg = "";
			}
		},
		
		async uploadPhoto() {
			this.errorMsg = "";
			this.successMsg = "";
			this.loadingPhoto = true;
			
			try {
				const formData = new FormData();
				formData.append("file", this.selectedFile);

				await this.$axios.put("/users/me/photo", formData);
				
				this.successMsg = "Foto profilo aggiornata con successo!";
				this.selectedFile = null;
				this.selectedPhoto = false;
				this.$refs.photoInput.value = ""; // Resetta l'input file
				
			} catch (e) {
				this.errorMsg = "Errore durante il caricamento della foto.";
				console.error(e);
			}
			this.loadingPhoto = false;
		}
	}
}
</script>

<style scoped>
.profile-container {
	background-color: #f8f9fa;
	padding: 2rem;
	border-radius: 8px;
}
</style>