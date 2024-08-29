package rpchandlers

import (
	"github.com/karlsen-network/karlsend/v2/app/appmessage"
	"github.com/karlsen-network/karlsend/v2/app/rpc/rpccontext"
	"github.com/karlsen-network/karlsend/v2/infrastructure/network/netadapter/router"
)

// HandleNotifyFinalityConflicts handles the respectively named RPC command
func HandleNotifyFinalityConflicts(context *rpccontext.Context, router *router.Router, _ appmessage.Message) (appmessage.Message, error) {
	listener, err := context.NotificationManager.Listener(router)
	if err != nil {
		return nil, err
	}
	listener.PropagateFinalityConflictNotifications()
	listener.PropagateFinalityConflictResolvedNotifications()

	response := appmessage.NewNotifyFinalityConflictsResponseMessage()
	return response, nil
}
