package main

func initializeRoutes() {
	logger.Println("Initializing routes")

	// Group peer related routes together
	peerRoutes := router.Group("/voters")
	{
		peerRoutes.GET("/list/", GetVoters)
	}
	deleRoutes := router.Group("/delegate")
	{
		deleRoutes.GET("/", GetDelegate)
	}
}
