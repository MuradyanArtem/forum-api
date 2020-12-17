package app

type App struct {
	User   User
	Forum  Forum
	Thread Thread
	Post   Post
}

func New(repository persistence.App) {

}
