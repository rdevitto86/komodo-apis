package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Returns the authenticated user's profile
// POST /users/me with body: { "size": "basic" | "minimal" | "full" }
// Defaults to "basic" if size not specified
func GetMyProfile(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req struct {
		Size string `json:"size"` // "basic" | "minimal" | "full"
	}
	
	if err := g.ShouldBindJSON(&req); err != nil || req.Size == "" {
		req.Size = "basic"
	}
	if req.Size != "basic" && req.Size != "minimal" && req.Size != "full" {
		req.Size = "basic"
	}
	
	// ========================================
	// TODO: Production Implementation (DynamoDB)
	// ========================================
	// Same as GetUserByID, but use userID from session context
	// Query DynamoDB with projection based on size parameter
	
	// ========================================
	// Mock Implementation (Development)
	// ========================================
	
	// Determine which mock file to load based on size
	var mockFile string
	switch req.Size {
		case "basic":
			mockFile = "user_basic.json"
		case "minimal":
			mockFile = "user_minimal.json"
		case "full":
			mockFile = "user_full.json"
	}
	
	mockPath := filepath.Join("mocks", mockFile)
	data, err := os.ReadFile(mockPath)
	if err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to load mock data",
			"details": err.Error(),
		})
		return
	}
	
	var userData map[string]interface{}
	if err := json.Unmarshal(data, &userData); err != nil {
		g.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to parse mock data",
			"details": err.Error(),
		})
		return
	}
	
	userData["user_id"] = userID
	
	g.JSON(http.StatusOK, userData)
}

// Updates the authenticated user's profile
func UpdateMyProfile(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var updateData map[string]interface{}
	if err := g.BindJSON(&updateData); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// ========================================
	// TODO: Production Implementation (DynamoDB)
	// ========================================
	// import "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	//
	// ctx := context.Background()
	// tableName := "komodo-users"
	//
	// // Build update expression from updateData
	// updateExpr := "SET "
	// exprAttrNames := make(map[string]string)
	// exprAttrValues := make(map[string]types.AttributeValue)
	//
	// i := 0
	// for key, value := range updateData {
	//     if i > 0 {
	//         updateExpr += ", "
	//     }
	//     attrName := fmt.Sprintf("#attr%d", i)
	//     attrValue := fmt.Sprintf(":val%d", i)
	//     exprAttrNames[attrName] = key
	//     exprAttrValues[attrValue] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%v", value)}
	//     updateExpr += fmt.Sprintf("%s = %s", attrName, attrValue)
	//     i++
	// }
	//
	// // Add updated_at timestamp
	// updateExpr += ", #updated_at = :updated_at"
	// exprAttrNames["#updated_at"] = "updated_at"
	// exprAttrValues[":updated_at"] = &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)}
	//
	// _, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
	//     TableName: aws.String(tableName),
	//     Key: map[string]types.AttributeValue{
	//         "user_id": &types.AttributeValueMemberS{Value: userID.(string)},
	//     },
	//     UpdateExpression: aws.String(updateExpr),
	//     ExpressionAttributeNames: exprAttrNames,
	//     ExpressionAttributeValues: exprAttrValues,
	// })
	//
	// if err != nil {
	//     logger.Error("failed to update user in DynamoDB", err)
	//     g.JSON(http.StatusInternalServerError, gin.H{
	//         "error": "failed to update profile",
	//         "error_code": errCodes.ERR_DB_QUERY_FAILED,
	//     })
	//     return
	// }
	
	// ========================================
	// Mock Implementation (Development)
	// ========================================
	g.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "profile updated successfully",
		"data":    updateData,
	})
}

// Deletes the authenticated user's account
func DeleteMyAccount(g *gin.Context) {
	userID, exists := g.Get("user_id")
	if !exists {
		g.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// ========================================
	// TODO: Production Implementation (DynamoDB)
	// ========================================
	// Soft delete: Set deleted_at timestamp and account_status = "deleted"
	//
	// ctx := context.Background()
	// tableName := "komodo-users"
	//
	// _, err := dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
	//     TableName: aws.String(tableName),
	//     Key: map[string]types.AttributeValue{
	//         "user_id": &types.AttributeValueMemberS{Value: userID.(string)},
	//     },
	//     UpdateExpression: aws.String("SET account_status = :status, deleted_at = :deleted_at"),
	//     ExpressionAttributeValues: map[string]types.AttributeValue{
	//         ":status": &types.AttributeValueMemberS{Value: "deleted"},
	//         ":deleted_at": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	//     },
	// })
	//
	// if err != nil {
	//     logger.Error("failed to delete user in DynamoDB", err)
	//     g.JSON(http.StatusInternalServerError, gin.H{
	//         "error": "failed to delete account",
	//         "error_code": errCodes.ERR_DB_QUERY_FAILED,
	//     })
	//     return
	// }
	
	// ========================================
	// Mock Implementation (Development)
	// ========================================
	g.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"message": "account deleted successfully",
	})
}
