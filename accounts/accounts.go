package accounts

import (
	"crypto/tls"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/nanobox-io/golang-scribble"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type Account struct {
	Name              string
	Cluster_name      string
	Consolegw         string
	Db_backups_bucket string
	Dbseeds_s3_bucket string
	Graylog_api_url   string
	Graylog_pass      string
	Graylog_user      string
	Livelogs_token    string
	Livelogs_url      string
	Livelogs_address  string
}

var Acc Account

func init() {
	// verbose := viper.GetBool("verbose")

	Acc = Account{}
	Acc.Load()
	// fmt.Println("Account Init")
}

func (a *Account) Load() (bool, error) {
	verbose := viper.GetBool("verbose")

	db := getStorage()

	if err := db.Read("Account", "active", &a); err != nil || a.Name == "" {
		if verbose {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}
		return false, err
	}
	return true, nil
}

func (a *Account) Verify() bool {
	verbose := viper.GetBool("verbose")

	if ok, err := a.Load(); ok {
		return true
	} else {
		fmt.Println("You need to login first.")
		a.ListAccounts()
		if verbose {
			fmt.Println(err)
		}
		os.Exit(2)
		return false
	}
}

func getConsul() *consulapi.Client {
	consul_address := viper.GetString("consul_address")
	consul_scheme := viper.GetString("consul_scheme")
	consul_user := viper.GetString("consul_user")
	consul_pass := viper.GetString("consul_pass")
	consul_token := viper.GetString("consul_token")

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	config := consulapi.DefaultConfig()
	config.Address = consul_address
	config.Scheme = consul_scheme
	config.HttpClient = &http.Client{Transport: transport}
	config.HttpAuth = &consulapi.HttpBasicAuth{Username: consul_user, Password: consul_pass}
	config.Token = consul_token

	consul, _ := consulapi.NewClient(config)
	return consul
}

func getStorage() *scribble.Driver {
	dataDir, err := filepath.Abs(filepath.Join(os.Getenv("HOME"), ".cnvy/db"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading filepath: ", err.Error())
		os.Exit(2)
	}
	db, _ := scribble.New(dataDir, nil)
	return db
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getAccounts() []string {
	consul := getConsul()

	kv := consul.KV()

	keys, qm, err := kv.Keys("cnvycli/", "/", nil)
	listAccounts := []string{}
	if err != nil {
		fmt.Println(err, qm)
	} else {
		cFolders := keys[1:]
		for _, cFolderABS := range cFolders {
			cFolderParts := strings.Split(cFolderABS, "/")
			listAccounts = append(listAccounts, cFolderParts[1])
		}
	}
	return listAccounts
}

func (a *Account) ListAccounts() {
	fmt.Println("Accounts:", a.Name)
	for _, AccountName := range getAccounts() {
		if a.Name == AccountName {
			fmt.Println("*", AccountName)
		} else {
			fmt.Println(" ", AccountName)
		}
	}
}

func (a *Account) SUAccount(AccountName string) {
	if stringInSlice(AccountName, getAccounts()) == false {
		fmt.Fprintln(os.Stderr, "Account does not exist.")
		os.Exit(2)
	}
	db := getStorage()

	Acc = Account{Name: AccountName}

	consul := getConsul()

	kv := consul.KV()

	kvPairs, qm, err := kv.List(fmt.Sprintf("cnvycli/%s/", AccountName), nil)
	if err != nil {
		fmt.Println(err, qm)
	} else {
		for _, kvPair := range kvPairs {
			keyName := strings.Split(kvPair.Key, "/")[2]
			if keyName != "" {
				val := string(kvPair.Value)
				if val != "" {
					field := reflect.ValueOf(&Acc).Elem().FieldByName(strings.Title(keyName))
					if field.CanSet() {
						field.SetString(string(val))
					}
				}
			}
		}
		db.Write("Account", "active", Acc)
		fmt.Println("Logged in to:", AccountName)
	}
}
