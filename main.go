package main

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"go.elara.ws/go-lemmy"
	"go.elara.ws/go-lemmy/types"
)

type Flemmy struct {
	app    fyne.App
	client *lemmy.Client
	ctx    context.Context
}

var flemmy Flemmy

func login() {
}

func init() {
	flemmy = Flemmy{
		app: app.New(),
		ctx: context.Background(),
	}
}

func main() {
	w := flemmy.app.NewWindow("Flemmy")
	posts := flemmy.app.NewWindow("Flemmy Posts")

	instance := widget.NewEntry()
	user := widget.NewEntry()
	pass := widget.NewPasswordEntry()

	instance.Text = "https://lemmy.ml"

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Instance", Widget: instance},
			{Text: "Login", Widget: user},
			{Text: "Password", Widget: pass},
		},
		OnSubmit: func() {
			var err error
			flemmy.client, err = lemmy.New(instance.Text)
			if err != nil {
				log.Fatal(err)
			}
			err = flemmy.client.ClientLogin(flemmy.ctx, types.Login{
				UsernameOrEmail: user.Text,
				Password:        pass.Text,
			})
			if err != nil {
				log.Fatal(err)
			}

			postList, err := flemmy.client.Posts(flemmy.ctx, types.GetPosts{
				CommunityName: types.NewOptional("!golang"),
			})
			if err != nil {
				log.Fatal(err)
			}

			list := widget.NewList(
				func() int {
					return len(postList.Posts)
				},
				func() fyne.CanvasObject {
					return widget.NewLabel("template")
				},
				func(i widget.ListItemID, o fyne.CanvasObject) {
					p := postList.Posts[i]
					o.(*widget.Label).SetText(p.Post.Name)
				})

			posts.SetContent(list)
			posts.Show()
		},
	}

	w.SetContent(form)
	w.ShowAndRun()
}
