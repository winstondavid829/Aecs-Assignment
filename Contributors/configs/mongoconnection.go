package configs

// // DBName of the database.
// const (
// 	///////////// Dev Mongo Environemnt Start////////////////////////////////////

// 	URI = "mongodb+srv://winston:Win12345@cluster0.ggapjbv.mongodb.net/?retryWrites=true&w=majority"

// ///////////// Dev Mongo Environemnt End////////////////////////////////////

// ///////////////// Production Mongo Environemnt Start//////////////////

// // URI    = "mongodb+srv://winston:Win12345@cluster0.ggapjbv.mongodb.net/?retryWrites=true&w=majority"

// ///////////////// Production Mongo Environemnt End //////////////////

// )

// func ConnectDB() *mongo.Client {
// 	// Set up MongoDB connection options
// 	clientOptions := options.Client().ApplyURI(URI)

// 	// Create a context with a timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Create a MongoDB client
// 	client, err := mongo.Connect(ctx, clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 		log.Println("[(ConnectDB): Cannot create a client to connect to MongoDB (err):]", err)
// 	}

// 	// Ping the MongoDB server to check if the connection is successful
// 	err = client.Ping(ctx, nil)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		log.Println("[(ConnectDB): Cannot ping MongoDB (err):]", err)
// 	}

// 	log.Println("[(ConnectDB): Connection Established to MongoDB]")
// 	return client
// }

// // Client instance
// var DB *mongo.Client = ConnectDB()

// // getting database collections
// func GetCollection(client *mongo.Client, DBName string, collectionName string) *mongo.Collection {
// 	collection := client.Database(DBName).Collection(collectionName)
// 	return collection
// }
