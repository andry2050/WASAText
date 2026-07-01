<script setup>
import { RouterView } from 'vue-router'
</script>

<script>
export default {

	created() {
		// Controlla se c'è un token salvato quando si apre l'app
		const token = localStorage.getItem("token");
		if (token) {
			this.$axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
		}
	},

	methods: {
		logout() {
			localStorage.removeItem("token");
			localStorage.removeItem("username");
			this.$router.push("/login");
		}
	},
	computed: {
		isLoggedIn() {
			return !!localStorage.getItem("token");
		},
		username() {
			return localStorage.getItem("username");
		}
	}
}
</script>

<template>
	<header class="navbar navbar-dark sticky-top bg-primary flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-5" href="#/">
			💬 WASAText
		</a>
		<div class="navbar-nav" v-if="isLoggedIn">
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