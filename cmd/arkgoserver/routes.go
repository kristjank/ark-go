package main

func initializeRoutes() {
	logger.Println("Initializing routes")

	// Group peer related routes together
	peerRoutes := router.Group("/voters")
	{
		peerRoutes.GET("/rewards/", GetVoters)
		peerRoutes.GET("/blocked/", GetBlocked)
	}
	deleRoutes := router.Group("/delegate")
	{
		deleRoutes.GET("/", GetDelegate)
		deleRoutes.GET("/config", GetDelegateSharingConfig)
	}
}
