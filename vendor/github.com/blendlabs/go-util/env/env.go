package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	exception "github.com/blendlabs/go-exception"
	util "github.com/blendlabs/go-util"
)

var (
	_env     Vars
	_envLock = sync.Mutex{}
)

const (
	// VarServiceEnv is a common env var name.
	VarServiceEnv = "SERVICE_ENV"
	// VarServiceName is a common env var name.
	VarServiceName = "SERVICE_NAME"
	// VarServiceSecret is a common env var name.
	VarServiceSecret = "SERVICE_SECRET"
	// VarPort is a common env var name.
	VarPort = "PORT"
	// VarSecurePort is a common env var name.
	VarSecurePort = "SECURE_PORT"
	// VarTLSCertPath is a common env var name.
	VarTLSCertPath = "TLS_CERT_PATH"
	// VarTLSKeyPath is a common env var name.
	VarTLSKeyPath = "TLS_KEY_PATH"
	// VarTLSCert is a common env var name.
	VarTLSCert = "TLS_CERT"
	// VarTLSKey is a common env var name.
	VarTLSKey = "TLS_KEY"

	// VarPGIdleConns is a common env var name.
	VarPGIdleConns = "PG_IDLE_CONNS"
	// VarPGMaxConns is a common env var name.
	VarPGMaxConns = "PG_MAX_CONNS"

	// ServiceEnvDev is a service environment.
	ServiceEnvDev = "dev"
	// ServiceEnvCI is a service environment.
	ServiceEnvCI = "ci"
	// ServiceEnvPreprod is a service environment.
	ServiceEnvPreprod = "preprod"
	// ServiceEnvBeta is a service environment.
	ServiceEnvBeta = "beta"
	// ServiceEnvProd is a service environment.
	ServiceEnvProd = "prod"
)

// Env returns the current env var set.
func Env() Vars {
	if _env == nil {
		_envLock.Lock()
		defer _envLock.Unlock()
		if _env == nil {
			_env = NewVarsFromEnvironment()
		}
	}
	return _env
}

// SetEnv sets the env vars.
func SetEnv(vars Vars) {
	_envLock.Lock()
	_env = vars
	_envLock.Unlock()
}

// NewVars returns a new env var set.
func NewVars() Vars {
	return Vars{}
}

