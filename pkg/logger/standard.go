package logger

//var (
//	std = New()
//)
//
//func StandardLogger() *Logger {
//	return std
//}
//
//func Debugf(out io.Writer, format string, args ...interface{}) (int, error) {
//	return std.Debugf(out, format, args...)
//}
//
//func Infof(out io.Writer, format string, args ...interface{}) (int, error) {
//	return std.Infof(out, format, args...)
//}
//
//func Printf(out io.Writer, format string, args ...interface{}) (int, error) {
//	return std.Printf(out, format, args...)
//}
//
//func Warnf(out io.Writer, format string, args ...interface{}) (int, error) {
//	return std.Warnf(out, format, args...)
//}
//
//func Errorf(out io.Writer, format string, args ...interface{}) (int, error) {
//	return std.Errorf(out, format, args...)
//}
//
//func Debug(out io.Writer, args ...interface{}) (int, error) {
//	return std.Debug(out, args...)
//}
//
//func Info(out io.Writer, args ...interface{}) (int, error) {
//	return std.Info(out, args...)
//}
//
//func Print(out io.Writer, args ...interface{}) (int, error) {
//	return std.Print(out, args...)
//}
//
//func Warn(out io.Writer, args ...interface{}) (int, error) {
//	return std.Warn(out, args...)
//}
//
//func Error(out io.Writer, args ...interface{}) (int, error) {
//	return std.Error(out, args...)
//}
//
//func Debugln(out io.Writer, args ...interface{}) (int, error) {
//	return std.Debugln(out, args...)
//}
//
//func Infoln(out io.Writer, args ...interface{}) (int, error) {
//	return std.Infoln(out, args...)
//}
//
//func Println(out io.Writer, args ...interface{}) (int, error) {
//	return std.Println(out, args...)
//}
//
//func Warnln(out io.Writer, args ...interface{}) (int, error) {
//	return std.Warnln(out, args...)
//}
//
//func Errorln(out io.Writer, args ...interface{}) (int, error) {
//	return std.Errorln(out, args...)
//}
//
//func GetLevel() Level {
//	return std.GetLevel()
//}
//
//func WithLevel(level Level) *Logger {
//	return std.WithLevel(level)
//}
//
//func IsLevelEnabled(level Level) bool {
//	return std.IsLevelEnabled(level)
//}
