package vdo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

type VideoCipherClient struct {
	url    string
	secret string
}

// UploadCredentials holds the response from VdoCipher for upload credentials
// Define the UploadCredentials struct with all required fields
type UploadCredentials struct {
	UploadURL      string `json:"uploadLink"`
	FileName       string `json:"key"`
	XAmzAlgorithm  string `json:"x-amz-algorithm"`
	XAmzCredential string `json:"x-amz-credential"`
	XAmzDate       string `json:"x-amz-date"`
	XAmzSignature  string `json:"x-amz-signature"`
	Policy         string `json:"policy"`
}

// UploadResponse holds the response from VdoCipher after uploading the video
type UploadResponse struct {
	VideoID string `json:"videoId"`
}

// UploadRequest holds the request body for getting upload credentials
type UploadRequest struct {
	Title string `json:"title"`
}

// NewVideoCipherClient initializes a new VdoCipher client using environment variables
func NewVideoCipherClient() *VideoCipherClient {
	return &VideoCipherClient{
		url:    os.Getenv("VDOCIPHER_URL"),
		secret: os.Getenv("VDOCIPHER_SECRET"),
	}
}

// CreateFolderRequest represents the request body for creating a folder
type CreateFolderRequest struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

// Folder represents the structure of a created folder
type Folder struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Parent       *string `json:"parent"`
	VideosCount  int     `json:"videosCount"`
	FoldersCount int     `json:"foldersCount"`
}

// CreateFolderRoot creates a new folder in the root directory
func (v *VideoCipherClient) CreateFolderRoot(name, parent string) (*Folder, error) {
	// Define the request payload
	createFolderReq := CreateFolderRequest{
		Name:   name,
		Parent: "root",
	}

	// Marshal the request body into JSON
	reqBody, err := json.Marshal(createFolderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	fmt.Println("Request Body:", string(reqBody))

	// Set the correct URL without a trailing slash
	url := v.url + "/videos/folders"
	fmt.Println("Request URL:", url) // Debugging URL

	// Prepare the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers for VdoCipher API request
	req.Header.Set("Authorization", "Apisecret "+v.secret)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Debug response status and body if not successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body))
		return nil, fmt.Errorf("unexpected status code from VdoCipher: %v", resp.Status)
	}

	// Decode the response body into the Folder struct
	var folder Folder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &folder, nil
}

// FolderListResponse represents the structure of the response for folder list

type FolderListResponse struct {
	FolderList []Folder `json:"folderList"`
}

// GetAllFolders fetches all folders available in VdoCipher
func (v *VideoCipherClient) GetAllFolders() (*FolderListResponse, error) {
	// Set the correct URL for fetching folders
	url := v.url + "/videos/folders/root"
	fmt.Println("Request URL:", url) // Debugging URL

	// Prepare the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the internal VdoCipher API secret in the headers
	req.Header.Set("Authorization", "Apisecret "+v.secret)
	req.Header.Set("Accept", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-OK status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body)) // Log the raw response body
		return nil, fmt.Errorf("unexpected status code from VdoCipher: %v", resp.Status)
	}

	// Log the full response body to check what exactly is returned
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw Response Body:", string(body))

	// Decode the response body into the FolderListResponse struct
	var folderList FolderListResponse
	if err := json.Unmarshal(body, &folderList); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &folderList, nil
}

func (v *VideoCipherClient) DeleteFolder(folderID string) error {
	// Set the correct URL for deleting the folder
	url := fmt.Sprintf("%s/videos/folders/%s", v.url, folderID)
	fmt.Println("Request URL:", url) // Debugging URL

	// Prepare the DELETE request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set the internal VdoCipher API secret in the headers
	req.Header.Set("Authorization", "Apisecret "+v.secret)
	req.Header.Set("Accept", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-OK status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body)) // Log the raw response body
		return fmt.Errorf("unexpected status code from VdoCipher: %v", resp.Status)
	}

	return nil
}

