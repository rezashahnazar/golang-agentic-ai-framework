package provider

import "fmt"

func ValidateModel(availableModels []Model, modelName string, providerName string) error {
	for _, m := range availableModels {
		if m.Name() == modelName {
			return nil
		}
	}
	return fmt.Errorf("model %s is not available in provider %s", modelName, providerName)
}

func ValidateRequestParameters(availableParams []string, requestParameters map[string]any, modelName string) error {
	availableParamsMap := make(map[string]bool)
	for _, param := range availableParams {
		availableParamsMap[param] = true
	}

	for key := range requestParameters {
		if !availableParamsMap[key] {
			return fmt.Errorf("request parameter '%s' is not available for model %s. Available parameters: %v", key, modelName, availableParams)
		}
	}
	return nil
}

