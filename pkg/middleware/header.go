package middleware

import (
	"github.com/gin-gonic/gin"
)

// Header names
const (
	ContentSecurityPolicy = "Content-Security-Policy"
	XContentTypeOptions   = "X-Content-Type-Options"
	XFrameOptions         = "X-Frame-Options"
	XXSSProtection        = "X-XSS-Protection"
)

// Header values
const (
	ContentSecurityPolicyValue = "default-src 'self'"
	XContentTypeOptionsValue   = "nosniff"
	XFrameOptionsValue         = "DENY"
	XXSSProtectionValue        = "1; mode=block"
)

// SecureHeaders adds security-related headers to the response to enhance security.
func SecureHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content-Security-Policy helps prevent XSS attacks by specifying which dynamic resources are allowed to load.
		c.Header(ContentSecurityPolicy, ContentSecurityPolicyValue)
		
		// X-Content-Type-Options prevents MIME type sniffing which can lead to XSS attacks.
		c.Header(XContentTypeOptions, XContentTypeOptionsValue)
		
		// X-Frame-Options prevents clickjacking by not allowing the page to be framed.
		c.Header(XFrameOptions, XFrameOptionsValue)
		
		// X-XSS-Protection enables the Cross-Site Scripting (XSS) filter built into most browsers.
		c.Header(XXSSProtection, XXSSProtectionValue)
		
		c.Next()
	}
}
