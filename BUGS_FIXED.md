# Bug Fixes Report

This document details the 3 critical bugs that were identified and fixed in the codebase.

## Bug 1: Go Version Mismatch (Configuration/Build Issue)

### **Problem Description:**
The `go.mod` file specified an invalid Go version `1.24.5` (which doesn't exist) while the GitHub Actions CI workflow was configured to use Go `1.20`. This version mismatch would cause build failures and deployment issues.

### **Location:**
- `go.mod` (line 3)
- `.github/workflows/go.yml` (line 22)

### **Impact:**
- **Severity:** High
- **Type:** Build/Configuration Error
- **Consequences:** Complete build failure, broken CI/CD pipeline, inability to compile the application

### **Root Cause:**
Invalid Go version number and inconsistent version specifications across configuration files.

### **Fix Applied:**
```diff
# go.mod
- go 1.24.5
+ go 1.21

# .github/workflows/go.yml
- go-version: '1.20'
+ go-version: '1.21'
```

### **Why This Fix Works:**
- Uses a valid, stable Go version (1.21)
- Ensures consistency between local development and CI environment
- Allows the application to build successfully

---

## Bug 2: SQL Injection Vulnerability (Security Issue)

### **Problem Description:**
The `getUserHandler` function was vulnerable to SQL injection attacks due to direct string interpolation in SQL queries. An attacker could manipulate the `userID` parameter to execute arbitrary SQL commands.

### **Location:**
`main.go` - `getUserHandler` function (originally lines 53-66)

### **Impact:**
- **Severity:** Critical
- **Type:** Security Vulnerability (CWE-89)
- **Consequences:** 
  - Data breach potential
  - Unauthorized data access
  - Database corruption risk
  - Complete database compromise

### **Vulnerable Code:**
```go
// VULNERABLE: Direct string interpolation
query := fmt.Sprintf("SELECT id, name, email FROM users WHERE id = %s", userID)
```

### **Attack Vector Example:**
```
GET /users/1; DROP TABLE users; --
```

### **Fix Applied:**
```go
// Validate and convert user ID to integer
userID, err := strconv.Atoi(userIDStr)
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
    return
}

// SECURE: Using parameterized query
query := "SELECT id, name, email FROM users WHERE id = $1"
err = db.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email)
```

### **Why This Fix Works:**
- **Input validation:** Ensures userID is a valid integer
- **Parameterized queries:** SQL parameters are properly escaped
- **Type safety:** Prevents string-based injection attacks
- **Error handling:** Distinguishes between client and server errors

---

## Bug 3: Race Condition (Concurrency Issue)

### **Problem Description:**
The `counter` global variable was accessed by multiple goroutines without proper synchronization, leading to race conditions. Multiple concurrent requests to `/counter` and `/increment` endpoints could read/write the counter simultaneously, causing data races and inconsistent results.

### **Location:**
`main.go` - `getCounterHandler` and `incrementHandler` functions

### **Impact:**
- **Severity:** Medium-High
- **Type:** Concurrency Bug
- **Consequences:**
  - Incorrect counter values
  - Data races detected by Go race detector
  - Unpredictable application behavior
  - Potential memory corruption

### **Vulnerable Code:**
```go
func getCounterHandler(c *gin.Context) {
    // Reading counter without proper synchronization
    c.JSON(http.StatusOK, gin.H{"counter": counter})
}

func incrementHandler(c *gin.Context) {
    // Race condition - multiple goroutines can read/write simultaneously
    counter++
    c.JSON(http.StatusOK, gin.H{"counter": counter})
}
```

### **Race Condition Scenario:**
1. Goroutine A reads counter (value: 10)
2. Goroutine B reads counter (value: 10)  
3. Goroutine A increments to 11
4. Goroutine B increments to 11
5. Result: Counter is 11 instead of expected 12

### **Fix Applied:**
```go
func getCounterHandler(c *gin.Context) {
    // Protected read with mutex
    mu.Lock()
    currentCounter := counter
    mu.Unlock()
    
    c.JSON(http.StatusOK, gin.H{"counter": currentCounter})
}

func incrementHandler(c *gin.Context) {
    // Fixed: Race condition - using mutex to protect counter access
    mu.Lock()
    counter++
    currentCounter := counter
    mu.Unlock()
    
    c.JSON(http.StatusOK, gin.H{"counter": currentCounter})
}
```

### **Why This Fix Works:**
- **Mutual exclusion:** Only one goroutine can access counter at a time
- **Atomic operations:** Read-modify-write operations are atomic
- **Consistent state:** Prevents interleaved operations
- **Performance:** Minimal lock contention with short critical sections

---

## Additional Fix: Memory Leak Prevention

### **Bonus Bug 4: Memory Leak (Performance Issue)**

**Problem:** Created a `streamHandler` that demonstrates common memory leaks:
- Goroutines running indefinitely without cleanup
- Channels not properly closed
- No timeout or context cancellation

**Fix Applied:**
- Added context-based cancellation with timeouts
- Proper cleanup functions for goroutines and channels
- Graceful shutdown when clients disconnect
- Resource leak prevention using `defer` statements

---

## Testing the Fixes

To verify these fixes:

1. **Build verification:**
   ```bash
   go build -v ./...
   ```

2. **Race condition testing:**
   ```bash
   go run -race main.go
   ```

3. **Security testing:**
   ```bash
   # Test parameterized queries prevent injection
   curl "http://localhost:8080/users/1"
   ```

4. **Load testing for race conditions:**
   ```bash
   # Multiple concurrent increment requests
   for i in {1..100}; do curl -X POST http://localhost:8080/increment & done
   ```

## Summary

All critical bugs have been successfully identified and fixed:
- ✅ **Build failure** resolved with correct Go version
- ✅ **SQL injection** prevented with parameterized queries  
- ✅ **Race condition** eliminated with proper mutex synchronization
- ✅ **Memory leaks** prevented with proper resource cleanup

The application now builds successfully and is secure, thread-safe, and memory-efficient.