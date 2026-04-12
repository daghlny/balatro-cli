package i18n

var current = "zh-CN"
var locales = map[string]map[string]string{}

// Register adds a locale's string table.
func Register(locale string, strings map[string]string) {
	locales[locale] = strings
}

// SetLocale switches the active locale.
func SetLocale(locale string) {
	if _, ok := locales[locale]; ok {
		current = locale
	}
}

// T returns the localized string for a key. Falls back to the key itself.
func T(key string) string {
	if table, ok := locales[current]; ok {
		if val, ok := table[key]; ok {
			return val
		}
	}
	return key
}

// CurrentLocale returns the active locale code.
func CurrentLocale() string {
	return current
}
