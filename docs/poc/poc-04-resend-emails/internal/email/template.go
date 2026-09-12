package email

import "fmt"

func BusinessInvitation(name string) (subject, htmlBody string) {
	subject = "¡Has recibido una invitación para unirte a Ruta de Origen!"
	htmlBody = fmt.Sprintf(`
	<h1> ¡Hola, %s! 👋</h1>
	<p>¡Soy Ori! Tengo una invitación especial para ti.</p>
	<p>Has recibido una invitación para unirte a un negocio en <strong>Ruta de Origen</strong>.</p>
	<p>Al aceptar la invitación, podrás formar parte del equipo y acceder a las herramientas y funcionalidades que el negocio tenga habilitadas para gestionar sus lotes y procesos de trazabilidad.</p>
	<p> ¡Espero verte pronto por Ruta de Origen! </p> 
	<p> <a href="https://rutadeorigen.co" style="display:inline-block; padding:12px 24px; background-color:#6B4F3A; color:#ffffff; text-decoration:none; border-radius:6px; font-weight:bold;"> Abrir Invitación </a> </p>
	<p> — Ori </p>
	`, name)

	return subject, htmlBody
}
