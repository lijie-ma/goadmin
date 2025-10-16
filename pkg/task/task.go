package task

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Service 统一接口
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
}

// ServiceManager 统一管理多个服务
type ServiceManager struct {
	services []Service
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	errChan  chan error
}

// NewServiceManager 创建
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		errChan: make(chan error, 1),
	}
}

func (sm *ServiceManager) AddService(s ...Service) {
	sm.services = append(sm.services, s...)
}

func (sm *ServiceManager) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	sm.cancel = cancel

	// 启动所有服务
	for _, s := range sm.services {
		sm.wg.Add(1)
		go sm.runService(ctx, s)
	}

	// 信号监听
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		log.Printf("[Manager] received signal: %v", sig)
	case err := <-sm.errChan:
		log.Printf("[Manager] service failed: %v", err)
	}

	sm.stopAll()
}

func (sm *ServiceManager) runService(ctx context.Context, s Service) {
	defer sm.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			sm.errChan <- &ServiceError{Name: s.Name(), Err: r}
		}
	}()

	log.Printf("[Manager] starting service: %s", s.Name())
	if err := s.Start(ctx); err != nil {
		sm.errChan <- &ServiceError{Name: s.Name(), Err: err}
	}
	log.Printf("[Manager] service exited: %s", s.Name())
}

func (sm *ServiceManager) stopAll() {
	log.Println("[Manager] stopping all services...")
	sm.cancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, s := range sm.services {
		log.Printf("[Manager] stopping service: %s", s.Name())
		if err := s.Stop(ctx); err != nil {
			log.Printf("[Manager] error stopping %s: %v", s.Name(), err)
		}
	}

	sm.wg.Wait()
	log.Println("[Manager] all services stopped.")
}

// 错误包装
type ServiceError struct {
	Name string
	Err  any
}

func (e *ServiceError) Error() string {
	return "[" + e.Name + "] crashed: " + e.toString()
}

func (e *ServiceError) toString() string {
	switch v := e.Err.(type) {
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}
