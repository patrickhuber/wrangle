package commands

import (
	"strings"
	"github.com/patrickhuber/wrangle/store"
	"fmt"
	"github.com/patrickhuber/wrangle/services"
	"github.com/urfave/cli"
)

// CreateSetCommand creates the set command with the cli app and credential service factory
func CreateSetCommand(app *cli.App, credentialServiceFactory services.CredentialServiceFactory) *cli.Command {
	command := &cli.Command{
		Name: "set",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "store, s",
				Usage: "the store where the credential will be set",
			},
			cli.StringFlag{
				Name: "path, p",
				Usage: "the path or key to store the credential",
			},
			cli.StringFlag{
				Name: "type, t",
				Usage: "Sets the credential type. Valid types include 'value', 'structured', 'password', 'user', 'certificate', 'ssh' and 'rsa'.",
			},
			cli.StringFlag{
				Name: "value, v",
				Usage: "[Value, Structured] Sets the value for the credential",
			},
			cli.StringFlag{
				Name: "password, p",
				Usage: "[Password, User] Sets the password value of the	credential",
			},
			cli.StringFlag{
				Name: "username, u",
				Usage: "[User] Sets the username value of the credential",
			},
			cli.StringFlag{
				Name: "private-key",
				Usage: "[Certificate, SSH, RSA] Sets the private key from file	or value",
			},
			cli.StringFlag{
				Name: "public-key",
				Usage: "[SSH, RSA, Certificate] Sets the public key from file or value",
			},
			cli.StringFlag{
				Name: "ca",
				Usage: "[Certificate] Sets the root CA from file or value",
			},
			cli.BoolFlag{
				Name: "encrypt, e",
				Usage: "if the credential should be encrypted",
			},
		},
		Action: func(context cli.Context) error{
			credentialType := context.String("type")
			if credentialType == ""{
				return fmt.Errorf("missing required flag 'type'")
			}

			path := context.String("path")
			if path == ""{
				return fmt.Errorf("missing required flag 'path'")
			}

			storeName := context.String("store")
			if storeName == ""{
				return fmt.Errorf("missing required flag 'store'")
			}

			var item store.Item
			switch(strings.ToLower(credentialType)){
			case string(store.Value):
				v := context.String("value")
				item = store.NewValueItem(path, v)
			case string(store.Certificate):
				private := context.String("private-key")
				public := context.String("public-key")
				ca := context.String("ca")
				item = store.NewCertificateItem(path, private, public, ca)
			case string(store.RSA):
				private := context.String("private-key")
				public := context.String("public-key")
				item = store.NewRSAItem(path, private, public)
			case string(store.SSH):
				private := context.String("private-key")
				public := context.String("public-key")
				item = store.NewSSHItem(path, private, public)
			case string(store.Password):
				password := context.String("password")
				item = store.NewPasswordItem(path, password)
			case string(store.User):
				password := context.String("password")
				username := context.String("username")
				item = store.NewUserItem(path, username, password)
			case string(store.Structured):
			default:
				return fmt.Errorf("unsupported credential type %s", credentialType)
			}
			
			configFile := context.GlobalString("config")

			credentialService, err := credentialServiceFactory.Create(configFile)
			if err != nil{
				return err
			}

			return credentialService.Set(storeName, item)
		},
	}
	setCommandCustomHelpTemplateWithGlobalOptions(app, command)
	return command
}
