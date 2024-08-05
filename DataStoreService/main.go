package main

import (
	"log/slog"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net/http"
	"strings"
)

func main() {
	
	slog.Info("Data service started")
	
	//Read promotion data
	readSourceData()	
	
	//Start listening for search quest
	http.ListenAndServe(":1323", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Client connected")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			slog.Error("Error starting socket server ", err.Error())
		}
		go func() {
			defer conn.Close()
			for {
					msgReq, op, err := wsutil.ReadClientData(conn)
					if err != nil {
						slog.Error("Error receiving data: " + err.Error())
						return
				}
				//Check and refresh data store, there could be long duration (>30 min) betweeen consecutive search requestes
				refreshDataStore()
				
				apiCall := string(msgReq[:])
				trimedAPICall := strings.ReplaceAll(apiCall,"find/","")
				promoId := strings.TrimSpace(trimedAPICall)
				slog.Info("Search request received with promotion ID: " + promoId)
				
				
				var response []byte
				
				//Find the promotion ID in the data store
				found := getPromotionData(&promoId,&response)
				if !found {
						// Write response status code
						slog.Info("Promotion: " + promoId + " not found")
						w.WriteHeader(http.StatusNotFound)
						return
				}				

				//Send data search response to the promotion service
				err = wsutil.WriteServerMessage(conn, op, []byte(response))
		
				if err != nil {
					slog.Error("Error sending data " + err.Error())
					return
				}
				
				//log the response for debugging!
				slog.Info("Server sent response: " + string(response[:]))
				
				//Check if the refresh of data store is needed
				refreshDataStore()
			}
		}()
	}))
}
