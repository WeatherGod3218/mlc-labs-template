package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/WeatherGod3218/mlc-project-template/internal/airtable"
	"github.com/WeatherGod3218/mlc-project-template/internal/firebase"

	"github.com/WeatherGod3218/mlc-project-template/internal/logging"
	"github.com/sirupsen/logrus"
)

func replaceNullWithString(obj map[string]interface{}) {
	for key, value := range obj {
		if value == nil {
			obj[key] = "null"
		} else if nested, ok := value.(map[string]interface{}); ok {
			replaceNullWithString(nested)
		}
	}
}

func getAirtableData(c *gin.Context) ([]airtable.AirtableRecord, error) {
	tableName, err := firebase.GetLowestTable(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtableData"}).Warn("error deciding which airtable to use!")
		return nil, err
	}

	tableURI := os.Getenv("AIRTABLE_TABLE" + tableName)

	airTable, err := airtable.GetAirtableURI(tableURI)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "getAirtab"}).Warn("error fetching airtable!")
		return nil, err
	}

	return airTable, nil
}

func GetHomepage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GetData(c *gin.Context) {
	data, err := getAirtableData(c)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "chooseAirtableMiddleware"}).Fatal("error fetching airtable!")

	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func SubmitResults(c *gin.Context) {
	var results map[string]interface{}

	err := c.ShouldBindJSON(&results)
	if err != nil {
		c.JSON(400, gin.H{"message": "Invalid request body"})
		return
	}

	replaceNullWithString(results)

	err = firebase.PushToDatabase(c, results)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{"error": err, "module": "api", "method": "SubmitResults"}).Warn("error updating database!")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving results.", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully saved results!"})
}
