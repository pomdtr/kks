package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kkga/kks/kak"
)

func NewListCmd() *ListCmd {
	c := &ListCmd{Cmd: Cmd{
		fs:        flag.NewFlagSet("list", flag.ExitOnError),
		alias:     []string{"ls", "l"},
		shortDesc: "List Kakoune sessions and clients.",
		usageLine: "[options]",
	}}
	c.fs.BoolVar(&c.json, "json", false, "json output")
	return c
}

type ListCmd struct {
	Cmd
	json bool
}

func (c *ListCmd) Run() error {
	kakSessions, err := kak.Sessions()
	if err != nil {
		return err
	}

	switch c.json {

	case true:
		type session struct {
			Name    string   `json:"name"`
			Clients []string `json:"clients"`
			Dir     string   `json:"dir"`
		}

		sessions := []session{}

		for i, s := range kakSessions {
			d, err := s.Dir()
			sessions = append(sessions, session{Name: s.Name, Clients: []string{}, Dir: d})

			c, err := s.Clients()
			for _, c := range c {
				if c.Name != "" {
					sessions[i].Clients = append(sessions[i].Clients, c.Name)
				}
			}

			if err != nil {
				return err
			}
		}

		j, err := json.MarshalIndent(sessions, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(j))

	case false:
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 1, '\t', 0)

		for _, s := range kakSessions {
			c, err := s.Clients()
			d, err := s.Dir()
			if len(c) == 0 {
				fmt.Fprintf(w, "%s\t: %s\t: %s\n", s.Name, " ", d)
			} else {
				for _, cl := range c {
					fmt.Fprintf(w, "%s\t: %s\t: %s\n", s.Name, cl.Name, d)
				}
			}
			if err != nil {
				return err
			}
		}

		w.Flush()
	}

	return nil
}
