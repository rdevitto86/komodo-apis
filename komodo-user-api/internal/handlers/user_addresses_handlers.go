package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMyAddresses returns all addresses for the authenticated user
func GetMyAddresses(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}


	// ========================================
	// TODO: Production Implementation (DynamoDB)
	// ========================================
	// Addresses are stored as a nested attribute in the user record
	//
	// ctx := context.Background()
	// tableName := "komodo-users"
	//
	// result, err := dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
	//     TableName: aws.String(tableName),
	//     Key: map[string]types.AttributeValue{
	//         "user_id": &types.AttributeValueMemberS{Value: userID.(string)},
	//     },
	//     ProjectionExpression: aws.String("addresses"),
	// })
	//
	// if err != nil || result.Item == nil {
	//     g.JSON(http.StatusOK, gin.H{
	//         "user_id": userID,
	//         "addresses": []interface{}{},
	//     })
	//     return
	// }
	//
	// var user struct {
	//     Addresses []map[string]interface{} `dynamodbav:"addresses"`
	// }
	// attributevalue.UnmarshalMap(result.Item, &user)
	//
	// g.JSON(http.StatusOK, gin.H{
	//     "user_id": userID,
	//     "addresses": user.Addresses,
	// })
	
	// ========================================
	// Mock Implementation (Development)
	// ========================================
	g.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"addresses": []interface{}{
			gin.H{
				"address_id": "addr_9h3k7j2m5n8p4q1r",
				"alias":      "Home",
				"line1":      "742 Evergreen Terrace",
				"line2":      "Apt 3B",
				"city":       "San Francisco",
				"state":      "CA",
				"zip_code":   "94102",
				"country":    "USA",
			},
		},
	})
}

// AddMyAddress adds a new address for the authenticated user
func AddMyAddress(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var addressData map[string]interface{}
	if err := g.BindJSON(&addressData); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// TODO: Validate and save address to database
	g.JSON(http.StatusCreated, gin.H{
		"user_id":   userID,
		"message":   "address added successfully",
		"address":   addressData,
	})
}

// UpdateMyAddress updates an address for the authenticated user
func UpdateMyAddress(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var updateData struct {
		AddrID string                 `json:"addr_id" binding:"required"`
		Data   map[string]interface{} `json:"data"`
	}
	
	if err := g.BindJSON(&updateData); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// TODO: Validate ownership and update address in database
	g.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"addr_id":   updateData.AddrID,
		"message":   "address updated successfully",
		"data":      updateData.Data,
	})
}

// DeleteMyAddress deletes an address for the authenticated user
func DeleteMyAddress(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		AddrID string `json:"addr_id" binding:"required"`
	}
	
	if err := g.BindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body or missing addr_id"})
		return
	}

	// TODO: Validate ownership and delete address from database
	g.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"addr_id": req.AddrID,
		"message": "address deleted successfully",
	})
}
