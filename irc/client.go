package irc

type Client interface {
  Loop()
  GetInbox() chan<- Message
}
