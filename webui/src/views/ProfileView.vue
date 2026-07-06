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

		<div class="text-center mb-4">
			<img v-if="photo_url" :src="getPhotoUrl(photo_url)" class="rounded-circle shadow" style="width: 150px; height: 150px; object-fit: cover;" />
			<div v-else class="rounded-circle bg-secondary d-flex justify-content-center align-items-center mx-auto text-white fs-1 shadow" style="width: 150px; height: 150px;">
				👤
			</div>
		</div>

		<div class="card mb-4 shadow-sm">
			<div class="card-body">
				<h5 class="card-title">Cambia Username</h5>
				<p class="text-muted small">Il tuo username attuale è: <strong>{{ currentUsername }}</strong></p>
				
				<div class="input-group">
					<input v-model="newUsername" type="text" class="form-control" placeholder="Nuovo username" />
					<button class="btn btn-primary" @click="changeUsername" :disabled="!newUsername || newUsername === currentUsername || loadingName">
						{{ loadingName ? 'Salvataggio...' : 'Salva Nome' }}
					</button>
				</div>
			</div>
		</div>

		<div class="card shadow-sm">
			<div class="card-body">
				<h5 class="card-title">Cambia Foto Profilo</h5>
				
				<input type="file" ref="photoInput" class="d-none" accept="image/*" @change="onPhotoSelected" />
				
				<div class="d-flex align-items-center gap-3">
					<button class="btn btn-outline-secondary" @click="triggerPhotoInput">Scegli Immagine...</button>
					<span v-if="selectedPhoto" class="text-success small">📎 {{ selectedFile.name }}</span>
				</div>

				<button v-if="selectedPhoto" class="btn btn-success mt-3 w-100" @click="uploadPhoto" :disabled="loadingPhoto">
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
			photo_url: null, // Aggiunto per salvare la foto
			selectedFile: null,
			selectedPhoto: false,
			loadingName: false,
			loadingPhoto: false,
			successMsg: "",
			errorMsg: ""
		}
	},
	methods: {
		getPhotoUrl(path) {
			if (!path) return "";
			const baseUrl = this.$axios.defaults.baseURL || "http://localhost:3000";
			// Assicura che path non abbia doppie slash e si componga bene
			const cleanPath = path.startsWith("/") ? path : "/" + path;
			return path.startsWith("http") ? path : baseUrl + cleanPath;
		},

		async loadMyProfile() {
			try {
				// Cerca il tuo stesso utente per ottenere la foto
				let res = await this.$axios.get("/users", { params: { username: this.currentUsername } });
				if (res.data && res.data.length > 0) {
					this.photo_url = res.data[0].photo_url;
				}
			} catch (e) { console.error("Errore caricamento profilo", e); }
		},
		async changeUsername() {
			this.errorMsg = ""; this.successMsg = "";
			const cleanName = this.newUsername.trim();
			if (cleanName.length < 3 || cleanName.length > 16) return;
			this.loadingName = true;
			try {
				await this.$axios.put("/users/me/username", { name: cleanName });
				localStorage.setItem("username", cleanName);
				this.currentUsername = cleanName; this.newUsername = "";
				this.successMsg = "Username aggiornato con successo!";
			} catch (e) {
				this.errorMsg = "Errore durante il cambio del nome o nome già in uso.";
			}
			this.loadingName = false;
		},
		triggerPhotoInput() { this.$refs.photoInput.click(); },
		onPhotoSelected(event) {
			const file = event.target.files[0];
			if (file) { this.selectedFile = file; this.selectedPhoto = true; }
		},
		async uploadPhoto() {
			this.loadingPhoto = true;
			try {
				const formData = new FormData();
				formData.append("file", this.selectedFile);
				await this.$axios.put("/users/me/photo", formData);
				
				this.successMsg = "Foto profilo aggiornata con successo!";
				this.selectedFile = null; this.selectedPhoto = false;
				this.$refs.photoInput.value = ""; 
				this.loadMyProfile(); // Ricarica la foto appena caricata
			} catch (e) {
				this.errorMsg = "Errore durante il caricamento della foto.";
			}
			this.loadingPhoto = false;
		}
	},
	mounted() {
		this.loadMyProfile(); // Carica la foto all'avvio
	}
}
</script>