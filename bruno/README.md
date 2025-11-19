# Fitness Coach API - Bruno Collection

This Bruno collection provides a complete testing suite for the Fitness Coach API.

## 🚀 Quick Start

### 1. Install Bruno
Download from: https://www.usebruno.com/downloads

### 2. Open Collection
1. Open Bruno
2. Click "Open Collection"
3. Navigate to this `bruno` folder
4. Select the folder

### 3. Select Environment
- **Local**: For testing against `localhost:8080`
- **Railway**: For testing against Railway deployment

### 4. Update Railway URL (if using Railway)
Edit `environments/Railway.bru` and update:
```
baseUrl: https://your-actual-app.railway.app
```

## 📋 Testing Sequence

### **Follow this order for first-time setup:**

1. **Health Check** (`Health/Health Check.bru`)
   - ✅ Verify API is running
   - No authentication needed
   - Should return `200 OK`

2. **Register User** (`Auth/1. Register User.bru`)
   - ✅ Create a new user account
   - Automatically saves `authToken` and `userId`
   - Should return `201 Created`
   - **Note**: Will fail if email already exists

3. **Login User** (`Auth/2. Login User.bru`)
   - ✅ Authenticate with existing credentials
   - Automatically saves new `authToken` and `userId`
   - Should return `200 OK`
   - Use this if registration fails (user already exists)

4. **Refresh Token** (`Auth/3. Refresh Token.bru`)
   - ℹ️ Currently not implemented (returns `501`)
   - Will be available in future updates

## 🔑 Authentication Flow

### Automatic Token Management

The collection automatically manages authentication:

1. **Register** or **Login** → Token saved to `{{authToken}}`
2. All subsequent requests use `{{authToken}}` automatically
3. Token is valid for **24 hours**

### Manual Token Override

If you need to use a different token:
1. Go to Environments
2. Update `authToken` variable
3. All authenticated requests will use the new token

## 📊 Environment Variables

### Available Variables

| Variable | Description | Auto-Set |
|----------|-------------|----------|
| `baseUrl` | API base URL | No |
| `apiVersion` | API version (v1) | No |
| `authToken` | JWT authentication token | Yes ✅ |
| `userId` | Current user ID | Yes ✅ |
| `testEmail` | Email for testing | No |
| `testPassword` | Password for testing | No |
| `testName` | Name for testing | No |

### Customizing Test Data

Edit environment files to change test credentials:

```
testEmail: your-email@example.com
testPassword: your-password
testName: Your Name
```

## 🧪 Running Tests

### Individual Request
1. Click on any request
2. Click "Send" button
3. View response and test results

### Run All Requests (Sequential)
1. Right-click on the collection root
2. Select "Run Collection"
3. Tests run in numbered order

### Expected Results

✅ **All Green** = Everything working!
- Health Check: 200 OK
- Register: 201 Created (or 409 if user exists)
- Login: 200 OK
- Refresh: 501 Not Implemented (expected)

## 📖 Request Documentation

Each request includes:
- **Docs tab**: Detailed explanation
- **Tests tab**: Automated assertions
- **Scripts tab**: Post-request automation

## 🔍 Troubleshooting

### Registration Fails with 409
- **Cause**: Email already registered
- **Solution**: Use "Login User" instead, or change `testEmail` in environment

### All Requests Fail
- **Check**: Is the server running?
- **Local**: `go run cmd/api/main.go` or check Docker
- **Railway**: Check deployment logs

### Authentication Failures (401)
- **Cause**: Token expired or invalid
- **Solution**: Run "Login User" again to get a new token

### Connection Refused
- **Check**: `baseUrl` in environment
- **Local**: Should be `http://localhost:8080`
- **Railway**: Should be your Railway deployment URL

## 🎯 What's Working Now

✅ **Authentication System**
- User Registration
- User Login
- JWT Token Generation
- Password Hashing

✅ **Infrastructure**
- Health Check
- Database Migrations
- CORS Configuration

⏳ **Coming Soon**
- Meals tracking
- Food database
- Activities logging
- Workouts tracking
- Metrics & goals
- Chat/AI features

## 📝 Adding New Requests

When new endpoints are enabled:

1. Create new folder (e.g., `Meals`)
2. Create `.bru` file with this structure:
```
meta {
  name: Request Name
  type: http
  seq: 5
}

get {
  url: {{baseUrl}}/api/{{apiVersion}}/endpoint
  auth: bearer
}

auth:bearer {
  token: {{authToken}}
}
```

## 🛠️ Development Tips

### Testing Against Local Development
```bash
# Start the server
cd backend
go run cmd/api/main.go

# Select "Local" environment in Bruno
# Run requests
```

### Testing Against Railway
```bash
# Deploy to Railway (automatic on git push)
git push origin main

# Update Railway.bru with your URL
# Select "Railway" environment in Bruno
# Run requests
```

## 📞 Support

If you encounter issues:
1. Check server logs
2. Verify environment variables
3. Review request documentation
4. Check response body for error details

Happy Testing! 🚀