func (v *VideoCipherClient) GetUploadCredentials(title string, folderID string) (*UploadCredentials, error) {
	// Construct the request URL with title and optional folder ID
	reqURL := fmt.Sprintf("%s/videos?title=%s", v.url, url.QueryEscape(title))
	if folderID != "" {
		reqURL += fmt.Sprintf("&folderId=%s", url.QueryEscape(folderID))
	}

	req, err := http.NewRequest("PUT", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Apisecret "+v.secret)
	req.Header.Set("Accept", "application/json") // Set Accept header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	log.Printf("Response Body: %s\n", string(body))

	// Check for non-OK status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v, response: %s", resp.Status, string(body))
	}

	// Parse the response body
	var response struct {
		ClientPayload UploadCredentials `json:"clientPayload"`
		VideoId       string            `json:"videoId"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Create UploadCredentials struct
	credentials := &UploadCredentials{
		UploadURL:      response.ClientPayload.UploadURL,
		XAmzAlgorithm:  response.ClientPayload.XAmzAlgorithm,
		XAmzCredential: response.ClientPayload.XAmzCredential,
		XAmzDate:       response.ClientPayload.XAmzDate,
		XAmzSignature:  response.ClientPayload.XAmzSignature,
		Policy:         response.ClientPayload.Policy,
		FileName:       response.ClientPayload.FileName,
	}

	log.Printf("Received credentials: %+v\n", credentials)

	return credentials, nil
}

// UploadFile uploads a file to S3 using the provided credentials
func (client *VideoCipherClient) UploadFile(credentials UploadCredentials, filePath string) error {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Prepare a new buffer for multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the form fields (based on your curl example)
	_ = writer.WriteField("key", credentials.FileName)
	_ = writer.WriteField("policy", credentials.Policy)
	_ = writer.WriteField("x-amz-algorithm", credentials.XAmzAlgorithm)
	_ = writer.WriteField("x-amz-credential", credentials.XAmzCredential)
	_ = writer.WriteField("x-amz-date", credentials.XAmzDate)
	_ = writer.WriteField("x-amz-signature", credentials.XAmzSignature)
	_ = writer.WriteField("success_action_status", "200")
	_ = writer.WriteField("success_action_redirect", "") // Add the empty value for success_action_redirect as required by the policy

	// Add the file part
	fileWriter, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the file content to the form
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", credentials.UploadURL, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set the Content-Type to multipart/form-data with the boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute the request using a new HTTP client
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	// Check for successful upload
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read the response body for error details
		return fmt.Errorf("failed to upload file, unexpected status: %v, response: %s", resp.Status, string(body))
	}

	return nil
}

func (v *VideoCipherClient) CreateSubFolder(name, parent string) (*Folder, error) {
	// Define the request payload
	createFolderReq := CreateFolderRequest{
		Name:   name,
		Parent: parent, // This will now be the provided parent folder ID for subfolders
	}

	// Marshal the request body into JSON
	reqBody, err := json.Marshal(createFolderReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	fmt.Println("Request Body:", string(reqBody))

	// Set the correct URL without a trailing slash
	url := v.url + "/videos/folders"
	fmt.Println("Request URL:", url) // Debugging URL

	// Prepare the POST request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers for VdoCipher API request
	req.Header.Set("Authorization", "Apisecret "+v.secret)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Debug response status and body if not successful
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body))
		return nil, fmt.Errorf("unexpected status code from VdoCipher: %v", resp.Status)
	}

	// Decode the response body into the Folder struct
	var folder Folder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &folder, nil
}

type SubFolder struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ParentID     string `json:"parent"`
	VideosCount  int    `json:"videosCount"`
	FoldersCount int    `json:"foldersCount"`
}

type FolderResponse struct {
	FolderList []SubFolder `json:"folderList"` // List of subfolders
	Current    SubFolder   `json:"current"`    // The current folder
	Parent     SubFolder   `json:"parent"`     // The parent folder
}

func (v *VideoCipherClient) GetSubFolders(folderID string) (*FolderResponse, error) {
	// Use the folderID directly in the path
	url := fmt.Sprintf("%s/videos/folders/%s", v.url, folderID)

	// Prepare the GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the necessary headers
	req.Header.Set("Authorization", "Apisecret "+v.secret)
	req.Header.Set("Accept", "application/json")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-OK status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response status:", resp.Status)
		fmt.Println("Response body:", string(body))
		return nil, fmt.Errorf("unexpected status code from VdoCipher: %v", resp.Status)
	}

	// Log raw response body for debugging
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response body:", string(body))

	// Decode the response body into the FolderResponse structure
	var folderResponse FolderResponse
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&folderResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Return the fetched folder data
	return &folderResponse, nil
}
