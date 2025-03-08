package utils

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindJSONAllowEmptyBody should be used when a request is fully optional, and the request body is allowed to be completely empty.
// The correct usage of this, and DisallowEmptyBody prevents certain cases where errors are ambiguous and nonsensical.
func BindJSONAllowEmptyBody(c *gin.Context, request any) error {
	return bindJSONtest(c, request, true)
}

// BindJSONDisallowEmptyBody should be a request is not fully optional, and the request body is never allowed to be completely empty.
// The correct usage of this, and AllowEmptyBody prevents certain cases where errors are ambiguous and nonsensical.
func BindJSONtestc(c *gin.Context, request any) error {
	return bindJSONtest(c, request, false)
}

func bindJSONtest(c *gin.Context, obj any, allowEmptyBody bool) error {
	// If empty bodies are allowed, ergo "", simply return nil.
	// Empty json object, "{}" will not be cought here, as they do not result in the "EOF" error.
	if allowEmptyBody {
		return nil
	}

	// Base case where bind is successful
	err := c.ShouldBindJSON(obj)
	if err == nil {
		return nil
	}

	// Empty request body case, where the ambiguous error "EOF" is
	// converted to something a potential user would understand
	if err.Error() == "EOF" {
		val := reflect.ValueOf(obj)

		// If obj is a pointer, get the element the pointer is referring to.
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		typ := val.Type()

		var requiredFields []string
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")
			// By convention, the fields in the struct with a "-" in the json tag, should be ignored when binding json
			if jsonTag == "-" {
				continue
			}

			jsonFieldName := strings.Split(jsonTag, ",")[0]
			validateTag := field.Tag.Get("binding")
			if !strings.Contains(validateTag, "required") {
				continue
			}
			requiredFields = append(requiredFields, jsonFieldName)
		}

		if len(requiredFields) > 0 {
			return errors.New(
				"request body cannot be empty, required fields: " + strings.Join(
					requiredFields,
					", ",
				),
			)
		}
		return errors.New(
			"request body cannot be completely empty, send empty json object instead.",
		)
	}

	if validatorErrors, ok := err.(validator.ValidationErrors); ok {
		jsonTagMap := make(map[string]string)
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

func BindJSON(c *gin.Context, obj any) error {
	// Attempt to bind the JSON to the provided struct
	err := c.ShouldBindJSON(obj)
	if err == nil {
		return nil
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	if err.Error() == "EOF" {
		// Reflect on the obj to find required fields
		var requiredFields []string
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" || len(jsonTag) <= 0 {
				continue
			}

			jsonFieldName := strings.Split(jsonTag, ",")[0]
			validateTag := field.Tag.Get("binding")
			if !strings.Contains(validateTag, "required") {
				continue
			}
			requiredFields = append(requiredFields, jsonFieldName)
		}

		if len(requiredFields) <= 0 {
			return errors.New("Request body cannot be empty")
		}

		return errors.New(
			"request body cannot be empty, required fields: " + strings.Join(
				requiredFields,
				", ",
			),
		)
	}

	validatorErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return errors.New("invalid request: " + err.Error())
	}
	// Map to hold the json field names
	jsonTagMap := make(map[string]string)
	// Reflect the obj to find the json tag names
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" || len(jsonTag) <= 0 {
			continue
		}
		jsonFieldName := strings.Split(jsonTag, ",")[0]
		jsonTagMap[field.Name] = jsonFieldName
	}

	out := make([]string, len(validatorErrors))
	for i, fieldError := range validatorErrors {
		fieldName := fieldError.Field()
		if jsonName, exists := jsonTagMap[fieldName]; exists {
			fieldName = jsonName
		}
		out[i] = fieldName + ": " + fieldError.Tag()
	}
	return errors.New("invalid request: " + strings.Join(out, ", "))
}

// BindJSON is a wrapper for gin.Context.ShouldBindJSON with enhanced error handling
// func BindJSON(c *gin.Context, obj any) error {
// 	// Attempt to bind the JSON to the provided struct
// 	if err := c.ShouldBindJSON(obj); err != nil {
// 		if err.Error() == "EOF" {
// 			// Reflect on the obj to find required fields
// 			val := reflect.ValueOf(obj)
// 			if val.Kind() == reflect.Ptr {
// 				val = val.Elem()
// 			}
// 			typ := val.Type()
//
// 			var requiredFields []string
// 			for i := 0; i < typ.NumField(); i++ {
// 				field := typ.Field(i)
// 				if jsonTag := field.Tag.Get("json"); jsonTag != "-" {
// 					jsonFieldName := strings.Split(jsonTag, ",")[0]
// 					if validateTag := field.Tag.Get("binding"); strings.Contains(
// 						validateTag,
// 						"required",
// 					) {
// 						requiredFields = append(requiredFields, jsonFieldName)
// 					}
// 				}
// 			}
// 			if len(requiredFields) > 0 {
// 				return errors.New(
// 					"request body cannot be empty, required fields: " + strings.Join(
// 						requiredFields,
// 						", ",
// 					),
// 				)
// 			}
// 			return errors.New("request body cannot be empty")
// 		}
//
// 		if validatorErrors, ok := err.(validator.ValidationErrors); ok {
// 			// Map to hold the json field names
// 			jsonTagMap := make(map[string]string)
// 			// Reflect the obj to find the json tag names
// 			val := reflect.ValueOf(obj)
// 			if val.Kind() == reflect.Ptr {
// 				val = val.Elem()
// 			}
// 			typ := val.Type()
// 			for i := 0; i < typ.NumField(); i++ {
// 				field := typ.Field(i)
// 				jsonTag := field.Tag.Get("json")
// 				if jsonTag == "-" {
// 					continue
// 				}
// 				jsonFieldName := strings.Split(jsonTag, ",")[0]
// 				jsonTagMap[field.Name] = jsonFieldName
// 			}
//
// 			out := make([]string, len(validatorErrors))
// 			for i, fe := range validatorErrors {
// 				fieldName := fe.Field()
// 				if jsonName, exists := jsonTagMap[fieldName]; exists {
// 					fieldName = jsonName
// 				}
// 				out[i] = fieldName + ": " + fe.Tag()
// 			}
// 			return errors.New("invalid request: " + strings.Join(out, ", "))
// 		}
//
// 		return errors.New("invalid request: " + err.Error())
// 	}
// 	return nil
// }
