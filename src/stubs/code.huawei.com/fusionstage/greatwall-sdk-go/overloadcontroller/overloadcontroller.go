package overloadcontroller

import (
	"encoding/json"
	"strconv"
	"strings"
)

var currentRateLimit int = 10

func Init() error { return nil }

func InitWithPolicies(policies interface{}) error {
	str, ok := policies.(string)
	if !ok {
		return nil
	}
	currentRateLimit = parseRateLimit(str)
	return nil
}

func Process(dimNameValues map[string]string) (bool, error) {
	if currentRateLimit <= 0 {
		return false, nil
	}
	return true, nil
}

func parseRateLimit(policyStr string) int {
	var policy struct {
		APIService struct {
			Default []struct {
				RateLimit int `json:"rate_limit"`
			} `json:"default"`
		} `json:"APIService"`
	}
	if err := json.Unmarshal([]byte(policyStr), &policy); err == nil {
		if len(policy.APIService.Default) > 0 {
			return policy.APIService.Default[0].RateLimit
		}
	}
	if strings.Contains(policyStr, "rate_limit") {
		start := strings.Index(policyStr, "rate_limit")
		if start > 0 {
			sub := policyStr[start+10:]
			sub = strings.TrimSpace(sub)
			if strings.HasPrefix(sub, ":") || strings.HasPrefix(sub, ",") {
				sub = strings.TrimLeft(sub, ":, ")
				end := strings.IndexFunc(sub, func(r rune) bool {
					return r < '0' || r > '9'
				})
				if end > 0 {
					val, _ := strconv.Atoi(sub[:end])
					return val
				}
			}
		}
	}
	return 10
}