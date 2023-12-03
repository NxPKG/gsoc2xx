package visualize

import "github.com/Gsoc2/gsoc2-merge/packages/models"

func PrintAllSecretDetails(secrets []models.SingleEnvironmentVariable) {
	rows := [][3]string{}
	for _, secret := range secrets {
		rows = append(rows, [...]string{secret.Key, secret.Value, secret.Type})
	}

	headers := [...]string{"SECRET NAME", "SECRET VALUE", "SECRET TYPE"}

	Table(headers, rows)
}
