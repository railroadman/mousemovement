package main

import (
	"fmt"
	hook "github.com/robotn/gohook"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

var mouseMovementEnabled int32 = 0
var moveRight, moveDown bool = true, true

func moveMouseSlowly() {
	for {
		if atomic.LoadInt32(&mouseMovementEnabled) == 1 {
			x, y := robotgo.GetMousePos()
			screenWidth, screenHeight := robotgo.GetScreenSize()

			if moveRight {
				x += 10
				if x >= screenWidth {
					moveRight = false
				}
			} else {
				x -= 10
				if x <= 0 {
					moveRight = true
				}
			}

			if moveDown {
				y += 10
				if y >= screenHeight {
					moveDown = false
				}
			} else {
				y -= 10
				if y <= 0 {
					moveDown = true
				}
			}

			robotgo.MoveMouseSmooth(x, y, 0.5, 0.5)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func toggleMouseMovement() {
	if atomic.LoadInt32(&mouseMovementEnabled) == 1 {
		atomic.StoreInt32(&mouseMovementEnabled, 0)
		fmt.Println("Движение мышки отключено")
	} else {
		atomic.StoreInt32(&mouseMovementEnabled, 1)
		fmt.Println("Мышка начинает двигаться")
	}
}

func main() {
	// Инициализация глобального перехвата клавиш с использованием gohook
	fmt.Println("Программа запущена. Нажмите Alt+A+0 для управления движением мышки.")

	// Запуск потока для перемещения мышки
	go moveMouseSlowly()

	// Установка глобальных горячих клавиш
	hook.Register(hook.KeyDown, []string{"alt", "a", "0"}, func(e hook.Event) {
		toggleMouseMovement()
	})

	// Запуск прослушивания событий клавиатуры в отдельном потоке
	go func() {
		s := hook.Start()
		<-hook.Process(s)
	}()

	// Обработка завершения программы
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Завершение программы.")
}
