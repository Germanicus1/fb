package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Germanicus1/fb/models"
)

const (
	restDirectoryBaseURL = "https://fb.mauvable.com/rest-directory/2"
	httpTimeout          = 30 * time.Second
)

// HTTP constants
const (
	httpMethodGET        = "GET"
	headerAuthorization  = "Authorization"
	headerContentType    = "Content-Type"
	contentTypeJSON      = "application/json"
	authorizationPrefix  = "bearer "
	httpStatusOK         = 200
	httpStatusMultipleOK = 300
)

// Client is the Flow Boards API client
type Client struct {
	authKey    string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client with the provided authentication key
func NewClient(authKey string) *Client {
	return &Client{
		authKey:    authKey,
		httpClient: createHTTPClient(),
	}
}

// createHTTPClient creates a configured HTTP client with timeout
func createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: httpTimeout,
	}
}

// DiscoverRestPrefix discovers the REST API prefix for the organization
func (c *Client) DiscoverRestPrefix(orgID string) error {
	discoveryURL := buildRestDirectoryURL(orgID)

	resp, err := c.doRequestWithoutBase(httpMethodGET, discoveryURL, nil)
	if err != nil {
		return fmt.Errorf("failed to discover REST prefix: %w", err)
	}

	prefixResp, err := parseRestPrefixResponse(resp)
	if err != nil {
		return err
	}

	if prefixResp.RestPrefix == "" {
		return fmt.Errorf("REST prefix not found in response")
	}

	c.baseURL = prefixResp.RestPrefix
	return nil
}

// buildRestDirectoryURL constructs the REST directory discovery URL
func buildRestDirectoryURL(orgID string) string {
	return fmt.Sprintf("%s/%s", restDirectoryBaseURL, orgID)
}

// parseRestPrefixResponse parses the REST prefix discovery response
func parseRestPrefixResponse(data []byte) (*models.RestPrefixResponse, error) {
	var prefixResp models.RestPrefixResponse
	if err := json.Unmarshal(data, &prefixResp); err != nil {
		return nil, fmt.Errorf("failed to parse REST prefix response: %w", err)
	}
	return &prefixResp, nil
}

