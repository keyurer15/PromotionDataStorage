package main

import(
 "encoding/json"
 "log/slog"
)

//Method to form JSON response from promotionData struct
func formPromotionJsonResponse(result *promotionData) ([]byte, error){
    jsonData, err := json.Marshal(result)
    if err != nil {

		slog.Error("Error in marshalling JSON ", err.Error())

    }
	return jsonData, err
}