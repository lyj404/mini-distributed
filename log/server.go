package log

import (
	"io"
	stlog "log" // 将系统的日志取别名为stlog
	"net/http"
	"os"
)

// 使用系统日志定义一个全局变量log
var log *stlog.Logger

// string类型的别名
type fileLog string

// 自定义日志类型
// 实现 io.Writer 接口中的 Write() 方法 将日志写入文件中
// 它打开文件以进行写入，如果文件不存在，则创建文件。
// 参数:
//   data []byte: 需要写入文件的日志数据。
// 返回值:
//   int: 写入文件的字节数。
//   error: 如果在写入过程中发生错误，则返回错误信息；否则返回nil。
func (fl fileLog) Write(data []byte) (int, error) {
    // 打开文件以进行写入。如果文件不存在，则创建文件。
    // 使用 O_CREATE 标志确保文件不存在时可以创建，
    // O_WRONLY 表明只进行写操作，O_APPEND 确保新数据追加到文件末尾。
    f, err := os.OpenFile(string(fl), os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0600)
    if err != nil {
        // 如果打开文件失败，则返回错误。
        return 0, err
    }
    // 确保在函数返回前关闭文件，以避免资源泄漏。
    defer f.Close()
    
    // 实际写入数据到文件。
    return f.Write(data)
}

// Run 函数初始化日志系统，将日志输出到指定的destination路径。
// 它使用stlog包中的New函数创建一个日志器，将文件日志记录器与"go"格式和标准标志LstdFlags组合。
// destination: 日志文件的目标路径。
func Run(destination string) {
	log = stlog.New(fileLog(destination), "[go]: ", stlog.LstdFlags)
}

// RegisterHandlers 向HTTP服务器注册处理日志请求的处理器。
// 该函数设置了一个路由处理器，用于接收和处理发往/log路径的HTTP请求。
func RegisterHandlers() {
    // 为/log路径注册一个处理函数，用于处理所有发送到该路径的HTTP请求。
    http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
        // 根据请求的方法进行不同的处理。
        switch r.Method {
        case http.MethodPost:
            // 读取请求体中的消息。
            msg, err := io.ReadAll(r.Body)
            // 如果读取请求体发生错误或消息为空，则返回400 Bad Request。
            if err != nil || len(msg) == 0 {
                w.WriteHeader(http.StatusBadRequest)
                return
            }
            // 将读取到的消息写入日志系统。
            write(string(msg))
        default:
            // 如果请求方法不是POST，则返回405 Method Not Allowed。
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }
    })
}

// write 函数记录给定的消息。
// 它使用日志记录功能来输出消息，以便于问题追踪和日志审计。
// 参数:
//   message - 需要记录的消息字符串。
func write(message string) {
    // 将日志写到目标位置
	log.Printf("%v\n", message)
}