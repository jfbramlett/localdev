package mockgen

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/bxcodec/faker"
	"google.golang.org/protobuf/encoding/protojson"
)

var registrar = &serviceRegistrar{services: make([]interface{}, 0)}

type serviceRegistrar struct {
	services []interface{}
	lock     sync.Mutex
}

func (s *serviceRegistrar) AddService(service interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.services = append(s.services, service)
}

func NewMockGenerator() *MockGen {
	mocks := make(map[string]proto.Message)

	for _, service := range registrar.services {
		newMocks := protosFromService(reflect.ValueOf(service))
		for k, v := range newMocks {
			mocks[k] = v
		}
	}

	return &MockGen{
		mocks: mocks,
	}
}

func protosFromService(service reflect.Value) map[string]proto.Message {
	result := make(map[string]proto.Message)

	for i := 0; i < service.NumMethod(); i++ {
		meth := service.Method(i)
		methType := meth.Type()
		if methType.NumIn() != 2 && methType.NumOut() != 2 {
			continue
		}
		argType := methType.In(1)
		input := reflect.New(argType.Elem())
		inputVal, ok := input.Interface().(proto.Message)
		if !ok {
			continue
		}
		result[string(inputVal.ProtoReflect().Type().Descriptor().FullName())] = inputVal

		argType = methType.Out(0)
		output := reflect.New(argType.Elem())
		outputVal, ok := output.Interface().(proto.Message)
		if !ok {
			continue
		}
		result[string(outputVal.ProtoReflect().Type().Descriptor().FullName())] = outputVal

	}

	return result
}

type MockGen struct {
	mocks map[string]proto.Message
}

func (m *MockGen) ListTypes() []string {
	result := make([]string, 0, len(m.mocks))
	for k := range m.mocks {
		result = append(result, k)
	}

	return result
}

func (m *MockGen) GenerateMock(typ string) (string, error) {
	mockType, found := m.mocks[typ]
	if !found {
		return "", errors.New("not found")
	}

	if err := faker.FakeData(mockType); err != nil {
		return "", err
	}

	json := protojson.Format(mockType)
	fmt.Println(json)

	return json, nil
}
