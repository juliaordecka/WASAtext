package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	// rt.router.GET("/", rt.getHelloWorld)
	// rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/user/:username/setmyusername", rt.wrap(rt.setMyUsername))
	//    rt.router.PUT("/user/:username/photo", rt.wrap(rt.setMyPhoto))
	//rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	//    rt.router.GET("/conversation/:conversation_id", rt.wrap(rt.getConversation))
	rt.router.POST("/message", rt.wrap(rt.sendMessage))
	//    rt.router.POST("/message/:message_id/forward", rt.wrap(rt.forwardMessage))
	//    rt.router.POST("/message/:message_id/comment", rt.wrap(rt.commentMessage))
	//    rt.router.DELETE("/message/:message_id/uncomment", rt.wrap(rt.uncommentMessage))
	//    rt.router.DELETE("/message/:message_id", rt.wrap(rt.deleteMessage))
	rt.router.POST("/group", rt.wrap(rt.createGroup))
	rt.router.POST("/group/:group_id/add", rt.wrap(rt.addToGroup))
	//    rt.router.DELETE("/group/:group_id/leave", rt.wrap(rt.leaveGroup))
	//    rt.router.PUT("/group/:group_id/name", rt.wrap(rt.setGroupName))
	//    rt.router.PUT("/group/:group_id/photo", rt.wrap(rt.setGroupPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
