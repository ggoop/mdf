package glog

// Default is the package-level ready-to-use logger,
// level had set to "info", is changeable.
var Default = createLogger()

// Reset re-sets the default logger to an empty one.
func Reset() {
	Default = createLogger()
}

// Print prints a log message without levels and colors.
func Print(v ...interface{}) {
	Default.Print(v...)
}
func CheckAndPrintError(flag string, err error) {
	Default.CheckAndPrintError(flag, err)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(msg string, fields ...Field) {
	Default.Fatal(msg, fields...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(msg interface{}, fields ...Field) error {
	return Default.Error(msg, fields...)

}

// Warn will print when logger's Level is warn, info or debug.
func Warn(msg string, fields ...Field) {
	Default.Warn(msg, fields...)
}

// Info will print when logger's Level is info or debug.
func Info(msg string, fields ...Field) {
	Default.Info(msg, fields...)
}

// Debug will print when logger's Level is debug.
func Debug(msg string, fields ...Field) {
	Default.Debug(msg, fields...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	Default.Fatalf(format, args...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) error {
	return Default.Errorf(format, args...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}

//w interfce

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Errorw(msg string, keysAndValues ...interface{}) {
	Default.Errorw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	Default.Warnw(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	Default.Infow(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Debugw(msg string, keysAndValues ...interface{}) {
	Default.Debugw(msg, keysAndValues...)
}
