package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

func Livelogs(url string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		panic(err)
	}
	defer c.Close()

	done := make(chan struct{})

	cPick := colorPicker()
	go func() {
		defer c.Close()
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Fprintln(os.Stderr, "read:", err)
				os.Exit(2)
				return
			}

			var f []interface{}
			err = json.Unmarshal(message, &f)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error parsing JSON: ", err)
				fmt.Fprintln(os.Stderr, string(message))
				continue
			}

			l := f[1].(map[string]interface{})

			printLog(l, cPick)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				fmt.Fprintln(os.Stderr, "write:", err)
				return
			}
		case <-interrupt:
			fmt.Fprintln(os.Stderr, "interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Fprintln(os.Stderr, "write close:", err)
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

func colorPicker() func(string) string {
	arrayColors := []func(...interface{}) string{
		color.New(color.FgRed).SprintFunc(),
		color.New(color.FgGreen).SprintFunc(),
		color.New(color.FgYellow).SprintFunc(),
		color.New(color.FgBlue).SprintFunc(),
		color.New(color.FgMagenta).SprintFunc(),
		color.New(color.FgCyan).SprintFunc(),
		color.New(color.FgWhite).SprintFunc(),
	}

	idents := map[string]string{}

	i := 0
	picker := func(ident string) string {
		if _, ok := idents[ident]; ok {
			return idents[ident]
		}
		idents[ident] = arrayColors[i](ident)
		if i >= (len(arrayColors) - 1) {
			i = 0
		} else {
			i++
		}
		return idents[ident]
	}
	return picker
}

func printLog(l map[string]interface{}, cPick func(string) string) {
	rawOutput := viper.GetBool("rawOutput")
	if rawOutput {
		jsonString, err := json.Marshal(l)
		if err != nil {
			fmt.Fprintln(os.Stderr, "json:", err)
			return
		}
		fmt.Println(string(jsonString))
	} else {
		timestamp, _ := time.Parse(time.RFC3339Nano, l["timestamp"].(string))
		fmt.Printf("%s ", timestamp.In(time.Now().Location()))
		if l["container_name"] != nil && l["stack_name"] != nil {
			fmt.Printf("%s %s",
				cPick(l["container_name"].(string)),
				l["stack_name"],
			)
		} else if l["container_name"] != nil && l["ecs_cluster"] != nil && l["task_definition"] != nil {
			fmt.Printf("%s %s %s",
				cPick(l["container_name"].(string)),
				l["ecs_cluster"],
				l["task_definition"],
			)
		} else if l["command"] != nil && l["image_name"] != nil {
			fmt.Printf("%s %s",
				cPick(l["command"].(string)),
				l["image_name"],
			)
		} else {
			fmt.Printf("%s",
				cPick(l["source"].(string)),
			)
		}
		fmt.Printf(" %s\n", l["message"])
	}
}
