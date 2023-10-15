package configs

import "os"

func EnvironmentVariables(Environment string) {

	if Environment == "Development" {
		os.Setenv("Github_Key", "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")

	} else if Environment == "Staging" {
		os.Setenv("Github_Key", "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")

	} else if Environment == "Production" {
		os.Setenv("Github_Key", "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")

	}
}
