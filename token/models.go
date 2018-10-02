package token

import "github.com/ONSdigital/dp-identity-api/persistence"

//Service encapsulates the logic for creating, updating and deleting tokens
type Service struct {
	persistence.TokenStore
}
