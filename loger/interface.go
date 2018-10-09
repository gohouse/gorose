package loger

type LogHandler interface {
	File(filepath string) *LogHandler
	Write(string)
}
