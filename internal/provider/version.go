package provider

import (
	"strconv"
	"strings"
)

func parseVersion(version string) (major, minor, patch, beta int) {
	parts := strings.Split(version, "-")
	versionParts := strings.Split(parts[0], ".")

	major, _ = strconv.Atoi(versionParts[0])
	minor, _ = strconv.Atoi(versionParts[1])
	patch, _ = strconv.Atoi(versionParts[2])

	if len(parts) > 1 && strings.HasPrefix(parts[1], "beta.") {
		beta, _ = strconv.Atoi(strings.TrimPrefix(parts[1], "beta."))
	}

	return
}

func isVersionCompatible(currentVersion, minVersion string) bool {
	currentMajor, currentMinor, currentPatch, currentBeta := parseVersion(currentVersion)
	minMajor, minMinor, minPatch, minBeta := parseVersion(minVersion)

	if currentMajor > minMajor {
		return true
	}
	if currentMajor < minMajor {
		return false
	}

	if currentMinor > minMinor {
		return true
	}
	if currentMinor < minMinor {
		return false
	}

	if currentPatch > minPatch {
		return true
	}
	if currentPatch < minPatch {
		return false
	}

	return currentBeta >= minBeta
}
