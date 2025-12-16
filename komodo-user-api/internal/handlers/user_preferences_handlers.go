package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Returns preferences for the authenticated user
func GetMyPreferences(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}


	// ========================================
	// TODO: Production Implementation (DynamoDB)
	// ========================================
	// Preferences are stored as a nested attribute in the user record
	//
	// ctx := context.Background()
	// tableName := "komodo-users"
	//
	// result, err := dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
	//     TableName: aws.String(tableName),
	//     Key: map[string]types.AttributeValue{
	//         "user_id": &types.AttributeValueMemberS{Value: userID.(string)},
	//     },
	//     ProjectionExpression: aws.String("preferences"),
	// })
	//
	// if err != nil || result.Item == nil {
	//     g.JSON(http.StatusOK, gin.H{
	//         "user_id": userID,
	//         "preferences": map[string]interface{}{},
	//     })
	//     return
	// }
	//
	// var user struct {
	//     Preferences map[string]interface{} `dynamodbav:"preferences"`
	// }
	// attributevalue.UnmarshalMap(result.Item, &user)
	//
	// g.JSON(http.StatusOK, gin.H{
	//     "user_id": userID,
	//     "preferences": user.Preferences,
	// })
	
	// ========================================
	// Mock Implementation (Development)
	// ========================================
	g.JSON(http.StatusOK, gin.H{
		"user_id":     userID,
		"preferences": gin.H{
			"language": "en-US",
			"timezone": "America/Los_Angeles",
		},
	})
}

// Updates preferences for the authenticated user
func UpdateMyPreferences(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var preferencesData map[string]interface{}
	if err := g.BindJSON(&preferencesData); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// TODO: Validate and update user preferences in database
	g.JSON(http.StatusOK, gin.H{
		"user_id":      userID,
		"message":      "preferences updated successfully",
		"preferences":  preferencesData,
	})
}