// NewVarsFromEnvironment reads an EnvVar set from the environment.
func NewVarsFromEnvironment() Vars {
	vars := Vars{}
	envVars := os.Environ()
	for _, ev := range envVars {
		parts := strings.SplitN(ev, "=", 2)
		if len(parts) > 1 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars
}

// Vars is a set of environment variables.
type Vars map[string]string

// Set sets a value for a key.
func (ev Vars) Set(envVar, value string) {
	ev[envVar] = value
}

// Restore resets an environment variable to it's environment value.
func (ev Vars) Restore(key string) {
	ev[key] = os.Getenv(key)
}

// String returns a string value for a given key.
func (ev Vars) String(envVar string, defaults ...string) string {
	if value, hasValue := ev[envVar]; hasValue {
		return value
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// Bool returns a boolean value for a key, defaulting to false.
// Valid "truthy" values are `true`, `yes`, and `1`.
// Everything else is false, including `REEEEEEEEEEEEEEE`.
func (ev Vars) Bool(envVar string, defaults ...bool) bool {
	if value, hasValue := ev[envVar]; hasValue {
		if len(value) > 0 {
			return util.String.CaseInsensitiveEquals(value, "true") || util.String.CaseInsensitiveEquals(value, "yes") || value == "1"
		}
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return false
}

// Int returns an integer value for a given key.
func (ev Vars) Int(envVar string, defaults ...int) int {
	if value, hasValue := ev[envVar]; hasValue {
		return util.String.ParseInt(value)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// Int64 returns an int64 value for a given key.
func (ev Vars) Int64(envVar string, defaults ...int64) int64 {
	if value, hasValue := ev[envVar]; hasValue {
		return util.String.ParseInt64(value)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

// Bytes returns a []byte value for a given key.
func (ev Vars) Bytes(key string, defaults ...[]byte) []byte {
	if value, hasValue := ev[key]; hasValue && len(value) > 0 {
		return []byte(value)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}

// Base64 returns a []byte value for a given key whose value is encoded in base64.
func (ev Vars) Base64(key string, defaults ...[]byte) []byte {
	if value, hasValue := ev[key]; hasValue && len(value) > 0 {
		result, _ := util.Base64.Decode(value)
		return result
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}

// HasVar returns if a key is present in the set.
func (ev Vars) HasVar(envVar string) bool {
	_, hasKey := ev[envVar]
	return hasKey
}

// HasAllVars returns if all of the given vars are present in the set.
func (ev Vars) HasAllVars(envVars ...string) bool {
	if len(envVars) == 0 {
		return false
	}
	for _, envVar := range envVars {
		if !ev.HasVar(envVar) {
			return false
		}
	}
	return true
}

// HasAnyVars returns if any of the given vars are present in the set.
func (ev Vars) HasAnyVars(envVars ...string) bool {
	for _, envVar := range envVars {
		if ev.HasVar(envVar) {
			return true
		}
	}
	return false
}

// Union returns the union of the two sets, other replacing conflicts.
func (ev Vars) Union(other Vars) Vars {
	newSet := NewVars()
	for key, value := range ev {
		newSet[key] = value
	}
	for key, value := range other {
		newSet[key] = value
	}
	return newSet
}

// Vars returns all the vars stored in the env var set.
func (ev Vars) Vars() []string {
	var envVars = make([]string, len(ev))
	var index int
	for envVar := range ev {
		envVars[index] = envVar
		index++
	}
	return envVars
}

// ServiceEnv is a common environment variable for the services environment.
// Common values include "dev", "ci", "sandbox", "preprod", "beta", and "prod".
func (ev Vars) ServiceEnv(defaults ...string) string {
	return ev.String(VarServiceEnv, defaults...)
}

// IsProduction returns if the ServiceEnv is a production environment.
func (ev Vars) IsProduction() bool {
	return ev.ServiceEnv() == ServiceEnvPreprod ||
		ev.ServiceEnv() == ServiceEnvProd
}

// IsProdLike returns if the ServiceEnv is "prodlike".
func (ev Vars) IsProdLike() bool {
	return ev.ServiceEnv() == ServiceEnvPreprod ||
		ev.ServiceEnv() == ServiceEnvBeta ||
		ev.ServiceEnv() == ServiceEnvProd
}

// ServiceName is a common environment variable for the service's name.
func (ev Vars) ServiceName(defaults ...string) string {
	return ev.String(VarServiceName, defaults...)
}

// ServiceSecret is the main secret for the app.
// It is typically a 32 byte / 256 bit key.
func (ev Vars) ServiceSecret(defaults ...[]byte) []byte {
	return ev.Base64(VarServiceSecret, defaults...)
}

// Port is a common environment variable.
// It is what TCP port to bind to for the HTTP server.
func (ev Vars) Port(defaults ...string) string {
	return ev.String(VarPort, defaults...)
}

// SecurePort is a common environment variable.
// It is what TCP port to bind to for the HTTPS server.
func (ev Vars) SecurePort(defaults ...string) string {
	return ev.String(VarSecurePort, defaults...)
}

// TLSCertFilepath is a common environment variable for the (whole) TLS cert to use with https.
func (ev Vars) TLSCertFilepath(defaults ...string) string {
	return ev.String(VarTLSCertPath, defaults...)
}

// TLSKeyFilepath is a common environment variable for the (whole) TLS key to use with https.
func (ev Vars) TLSKeyFilepath(defaults ...string) string {
	return ev.String(VarTLSKeyPath, defaults...)
}

// TLSCert is a common environment variable for the (whole) TLS cert to use with https.
func (ev Vars) TLSCert(defaults ...[]byte) []byte {
	return ev.Bytes(VarTLSCert, defaults...)
}

// TLSKey is a common environment variable for the (whole) TLS key to use with https.
func (ev Vars) TLSKey(defaults ...[]byte) []byte {
	return ev.Bytes(VarTLSKey, defaults...)
}

// RequireVars enforces that a given set of environment variables are present.
func (ev Vars) RequireVars(keys ...string) error {
	for _, key := range keys {
		if !ev.HasVar(key) {
			return fmt.Errorf("the following environment variables are required: `%s`", strings.Join(keys, ","))
		}
	}
	return nil
}

// MustVars enforces that a given set of environment variables are present and panics
// if they're not present.
func (ev Vars) MustVars(keys ...string) {
	for _, key := range keys {
		if !ev.HasVar(key) {
			panic(fmt.Sprintf("the following environment variables are required: `%s`", strings.Join(keys, ",")))
		}
	}
}

const (
	// TagNameEnvironmentVariableName is the struct tag for what environment variable to use to populate a field.
	TagNameEnvironmentVariableName = "env"
	// TagNameEnvironmentVariableDefault is the struct tag for what to use if the environment variable is empty.
	TagNameEnvironmentVariableDefault = "env_default"
)

// ReadInto reads the environment into tagged fields on the `obj`.
func (ev Vars) ReadInto(obj interface{}) error {
	objMeta := reflectType(obj)

	var field reflect.StructField
	var tag string
	var envValue string
	var defaultValue string
	var err error

	for x := 0; x < objMeta.NumField(); x++ {
		field = objMeta.Field(x)
		tag = field.Tag.Get(TagNameEnvironmentVariableName)

		if len(tag) > 0 {
			envValue = ev.String(tag)
			if len(envValue) > 0 {
				err = setValueByName(obj, field.Name, envValue)
				if err != nil {
					return err
				}
			} else {
				defaultValue = field.Tag.Get(TagNameEnvironmentVariableDefault)
				if len(defaultValue) > 0 {
					err = setValueByName(obj, field.Name, defaultValue)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

// ReflectType returns the integral type for an object.
func reflectType(obj interface{}) reflect.Type {
	t := reflect.TypeOf(obj)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}

	return t
}

// ReflectValue returns the integral reflect.Value for an object.
func reflectValue(obj interface{}) reflect.Value {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	return v
}

// SetValueByName sets a value on an object by its field name.
func setValueByName(target interface{}, fieldName string, fieldValue interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = exception.Newf("Error setting field: %v", r)
		}
	}()
	typeCheck := reflect.TypeOf(target)
	if typeCheck.Kind() != reflect.Ptr {
		return exception.New("Cannot modify non-pointer target")
	}

	targetValue := reflectValue(target)
	targetType := reflectType(target)
	relevantField, hasField := targetType.FieldByName(fieldName)

	if !hasField {
		return exception.Newf("Field not found  %s.%s", targetType.Name(), fieldName)
	}

	field := targetValue.FieldByName(relevantField.Name)
	fieldType := field.Type()
	if !field.CanSet() {
		return exception.Newf("Cannot set field %s", fieldName)
	}

	valueReflected := reflectValue(fieldValue)
	if !valueReflected.IsValid() {
		return exception.New("Reflected value is invalid, cannot continue.")
	}

	if valueReflected.Type().AssignableTo(fieldType) {
		field.Set(valueReflected)
		return nil
	}

	if field.Kind() == reflect.Ptr {
		if valueReflected.CanAddr() {
			convertedValue := valueReflected.Convert(fieldType.Elem())
			if convertedValue.CanAddr() {
				field.Set(convertedValue.Addr())
				return nil
			}
		}
		return exception.New("Cannot take address of value for assignment to field pointer")
	}

	if fieldAsString, isString := valueReflected.Interface().(string); isString {
		var parsedValue reflect.Value
		handledType := true
		switch fieldType.Kind() {
		case reflect.Int:
			intValue, err := strconv.Atoi(fieldAsString)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(intValue)
		case reflect.Int64:
			int64Value, err := strconv.ParseInt(fieldAsString, 10, 64)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(int64Value)
		case reflect.Uint16:
			intValue, err := strconv.Atoi(fieldAsString)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(uint16(intValue))
		case reflect.Uint: //a.k.a. uint32
			intValue, err := strconv.Atoi(fieldAsString)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(uint(intValue))
		case reflect.Uint32:
			intValue, err := strconv.Atoi(fieldAsString)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(uint32(intValue))
		case reflect.Uint64:
			intValue, err := strconv.Atoi(fieldAsString)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(uint64(intValue))
		case reflect.Float32:
			floatValue, err := strconv.ParseFloat(fieldAsString, 32)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(floatValue)
		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(fieldAsString, 64)
			if err != nil {
				return exception.Wrap(err)
			}
			parsedValue = reflect.ValueOf(floatValue)
		default:
			handledType = false
		}
		if handledType {
			field.Set(parsedValue)
			return nil
		}
	}

	convertedValue := valueReflected.Convert(fieldType)
	if convertedValue.IsValid() && convertedValue.Type().AssignableTo(fieldType) {
		field.Set(convertedValue)
		return nil
	}

	return exception.New("Couldnt set field %s.%s", targetType.Name(), fieldName)
}
