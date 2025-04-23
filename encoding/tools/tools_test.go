package tools

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/triple/pkg/common/constant"
	"github.com/dubbogo/triple/pkg/config"
)

func TestValidate(t *testing.T) {
	opt := config.NewTripleOption()
	opt.Validate()
	assert.Equal(t, constant.DefaultTimeout, opt.Timeout)
	assert.Equal(t, uint32(constant.DefaultHttp2ControllerReadBufferSize), opt.BufferSize)
	assert.Equal(t, constant.DefaultListeningAddress, opt.Location)
	assert.Equal(t, constant.TRIPLE, opt.Protocol)
	assert.Equal(t, constant.PBCodecName, opt.CodecType)
}
