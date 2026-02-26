package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Profilo Utente e Ricerca
	rt.router.GET("/users", rt.wrap(rt.searchUsers))
	rt.router.PUT("/users/me/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/me/photo", rt.wrap(rt.setMyPhoto))

	// Conversazioni (lista e singola conversazione)
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.GET("/conversations/:conversation_id", rt.wrap(rt.getConversation))

	// Gestione Messaggi
	rt.router.POST("/conversations/:conversation_id/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversations/:conversation_id/forward", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/messages/:message_id", rt.wrap(rt.deleteMessage))

	// Gestione Reazioni (Commenti)
	rt.router.POST("/messages/:message_id/comments", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/messages/:message_id/comments/:comment_id", rt.wrap(rt.uncommentMessage))

	// Gestione Gruppi
	rt.router.POST("/groups", rt.wrap(rt.createGroup))
	rt.router.POST("/groups/:group_id/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/groups/:group_id/members/:user_id", rt.wrap(rt.leaveGroup))
	rt.router.PUT("/groups/:group_id/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/groups/:group_id/photo", rt.wrap(rt.setGroupPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
