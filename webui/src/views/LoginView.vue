<template>
	<div class="login-page">
		<div class="login-box">
			<h1>WASAText 💬</h1>
			<p>Inserisci il tuo username per iniziare a chattare.</p>
			
			<input 
				v-model="username" 
				type="text" 
				placeholder="Username (3-16 caratteri)" 
				@keyup.enter="login"
			/>
			
			<button @click="login" :disabled="isLoggingIn">
				{{ isLoggingIn ? 'Accesso in corso...' : 'Entra' }}
			</button>

			<p v-if="errorMessage" class="error-msg">{{ errorMessage }}</p>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			username: "",
			errorMessage: "",
			isLoggingIn: false,
		};
	},
	methods: {
		async login() {
			// 1. Pulizia e Validazione locale
			const cleanName = this.username.trim();
			if (cleanName.length < 3 || cleanName.length > 16) {
				this.errorMessage = "L'username deve avere tra 3 e 16 caratteri.";
				return;
			}

			this.errorMessage = "";
			this.isLoggingIn = true;

			try {
				// 2. Chiamata al backend in Go
				const response = await this.$axios.post("/session", {
					name: cleanName,
				});

				// 3. Il backend ci risponde con l'identificatore (l'ID utente)
				const userID = response.data.identifier;

				// 4. Salva l'ID e il nome utente nella memoria del browser
				localStorage.setItem("token", userID);
				localStorage.setItem("username", cleanName);

				this.$axios.defaults.headers.common['Authorization'] = `Bearer ${userID}`;

				// 5.  Porta alla pagina principale delle chat
				this.$router.push("/chat"); 
				
			} catch (error) {
				// Se il server restituisce errore (es. 400 Bad Request o 500)
				this.errorMessage = "Errore durante l'accesso. Riprova.";
				console.error(error);
			} finally {
				this.isLoggingIn = false;
			}
		},
	},
};
</script>

<style scoped>

.login-page {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100vh;
	background-color: #f0f2f5;
}
.login-box {
	background: white;
	padding: 2rem;
	border-radius: 8px;
	box-shadow: 0 4px 6px rgba(0,0,0,0.1);
	text-align: center;
	display: flex;
	flex-direction: column;
	gap: 1rem;
	width: 100%;
	max-width: 350px;
}
input {
	padding: 0.8rem;
	border: 1px solid #ccc;
	border-radius: 4px;
	font-size: 1rem;
}
button {
	padding: 0.8rem;
	background-color: #008cba;
	color: white;
	border: none;
	border-radius: 4px;
	font-size: 1rem;
	cursor: pointer;
}
button:disabled {
	background-color: #ccc;
}
.error-msg {
	color: red;
	font-size: 0.9rem;
}
</style>