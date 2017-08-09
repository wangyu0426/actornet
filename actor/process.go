package actor

import (
	"github.com/davyxu/actornet/mailbox"
	"github.com/davyxu/actornet/proto"
)

type Process interface {
	Send(sender *PID, data interface{})
	Stop()

	PID() *PID
}

type localProcess struct {
	mailbox mailbox.MailBox

	pid PID

	a Actor
}

func (self *localProcess) PID() *PID {
	return &self.pid
}

func (self *localProcess) Send(sender *PID, data interface{}) {

	self.mailbox.Push(&mailContext{
		msg:  data,
		src:  sender,
		self: &self.pid,
	})
}

func (self *localProcess) Stop() {

	self.Send(&self.pid, &proto.Stop{})
}

func (self *localProcess) Recv(msg interface{}) {

	self.a.Receive(msg.(Context))
}

func NewLocalProcess(a Actor, pid PID) *localProcess {

	self := &localProcess{
		mailbox: mailbox.NewBounded(10),
		a:       a,
		pid:     pid,
	}

	self.pid.proc = self

	self.mailbox.Start(self)

	self.Send(&self.pid, &proto.Start{})

	return self
}
