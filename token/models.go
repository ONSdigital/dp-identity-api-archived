package token

import "github.com/ONSdigital/dp-identity-api/persistence"

//Service encapsulates the logic for storing and getting tokens
type Service struct {
	persistence.TokenStore
}
