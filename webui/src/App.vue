<script setup>
import { RouterView } from 'vue-router'
</script>

<script>
export default {
	data() {
		return {
			token: localStorage.getItem("token"),
			username: localStorage.getItem("username")
		}
	},
	created() {
		if (this.token) {
			this.$axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`;
		}
	},
	watch: {
		$route() {
			this.token = localStorage.getItem("token");
			this.username = localStorage.getItem("username");
		}
	},
	methods: {
		logout() {
			localStorage.removeItem("token");
			localStorage.removeItem("username");
			this.token = null;
			this.username = null;
            // Rimuoviamo l'autorizzazione di Axios per sicurezza
            delete this.$axios.defaults.headers.common['Authorization'];
			this.$router.push("/login");
		}
	}
}
</script>

<template>
	<header class="navbar navbar-dark sticky-top bg-primary flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-5" href="#/">
			💬 WASAText
		</a>
		<div class="navbar-nav" v-if="token">
			<div class="nav-item text-nowrap d-flex align-items-center">
				<span class="text-white me-3">Ciao, {{ username }}!</span>
				
				<button class="nav-link px-3 bg-transparent border-0" @click="$router.push('/profile')">
					👤 Profilo
				</button>
				
				<button class="nav-link px-3 bg-transparent border-0" @click="logout">
					Logout
				</button>
			</div>
		</div>
	</header>

	<div class="container-fluid mt-3">
		<RouterView />
	</div>
</template>

<style>

.bg-primary {
	background-color: #008cba !important;
}
.navbar-brand {
	font-weight: bold;
}
</style>