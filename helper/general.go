package helper

import (
	"github.com/gin-gonic/gin"
)

func JsonResponse(stringSlice []string, message string, status int, c *gin.Context) {

	lenSlice := len(stringSlice)
	stringMap := make(map[string]string)

	for i := 0; i < lenSlice-1; i += 2 {
		// Use pairs of elements from the slice as key-value pairs in the map
		stringMap[stringSlice[i]] = stringSlice[i+1]
	}

	c.JSON(status, gin.H{
		"status":  status,
		"message": message,
		"data":    stringMap,
	})
}

func JsonResponseMap(stringSlice map[string]string, message string, status int, c *gin.Context) {

}
