package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
	"github.com/alexandrevicenzi/go-sse"
)


var (
    //go:embed public
    embedRootFS embed.FS
    publicRootFS, _ = fs.Sub(embedRootFS, "public")
    chatChannel = make(chan []byte)
    parser = argparse.NewParser("picopear", "Fresh and tiny tool for aggregating livestream chats into a browser source")
)

func relayChatEndpoint(w http.ResponseWriter, req *http.Request) {
    // When POST with chat message json string, forward to chatChannel
    chatMsg, _ := io.ReadAll(req.Body)
    fmt.Println(string(chatMsg))
    chatChannel <- chatMsg
}

func mergedChatEndpoint() *sse.Server{
    s := sse.NewServer(nil)
    
    go func() {
        for {
            // Auto send live chat messages to subscribed clients
            s.SendMessage("/mergedchat", sse.SimpleMessage(string(<-chatChannel)))
        }
    }()

    return s
}

func setupRoutes() *http.ServeMux{

    mux := http.NewServeMux()

    // Root of server is mapped to root of embedded public folder
    mux.Handle("/", http.FileServer(http.FS(publicRootFS)))

    // Subpath for forwarding merged chat messages to browser
    mux.Handle("/mergedchat", mergedChatEndpoint())

    // Subpath where 3rd party services POST json chat messages
    mux.HandleFunc("/relay", relayChatEndpoint)
    
    return mux

}

func main() {
    // Build and parse cli arguments
    listenPort := parser.String("p", "port", &argparse.Options{Required: false, Default: "4554", Help: "Port to bind"})
    parser.Parse(os.Args)

    addr := fmt.Sprintf(":%s", *listenPort)

    http.ListenAndServe(addr, setupRoutes())
}
