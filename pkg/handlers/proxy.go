package handlers

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"rbac/pkg/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func ProxyHandler(c echo.Context) error {
	cluster := c.Param("cluster")
	target := getClusterURL(cluster)
	if target == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cluster"})
	}

	url, err := url.Parse(target)
	if err != nil {
		utils.Logger.Error("Failed to parse target URL", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse target URL"})
	}

	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

func getClusterURL(cluster string) string {
	// Map cluster names to their respective URLs
	clusterURLs := map[string]string{
		"cluster1": "https://kubernetes.cluster1.svc",
		"cluster2": "https://kubernetes.cluster2.svc",
		// Add more clusters as needed
	}
	return clusterURLs[cluster]
}
