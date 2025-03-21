package appwrite

import (
	"github.com/abhisheksharm-3/shrtn/internal/config"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
)

func GetClient(cfg *config.Config) client.Client {
	var appwriteClient = appwrite.NewClient(
		appwrite.WithProject(cfg.AppwriteProjectID),
		appwrite.WithKey(cfg.AppwriteAPIKey),
	)
	return appwriteClient
}
