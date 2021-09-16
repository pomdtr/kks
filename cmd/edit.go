package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/kkga/kks/kak"
)

func NewEditCmd() *EditCmd {
	c := &EditCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("edit", flag.ExitOnError),
		alias:     []string{"e"},
		shortDesc: "Edit file. In session and client, if set.",
		usageLine: "[options] [file] [+<line>[:<col>]]",
	}}
	// TODO add flag that allows creating new files (removes -existing)
	c.fs.StringVar(&c.session, "s", "", "session")
	c.fs.StringVar(&c.client, "c", "", "client")
	return c
}

type EditCmd struct {
	Cmd
}

func (c *EditCmd) Run() error {
	fp := kak.NewFilepath(c.fs.Args())

	switch c.kakContext.Session.Name {

	case "":
		kctx := &kak.Context{}

		if c.useGitDirSessions {
			kctx.Session = kak.Session{Name: fp.ParseGitDir()}

			if kctx.Session.Name != "" {
				if exists, _ := kctx.Session.Exists(); !exists {
					sessionName, err := kak.Start(kctx.Session.Name)
					if err != nil {
						return err
					}
					fmt.Println("git-dir session started:", sessionName)
				}
			}
		}

		if kctx.Session.Name == "" {
			kctx.Session = kak.Session{Name: c.defaultSession}
		}

		sessionExists, err := kctx.Session.Exists()
		if err != nil {
			return err
		}

		switch sessionExists {
		case true:
			if err := kak.Connect(kctx, fp); err != nil {
				return err
			}
		case false:
			if err := kak.Run(&kak.Context{}, []string{}, fp); err != nil {
				return err
			}
		}

	default:
		switch c.kakContext.Client.Name {
		case "":
			// if no client, attach to session with new client
			if err := kak.Connect(c.kakContext, fp); err != nil {
				return err
			}
		default:
			// if client set, send 'edit [file]' to client
			sb := strings.Builder{}
			sb.WriteString(fmt.Sprintf("edit -existing %s", fp.Name))
			if fp.Line != 0 {
				sb.WriteString(fmt.Sprintf(" %d", fp.Line))
			}
			if fp.Column != 0 {
				sb.WriteString(fmt.Sprintf(" %d", fp.Column))
			}

			if err := kak.Send(c.kakContext, sb.String()); err != nil {
				return err
			}
		}
	}

	return nil
}
