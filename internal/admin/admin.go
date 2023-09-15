package internal

import (
	"bufio"
	"fmt"
	internal "net-cat/internal/server"
	"os"
	"strings"
)

func Admin(server *internal.Server) {
	for {
		fmt.Print("[ADMIN]: ")
		input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		inputArr := strings.Split(strings.Replace(input, "\n", "", 1), " ")

		switch inputArr[0] {
		case "history":
			fmt.Print(server.AllMessages)
		case "delete":
			if len(inputArr) != 2 {
				continue
			}
			delete(server.Connections, inputArr[1])
			fmt.Println(internal.AdminDeletedUser)
			server.NewMsg(&internal.Message{
				Text: internal.AdminDeletedUser + inputArr[1] + "\n",
				Time: "",
			})
		}
	}
}
