// В desktop/structs.go
package desktop

type Desktop struct {
}

func (d *Desktop) Init() {
	// Инициализация для desktop
}

func NewDesktop() *Desktop {
	return &Desktop{}
}
