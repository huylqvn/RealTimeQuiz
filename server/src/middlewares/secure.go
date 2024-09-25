package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

type Config struct {
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter func(*chi.Context) bool
	// XSSProtection
	// Optional. Default value "0".
	XSSProtection string
	// ContentTypeNosniff
	// Optional. Default value "nosniff".
	ContentTypeNosniff string
	// XFrameOptions
	// Optional. Default value "SAMEORIGIN".
	// Possible values: "SAMEORIGIN", "DENY", "ALLOW-FROM uri"
	XFrameOptions string
	// HSTSMaxAge
	// Optional. Default value 0.
	HSTSMaxAge int
	// HSTSExcludeSubdomains
	// Optional. Default value false.
	HSTSExcludeSubdomains bool
	// ContentSecurityPolicy
	// Optional. Default value "".
	ContentSecurityPolicy string
	// CSPReportOnly
	// Optional. Default value false.
	CSPReportOnly bool
	// HSTSPreloadEnabled
	// Optional. Default value false.
	HSTSPreloadEnabled bool
	// ReferrerPolicy
	// Optional. Default value "no-referrer".
	ReferrerPolicy string
	// Permissions-Policy
	// Optional. Default value "".
	PermissionPolicy string
	// Cross-Origin-Embedder-Policy
	// Optional. Default value "require-corp".
	CrossOriginEmbedderPolicy string
	// Cross-Origin-Opener-Policy
	// Optional. Default value "same-origin".
	CrossOriginOpenerPolicy string
	// Cross-Origin-Resource-Policy
	// Optional. Default value "same-origin".
	CrossOriginResourcePolicy string
	// Origin-Agent-Cluster
	// Optional. Default value "?1".
	OriginAgentCluster string
	// X-DNS-Prefetch-Control
	// Optional. Default value "off".
	XDNSPrefetchControl string
	// X-Download-Options
	// Optional. Default value "noopen".
	XDownloadOptions string
	// X-Permitted-Cross-Domain-Policies
	// Optional. Default value "none".
	XPermittedCrossDomain string
}

func WrapEncodeWithSercureHeader(f func(ctx context.Context, w http.ResponseWriter, response interface{}) error, config ...Config) func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	// Set config default values
	if cfg.XSSProtection == "" {
		cfg.XSSProtection = "0"
	}
	if cfg.ContentTypeNosniff == "" {
		cfg.ContentTypeNosniff = "nosniff"
	}
	if cfg.XFrameOptions == "" {
		cfg.XFrameOptions = "SAMEORIGIN"
	}
	if cfg.ReferrerPolicy == "" {
		cfg.ReferrerPolicy = "no-referrer"
	}
	if cfg.CrossOriginEmbedderPolicy == "" {
		cfg.CrossOriginEmbedderPolicy = "require-corp"
	}
	if cfg.CrossOriginOpenerPolicy == "" {
		cfg.CrossOriginOpenerPolicy = "same-origin"
	}
	if cfg.CrossOriginResourcePolicy == "" {
		cfg.CrossOriginResourcePolicy = "same-origin"
	}
	if cfg.OriginAgentCluster == "" {
		cfg.OriginAgentCluster = "?1"
	}
	if cfg.XDNSPrefetchControl == "" {
		cfg.XDNSPrefetchControl = "off"
	}
	if cfg.XDownloadOptions == "" {
		cfg.XDownloadOptions = "noopen"
	}
	if cfg.XPermittedCrossDomain == "" {
		cfg.XPermittedCrossDomain = "none"
	}
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if cfg.XSSProtection != "" {
			w.Header().Set(HeaderXXSSProtection, cfg.XSSProtection)
		}
		if cfg.ContentTypeNosniff != "" {
			w.Header().Set(HeaderXContentTypeOptions, cfg.ContentTypeNosniff)
		}
		if cfg.XFrameOptions != "" {
			w.Header().Set(HeaderXFrameOptions, cfg.XFrameOptions)
		}
		if cfg.CrossOriginEmbedderPolicy != "" {
			w.Header().Set("Cross-Origin-Embedder-Policy", cfg.CrossOriginEmbedderPolicy)
		}
		if cfg.CrossOriginOpenerPolicy != "" {
			w.Header().Set("Cross-Origin-Opener-Policy", cfg.CrossOriginOpenerPolicy)
		}
		if cfg.CrossOriginResourcePolicy != "" {
			w.Header().Set("Cross-Origin-Resource-Policy", cfg.CrossOriginResourcePolicy)
		}
		if cfg.OriginAgentCluster != "" {
			w.Header().Set("Origin-Agent-Cluster", cfg.OriginAgentCluster)
		}
		if cfg.ReferrerPolicy != "" {
			w.Header().Set(HeaderReferrerPolicy, cfg.ReferrerPolicy)
		}
		if cfg.XDNSPrefetchControl != "" {
			w.Header().Set("X-DNS-Prefetch-Control", cfg.XDNSPrefetchControl)
		}
		if cfg.XDownloadOptions != "" {
			w.Header().Set("X-Download-Options", cfg.XDownloadOptions)
		}
		if cfg.XPermittedCrossDomain != "" {
			w.Header().Set("X-Permitted-Cross-Domain-Policies", cfg.XPermittedCrossDomain)
		}

		if cfg.ContentSecurityPolicy != "" {
			if cfg.CSPReportOnly {
				w.Header().Set(HeaderContentSecurityPolicyReportOnly, cfg.ContentSecurityPolicy)
			} else {
				w.Header().Set(HeaderContentSecurityPolicy, cfg.ContentSecurityPolicy)
			}
		}
		return f(ctx, w, response)
	}
}