// GetCurrentUser retrieves the user information by email
func (c *Client) GetCurrentUser(email string) (*models.User, error) {
	if err := c.requireBaseURL(); err != nil {
		return nil, err
	}

	path := buildUserPath(email)
	resp, err := c.doRequest(httpMethodGET, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := parseUserResponse(resp)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// requireBaseURL checks if the base URL has been discovered
func (c *Client) requireBaseURL() error {
	if c.baseURL == "" {
		return fmt.Errorf("REST prefix not discovered, call DiscoverRestPrefix first")
	}
	return nil
}

// buildUserPath constructs the user lookup API path
func buildUserPath(email string) string {
	return fmt.Sprintf("/users/%s", url.PathEscape(email))
}

// parseUserResponse parses the user API response (single user object)
func parseUserResponse(data []byte) (*models.User, error) {
	var user models.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}
	return &user, nil
}

// SearchTickets searches for tickets assigned to the given user IDs
func (c *Client) SearchTickets(userIDs []string) ([]models.Ticket, error) {
	return c.SearchTicketsWithFilters(userIDs, "", "")
}

// SearchTicketsWithFilters searches for tickets with optional bin and board filters
func (c *Client) SearchTicketsWithFilters(userIDs []string, binID, boardID string) ([]models.Ticket, error) {
	if err := c.requireBaseURL(); err != nil {
		return nil, err
	}

	path := buildTicketSearchPathWithFilters(userIDs, binID, boardID)

	resp, err := c.doRequest(httpMethodGET, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search tickets: %w", err)
	}

	tickets, err := parseTicketSearchResponse(resp)
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// buildTicketSearchPath constructs the ticket search API path with comma-separated user IDs
func buildTicketSearchPath(userIDs []string) string {
	return buildTicketSearchPathWithFilters(userIDs, "", "")
}

// buildTicketSearchPathWithFilters constructs the ticket search API path with filters
func buildTicketSearchPathWithFilters(userIDs []string, binID, boardID string) string {
	params := []string{}

	if len(userIDs) > 0 {
		usersParam := strings.Join(userIDs, ",")
		params = append(params, fmt.Sprintf("users=%s", url.QueryEscape(usersParam)))
	}

	if binID != "" {
		params = append(params, fmt.Sprintf("bins=%s", url.QueryEscape(binID)))
	}

	if boardID != "" {
		params = append(params, fmt.Sprintf("boards=%s", url.QueryEscape(boardID)))
	}

	return "/ticket-search?" + strings.Join(params, "&")
}

// parseTicketSearchResponse parses the ticket search API response
func parseTicketSearchResponse(data []byte) ([]models.Ticket, error) {
	// The API returns an array of tickets directly
	var tickets []models.Ticket
	if err := json.Unmarshal(data, &tickets); err != nil {
		return nil, fmt.Errorf("failed to parse ticket response: %w", err)
	}
	return tickets, nil
}

// GetBins retrieves all bins from the API
func (c *Client) GetBins() ([]models.Bin, error) {
	if err := c.requireBaseURL(); err != nil {
		return nil, err
	}

	var allBins []models.Bin
	pageToken := ""

	for {
		path := buildPaginatedPath("/bins", pageToken)

		resp, err := c.doRequest(httpMethodGET, path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get bins: %w", err)
		}

		bins, nextToken, err := parseBinsPage(resp)
		if err != nil {
			return nil, err
		}

		allBins = append(allBins, bins...)

		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}

	return allBins, nil
}

// LookupBinIDByName looks up a bin ID by name (case-insensitive)
func (c *Client) LookupBinIDByName(binName string) (string, error) {
	bins, err := c.GetBins()
	if err != nil {
		return "", err
	}

	lowerBinName := strings.ToLower(binName)
	for _, bin := range bins {
		if strings.ToLower(bin.Name) == lowerBinName {
			return bin.ID, nil
		}
	}

	return "", fmt.Errorf("bin not found: %s", binName)
}

// GetBoards retrieves all boards from the API
func (c *Client) GetBoards() ([]models.Board, error) {
	if err := c.requireBaseURL(); err != nil {
		return nil, err
	}

	var allBoards []models.Board
	pageToken := ""

	for {
		path := buildPaginatedPath("/boards", pageToken)

		resp, err := c.doRequest(httpMethodGET, path, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to get boards: %w", err)
		}

		boards, nextToken, err := parseBoardsPage(resp)
		if err != nil {
			return nil, err
		}

		allBoards = append(allBoards, boards...)

		if nextToken == "" {
			break
		}
		pageToken = nextToken
	}

	return allBoards, nil
}

// LookupBoardIDByName looks up a board ID by name (case-insensitive)
func (c *Client) LookupBoardIDByName(boardName string) (string, error) {
	boards, err := c.GetBoards()
	if err != nil {
		return "", err
	}

	lowerBoardName := strings.ToLower(boardName)
	for _, board := range boards {
		if strings.ToLower(board.Name) == lowerBoardName {
			return board.ID, nil
		}
	}

	return "", fmt.Errorf("board not found: %s", boardName)
}

// doRequest makes an HTTP request with authentication using the base URL
func (c *Client) doRequest(method, path string, body io.Reader) ([]byte, error) {
	fullURL := c.baseURL + path
	return c.doRequestWithoutBase(method, fullURL, body)
}

// doRequestWithoutBase makes an HTTP request with authentication without using the base URL
func (c *Client) doRequestWithoutBase(method, fullURL string, body io.Reader) ([]byte, error) {
	req, err := c.createRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	if err := checkStatusCode(resp.StatusCode, respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

// createRequest creates an HTTP request with authentication headers
func (c *Client) createRequest(method, fullURL string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setAuthHeaders(req)
	return req, nil
}

// setAuthHeaders sets the authorization and content-type headers
func (c *Client) setAuthHeaders(req *http.Request) {
	req.Header.Set(headerAuthorization, authorizationPrefix+c.authKey)
	req.Header.Set(headerContentType, contentTypeJSON)
}

// executeRequest executes an HTTP request
func (c *Client) executeRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

// readResponseBody reads the response body into a byte slice
func readResponseBody(resp *http.Response) ([]byte, error) {
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	return respBody, nil
}

// checkStatusCode validates the HTTP status code is in the 2xx range
func checkStatusCode(statusCode int, respBody []byte) error {
	if statusCode < httpStatusOK || statusCode >= httpStatusMultipleOK {
		return fmt.Errorf("API request failed with status %d: %s", statusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}

// buildPaginatedPath constructs a paginated API path with max-results and optional page-token
func buildPaginatedPath(basePath string, pageToken string) string {
	path := basePath + "?max-results=1000"
	if pageToken != "" {
		path += "&page-token=" + url.QueryEscape(pageToken)
	}
	return path
}

// parseBinsPage attempts to parse a page of bins from the API response
// Returns the bins, next page token, and error
func parseBinsPage(data []byte) ([]models.Bin, string, error) {
	// Try parsing as paginated response first
	var paginatedResp models.BinsResponse
	if err := json.Unmarshal(data, &paginatedResp); err == nil && paginatedResp.Results != nil {
		return paginatedResp.Results, paginatedResp.PageToken, nil
	}

	// Fall back to old format (direct array)
	var bins []models.Bin
	if err := json.Unmarshal(data, &bins); err != nil {
		return nil, "", fmt.Errorf("failed to parse bins response: %w", err)
	}
	return bins, "", nil
}

// parseBoardsPage attempts to parse a page of boards from the API response
// Returns the boards, next page token, and error
func parseBoardsPage(data []byte) ([]models.Board, string, error) {
	// Try parsing as paginated response first
	var paginatedResp models.BoardsResponse
	if err := json.Unmarshal(data, &paginatedResp); err == nil && paginatedResp.Results != nil {
		return paginatedResp.Results, paginatedResp.PageToken, nil
	}

	// Fall back to old format (direct array)
	var boards []models.Board
	if err := json.Unmarshal(data, &boards); err != nil {
		return nil, "", fmt.Errorf("failed to parse boards response: %w", err)
	}
	return boards, "", nil
}

// PostComment posts a comment to a ticket
func (c *Client) PostComment(payload models.CommentPayload) error {
	if err := c.requireBaseURL(); err != nil {
		return err
	}

	path := fmt.Sprintf("/ticket-comments/%s", url.PathEscape(payload.ID))

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal comment payload: %w", err)
	}

	_, err = c.doRequest("POST", path, strings.NewReader(string(jsonData)))
	if err != nil {
		return fmt.Errorf("failed to post comment: %w", err)
	}

	return nil
}
