package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindJSON is a wrapper for gin.Context.ShouldBindJSON with enhanced error handling
func BindJSON(c *gin.Context, obj any) error {
	// Attempt to bind the JSON to the provided struct
	if err := c.ShouldBindJSON(obj); err != nil {
		if err.Error() == "EOF" {
			// Reflect on the obj to find required fields
			val := reflect.ValueOf(obj)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			typ := val.Type()

			var requiredFields []string
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				if jsonTag := field.Tag.Get("json"); jsonTag != "-" {
					jsonFieldName := strings.Split(jsonTag, ",")[0]
					if validateTag := field.Tag.Get("binding"); strings.Contains(
						validateTag,
						"required",
					) {
						requiredFields = append(requiredFields, jsonFieldName)
					}
				}
			}
			if len(requiredFields) > 0 {
				return errors.New(
					"request body cannot be empty, required fields: " + strings.Join(
						requiredFields,
						", ",
					),
				)
			}
			return errors.New("request body cannot be empty")
		}

		if validatorErrors, ok := err.(validator.ValidationErrors); ok {
			// Map to hold the json field names
			jsonTagMap := make(map[string]string)
			// Reflect the obj to find the json tag names
			val := reflect.ValueOf(obj)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}
			typ := val.Type()
			for i := 0; i < typ.NumField(); i++ {
				field := typ.Field(i)
				jsonTag := field.Tag.Get("json")
				if jsonTag == "-" {
					continue
				}
				jsonFieldName := strings.Split(jsonTag, ",")[0]
				jsonTagMap[field.Name] = jsonFieldName
			}

			out := make([]string, len(validatorErrors))
			for i, fe := range validatorErrors {
				fieldName := fe.Field()
				if jsonName, exists := jsonTagMap[fieldName]; exists {
					fieldName = jsonName
				}
				out[i] = fieldName + ": " + fe.Tag()
			}
			return errors.New("invalid request: " + strings.Join(out, ", "))
		}

		return errors.New("invalid request: " + err.Error())
	}
	return nil
}
