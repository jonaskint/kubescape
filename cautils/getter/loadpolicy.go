package getter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/armosec/armoapi-go/armotypes"
	"github.com/armosec/opa-utils/reporthandling"
)

// =======================================================================================================================
// ============================================== LoadPolicy =============================================================
// =======================================================================================================================
const DefaultLocalStore = ".kubescape"

// Load policies from a local repository
type LoadPolicy struct {
	filePaths []string
}

func NewLoadPolicy(filePaths []string) *LoadPolicy {
	return &LoadPolicy{
		filePaths: filePaths,
	}
}

// Return control from file
func (lp *LoadPolicy) GetControl(controlName string) (*reporthandling.Control, error) {

	control := &reporthandling.Control{}
	filePath := lp.filePath()
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(f, control); err != nil {
		return control, err
	}
	if controlName != "" && !strings.EqualFold(controlName, control.Name) && !strings.EqualFold(controlName, control.ControlID) {
		framework, err := lp.GetFramework(control.Name)
		if err != nil {
			return nil, fmt.Errorf("control from file not matching")
		} else {
			for _, ctrl := range framework.Controls {
				if strings.EqualFold(ctrl.Name, controlName) || strings.EqualFold(ctrl.ControlID, controlName) {
					control = &ctrl
					break
				}
			}
		}
	}
	return control, err
}

func (lp *LoadPolicy) GetFramework(frameworkName string) (*reporthandling.Framework, error) {
	framework := &reporthandling.Framework{}
	var err error
	for _, filePath := range lp.filePaths {
		f, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(f, framework); err != nil {
			return framework, err
		}
		if strings.EqualFold(frameworkName, framework.Name) {
			break
		}
	}
	if frameworkName != "" && !strings.EqualFold(frameworkName, framework.Name) {

		return nil, fmt.Errorf("framework from file not matching")
	}
	return framework, err
}

func (lp *LoadPolicy) GetExceptions(customerGUID, clusterName string) ([]armotypes.PostureExceptionPolicy, error) {
	filePath := lp.filePath()
	exception := []armotypes.PostureExceptionPolicy{}
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(f, &exception)
	return exception, err
}

func (lp *LoadPolicy) GetControlsInputs(customerGUID, clusterName string) (map[string][]string, error) {
	filePath := lp.filePath()
	accountConfig := &armotypes.CustomerConfig{}
	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(f, &accountConfig); err == nil {
		return accountConfig.Settings.PostureControlInputs, nil
	}
	return nil, err
}

// temporary support for a list of files
func (lp *LoadPolicy) filePath() string {
	if len(lp.filePaths) > 0 {
		return lp.filePaths[0]
	}
	return ""
}
