package models

type UserCredentials struct {
	Email        string `json:"email"`
	PrivateKey   string `json:"privateKey"`
	JTWToken     string `json:"JTWToken"`
	RefreshToken string `json:"RefreshToken"`
}

// The file struct for Gsoc2 config file
type ConfigFile struct {
	LoggedInUserEmail  string         `json:"loggedInUserEmail"`
	LoggedInUserDomain string         `json:"LoggedInUserDomain,omitempty"`
	LoggedInUsers      []LoggedInUser `json:"loggedInUsers,omitempty"`
	VaultBackendType   string         `json:"vaultBackendType,omitempty"`
}

type LoggedInUser struct {
	Email  string `json:"email"`
	Domain string `json:"domain"`
}

type SingleEnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
	ID    string `json:"_id"`
	Tags  []struct {
		ID        string `json:"_id"`
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		Workspace string `json:"workspace"`
	} `json:"tags"`
	Comment string `json:"comment"`
}

type SingleFolder struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Workspace struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	Plan         string `json:"plan,omitempty"`
	V            int    `json:"__v"`
	Organization string `json:"organization,omitempty"`
}

type WorkspaceConfigFile struct {
	WorkspaceId                   string            `json:"workspaceId"`
	DefaultEnvironment            string            `json:"defaultEnvironment"`
	GitBranchToEnvironmentMapping map[string]string `json:"gitBranchToEnvironmentMapping"`
}

type SymmetricEncryptionResult struct {
	CipherText []byte `json:"CipherText"`
	Nonce      []byte `json:"Nonce"`
	AuthTag    []byte `json:"AuthTag"`
}

type GetAllSecretsParameters struct {
	Environment              string
	EnvironmentPassedViaFlag bool
	Gsoc2Token           string
	TagSlugs                 string
	WorkspaceId              string
	SecretsPath              string
	IncludeImport            bool
}

type GetAllFoldersParameters struct {
	WorkspaceId    string
	Environment    string
	FoldersPath    string
	Gsoc2Token string
}

type CreateFolderParameters struct {
	FolderName     string
	WorkspaceId    string
	Environment    string
	FolderPath     string
	Gsoc2Token string
}

type DeleteFolderParameters struct {
	FolderName     string
	WorkspaceId    string
	Environment    string
	FolderPath     string
	Gsoc2Token string
}
