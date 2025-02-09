package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/Gsoc2/gsoc2-merge/packages/api"
	"github.com/Gsoc2/gsoc2-merge/packages/models"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

func GetAllFolders(params models.GetAllFoldersParameters) ([]models.SingleFolder, error) {

	if params.Gsoc2Token == "" {
		params.Gsoc2Token = os.Getenv(GSOC2_TOKEN_NAME)
	}

	var foldersToReturn []models.SingleFolder
	var folderErr error
	if params.Gsoc2Token == "" {

		log.Debug().Msg("GetAllFolders: Trying to fetch folders using logged in details")

		loggedInUserDetails, err := GetCurrentLoggedInUserDetails()
		if err != nil {
			return nil, err
		}

		if loggedInUserDetails.LoginExpired {
			PrintErrorMessageAndExit("Your login session has expired, please run [gsoc2 login] and try again")
		}

		workspaceFile, err := GetWorkSpaceFromFile()
		if err != nil {
			return nil, err
		}

		if params.WorkspaceId != "" {
			workspaceFile.WorkspaceId = params.WorkspaceId
		}

		folders, err := GetFoldersViaJTW(loggedInUserDetails.UserCredentials.JTWToken, workspaceFile.WorkspaceId, params.Environment, params.FoldersPath)
		folderErr = err
		foldersToReturn = folders
	} else {
		// get folders via service token
		folders, err := GetFoldersViaServiceToken(params.Gsoc2Token, params.WorkspaceId, params.Environment, params.FoldersPath)
		folderErr = err
		foldersToReturn = folders
	}
	return foldersToReturn, folderErr
}

func GetFoldersViaJTW(JTWToken string, workspaceId string, environmentName string, foldersPath string) ([]models.SingleFolder, error) {
	// set up resty client
	httpClient := resty.New()
	httpClient.SetAuthToken(JTWToken).
		SetHeader("Accept", "application/json")

	getFoldersRequest := api.GetFoldersV1Request{
		WorkspaceId: workspaceId,
		Environment: environmentName,
		FoldersPath: foldersPath,
	}

	apiResponse, err := api.CallGetFoldersV1(httpClient, getFoldersRequest)
	if err != nil {
		return nil, err
	}

	var folders []models.SingleFolder

	for _, folder := range apiResponse.Folders {
		folders = append(folders, models.SingleFolder{
			Name: folder.Name,
			ID:   folder.ID,
		})
	}

	return folders, nil
}

func GetFoldersViaServiceToken(fullServiceToken string, workspaceId string, environmentName string, foldersPath string) ([]models.SingleFolder, error) {
	serviceTokenParts := strings.SplitN(fullServiceToken, ".", 4)
	if len(serviceTokenParts) < 4 {
		return nil, fmt.Errorf("invalid service token entered. Please double check your service token and try again")
	}

	serviceToken := fmt.Sprintf("%v.%v.%v", serviceTokenParts[0], serviceTokenParts[1], serviceTokenParts[2])

	httpClient := resty.New()

	httpClient.SetAuthToken(serviceToken).
		SetHeader("Accept", "application/json")

	serviceTokenDetails, err := api.CallGetServiceTokenDetailsV2(httpClient)
	if err != nil {
		return nil, fmt.Errorf("unable to get service token details. [err=%v]", err)
	}

	// if multiple scopes are there then user needs to specify which environment and folder path
	if environmentName == "" {
		if len(serviceTokenDetails.Scopes) != 1 {
			return nil, fmt.Errorf("you need to provide the --env for multiple environment scoped token")
		} else {
			environmentName = serviceTokenDetails.Scopes[0].Environment
		}
	}

	getFoldersRequest := api.GetFoldersV1Request{
		WorkspaceId: serviceTokenDetails.Workspace,
		Environment: environmentName,
		FoldersPath: foldersPath,
	}

	apiResponse, err := api.CallGetFoldersV1(httpClient, getFoldersRequest)
	if err != nil {
		return nil, fmt.Errorf("unable to get folders. [err=%v]", err)
	}

	var folders []models.SingleFolder

	for _, folder := range apiResponse.Folders {
		folders = append(folders, models.SingleFolder{
			Name: folder.Name,
			ID:   folder.ID,
		})
	}

	return folders, nil
}

// CreateFolder creates a folder in Gsoc2
func CreateFolder(params models.CreateFolderParameters) (models.SingleFolder, error) {
	loggedInUserDetails, err := GetCurrentLoggedInUserDetails()
	if err != nil {
		return models.SingleFolder{}, err
	}

	if loggedInUserDetails.LoginExpired {
		PrintErrorMessageAndExit("Your login session has expired, please run [gsoc2 login] and try again")
	}

	// set up resty client
	httpClient := resty.New()
	httpClient.
		SetAuthToken(loggedInUserDetails.UserCredentials.JTWToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	createFolderRequest := api.CreateFolderV1Request{
		WorkspaceId: params.WorkspaceId,
		Environment: params.Environment,
		FolderName:  params.FolderName,
		Directory:   params.FolderPath,
	}

	apiResponse, err := api.CallCreateFolderV1(httpClient, createFolderRequest)
	if err != nil {
		return models.SingleFolder{}, err
	}

	folder := apiResponse.Folder

	return models.SingleFolder{
		Name: folder.Name,
		ID:   folder.ID,
	}, nil
}

func DeleteFolder(params models.DeleteFolderParameters) ([]models.SingleFolder, error) {
	loggedInUserDetails, err := GetCurrentLoggedInUserDetails()
	if err != nil {
		return nil, err
	}

	if loggedInUserDetails.LoginExpired {
		PrintErrorMessageAndExit("Your login session has expired, please run [gsoc2 login] and try again")
	}

	// set up resty client
	httpClient := resty.New()
	httpClient.
		SetAuthToken(loggedInUserDetails.UserCredentials.JTWToken).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json")

	deleteFolderRequest := api.DeleteFolderV1Request{
		WorkspaceId: params.WorkspaceId,
		Environment: params.Environment,
		FolderName:  params.FolderName,
		Directory:   params.FolderPath,
	}

	apiResponse, err := api.CallDeleteFolderV1(httpClient, deleteFolderRequest)
	if err != nil {
		return nil, err
	}

	var folders []models.SingleFolder

	for _, folder := range apiResponse.Folders {
		folders = append(folders, models.SingleFolder{
			Name: folder.Name,
			ID:   folder.ID,
		})
	}

	return folders, nil
}
