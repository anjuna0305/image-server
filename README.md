# Image Server

A secure image server built with Go (Gin framework) that provides image upload, retrieval, update, and deletion capabilities with HMAC-signed URL authentication.

## Features

- **Secure Image Upload** - Upload images with automatic UUID-based filename generation
- **Signed URL Authentication** - HMAC-SHA256 signed URLs with expiration times
- **Method-Specific Tokens** - Each HTTP method (GET, PUT, DELETE, POST) requires its own token for security
- **Image Management** - Support for GET, POST, PUT, and DELETE operations
- **Automatic MIME Type Detection** - Content-Type headers set automatically based on file extensions
- **Token Generation Tool** - JavaScript utility to generate signed URLs from command line

## Requirements

- Go 1.24.4 or later
- Node.js (for the signed URL generation script)

## Configuration

The server uses the following configuration (located in `main.go`):

- **Upload Directory**: `/home/anjuna/kethaka/imageServer/uploads` (default)
- **Secret Key**: `secret-key` (default - **CHANGE THIS IN PRODUCTION**)
- **Server Port**: `:8000` (default)
- **Base URL**: `http://localhost:8000` (for token generation)

**Important**: Change the `secretKey` constant in `main.go` before deploying to production!

## Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8000`

## API Endpoints

### Health Check
```
GET /
```
Returns server status.

### Upload Image
```
POST /images
```
Uploads a new image file. Requires a signed URL token.

**Request**: `multipart/form-data` with `file` field

**Response**:
```json
{
  "message": "File uploaded",
  "filename": "uuid-here.jpg",
  "original_filename": "original.jpg",
  "size": 12345
}
```

### Retrieve Image
```
GET /images/:filename
```
Retrieves an image file. Requires a signed URL token specific to GET method.

**Query Parameters**:
- `expires`: Unix timestamp for expiration
- `signature`: HMAC-SHA256 signature

### Update Image
```
PUT /images/:filename
```
Updates an existing image file. Requires a signed URL token specific to PUT method.

**Request**: `multipart/form-data` with `file` field

**Response**:
```json
{
  "message": "File updated",
  "size": 12345
}
```

### Delete Image
```
DELETE /images/:filename
```
Deletes an image file. Requires a signed URL token specific to DELETE method.

**Response**:
```json
{
  "message": "File removed"
}
```

## Signed URL Generation

Use the provided JavaScript script to generate signed URLs for secure access.

### Installation
The script is ready to use - it uses Node.js built-in modules (no npm install required).

### Usage

#### For GET requests (view only):
```bash
node generate-signed-url.js --get <image-name> <time-in-seconds>
# or short form:
node generate-signed-url.js -g <image-name> <time-in-seconds>
```

#### For PUT requests (update):
```bash
node generate-signed-url.js --put <image-name> <time-in-seconds>
# or short form:
node generate-signed-url.js -u <image-name> <time-in-seconds>
```

#### For DELETE requests:
```bash
node generate-signed-url.js --delete <image-name> <time-in-seconds>
# or short form:
node generate-signed-url.js -d <image-name> <time-in-seconds>
```

#### For POST requests (upload):
```bash
node generate-signed-url.js --post <time-in-seconds>
# or short form:
node generate-signed-url.js -p <time-in-seconds>
```

### Examples

```bash
# Generate GET token valid for 1 hour (3600 seconds)
node generate-signed-url.js --get 4d030191-3362-4491-8fab-f9c4c6ef17e0.jpg 3600

# Generate DELETE token valid for 30 minutes (1800 seconds)
node generate-signed-url.js -d myimage.jpg 1800

# Generate POST token valid for 1 hour
node generate-signed-url.js --post 3600
```

The script outputs a complete URL with `expires` and `signature` query parameters.

## Security Features

### Method-Specific Tokens
Each HTTP method requires its own signed token. A GET token cannot be reused for PUT or DELETE operations, preventing unauthorized modifications.

### Token Expiration
All tokens have an expiration time (Unix timestamp). Expired tokens are automatically rejected.

### HMAC-SHA256 Signing
All signed URLs use HMAC-SHA256 with a secret key. The signature includes:
- HTTP method (GET, PUT, DELETE, POST)
- Filename (empty for POST)
- Expiration timestamp

Format: `METHOD:filename:expires`

## Example Usage

### Upload an Image

1. Generate a POST token:
```bash
node generate-signed-url.js --post 3600
# Output: http://localhost:8000/images?expires=1234567890&signature=abc123...
```

2. Upload the image using curl:
```bash
curl -X POST "http://localhost:8000/images?expires=1234567890&signature=abc123..." \
  -F "file=@/path/to/image.jpg"
```

### Retrieve an Image

1. Generate a GET token:
```bash
node generate-signed-url.js --get myimage.jpg 3600
# Output: http://localhost:8000/images/myimage.jpg?expires=1234567890&signature=xyz789...
```

2. Access the image (browser, curl, etc.):
```bash
curl "http://localhost:8000/images/myimage.jpg?expires=1234567890&signature=xyz789..."
```

### Delete an Image

1. Generate a DELETE token:
```bash
node generate-signed-url.js --delete myimage.jpg 3600
```

2. Delete the image:
```bash
curl -X DELETE "http://localhost:8000/images/myimage.jpg?expires=1234567890&signature=def456..."
```

## Security Considerations

⚠️ **Production Deployment Checklist**:
- [ ] Change the hardcoded `secretKey` in `main.go` to a secure random string
- [ ] Use environment variables for sensitive configuration
- [ ] Implement additional authentication (API keys, JWT tokens, etc.)
- [ ] Add rate limiting to prevent abuse
- [ ] Validate file types and sizes
- [ ] Use HTTPS in production
- [ ] Implement proper logging and monitoring
- [ ] Set appropriate file permissions on the uploads directory
- [ ] Consider adding CORS headers if serving from a frontend

## License

[Add your license here]

## Author

[Add your name/info here]

