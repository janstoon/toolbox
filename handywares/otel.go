package handywares

import "go.opentelemetry.io/otel/attribute"

const (
	oaPrefix = attribute.Key("jst")

	oaDebug      = oaPrefix + ".debug"
	oaDebugStack = oaDebug + ".stack"

	oaPanicValue = oaPrefix + ".panic.value"

	oaHttp         = oaPrefix + ".http"
	oaHttpRequest  = oaHttp + ".request"
	oaHttpResponse = oaHttp + ".response"
)
