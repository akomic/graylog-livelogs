package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	livelogsCmd = &cobra.Command{
		Use:   "livelogs",
		Short: "App Live Logs",
		Long:  ``,

		Run: livelogs,
	}
)

var (
	cluster   string
	rawOutput bool
)

// init
func init() {
	livelogsCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "Container Cluster Name")
	livelogsCmd.Flags().StringSliceP("filter", "f", nil, "Filter e.g. -f stack_name=idea1")
	livelogsCmd.Flags().BoolVarP(&rawOutput, "raw", "r", false, "Dump complete messages as json")
	viper.BindPFlag("filter", livelogsCmd.Flags().Lookup("filter"))
}

func livelogs(ccmd *cobra.Command, args []string) {
	filters := viper.GetStringSlice("filter")

	livelogs_url := viper.GetString("livelogs_url")
	livelogs_token := viper.GetString("livelogs_token")

	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: livelogs_url, Path: "/filter"}
	querySlice := []string{}

	if livelogs_token != "" {
		querySlice = append(querySlice, "token="+livelogs_token)
	}

	for _, filter := range filters {
		querySlice = append(querySlice, filter)
	}
	u.RawQuery = strings.Join(querySlice, "&")

	// log.Printf("connecting to %s", s)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var f []interface{}
			err = json.Unmarshal(message, &f)
			if err != nil {
				fmt.Println("Error parsing JSON: ", err)
				fmt.Println(string(message))
				continue
			}

			l := f[1].(map[string]interface{})

			// for k := range l {
			// 	fmt.Printf("%s ", k)
			// }
			// fmt.Printf("\n")
			// fmt.Println(reflect.TypeOf(l))

			if rawOutput {
				jsonString, err := json.Marshal(l)
				if err != nil {
					log.Println("json:", err)
					continue
				}
				fmt.Println(string(jsonString))
			} else {
				fmt.Printf("%s ", l["timestamp"])
				if l["container_name"] != nil && l["stack_name"] != nil {
					fmt.Printf("%s %s",
						l["container_name"],
						l["stack_name"],
					)
				} else if l["container_name"] != nil && l["ecs_cluster"] != nil && l["task_definition"] != nil {
					fmt.Printf("%s %s %s",
						l["container_name"],
						l["ecs_cluster"],
						l["task_definition"],
					)
				} else if l["command"] != nil && l["image_name"] != nil {
					fmt.Printf("%s %s",
						l["command"],
						l["image_name"],
					)
				} else {
					fmt.Printf("%s",
						l["source"],
					)
				}
				fmt.Printf(" %s\n", l["message"])
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
