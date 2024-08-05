package main

import (
    "context"
	"log/slog"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"	
	"net/http"
	"strings"
)

//Hander method to process promotion search request
func handlerPromotionSearch(w http.ResponseWriter, r *http.Request) {
    // Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Access request URL
	reqURI := r.RequestURI
	
	//Extract promotion ID from the path
	trimedURI := strings.ReplaceAll(reqURI,"/promotions/","")
	promoId := strings.TrimSpace(trimedURI)
	
	slog.Info("reqURI: " + reqURI + " promoId: " + promoId)
	
	//Create a websocket connection to the datastore service to request data for the promotion ID
	conn, _, _, err := ws.DefaultDialer.Dial(context.Background(), "ws://127.0.0.1:1323/")
	if err != nil {
		slog.Error("Cannot connect to data service " + err.Error())
		// Write response status code
		w.WriteHeader(http.StatusNotFound)
		return		
	}
		
	//Initiate search call to the datastore service
	msg := []byte("find/" + promoId)
	err = wsutil.WriteClientMessage(conn, ws.OpText, msg)
	if err != nil {
		slog.Error("Cannot send: " + err.Error())
		//Write response status code
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	//Read the response sent by the datastore service
	response, _, err := wsutil.ReadServerData(conn)
	if err != nil {
		slog.Error("Cannot receive data: " + err.Error())
		// Write response status code
		w.WriteHeader(http.StatusNotFound)
		return
	}	

	// Write response status code
	w.WriteHeader(http.StatusOK)
	// Write response body
	w.Write(response)	
	
	slog.Info("response sent:" + string(response[:]))
}



//Listener method to serve incoming promotion search request
func listen() {

	//Register handler for processing promotion search service
	http.HandleFunc("/promotions/",handlerPromotionSearch)
	
	//Listen over port 1321
	err := http.ListenAndServe(":1321",nil)
	if err != nil {
		slog.Error ("Failed to start Promotion server: " + err.Error())
	}
}
