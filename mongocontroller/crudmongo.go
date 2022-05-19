package mongocontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

// struct for storing data
type user struct {
	Name string `json:name`
	Age  int    `json:age`
	City string `json:city`
}

type HttpMongo struct {
	db *mongo.Database
}

func New(db *mongo.Database) *HttpMongo {
	return &HttpMongo{db}
}

func (h *HttpMongo) Mount(echoGroup *echo.Group) {
	echoGroup.GET("/getAllUsers", h.GetAllUsers)
	echoGroup.POST("/createProfile", h.CreateProfile)
	echoGroup.POST("/getUserProfile", h.GetUserProfile)
	echoGroup.PUT("/updateProfile/:id", h.UpdateProfile)
	echoGroup.DELETE("/deleteProfile/:id", h.DeleteProfile)
}

// Create Profile or Signup
func (h *HttpMongo) GetAllUsers(c echo.Context) error {
	var results []primitive.M //slice for multiple documents
	var userCollection = h.db
	cur, err := userCollection.Collection("mahasiswa").Find(context.TODO(), bson.D{{}}) //returns a *mongo.Cursor
	if err != nil {
		fmt.Println(err)

	}
	for cur.Next(context.TODO()) { //Next() gets the next document for corresponding cursor

		var elem primitive.M
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem) // appending document pointed by Next()
	}
	cur.Close(context.TODO()) // close the cursor once stream of documents has exhausted
	jsonData, err := json.Marshal(results)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, string(jsonData))

}

// Get Profile of a particular User by Name

func (h *HttpMongo) GetUserProfile(c echo.Context) error {

	var body user
	e := json.NewDecoder(c.Request().Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}
	var result primitive.M //  an unordered representation of a BSON document which is a Map
	var userCollection = h.db
	err := userCollection.Collection("mahasiswa").FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)
	if err != nil {

		fmt.Println(err)

	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, string(jsonData))
}

func (h *HttpMongo) CreateProfile(c echo.Context) error {

	var person user
	var userCollection = h.db
	err := json.NewDecoder(c.Request().Body).Decode(&person) // storing in person variable of type user
	if err != nil {
		fmt.Print(err)
	}
	insertResult, err := userCollection.Collection("mahasiswa").InsertOne(context.TODO(), person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	jsonData, err := json.Marshal(insertResult.InsertedID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, string(jsonData))

}

//Update Profile of User

func (h *HttpMongo) UpdateProfile(c echo.Context) error {
	id := c.Param("id")
	var body user
	e := json.NewDecoder(c.Request().Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
	}
	after := options.After // for returning updated document
	returnOpt := options.FindOneAndUpdateOptions{

		ReturnDocument: &after,
	}
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	}
	var userCollection = h.db
	toDocBody, err := toDoc(body)
	if err != nil {
		log.Fatal(" ERROR:", err)
	}
	update := bson.D{{"$set", toDocBody}}
	//update := bson.M{"$set": bson.M{"name": body.Name, "age": body.Age, "city": body.City}}
	updateResult := userCollection.Collection("mahasiswa").FindOneAndUpdate(context.TODO(), bson.D{{"_id", idPrimitive}}, update, &returnOpt)

	var result primitive.M
	_ = updateResult.Decode(&result)

	jsonData, err := json.Marshal(result)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, string(jsonData))
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

//Delete Profile of User

func (h *HttpMongo) DeleteProfile(c echo.Context) error {
	id := c.Param("id")
	opts := options.Delete().SetCollation(&options.Collation{}) // to specify language-specific rules for string comparison, such as rules for lettercase
	var userCollection = h.db
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	}
	res, err := userCollection.Collection("mahasiswa").DeleteOne(context.TODO(), bson.D{{"_id", idPrimitive}}, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	jsonData, err := json.Marshal(res.DeletedCount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, string(jsonData))

}
