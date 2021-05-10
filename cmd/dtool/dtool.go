package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/matrix-org/dendrite/setup/config"
	"github.com/matrix-org/dendrite/userapi/storage/accounts"
	"golang.org/x/crypto/bcrypt"
)

const (
	addusercmd     = "adduser"
	passwdcmd      = "passwd"
	helpcmd        = "help"
	accountdatacmd = "accountdata"
)

func main() {

	var config, username, password string

	addUserSet := flag.NewFlagSet(addusercmd, flag.ExitOnError)
	addUserSet.StringVar(&config, "config", "./config.yml", "The config file")
	addUserSet.StringVar(&username, "name", "", "The name of the new user")
	addUserSet.StringVar(&password, "password", "", "Wanted password")

	passwdSet := flag.NewFlagSet(passwdcmd, flag.ExitOnError)
	passwdSet.StringVar(&config, "config", "./config.yml", "The config file")
	passwdSet.StringVar(&username, "name", "", "The name of the existing user")
	passwdSet.StringVar(&password, "password", "", "new password")

	getAccountDataSet := flag.NewFlagSet(accountdatacmd, flag.ExitOnError)
	getAccountDataSet.StringVar(&config, "config", "./config.yml", "The config file")
	getAccountDataSet.StringVar(&username, "name", "", "The name of the existing user")

	switch os.Args[1] {
	case addusercmd:
		addUserSet.Parse(os.Args[2:])
		if username == "" {
			fmt.Println("Name must not be empty.")
			passwdSet.Usage()
			os.Exit(0)
		}
		err := doAddUser(config, username, password)
		if err != nil {
			fmt.Printf("Failed to add user %v (%v)\n", username, err.Error())
			os.Exit(1)
		}
	case passwdcmd:
		passwdSet.Parse(os.Args[2:])
		if username == "" {
			fmt.Println("Name must not be empty.")
			passwdSet.Usage()
			os.Exit(0)
		}
		err := doPasswd(config, username, password)
		if err != nil {
			fmt.Printf("Failed to change password (%v)\n", err.Error())
			os.Exit(1)
		}
	case accountdatacmd:
		getAccountDataSet.Parse(os.Args[2:])
		if username == "" {
			fmt.Println("Name must not be empty.")
			passwdSet.Usage()
			os.Exit(0)
		}
		err := doGetAccountData(config, username)
		if err != nil {
			fmt.Printf("Failed to change password (%v)\n", err.Error())
			os.Exit(1)
		}

	case helpcmd:
		flag.Usage()
		os.Exit(0)
	default:
		fmt.Printf("FAIL\n")
		os.Exit(2)
	}
}

func doAddUser(configfile, username, password string) error {
	cfg, err := config.Load(configfile, true)
	if err != nil {
		return err
	}
	accountDB, err := accounts.NewDatabase(&config.DatabaseOptions{
		ConnectionString: cfg.UserAPI.AccountDatabase.ConnectionString,
	},
		cfg.Global.ServerName,
		bcrypt.DefaultCost,
		cfg.UserAPI.OpenIDTokenLifetimeMS)

	if err != nil {
		return err
	}

	_, err = accountDB.CreateAccount(context.Background(), username, password, "")
	if err != nil {
		return err
	}

	return nil
}

func doPasswd(configfile, username, password string) error {
	cfg, err := config.Load(configfile, true)
	if err != nil {
		return err
	}
	accountDB, err := accounts.NewDatabase(&config.DatabaseOptions{
		ConnectionString: cfg.UserAPI.AccountDatabase.ConnectionString,
	},
		cfg.Global.ServerName,
		bcrypt.DefaultCost,
		cfg.UserAPI.OpenIDTokenLifetimeMS)

	if err != nil {
		return err
	}

	err = accountDB.SetPassword(context.Background(), username, password)
	if err != nil {
		return err
	}
	return nil
}

func doGetAccountData(configfile, username string) error {
	cfg, err := config.Load(configfile, true)
	if err != nil {
		return err
	}
	accountDB, err := accounts.NewDatabase(&config.DatabaseOptions{
		ConnectionString: cfg.UserAPI.AccountDatabase.ConnectionString,
	},
		cfg.Global.ServerName,
		bcrypt.DefaultCost,
		cfg.UserAPI.OpenIDTokenLifetimeMS)

	if err != nil {
		return err
	}

	_, rooms, err := accountDB.GetAccountData(context.Background(), username)
	if err != nil {
		return err
	}

	fmt.Printf("Rooms:\n")
	for k := range rooms {
		fmt.Println(k)
	}

	return nil
}
