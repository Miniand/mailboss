package mailboss

type Handler interface {
	Handle(mail string)
}
