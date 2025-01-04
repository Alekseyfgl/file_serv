package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

// GetClientIP extracts the client's IP address from the request context.
// It checks the following headers in order of priority:
// 1. X-Real-IP: This header is often set by proxies and load balancers to indicate the real IP of the client.
// 2. X-Forwarded-For: This header contains a list of IPs, where the first IP represents the original client IP.
// 3. c.IP(): Falls back to the IP from the connection if the above headers are not set.
//
// Parameters:
// - c: The *fiber.Ctx object representing the current request context.
//
// Returns:
// - string: The determined client IP address.
func GetClientIP(c *fiber.Ctx) string {
	// Priority: X-Real-IP -> X-Forwarded-For -> c.IP()
	if realIP := c.Get("X-Real-IP"); realIP != "" {
		fmt.Println("X-Real-IP:", realIP)
		return realIP
	} else if forwardedFor := c.Get("X-Forwarded-For"); forwardedFor != "" {

		ip := strings.Split(forwardedFor, ",")[0]
		fmt.Println("X-Forwarded-For:", ip)
		return ip
	}
	return c.IP()
}
