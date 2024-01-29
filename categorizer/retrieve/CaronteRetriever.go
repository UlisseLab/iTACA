package retrieve

import (
	"categorizer/storage"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"slices"
)

// CaronteRetriever : implements Retriever for Caronte service, fetching reconstructed TCP streams from a service exposing a REST API
// address must be the address of the machine which is hosting the service, port must be the port which exposes the service
type CaronteRetriever struct {
	address string
	port    uint16
}

func NewCaronteRetriever(address string, port uint16) *CaronteRetriever {
	return &CaronteRetriever{address: address, port: port}
}

func (r *CaronteRetriever) Retrieve(ctx context.Context, results chan<- Result) {
	var visited []storage.RowID

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Retriever: task stopped")
			return
		default:
			addr := fmt.Sprintf("%s:%d/api/connections?limit=50", r.address, r.port)
			req, err := http.NewRequest(http.MethodGet, addr, nil)
			if err != nil {
				fmt.Printf("client: could not create request: %s\n", err)
				os.Exit(1)
			}

			req.Header.Add("Accept", "*/*")
			req.Header.Add("Accept-Language", "en-GB,en;q=0.9,it-IT;q=0.8,it;q=0.7,en-US;q=0.6")
			req.Header.Add("Cache-Control", "no-cache")
			req.Header.Add("Connection", "keep-alive")
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("DNT", "1")
			req.Header.Add("Pragma", "no-cache")
			req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("client: error making http request: %s\n", err)
				os.Exit(1)
			}

			var connections []storage.Connection
			err = json.NewDecoder(res.Body).Decode(&connections)
			if err != nil {
				fmt.Printf("client: error decoding response body: %s\n", err)
				os.Exit(1)
			}

			var Ids = map[storage.RowID]uint16{}
			for _, conn := range connections {
				Ids[conn.ID] = conn.SourcePort
			}

			for id, port := range Ids {
				if slices.Contains(visited, id) {
					continue
				}

				visited = append(visited, id)

				// id is converted to hex because Caronte expects a hex string as id
				req.URL.Path = fmt.Sprintf("/api/streams/%s/download?format=only_client", id.Hex())

				res, err = http.DefaultClient.Do(req)
				if err != nil {
					fmt.Printf("client: error making http request: %s\n", err)
					os.Exit(1)
				}

				resBody, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Printf("client: error reading response body: %s\n", err)
					os.Exit(1)
				}

				results <- Result{Stream: string(resBody), SrcPort: port}
			}
		}
	}
}
