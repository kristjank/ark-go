package arkgoserver

func initializeRoutes() {
	logger.Println("Initializing routes")

	// Group peer related routes together
	peerRoutes := router.Group("/voters")
	{
		// Handle the GET requests at /peer/list
		peerRoutes.GET("/list/", GetVoters)
	}

}
