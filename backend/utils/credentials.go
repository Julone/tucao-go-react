package utils

import (
	"fmt"

	"tuxiaocao/pkg/repository"
)

// GetCredentialsByRole func for getting credentials from a role name.
func GetCredentialsByRole(role string) ([]string, error) {
	// Define credentials variable.
	var credentials []string

	// Switch given role.
	switch role {
	case repository.AdminRoleName:
		// Admin credentials (all access).
		credentials = []string{
			repository.ProductCreateCredential,
			repository.ProductUpdateCredential,
			repository.ProductDeleteCredential,
		}
	case repository.ModeratorRoleName:
		// Moderator credentials (only Product creation and update).
		credentials = []string{
			repository.ProductCreateCredential,
			repository.ProductUpdateCredential,
		}
	case repository.UserRoleName:
		// Simple user credentials (only Product creation).
		credentials = []string{
			repository.ProductCreateCredential,
		}
	default:
		// Return error message.
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return credentials, nil
}
