package extensions

import (
	"errors"
	"strings"

	"github.com/ski7777/goextensioniser/pkg/extensioniser"
	"github.com/ski7777/gomsgqueue/pkg/interfaces"
	"github.com/ski7777/gomsgqueue/pkg/messagequeue"
	sm "github.com/ski7777/gomultiwa/internal/scopemanager"
	"github.com/ski7777/gomultiwa/internal/usermanager"
	"github.com/ski7777/gomultiwa/internal/webserver/websocketserver"
)

type Extension struct { // implements "github.com/ski7777/gomsgqueue/pkg/interfaces".Master
	ext      interfaces.Extension
	sm       *sm.ScopeManager
	ws       *websocketserver.WSServer
	um       *usermanager.UserManager
	mq       *messagequeue.MessageQueue
	mqm      *messagequeue.Master
	mqe      *messagequeue.Extension
	embedded bool
}

func (e *Extension) ConnectEmbedded(f func(*messagequeue.MessageQueue) interfaces.Extension) {
	e.mqe = messagequeue.NewExtension()
	e.mqe.ConnectToMaster(e.mqm)
	e.ext = f(e.mqe.GetMessageQueue())
	e.embedded = true
}

func (e *Extension) ConnectDedicated(cmd string) error {
	cm, err := extensioniser.NewDedicatedExtension(cmd)
	if err != nil {
		return err
	}
	cm.ConnectToMaster(e.mqm)
	return nil
}

func (e *Extension) start() {
	e.mq.Run()
	if e.embedded {
		e.mqe.Run()
	}
}

func (e *Extension) stop() {
	e.mq.Stop()
}

func (e *Extension) handleScopeRequest() {
	sl := e.sm.GetScopes()
	for _, s := range sl {
		if !s.Approved {
			var approve bool
			if e.embedded {
				approve = true
			} else {
				// Some more magic needed
			}
			var err error
			if approve {
				if s.Name == "http-serve" {
					d := strings.Split(s.Value, ":")
					if len(d) == 2 {
						e.ws.HandleExtensionFunc(d[0], d[1], e.mq)
					} else {
						err = errors.New("Invalid Value")
					}
				}
			}
			if err == nil {
				go func() {
					_ = e.sm.ApproveScope(s)
				}()
			}
		}
	}
}

func NewExtension(ws *websocketserver.WSServer, um *usermanager.UserManager) *Extension {
	e := new(Extension)
	e.ws = ws
	e.um = um
	e.mqm = messagequeue.NewMaster()
	e.mq = e.mqm.GetMessageQueue()
	e.sm = sm.NewScopeManager(e.mq)
	e.sm.SetRequestHandler(e.handleScopeRequest)
	return e
}
