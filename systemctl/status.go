package systemctl

import (
	"bufio"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/digiconvent/install_on_debian/utils"
)

type ActiveState string

const (
	ActiveStateActive       ActiveState = "active"
	ActiveStateInactive     ActiveState = "inactive"
	ActiveStateActivating   ActiveState = "activating"
	ActiveStateDeactivating ActiveState = "deactivating"
	ActiveStateFailed       ActiveState = "failed"
)

type LoadState string

const (
	LoadStateLoaded     LoadState = "loaded"
	LoadStateNotFound   LoadState = "not-found"
	LoadStateBadSetting LoadState = "bad-setting"
	LoadStateError      LoadState = "error"
	LoadStateMasked     LoadState = "masked"
)

type SubState string

const (
	SubStateRunning   SubState = "running"
	SubStateExited    SubState = "exited"
	SubStateDead      SubState = "dead"
	SubStateStartPre  SubState = "start-pre"
	SubStateStart     SubState = "start"
	SubStateStartPost SubState = "start-post"
)

type UnitFileState string

const (
	UnitFileStateEnabled        UnitFileState = "enabled"
	UnitFileStateEnabledRuntime UnitFileState = "enabled-runtime"
	UnitFileStateLinked         UnitFileState = "linked"
	UnitFileStateStatic         UnitFileState = "static"
	UnitFileStateDisabled       UnitFileState = "disabled"
	UnitFileStateIndirect       UnitFileState = "indirect"
	UnitFileStateGenerated      UnitFileState = "generated"
	UnitFileStateTransient      UnitFileState = "transient"
	UnitFileStateMasked         UnitFileState = "masked"
	UnitFileStateInvalid        UnitFileState = "invalid"
)

type ServiceStatus struct {
	ActiveState   ActiveState
	LoadState     LoadState
	SubState      SubState
	Result        string
	ExitCode      int
	MainPID       int
	MemoryCurrent uint64
	Loaded        bool
	UnitFileState string
}

func (s *SystemCtl) refreshStatus() (*ServiceStatus, error) {
	sStatus := s.status
	output, err := utils.Execute("systemctl show " + s.serviceName + ".service")
	if err != nil {
		return nil, errors.New(output + ": " + err.Error())
	}
	v := reflect.ValueOf(sStatus).Elem()
	t := v.Type()

	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		fieldMap[strings.ToLower(t.Field(i).Name)] = i
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "="); idx > 0 {
			key, val := line[:idx], line[idx+1:]
			if fieldIndex, exists := fieldMap[strings.ToLower(key)]; exists {
				field := v.Field(fieldIndex)
				switch field.Kind() {
				case reflect.String:
					field.SetString(val)
				case reflect.Bool:
					field.SetBool(val == "true")
				case reflect.Int, reflect.Int64:
					if intVal, err := strconv.Atoi(val); err == nil {
						field.SetInt(int64(intVal))
					}
				case reflect.Uint, reflect.Uint64:
					if uintVal, err := strconv.ParseUint(val, 10, 64); err == nil {
						field.SetUint(uintVal)
					}
				}
			}
		}
	}

	return sStatus, scanner.Err()
}
