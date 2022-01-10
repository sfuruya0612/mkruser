package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/urfave/cli/v2"
)

var (
	name      string
	date      string
	hash      string
	goversion string
)

type User struct {
	ID    string
	Email string
}

func main() {
	app := cli.NewApp()

	app.Name = name
	app.Usage = "CLI to get and delete mackerel users."
	app.Version = fmt.Sprintf("%s %s (Build by: %s)", date, hash, goversion)

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "key",
			Aliases:  []string{"k"},
			EnvVars:  []string{"MACKEREL_API_KEY"},
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "Get user list.",
			Action: func(c *cli.Context) error {
				findUsers(c)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"del"},
			Usage:   "Delete user.",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "user-id",
					Aliases:  []string{"id"},
					Usage:    "Specify the target user id.",
					Required: true,
				},
			},
			Action: func(c *cli.Context) error {
				deleteUsers(c)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		code := 1
		if c, ok := err.(cli.ExitCoder); ok {
			code = c.ExitCode()
		}
		log.Fatalln(err.Error())
		os.Exit(code)
	}
}

// TODO: Add bulk get method
func findUsers(c *cli.Context) {
	mkr := getClient(c.String("key"))

	users, err := mkr.FindUsers()
	if err != nil {
		log.Fatalln(err)
	}

	list := []User{}
	for _, u := range users {
		list = append(list, User{
			ID:    u.ID,
			Email: u.Email,
		})
	}

	header := []string{
		"ID",
		"Email",
	}

	if err := output(os.Stdout, header, list); err != nil {
		log.Fatalln(err)
	}
}

func output(wrt io.Writer, header []string, list []User) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return err
	}

	for _, l := range list {
		if _, err := fmt.Fprintln(w, l.tabString()); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func (u *User) tabString() string {
	fields := []string{
		u.ID,
		u.Email,
	}

	return strings.Join(fields, "\t")
}

// TODO: Add bulk delete method
func deleteUsers(c *cli.Context) {
	mkr := getClient(c.String("key"))
	id := c.String("user-id")

	// TODO: Confirm that it's okay to delete.
	user, err := mkr.DeleteUser(id)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s delete success.", user.Email)
}

func getClient(key string) *mackerel.Client {
	return mackerel.NewClient(key)
}
