#!/usr/bin/env node

const crypto = require('crypto');

// Load configuration from environment variables with fallback defaults
const secretKey = process.env.SECRET_KEY || 'secret-key';
const baseUrl = process.env.BASE_URL || 'http://localhost:8000';

function generateSignedUrl(method, filename, validForSeconds) {
    // Calculate expiration timestamp (current time + validForSeconds)
    const expires = Math.floor(Date.now() / 1000) + parseInt(validForSeconds);

    // Create the data string to sign: "METHOD:filename:expires" (empty filename for POST)
    // Including method prevents token reuse across different HTTP methods
    const data = `${method}:${filename || ''}:${expires}`;

    // Create HMAC-SHA256 signature
    const hmac = crypto.createHmac('sha256', secretKey);
    hmac.update(data);
    const signature = hmac.digest('hex');

    // Construct the signed URL
    if (filename) {
        // GET/PUT/DELETE requests with filename
        const signedUrl = `${baseUrl}/images/${filename}?expires=${expires}&signature=${signature}`;
        return signedUrl;
    } else {
        // POST request without filename
        const signedUrl = `${baseUrl}/images?expires=${expires}&signature=${signature}`;
        return signedUrl;
    }
}

// Get command line arguments
const args = process.argv.slice(2);

if (args.length < 1) {
    console.error('Usage:');
    console.error('  For GET:  node generate-signed-url.js --get <image-name> <time-in-seconds>');
    console.error('  For PUT:  node generate-signed-url.js --put <image-name> <time-in-seconds>');
    console.error('  For DELETE: node generate-signed-url.js --delete <image-name> <time-in-seconds>');
    console.error('  For POST: node generate-signed-url.js --post <time-in-seconds>');
    console.error('  Short forms: -g (GET), -u (PUT), -d (DELETE), -p (POST)');
    process.exit(1);
}

let method, imageName, timeInSeconds;

// Parse method flag
const methodFlag = args[0].toLowerCase();
switch (methodFlag) {
    case '--post':
    case '-p':
        method = 'POST';
        if (args.length < 2) {
            console.error('Usage: node generate-signed-url.js --post <time-in-seconds>');
            process.exit(1);
        }
        imageName = null;
        timeInSeconds = args[1];
        break;
    case '--get':
    case '-g':
        method = 'GET';
        if (args.length < 3) {
            console.error('Usage: node generate-signed-url.js --get <image-name> <time-in-seconds>');
            process.exit(1);
        }
        imageName = args[1];
        timeInSeconds = args[2];
        break;
    case '--put':
    case '-u':
        method = 'PUT';
        if (args.length < 3) {
            console.error('Usage: node generate-signed-url.js --put <image-name> <time-in-seconds>');
            process.exit(1);
        }
        imageName = args[1];
        timeInSeconds = args[2];
        break;
    case '--delete':
    case '-d':
        method = 'DELETE';
        if (args.length < 3) {
            console.error('Usage: node generate-signed-url.js --delete <image-name> <time-in-seconds>');
            process.exit(1);
        }
        imageName = args[1];
        timeInSeconds = args[2];
        break;
    default:
        console.error('Error: Invalid method flag. Use --get, --put, --delete, or --post');
        console.error('  Short forms: -g (GET), -u (PUT), -d (DELETE), -p (POST)');
        process.exit(1);
}

// Validate time is a number
if (isNaN(timeInSeconds) || parseInt(timeInSeconds) <= 0) {
    console.error('Error: time-in-seconds must be a positive number');
    process.exit(1);
}

// Generate and output the signed URL
const signedUrl = generateSignedUrl(method, imageName, timeInSeconds);
console.log(signedUrl);

