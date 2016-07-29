package log4go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexPathURI(t *testing.T) {
	uri := "/api/v1/doctors/1231231aaa/clinis/34234242/timetables/2229237498237489237849372"
	ouput := "_api_v1_doctors_id_clinis_id_timetables_id"
	tOutput := ReplaceIntegerPathParameters(uri, PathParamIntegerRegex, LastPathParamIntegerRegex)
	t.Log(tOutput)
	assert.Equal(t, ouput, tOutput)

	uri = "/doctors/clinic/doctors/1343"
	ouput = "_doctors_clinic_doctors_id"
	tOutput = ReplaceIntegerPathParameters(uri, PathParamIntegerRegex, LastPathParamIntegerRegex)
	assert.Equal(t, ouput, tOutput)
}
