package controllers

import (
	"context"
	"log"
	"net/http"
	"restaurant-management/databases"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemsPack struct {
	Table_id    string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = databases.OpenCollection(databases.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error listing ordered items"})
			return
		}
		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func GetOrderItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {

		orderId := c.Param("order_id")

		allOrderItem, err := ItemByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while listing order item by order id"})
			return
		}
		c.JSON(http.StatusOK, allOrderItem)
	}
}

func GetOrderItemsbyOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("order_item_id")
		var orderItem models.OrderItem

		err := orderCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error listing ordered item"})
			return
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func ItemByOrder(id string) (OrderItems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"food", "order"}, {"localfield", "order_id"}, {"foreign_field", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookUpTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localfield", "order.Table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food_price"},
			{"total_count", 1},
			{"food_name", "$food_name"},
			{"food_image", "$food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food_price"},
			{"quantity", 1},
		},
		},
	}
	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", ""}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1},
		}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		lookUpTableStage,
		unwindStage,
		unwindOrderStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}
	defer cancel()

	return OrderItems, err
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var (
			orderItemPack OrderItemsPack
			order         models.Order
		)

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}

		order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemstoInserted := []interface{}{}
		order.Table_id = &orderItemPack.Table_id
		// orderId := OrderItemCreator(order)

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order.Order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()
			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemstoInserted = append(orderItemstoInserted, orderItem)
		}
		insertedOrderItem, err := orderCollection.InsertMany(ctx, orderItemstoInserted)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.JSON(http.StatusOK, insertedOrderItem)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItems models.OrderItem
		orderItemId := c.Param("order_item_id")

		filter := bson.M{"order_item_id": orderItemId}

		var updateObj primitive.D

		if orderItems.Unit_price != nil {
			updateObj = append(updateObj, bson.E{"unit_price", *&orderItems.Unit_price})
		}
		if orderItems.Quantity != nil {
			updateObj = append(updateObj, bson.E{"quantity", *&orderItems.Quantity})
		}
		if orderItems.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", *&orderItems.Food_id})
		}

		orderItems.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItems.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)
		if err != nil {
			msg := "Order item create failed"
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
