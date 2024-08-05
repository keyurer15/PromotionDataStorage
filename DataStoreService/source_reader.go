package main

import(
 "encoding/csv"
 "os"
 "log/slog"
 "strings"
 "strconv"
 "regexp"
 "time"
)

//struct to hold data read from the promotion csv file
type promotionData struct {
	Id              string  `json:"id"`
	Price           float64 `json:"price"`
	Expiration_date string  `json:"expiration_date"`
}

//Map to quickly access promotion data
var promotionDataStore map[string]promotionData = make(map[string]promotionData) 

//Store time stamp of last refresh of data store
var lastRefreshedTime = time.Now()
//Read the promotion.csv file if last read was 1200 seconds(20 minutes*60) ago, data changes every 30 minutes
const refreshDuration = 1200

//Method to read and process data from promotion data file
func readSourceData(dataFile string) {
   // Open the CSV file
   slog.Info("Enter ReadSourceData\n")
   file, err := os.Open(dataFile)
   if err != nil {
       panic(err)
   }
   defer file.Close()
   // Read the CSV data
   csvReader := csv.NewReader(file)
   csvReader.FieldsPerRecord = -1 // Allow variable number of fields
   slog.Info("Reading file promotion data file\n")
   rawData, err := csvReader.ReadAll()
   if err != nil {
       panic(err)
   }
   
	//Now process/transform the read data and store

	//Transform of expiration date from UnixDate format to YYYY-MM-DD hh:mm:ss format is expected in the response
	//Utilizing regular expression to trim timezone details
	re := regexp.MustCompile(`[+-]\d{4}.*`)  //stored date example:2018-09-20 07:41:26 +0200 CEST 

    
	//TODO: Add benchmarking
	
	var dataFromSource promotionData
	for recordIndex := 0; recordIndex < len(rawData); recordIndex++ {
		for dataIndex := 0; dataIndex < len(rawData[recordIndex]); dataIndex++ {

			dataFromSource.Id = rawData[recordIndex][0]
			//Transform price string to float data type
			dataFromSource.Price,_ = strconv.ParseFloat(strings.TrimSpace(rawData[recordIndex][1]), 64)
			//Transform expiration date to YYYY-MM-DD hh:mm:ss format as expected in the response
			dataFromSource.Expiration_date = strings.TrimSpace(re.ReplaceAllString(rawData[recordIndex][2], ""))
			promotionDataStore[dataFromSource.Id] = dataFromSource
		}
	}
	slog.Info("Datastore has","entries",len(promotionDataStore))
	//Capture the time stamp of data store refresh
	lastRefreshedTime = time.Now()
	slog.Info("Exit ReadSourceData\n")
}

//Method to retrieve promotdata based upon giveb promotion ID
//The method return false if the promotdata record is not found
func getPromotionData(promotionID *string, responseData *[]byte) bool{

	promotdata,doExist := promotionDataStore[*promotionID]
	if !doExist {
		slog.Info(*promotionID," not found")
		return false
	}
    //Convert promotionData struct to promotionResult to json marshalling

	responseJson, err := formPromotionJsonResponse(&promotdata)
	if err != nil {
		slog.Error(err.Error())
		return false
	}
	*responseData = responseJson
	return true
}

//Method to read the csv file on demand and update the data store
func refreshDataStore() {
 
 timenow := time.Now()
 duration := timenow.Sub(lastRefreshedTime)
 timeDiff := duration.Seconds()
 
 
 if timeDiff > refreshDuration { 
		readSourceData(promotion_data_file)
		slog.Info("DataStore refreshed at " , lastRefreshedTime.Format("2006-01-02 15:04:05")) //YYYY-MM-DD HH:mm:ss format
   }
}



