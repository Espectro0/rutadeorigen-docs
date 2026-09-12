const content = document.getElementById("content");

async function cargarEstado() {
	try {
		const res = await fetch("/api/me");

		if (res.ok) {
			const user = await res.json();
			renderPerfil(user);
			return;
		}

		if (res.status === 401) {
			renderLogin();
			return;
		}

		renderError(`Error inesperado (${res.status})`);
	} catch (err) {
		renderError("No se pudo conectar con el servidor.");
	}
}

function renderLogin() {
	content.innerHTML = `
		<a class="button" href="/login">Iniciar sesión con Ruta de Origen</a>
	`;
}

function renderPerfil(user) {
	const nombre = [user.first_name, user.last_name].filter(Boolean).join(" ") || user.email;
	const inicial = nombre.charAt(0).toUpperCase();

	content.innerHTML = `
		<div class="avatar">${inicial}</div>
		<p class="user-name">${nombre}</p>
		<p class="user-email">${user.email}</p>
		<button class="button secondary" id="logout-btn">Cerrar sesión</button>
	`;

	document.getElementById("logout-btn").addEventListener("click", cerrarSesion);
}

function renderError(mensaje) {
	content.innerHTML = `<p class="error">${mensaje}</p>`;
}

async function cerrarSesion() {
	await fetch("/api/logout", { method: "POST" });
	cargarEstado();
}

cargarEstado();